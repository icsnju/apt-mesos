package core

import (
	"fmt"
	"strings"
	
	"github.com/golang/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

func createTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) *mesosproto.TaskInfo {
	taskInfo := mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("volt-task-%s", task.ID)),
		TaskId: &mesosproto.TaskID{
			Value: &task.ID,
		},
		SlaveId:   offer.SlaveId,
		Resources: resources,
		Command:   &mesosproto.CommandInfo{},
	}

	// Set value only if provided
	commands := strings.Split(task.Command, " ")
	if commands[0] != "" {
		taskInfo.Command.Value = &commands[0]
	}

	// Set args only if they exist
	if len(commands) > 1 {
		taskInfo.Command.Arguments = commands[1:]
	}
	
	// Set the docker image if specified
	if task.DockerImage != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &task.DockerImage,
			},
		}

		for _, v := range task.Volumes {
			var (
				vv   = v
				mode = mesosproto.Volume_RW
			)

			if vv.Mode == "ro" {
				mode = mesosproto.Volume_RO
			}

			taskInfo.Container.Volumes = append(taskInfo.Container.Volumes, &mesosproto.Volume{
				ContainerPath: &vv.ContainerPath,
				HostPath:      &vv.HostPath,
				Mode:          &mode,
			})
		}

		taskInfo.Command.Shell = proto.Bool(false)
	}

	return &taskInfo
}