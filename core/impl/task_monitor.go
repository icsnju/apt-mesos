package impl

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	client "github.com/google/cadvisor/client/v2"
	info "github.com/google/cadvisor/info/v2"
	comm "github.com/icsnju/apt-mesos/communication"
	"github.com/icsnju/apt-mesos/registry"
)

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
			}

			// fetch completed tasks
			for _, task := range framework.CompletedTasks {
				taskInfo, err := core.GetTask(task.ID)
				if err != nil {
					continue
				}
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
		task.LastUpdateTime = time.Now().Unix()
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

func (core *Core) updateTaskByDockerInfo(task *registry.Task, dockerInspectOutput []byte) {
	//	glog.Infof("Docker Inspect: %s", dockerInspectOutput)
	var dockerTasks []*registry.DockerTask
	err := json.Unmarshal(dockerInspectOutput, &dockerTasks)
	if err != nil {
		log.Errorf("UpdateTaskWithDockerInfo error: %v\n", err)
		return
	}
	task.DockerID = dockerTasks[0].DockerID
	task.DockerName = dockerTasks[0].DockerName[1:]

	var dockerState *registry.DockerState
	err = json.Unmarshal(dockerTasks[0].DockerState, &dockerState)
	if err != nil {
		log.Errorf("UpdateTaskWithDockerInfo error: %v\n", err)
		return
	}

	task.ProcessID = dockerState.Pid
	task.LastUpdateTime = time.Now().Unix()
}

func (core *Core) updateTasksByCAdvisor() {
	for _, task := range core.GetAllTasks() {
		if task.State == "TASK_RUNNING" && task.DockerID != "" && task.SlaveID != "" {
			node, _ := core.GetNode(task.SlaveID)
			client, err := client.NewClient("http://" + node.Host + ":" + core.GetAgentLisenPort())
			if err != nil {
				log.Errorf("Cannot connect to cadvisor agent: %v", err)
				continue
			}

			request := info.RequestOptions{
				IdType: "docker",
				Count:  15,
			}
			containerInfo, err := client.Stats(task.DockerID, &request)
			if err != nil {
				log.Errorf("Fetch container info failed: %v", err)
			}

			var cpuStats []*registry.Usage
			var memoryStats []*registry.Usage
			for _, containerInfo := range containerInfo {
				for _, containerStats := range containerInfo.Stats {
					cpuStats = append(cpuStats, &registry.Usage{
						Total:     containerStats.Cpu.Usage.Total,
						Timestamp: containerStats.Timestamp,
					})
					memoryStats = append(memoryStats, &registry.Usage{
						Total:     containerStats.Memory.Usage,
						Timestamp: containerStats.Timestamp,
					})
				}
			}
			task.CPUUsage = cpuStats
			task.MemoryUsage = memoryStats
		}
		task.LastUpdateTime = time.Now().Unix()
	}
}
