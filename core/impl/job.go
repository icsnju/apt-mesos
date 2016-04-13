package impl

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// CreateSingleTaskInfo build single taskInfo for task
func (core *Core) CreateSingleTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) *mesosproto.TaskInfo {
	portResources := []*mesosproto.Value_Range{}

	// Set the docker image if specified
	dockerInfo := &mesosproto.ContainerInfo_DockerInfo{
		Image: &task.DockerImage,
	}
	containerInfo := &mesosproto.ContainerInfo{
		Type:   mesosproto.ContainerInfo_DOCKER.Enum(),
		Docker: dockerInfo,
	}
	for _, volume := range task.Volumes {
		mode := mesosproto.Volume_RW
		if volume.Mode == "ro" {
			mode = mesosproto.Volume_RO
		}

		containerInfo.Volumes = append(containerInfo.Volumes, &mesosproto.Volume{
			ContainerPath: &volume.ContainerPath,
			HostPath:      &volume.HostPath,
			Mode:          &mode,
		})
	}

	for _, port := range task.Ports {
		dockerInfo.PortMappings = append(dockerInfo.PortMappings, &mesosproto.ContainerInfo_DockerInfo_PortMapping{
			ContainerPort: &port.ContainerPort,
			HostPort:      &port.HostPort,
		})
		portResources = append(portResources, &mesosproto.Value_Range{
			Begin: proto.Uint64(uint64(port.HostPort)),
			End:   proto.Uint64(uint64(port.HostPort)),
		})
	}

	if len(task.Ports) > 0 {
		// port mapping only works in bridge mode
		dockerInfo.Network = mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()
	} else if len(task.NetworkMode) > 0 {
		if task.NetworkMode == registry.NetworkModeBridge {
			dockerInfo.Network = mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()
		} else if task.NetworkMode == registry.NetworkModeHost {
			dockerInfo.Network = mesosproto.ContainerInfo_DockerInfo_HOST.Enum()
		} else if task.NetworkMode == registry.NetworkModeNone {
			dockerInfo.Network = mesosproto.ContainerInfo_DockerInfo_NONE.Enum()
		}
	}

	commandInfo := &mesosproto.CommandInfo{
		Shell: proto.Bool(false),
	}
	if len(task.Arguments) > 0 {
		for _, argument := range task.Arguments {
			commandInfo.Arguments = append(commandInfo.Arguments, argument)
		}
	}

	if len(task.Ports) > 0 {
		resources = append(resources,
			&mesosproto.Resource{
				Name:   proto.String("ports"),
				Ranges: &mesosproto.Value_Ranges{Range: portResources},
				Type:   mesosproto.Value_RANGES.Enum(),
			},
		)
	}

	taskInfo := &mesosproto.TaskInfo{
		Name:      proto.String(fmt.Sprintf("task-%s", task.ID)),
		TaskId:    &mesosproto.TaskID{Value: &task.ID},
		SlaveId:   offer.SlaveId,
		Container: containerInfo,
		Command:   commandInfo,
		Resources: resources,
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

	return taskInfo
}

// func createBuildImageTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, dockerfile *docker.Dockerfile) *mesosproto.TaskInfo {
// 	if dockerfile.HasLocalSources() {
// 		dockerfile.BuildContext()
// 	}
// }
