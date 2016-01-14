package registry

import (
	"github.com/icsnju/apt-mesos/mesosproto"
)
type Task struct {
	ID          	string   				`json:"id"`
	DockerImage 	string   				`json:"docker_image"`
	Command     	string   				`json:"cmd"`
	Cpus        	float64  				`json:"cpus,string"`
	Disk        	float64  				`json:"disk,string"`
	Mem         	float64  				`json:"mem,string"`
	Arguments     	[]string  				`json:"arguments,omitempty"`
	State         	*mesosproto.TaskState 	`json:"state,string"`

	DockerID		string
	DockerName		string
	SlaveID     	string
	CreatedTime     int64
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
