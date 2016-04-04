package impl

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	comm "github.com/icsnju/apt-mesos/communication"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ErrTaskNotExists defined errors
var (
	ErrTaskNotExists             = errors.New("Specific task not exist")
	ErrBasicResourceNotSatisfied = errors.New("Cpus and mem must be required in resources")
)

type filter func(task *registry.Task) bool

// AddTask is called when user submit a task and add the task to the registry
func (core *Core) AddTask(id string, task *registry.Task) error {
	if err := core.tasks.Add(id, task); err != nil {
		return err
	}
	return nil
}

// GetAllTasks Return all the tasks in registry
func (core *Core) GetAllTasks() []*registry.Task {
	rawList := core.tasks.List()
	tasks := make([]*registry.Task, len(rawList))

	for i, v := range rawList {
		tasks[i] = v.(*registry.Task)
	}
	return tasks
}

// GetTask : Get the task that specified id
func (core *Core) GetTask(id string) (*registry.Task, error) {
	if task := core.tasks.Get(id); task != nil {
		return task.(*registry.Task), nil
	}
	return nil, ErrTaskNotExists
}

// DeleteTask give an id and delete the task
func (core *Core) DeleteTask(id string) error {
	if err := core.tasks.Delete(id); err != nil {
		return err
	}
	return nil
}

// UpdateTask update task info
func (core *Core) UpdateTask(id string, task *registry.Task) error {
	return core.tasks.Update(id, task)
}

// KillTask kill the task
func (core *Core) KillTask(id string) error {
	if task := core.tasks.Get(id); task == nil {
		return ErrTaskNotExists
	}
	frameworkID := core.frameworkInfo.GetId().GetValue()
	message := &mesosproto.KillTaskMessage{
		FrameworkId: &mesosproto.FrameworkID{Value: &frameworkID},
		TaskId:      &mesosproto.TaskID{Value: &id},
	}
	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	return comm.SendMessageToMesos(core.coreUPID, messagePackage)
}

// FilterTask filter task by func
func (core *Core) FilterTask(choose filter) []*registry.Task {
	var result []*registry.Task
	for _, task := range core.tasks.List() {
		if choose(task.(*registry.Task)) {
			result = append(result, task.(*registry.Task))
		}
	}

	return result
}

func (core *Core) generateResource(task *registry.Task) {
	var resources = []*mesosproto.Resource{}
	resources = append(resources, &mesosproto.Resource{
		Name:   proto.String("cpus"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Cpus)},
	})
	resources = append(resources, &mesosproto.Resource{
		Name:   proto.String("mem"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Mem)},
	})
	resources = append(resources, &mesosproto.Resource{
		Name:   proto.String("disk"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Disk)},
	})
	resources = append(resources, core.MergePorts(task.Ports))
	task.Resources = resources
}

func (core *Core) MergePorts(ports []*registry.Port) *mesosproto.Resource {
	var used = [65536]bool{false}
	for _, port := range ports {
		used[port.HostPort] = true
	}
	var rangesPort = &mesosproto.Value_Ranges{}
	for i := 1; i < 65536; i++ {
		if used[i] {
			Begin := uint64(i)
			var rangePort = &mesosproto.Value_Range{
				Begin: &Begin,
			}
			for ; used[i]; i++ {
			}
			End := uint64(i - 1)
			rangePort.End = &End
			rangesPort.Range = append(rangesPort.Range, rangePort)
		}
	}
	return &mesosproto.Resource{
		Name:   proto.String("ports"),
		Type:   mesosproto.Value_RANGES.Enum(),
		Ranges: rangesPort,
	}
}

// GetUnScheduledTask return all un schedulered task
func (core *Core) GetUnScheduledTask() []*registry.Task {
	return core.FilterTask(func(task *registry.Task) bool {
		return task.State == ""
	})
}

func (core *Core) updateTasksByMetrics(metrics *registry.MetricsData) {
	// get executorID, slaveID, slaveHostname of task
	for _, framework := range metrics.Frameworks {
		if framework.ID == core.frameworkInfo.GetId().GetValue() {
			// fetch running task
			for _, task := range framework.Tasks {
				taskInfo, _ := core.GetTask(task.ID)
				// get executorID
				if taskInfo.ExecutorID == "" {
					taskInfo.ExecutorID = task.ExecutorID
					if task.ExecutorID == "" {
						taskInfo.ExecutorID = taskInfo.ID
					}
				}
				// get slaveID
				if taskInfo.SlaveID == "" {
					taskInfo.SlaveID = task.SlaveID
				}
				// get slave hostname
				if taskInfo.SlaveHostname == "" {
					node, _ := core.GetNode(task.SlaveID)
					taskInfo.SlaveHostname = node.Hostname
				}
				// update task's state
				taskInfo.State = task.State
				core.UpdateTask(taskInfo.ID, taskInfo)
			}

			// fetch completed tasks
			for _, task := range framework.CompletedTasks {
				taskInfo, _ := core.GetTask(task.ID)
				// get executorID
				if taskInfo.ExecutorID == "" {
					taskInfo.ExecutorID = task.ExecutorID
					if task.ExecutorID == "" {
						taskInfo.ExecutorID = taskInfo.ID
					}
				}
				// update state
				taskInfo.State = task.State
				// get slaveID
				if taskInfo.SlaveID == "" {
					taskInfo.SlaveID = task.SlaveID
				}
				// get slave hostname
				if taskInfo.SlaveHostname == "" {
					node, _ := core.GetNode(task.SlaveID)
					taskInfo.SlaveHostname = node.Hostname
				}
				core.UpdateTask(taskInfo.ID, taskInfo)
			}
		}
	}

	// get slavePID of task
	for _, slave := range metrics.Slaves {
		for _, task := range core.GetAllTasks() {
			if task.SlavePID != "" {
				continue
			} else {
				if task.SlaveID == slave.ID {
					task.SlavePID = slave.PID
					upid, err := comm.Parse(task.SlavePID)
					if err != nil {
						continue
					}
					task.SlaveHost = upid.Host
				}
			}
		}
	}

	// get directory of task
	for _, task := range core.GetAllTasks() {
		if task.Directory != "" || task.SlavePID == "" {
			continue
		}
		directory, err := core.getTaskDirectory(task.SlavePID, task.ExecutorID)
		if err != nil {
			log.Errorf("Cannot read task directory of %v: %v", task.ID, err)
		}
		task.Directory = directory
	}
}

// getTaskDirectory return the directory of a task
func (core *Core) getTaskDirectory(slavePID, executorID string) (string, error) {
	resp, err := http.Get("http://" + slavePID + "/state.json")
	if err != nil {
		return "", err
	}

	data := struct {
		Frameworks []struct {
			Executors []struct {
				ID        string
				Directory string
			}
			CompletedExecutors []struct {
				ID        string
				Directory string
			} `json:"completed_executors"`
			ID string
		}
		CompletedFrameworks []struct {
			CompletedExecutors []struct {
				ID        string
				Directory string
			} `json:"completed_executors"`
			ID string
		} `json:"completed_frameworks"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	for _, framework := range data.Frameworks {
		if framework.ID != core.frameworkInfo.GetId().GetValue() {
			continue
		}
		for _, executor := range framework.Executors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
	}

	for _, framework := range data.CompletedFrameworks {
		if framework.ID != core.frameworkInfo.GetId().GetValue() {
			continue
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
	}

	return "", nil
}
