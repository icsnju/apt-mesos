package impl

import (
	"time"

	log "github.com/Sirupsen/logrus"
	client "github.com/google/cadvisor/client/v2"
	info "github.com/google/cadvisor/info/v2"
	comm "github.com/icsnju/apt-mesos/communication"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	scheduler "github.com/icsnju/apt-mesos/scheduler/impl"
)

// Update node information by MESOS-OFFERS
// Only used when framework registered and get first offers
func (core *Core) updateNodesByOffer(offers []*mesosproto.Offer) {
	for _, offer := range offers {
		// update resources
		resources := make(map[string]*mesosproto.Resource)
		for _, resource := range offer.GetResources() {
			resources[resource.GetName()] = resource
		}

		// if it is a new node
		slaveID := offer.GetSlaveId().GetValue()

		if exists := core.ExistsNode(slaveID); !exists {
			node := &registry.Node{
				ID:             slaveID,
				Hostname:       offer.GetHostname(),
				LastUpdateTime: time.Now().Unix(),
				Resources:      resources,
			}
			core.RegisterNode(slaveID, node)
		}
	}
}

// Update node information by tasks
func (core *Core) updateNodeByTask(id string, task *registry.Task) {
	node, _ := core.GetNode(id)
	for _, resource := range task.Resources {
		// Update scalar
		if resource.GetType().String() == "SCALAR" {
			newScalar := node.OfferedResources[resource.GetName()].GetScalar().GetValue() - resource.GetScalar().GetValue()
			node.OfferedResources[resource.GetName()].Scalar.Value = &newScalar
		} else if resource.GetType().String() == "RANGES" {
			// Update ranges
			node.OfferedResources[resource.GetName()].Ranges = scheduler.RangeUsedUpdate(resource.GetRanges(), node.OfferedResources[resource.GetName()].GetRanges())
		}
	}
	node.Tasks = append(node.Tasks, task)
	node.LastUpdateTime = time.Now().Unix()
}

// Update node information by Metrics
func (core *Core) updateNodesByMetrics() {
	metrics, err := core.FetchMetricData()
	if err != nil {
		log.Errorf("Fetch metric data error: %v", err)
		return
	}

	// get slavePID of task
	for _, slave := range metrics.Slaves {
		for _, node := range core.GetAllNodes() {
			// find for master node
			if node.Hostname == metrics.Hostname {
				node.IsMaster = true
			}

			if node.ID == slave.ID {
				node.IsSlave = true
				// update node pid for the first time
				if node.PID == "" {
					node.PID = slave.PID
				}

				// update node host for the first time
				if node.Host == "" {
					upid, err := comm.Parse(slave.PID)
					if err != nil {
						continue
					}
					node.Host = upid.Host
				}

				offeredResources := scheduler.BuildResourcesFromMap(slave.OfferedResources)
				node.OfferedResources = offeredResources

				log.Debugf("Node %v status updated: %v", node.Hostname, node.OfferedResources["ports"])
				// node.CPURegistered = slave.Resources.Cpus
				// node.MemoryRegistered = uint64(slave.Resources.Mem)
				//find for slave node
				node.LastUpdateTime = time.Now().Unix()
			}
		}
	}
}

// Update node information by CAdvisor
func (core *Core) updateNodesByCAdvisor() {
	for _, node := range core.GetAllNodes() {
		if node.Host != "" {
			client, err := client.NewClient("http://" + node.Host + ":" + core.GetAgentLisenPort())
			if err != nil {
				log.Errorf("Cannot connect to cadvisor agent: %v", err)
				continue
			}

			// Fetch software versions and hardware information
			// One node fetches just one time
			if !node.MachineInfoFetched {
				attributes, err := client.Attributes()
				if err != nil {
					log.Errorf("Fetch machine info failed: %v", err)
				}
				node.NumCores = attributes.NumCores
				node.KernelVersion = attributes.KernelVersion
				node.CPUFrequency = attributes.CpuFrequency
				node.ContainerOsVersion = attributes.ContainerOsVersion
				node.DockerVersion = attributes.DockerVersion
				node.MemoryCapacity = attributes.MemoryCapacity
				node.MachineInfoFetched = true
			}

			request := info.RequestOptions{
				IdType:    "docker",
				Count:     5,
				Recursive: true,
			}
			containerInfo, err := client.Stats("", &request)
			if err != nil {
				log.Errorf("Fetch container info failed: %v", err)
			}

			node.Containers = make([]string, 0, len(containerInfo))
			for container, info := range containerInfo {
				node.Containers = append(node.Containers, container)
				if len(info.Stats) > 1 {
					lastStats := info.Stats[len(info.Stats)-2]
					currentStats := info.Stats[len(info.Stats)-1]

					// ms -> ns.
					timeInterval := float64((currentStats.Timestamp.Unix() - lastStats.Timestamp.Unix()) * 1000000)
					node.CPUUsage = float64(currentStats.Cpu.Usage.Total-lastStats.Cpu.Usage.Total) / timeInterval
					node.MemoryUsage = currentStats.Memory.Usage
				}
			}
		}
	}
}
