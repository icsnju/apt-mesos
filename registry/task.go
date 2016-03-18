package registry

import (
	"github.com/icsnju/apt-mesos/mesosproto"
)

const (
	SLA_ONE_PER_NODE = "one-per-node"
	SLA_SINGLETON = "singleton"

	NETWORK_MODE_BRIDGE = "bridge"
	NETWORK_MODE_HOST = "host"
	NETWORK_MODE_NONE = "none"
)

type Task struct {
	ID          	string   				`json:"id"`
	Name			string					`json:"name"`
	DockerImage 	string   				`json:"docker_image"`
	Command     	string   				`json:"cmd"`
	Cpus        	float64  				`json:"cpus,string"`
	Disk        	float64  				`json:"disk,string"`
	Mem         	float64  				`json:"mem,string"`
	Arguments     	[]string  				`json:"arguments,omitempty"`
	State         	*mesosproto.TaskState 	`json:"state,string"`
	Volumes         []*Volume    		    `json:"volumes,omitempty"`
	Ports         	[]*Port   				`json:"port_mappings,omitempty"`
	NetworkMode   	string    				`json:"network_mode"`

	DockerID		string
	DockerName		string
	SlaveId       string                `json:"slave_id,string"`
	SlaveHostname string                `json:"slave_hostname"`
	CreatedTime     int64				`json:"create_time"`
	TaskInfo		*mesosproto.TaskInfo
}

func TestTask(id string) *Task {
	return &Task{
		ID: 			id,
		DockerImage: 	"ubuntu",
		Command:		"echo `hello sher`",
		Cpus:			0.5,
		Disk:			0,
		Mem:			16,
	}
}
