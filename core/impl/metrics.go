package impl

import (
	"encoding/json"
	"net/http"

	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// GetMetrics calculate some needed metrics from raw metric data
func (core *Core) GetMetrics() (*registry.Metrics, map[string]mesosproto.TaskState, error) {
	states := make(map[string]mesosproto.TaskState)

	data, err := core.GetMetricsData()
	if err != nil {
		return nil, nil, err
	}

	var metrics registry.Metrics

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
func (core *Core) GetMetricsData() (*registry.MetricsData, error) {
	resp, err := http.Get("http://" + core.master + "/master/state.json")
	if err != nil {
		return nil, err
	}

	data := new(registry.MetricsData)
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
