package impl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/docker"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/utils"
)

var (
	BUILD_CPU float64 = 0.5
	BUILD_MEM float64 = 128
)

func (core *Core) StartJob(job *registry.Job) error {
	log.Infof("Starting job: %v", job.ID)
	if job.Image == "" && job.ContextDir == "" {
		return errors.New("Start job error: image, context_dir cannot be nil at the same time")
	}

	if job.ContextDir != "" {
		core.BuildImage(job, len(job.Tasks))
	}

	// TODO split input
	core.RunTask(job)

	return nil
}

func (core *Core) BuildImage(job *registry.Job, size int) error {
	// Build Images before run test task
	// TaskID: build-{JobID}-{randID}-{NumberOfScale}
	log.Infof("Create task for job(%v) to build image", job.ID)
	job.Image = "image-" + job.ID
	for index := 1; index <= size; index++ {
		task := &registry.Task{
			Cpus:       BUILD_CPU,
			Mem:        BUILD_MEM,
			ID:         "build-" + job.ID + "-" + strconv.Itoa(index),
			Name:       job.Name + " [BUILD IMAGE]",
			Type:       registry.TaskTypeBuild,
			CreateTime: time.Now().UnixNano(),
			JobID:      job.ID,
			State:      "TASK_WAITING",
			SLA:        registry.SLAOnePerNode,
		}

		err := core.AddTask(task.ID, task)
		if err != nil {
			log.Errorf("Error when add %d build image task: %v", index, err)
			continue
		}
	}
	return nil
}

func (core *Core) RunTask(job *registry.Job) {
	// Run test task of specified job
	// TaskID: test-{JobID}-{randID}-{scaleNumber}
	for _, task := range job.Tasks {
		randID, err := utils.Encode(6)
		if err != nil {
			log.Errorf("Error when generate id to task %d of job %v", task.ID, job.ID)
			continue
		}

		if task.Scale <= 0 {
			task.Scale = 1
		}

		for index := 1; index <= task.Scale; index++ {
			taskInstance := &registry.Task{
				JobID:       job.ID,
				ID:          "task-" + job.ID + "-" + randID + "-" + strconv.Itoa(index),
				Name:        job.Name + " [RUN TASK]",
				DockerImage: job.Image,
				Cpus:        task.Cpus,
				Mem:         task.Mem,
				Disk:        task.Disk,
				Ports:       task.Ports,
				Command:     task.Command,
				Resources:   task.Resources,
				Attributes:  task.Attributes,
				CreateTime:  time.Now().UnixNano(),
				Type:        registry.TaskTypeTest,
				State:       "TASK_WAITING",
			}

			// if task was build from dockerfile
			// add attribute to task
			if job.ContextDir != "" {
				taskInstance.Attributes = append(task.Attributes, &mesosproto.Attribute{
					Name: proto.String("Image"),
					Text: &mesosproto.Value_Text{
						Value: proto.String(job.Image),
					},
				})
			}

			err = core.AddTask(taskInstance.ID, taskInstance)
			if err != nil {
				task.State = "TASK_FAILED"
				log.Errorf("Error when running task %v: %v", task.ID, err)
				continue
			}
		}
	}
}

// CreateSingleTaskInfo build single taskInfo for task
func (core *Core) CreateSingleTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) (*mesosproto.TaskInfo, error) {
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
		Name:      proto.String(fmt.Sprintf("test-%s", task.ID)),
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

	return taskInfo, nil
}

func (core *Core) CreateBuildImageTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) (*mesosproto.TaskInfo, error) {
	log.Debugf("Build image taskInfo of task(%v)", task.ID)
	job, err := core.GetJob(task.JobID)
	if err != nil {
		return nil, err
	}

	if !job.HasContextDir() {
		return nil, errors.New("Context directory needed.")
	}

	if exists := job.DockerfileExists(); !exists {
		return nil, errors.New("Cannot found Dockerfile in context directory.")
	}

	job.Dockerfile = docker.NewDockerfile("dockerfile-"+job.ID, job.ContextDir)
	err = job.Dockerfile.BuildContext()
	if err != nil {
		return nil, err
	}

	contextServePath := "http://" + core.GetAddr() + "/context/" + job.Dockerfile.ID + ".tar"
	executorServePath := "http://" + core.GetAddr() + "/executor/image_builder"
	log.Debugf("Context file served on path: %v", contextServePath)

	executorUris := []*mesosproto.CommandInfo_URI{
		{
			Value:      &executorServePath,
			Executable: proto.Bool(true),
		},
	}

	executorInfo := &mesosproto.ExecutorInfo{
		ExecutorId: &mesosproto.ExecutorID{
			Value: proto.String(task.ID),
		},
		Name: proto.String("Build Image (APT-MESOS)"),
		Command: &mesosproto.CommandInfo{
			Uris:  executorUris,
			Value: proto.String("./image_builder"),
		},
	}
	return &mesosproto.TaskInfo{
		Executor:  executorInfo,
		Name:      proto.String("image-" + task.JobID),
		Resources: resources,
		SlaveId:   offer.SlaveId,
		TaskId: &mesosproto.TaskID{
			Value: proto.String(task.ID),
		},
		Data: []byte(contextServePath),
	}, nil
}
