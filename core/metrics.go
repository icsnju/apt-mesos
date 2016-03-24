package core

import (
	"encoding/json"
	"net/http"

	"github.com/icsnju/apt-mesos/mesosproto"
)

// Metrics provide some useful metrics to frontend
type Metrics struct {
	TotalCpus    float64 `json:"total_cpus"`
	TotalMem     float64 `json:"total_mem"`
	TotalDisk    float64 `json:"total_disk"`
	UsedCpus     float64 `json:"used_cpus"`
	UsedMem      float64 `json:"used_mem"`
	UsedDisk     float64 `json:"used_disk"`
	TaskRunning  int64   `json:"task_running"`
	TaskStaging  int64   `json:"task_staging"`
	TaskFinished int64   `json:"task_finished"`
	TaskKilled   int64   `json:"task_killed"`
	//TODO add customed resources
}

// MetricsData is a struct suit for json data from mesos-master
type MetricsData struct {
	Frameworks []struct {
		Tasks []struct {
			ExecutorID string `json:"executor_id"`
			ID         string
			SlaveID    string `json:"slave_id"`
			Resources  struct {
				Cpus float64
				Mem  float64
				Disk float64
			}
			State string `json:"state"`
		}
		CompletedTasks []struct {
			ExecutorID string `json:"executor_id"`
			ID         string
			SlaveID    string `json:"slave_id"`
			State      string `json:"state"`
		} `json:"completed_tasks"`
		ID string
	}
	CompletedFrameworks []struct {
		CompletedTasks []struct {
			ExecutorID string `json:"executor_id"`
			ID         string
			SlaveID    string `json:"slave_id"`
		} `json:"completed_tasks"`
		ID string
	} `json:"completed_frameworks"`
	Slaves []struct {
		ID        string
		Pid       string
		Hostname  string
		Resources struct {
			Cpus float64
			Mem  float64
			Disk float64
		}
	}
}

// Metrics calculate some needed metrics from raw metric data
// TODO Bad Solution to update task state
func (core *Core) Metrics() (*Metrics, map[string]mesosproto.TaskState, error) {
	states := make(map[string]mesosproto.TaskState)

	data, err := core.GetMetricsData()
	if err != nil {
		return nil, nil, err
	}

	var metrics Metrics

	for _, framework := range data.Frameworks {
		for _, task := range framework.Tasks {
			metrics.UsedMem += task.Resources.Mem
			metrics.UsedCpus += task.Resources.Cpus
			metrics.UsedDisk += task.Resources.Disk
		}

		if framework.ID == core.frameworkInfo.Id.GetValue() {
			for _, task := range framework.Tasks {
				switch task.State {
				case "TASK_RUNNING":
					metrics.TaskRunning++
					states[task.ID] = mesosproto.TaskState_TASK_RUNNING
				case "TASK_STAGING":
					metrics.TaskStaging++
					states[task.ID] = mesosproto.TaskState_TASK_STAGING
				}
			}
			for _, task := range framework.CompletedTasks {
				switch task.State {
				case "TASK_FINISHED":
					metrics.TaskFinished++
					states[task.ID] = mesosproto.TaskState_TASK_FINISHED
				case "TASK_KILLED":
					metrics.TaskKilled++
					states[task.ID] = mesosproto.TaskState_TASK_KILLED
				}
			}
		}
	}

	for _, slave := range data.Slaves {
		metrics.TotalMem += slave.Resources.Mem
		metrics.TotalCpus += slave.Resources.Cpus
		metrics.TotalDisk += slave.Resources.Disk
	}
	return &metrics, states, nil
}

// GetMetricsData connect to mesos-master and get raw metric data,
// and decode json data to struct
func (core *Core) GetMetricsData() (*MetricsData, error) {
	resp, err := http.Get("http://" + core.master + "/master/state.json")
	if err != nil {
		return nil, err
	}

	data := new(MetricsData)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}

// GetSlaveHostname return slave's hostname when given slave's id
func (core *Core) GetSlaveHostname(slaveID string) (string, error) {
	data, err := core.GetMetricsData()
	if err != nil {
		return "", err
	}

	for _, slave := range data.Slaves {
		if slave.ID == slaveID {
			return slave.Hostname, nil
		}
	}

	return "", nil
}

// getSlavePIDAndExecutorID extract slave's pid and executor's id from raw metric data
func (core *Core) getSlavePIDAndExecutorID(taskID string) (string, string, error) {
	data, err := core.GetMetricsData()
	if err != nil {
		return "", "", err
	}

	var (
		executorID string
		slaveID    string
	)

	// Iterates over all frameworks' tasks to find executorID and slaveID
SearchRunningFrameworks:
	for _, framework := range data.Frameworks {
		if framework.ID != *core.frameworkInfo.Id.Value {
			continue
		}
		for _, task := range framework.Tasks {
			if task.ID == taskID {
				executorID = task.ExecutorID
				slaveID = task.SlaveID
				break SearchRunningFrameworks
			}
		}
		for _, task := range framework.CompletedTasks {
			if task.ID == taskID {
				executorID = task.ExecutorID
				slaveID = task.SlaveID
				break SearchRunningFrameworks
			}
		}
	}

SearchCompletedFrameworks:
	for _, framework := range data.CompletedFrameworks {
		if framework.ID != *core.frameworkInfo.Id.Value {
			continue
		}
		for _, task := range framework.CompletedTasks {
			if task.ID == taskID {
				executorID = task.ExecutorID
				slaveID = task.SlaveID
				break SearchCompletedFrameworks
			}
		}
	}

	for _, slave := range data.Slaves {
		if slave.ID == slaveID {
			return slave.Pid, executorID, nil
		}
	}

	return "", "", nil
}
