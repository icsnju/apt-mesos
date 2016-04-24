package registry

import (
	"encoding/json"
	"strings"
	"time"

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
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	Cpus       float64                 `json:"cpus,string"`
	Mem        float64                 `json:"mem,string"`
	Disk       float64                 `json:"disk,string"`
	Resources  []*mesosproto.Resource  `json:"resources,omitempty"`
	Attributes []*mesosproto.Attribute `json:attributes,omitempty`
	SLA        string                  `json:"sla"`

	// Monitoring
	State          string   `json:"state"`
	CPUUsage       []*Usage `json:"cpu_usage"`
	MemoryUsage    []*Usage `json:"memory_usage"`
	LastUpdateTime int64    `json:"last_update_time"`

	// Docker settings
	Command     string       `json:"cmd"`
	Arguments   []string     `json:"arguments,omitempty"`
	DockerImage string       `json:"docker_image"`
	Volumes     []*Volume    `json:"volumes,omitempty"`
	Ports       []*Port      `json:"port_mappings,omitempty"`
	NetworkMode string       `json:"network_mode"`
	Privileged  bool         `json:"privileged"`
	Parameters  []*Parameter `json:"parameters,omitempty"`

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
	CreateTime    int64  `json:"create_time"`
	RunTime       int64  `json:"run_time"`

	TaskInfo *mesosproto.TaskInfo
	Type     TaskType `enum=TaskType,json:"type,omitempty"`
	JobID    string   `json:"job_id"`
	Scale    int      `json:"scale"`
}

// DockerTask is docker information struct
type DockerTask struct {
	DockerID    string          `json:"Id"`
	DockerName  string          `json:"Name"`
	DockerState json.RawMessage `json:"State"`
}

// Parameter of docker
type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DockerState is the state of docker
type DockerState struct {
	Pid uint `json:"Pid"`
}

type Usage struct {
	Total     uint64    `json:"total"`
	Timestamp time.Time `json:"timestamp"`
}

type TaskType int32

const (
	TaskTypeTest  TaskType = 0
	TaskTypeBuild TaskType = 1
)

func (task *Task) Parse() string {
	if task.JobID == "" {
		return task.ID
	}
	index := strings.LastIndex(task.ID, "-")
	if index < 0 {
		return task.ID
	}
	return task.ID[0:index]
}
