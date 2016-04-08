package registry

import (
	"encoding/json"

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
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Cpus      float64                `json:"cpus"`
	Mem       float64                `json:"mem"`
	Disk      float64                `json:"disk"`
	Resources []*mesosproto.Resource `json:"resources,omitempty"`
	SLA       string                 `json:"sla"`

	// Monitoring
	State       string  `json:"state"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage uint64  `json:"memory_usage"`

	// Docker settings
	Command     string    `json:"cmd"`
	Arguments   []string  `json:"arguments,omitempty"`
	DockerImage string    `json:"docker_image"`
	Volumes     []*Volume `json:"volumes,omitempty"`
	Ports       []*Port   `json:"port_mappings,omitempty"`
	NetworkMode string    `json:"network_mode"`

	// Docker inspect
	DockerID   string `json:"docker_id"`
	DockerName string `json:"docker_name"`
	ProcessID  uint   `json:"process_id"`

	// Node
	SlaveID       string `json:"slave_id"`
	SlaveHost     string `json:"slave_host"`
	SlaveHostname string `json:"slave_hostname"`
	SlavePID      string `json:"slave_pid"`
	ExecutorID    string `json:"executor_id"`
	Directory     string `json:"directory"`
	CreatedTime   int64  `json:"create_time"`
	TaskInfo      *mesosproto.TaskInfo
}

// DockerTask is docker information struct
type DockerTask struct {
	DockerID    string          `json:"Id"`
	DockerName  string          `json:"Name"`
	DockerState json.RawMessage `json:"State"`
}

// DockerState is the state of docker
type DockerState struct {
	Pid uint `json:"Pid"`
}

// TestTask returns a task for testing
func TestTask(id string) *Task {
	return &Task{
		ID:          id,
		DockerImage: "ubuntu",
		Command:     "echo `hello sher`",
	}
}
