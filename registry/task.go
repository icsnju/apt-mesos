package registry

import (
	"github.com/icsnju/apt-mesos/mesosproto"
)

const (
	// SLAOnePerNode is a mode that a task run only one time on a machine
	SLAOnePerNode = "one-per-node"
	// SLASingleton is a mode that a task run only one time on the whole cluster
	SLASingleton = "singleton"

	// NetworkModeBridge means set up docker containers with network mode of bridge
	NetworkModeBridge = "bridge"
	// NetworkModeHost means set up docker containers with network mode of host
	NetworkModeHost = "host"
	// NetworkModeNone means do not use any mode of network
	NetworkModeNone = "none"
)

// Task struct
type Task struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Cpus        float64                `json:"cpus"`
	Mem         float64                `json:"mem"`
	Disk        float64                `json:"disk"`
	DockerImage string                 `json:"docker_image"`
	Command     string                 `json:"cmd"`
	Resources   []*mesosproto.Resource `json:"resources,omitempty"`
	SLA         string                 `json:"sla"`
	Arguments   []string               `json:"arguments,omitempty"`
	State       string                 `json:"state"`
	Volumes     []*Volume              `json:"volumes,omitempty"`
	Ports       []*Port                `json:"port_mappings,omitempty"`
	NetworkMode string                 `json:"network_mode"`

	DockerID      string
	DockerName    string
	SlaveID       string `json:"slave_id"`
	SlaveHost     string `json:"slave_host"`
	SlaveHostname string `json:"slave_hostname"`
	SlavePID      string `json:"slave_pid"`
	ExecutorID    string `json:"executor_id"`
	Directory     string `json:"directory"`
	CreatedTime   int64  `json:"create_time"`
	TaskInfo      *mesosproto.TaskInfo
}

// TestTask returns a task for testing
func TestTask(id string) *Task {
	return &Task{
		ID:          id,
		DockerImage: "ubuntu",
		Command:     "echo `hello sher`",
	}
}
