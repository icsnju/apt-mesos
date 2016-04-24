package impl

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/registry"
)

// FetchMetricData connect to mesos-master and get raw metric data,
func (core *Core) FetchMetricData() (*registry.MetricsData, error) {
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

func (core *Core) GetSystemUsage() *registry.Metrics {
	var metrics registry.Metrics
	for _, node := range core.GetAllNodes() {
		for name, resource := range node.Resources {
			if name == "cpus" {
				metrics.FreeCpus += resource.GetScalar().GetValue()
			}
			if name == "mem" {
				metrics.FreeMem += resource.Scalar.GetValue()
			}
			if name == "disk" {
				metrics.FreeDisk += resource.Scalar.GetValue()
			}
		}
	}

	return &metrics
}

func (core *Core) metricMonitor() {
	for {
		// Fetch mesos state and update tasks and nodes
		data, err := core.FetchMetricData()
		if err != nil {
			log.Errorf("Fetch metric data error: %v", err)
			return
		}
		core.updateTasksByMetrics(data)

		// Fetch agent data and update
		core.updateNodesByCAdvisor()
		core.updateTasksByCAdvisor()
		time.Sleep(500 * time.Millisecond)
	}
}

func (core *Core) addFailureMetric(value float32) {
	if len(core.metric.FailureRate) > 60 {
		core.metric.FailureRate = core.metric.FailureRate[1:]
	}
	core.metric.FailureRate = append(core.metric.FailureRate, registry.SystemMetricItem{
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	})
}

func (core *Core) addWaittimeMetric(value float32) {
	if len(core.metric.WaitTime) > 60 {
		core.metric.WaitTime = core.metric.WaitTime[1:]
	}
	core.metric.WaitTime = append(core.metric.WaitTime, registry.SystemMetricItem{
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	})
}

func (core *Core) taskHealthCheck(timeInterval time.Duration) {
	for {
		totalWaitTime := int64(0)
		taskCount := int64(0)
		taskTotal := 0
		taskFailed := 0
		for _, task := range core.GetAllTasks() {
			if task.RunTime > time.Now().UnixNano()-int64(1*time.Minute) {
				totalWaitTime += (task.RunTime - task.CreateTime)
				taskCount++
			}
			if task.State != "TASK_RUNNING" && task.State != "TASK_STAGING" && task.State != "TASK_WAITING" {
				taskCount++
				if task.State != "TASK_FINISHED" {
					taskFailed++
				}
			}
		}
		if taskTotal == 0 {
			core.addFailureMetric(float32(0))
		} else {
			core.addFailureMetric(float32(taskFailed / taskTotal))
		}

		if taskCount == 0 {
			core.addWaittimeMetric(float32(0))
		} else {
			core.addWaittimeMetric(float32(totalWaitTime / taskCount))
		}
		time.Sleep(timeInterval)
	}
}

func (core *Core) GetSystemMetric() *registry.SystemMetric {
	return &core.metric
}

// Monitor fetch mesos state and update task info
func (core *Core) monitor() {
	go core.metricMonitor()
	go core.taskHealthCheck(1 * time.Minute)
}
