package core

import (
	"net/http"
	"encoding/json"
)

type Metrics struct {
	TotalCpus float64 `json:"total_cpus"`
	TotalMem  float64 `json:"total_mem"`
	TotalDisk float64 `json:"total_disk"`
	UsedCpus  float64 `json:"used_cpus"`
	UsedMem   float64 `json:"used_mem"`
	UsedDisk  float64 `json:"used_disk"`
	//TODO add customed resources
}

type MetricsData struct {
	Frameworks []struct {
		Tasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
			Resources  struct {
				Cpus float64
				Mem  float64
				Disk float64
			}
		}
		CompletedTasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
		} `json:"completed_tasks"`
		Id string
	}
	CompletedFrameworks []struct {
		CompletedTasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
		} `json:"completed_tasks"`
		Id string
	} `json:"completed_frameworks"`
	Slaves []struct {
		Id        string
		Pid       string
		Hostname  string
		Resources struct {
			Cpus float64
			Mem  float64
			Disk float64
		}
	}
}

func (core *Core) Metrics() (*Metrics, error) {
	data, err := core.GetMetricsData()
	if err != nil {
		return nil, err
	}

	var metrics Metrics

	for _, framework := range data.Frameworks {
		for _, task := range framework.Tasks {
			metrics.UsedMem += task.Resources.Mem
			metrics.UsedCpus += task.Resources.Cpus
			metrics.UsedDisk += task.Resources.Disk
		}
	}

	for _, slave := range data.Slaves {
		metrics.TotalMem += slave.Resources.Mem
		metrics.TotalCpus += slave.Resources.Cpus
		metrics.TotalDisk += slave.Resources.Disk
	}

	return &metrics, nil
}

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
