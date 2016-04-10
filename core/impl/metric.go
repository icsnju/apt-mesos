package impl

import (
	"encoding/json"
	"net/http"
	"time"

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

// Monitor fetch mesos state and update task info
func (core *Core) monitor() {
	for {
		// Fetch mesos state and update tasks and nodes
		data, err := core.FetchMetricData()
		if err != nil {
			return
		}
		core.updateTasksByMetrics(data)
		core.updateNodesByMetrics(data)

		// Fetch agent data and update
		core.updateNodesByCAdvisor()
		core.updateTasksByCAdvisor()
		time.Sleep(1 * time.Second)
	}
}
