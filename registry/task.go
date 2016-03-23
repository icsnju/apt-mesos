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
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	DockerImage string                `json:"docker_image"`
	Command     string                `json:"cmd"`
	Cpus        float64               `json:"cpus,string"`
	Disk        float64               `json:"disk,string"`
	Mem         float64               `json:"mem,string"`
	Arguments   []string              `json:"arguments,omitempty"`
	State       *mesosproto.TaskState `json:"state,string"`
	Volumes     []*Volume             `json:"volumes,omitempty"`
	Ports       []*Port               `json:"port_mappings,omitempty"`
	NetworkMode string                `json:"network_mode"`

	DockerID      string
	DockerName    string
	SlaveID       string `json:"slave_id,string"`
	SlaveHostname string `json:"slave_hostname"`
	CreatedTime   int64  `json:"create_time"`
	TaskInfo      *mesosproto.TaskInfo
}

// TestTask returns a task for testing
func TestTask(id string) *Task {
	return &Task{
		ID:          id,
		DockerImage: "ubuntu",
		Command:     "echo `hello sher`",
		Cpus:        0.5,
		Disk:        0,
		Mem:         16,
	}
}
