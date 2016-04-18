package registry

import "github.com/icsnju/apt-mesos/mesosproto"

// Node is one node of cluster
type Node struct {
	ID               string                          `json:"id"`
	Host             string                          `json:"host"`
	Hostname         string                          `json:"hostname"`
	IsMaster         bool                            `json:"is_master"`
	IsSlave          bool                            `json:"is_slave"`
	Resources        map[string]*mesosproto.Resource `json:"resources"`
	OfferedResources map[string]*mesosproto.Resource `json:"offered_resources"`
	LastUpdateTime   int64                           `json:"last_update_time"`

	MachineInfoFetched bool
	NumCores           int    `json:"num_cores"`
	CPUFrequency       uint64 `json:"cpu_frequency_khz"`
	MemoryCapacity     uint64 `json:"memory_capacity"`
	KernelVersion      string `json:"kernel_version"`
	ContainerOsVersion string `json:"container_os_version"`
	DockerVersion      string `json:"docker_version"`
	DockerDaemonHealth string `json:"docker_daemon_health"`

	CPUUsage    float64  `json:"cpu_usage"`
	MemoryUsage uint64   `json:"memory_usage"`
	Containers  []string `json:"containers"`

	Tasks []*Task `json:"tasks"`
	PID   string
}

var (
	DockerDaemonUp   = "up"
	DockerDaemonDown = "down"
)

func (node *Node) GetTasks() []*Task {
	return node.Tasks
}

func (node *Node) GetSLATasks() []*Task {
	var result []*Task
	for _, task := range node.Tasks {
		if task.SLA != "" {
			result = append(result, task)
		}
	}
	return result
}
