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
	"github.com/icsnju/apt-mesos/fs"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/scheduler/impl/resource"
	"github.com/icsnju/apt-mesos/utils"
)

var (
	BuildCPU   float64 = 0.5
	BuildMem   float64 = 128
	CollectCPU float64 = 0.1
	CollectMem float64 = 16
)

func (core *Core) StartJob(job *registry.Job) error {
	log.Infof("Starting job: %v", job.ID)
	if job.Image == "" && job.ContextDir == "" {
		return errors.New("Start job error: image, context_dir cannot be nil at the same time")
	}

	job.StartTime = time.Now().UnixNano()
	if job.ContextDir != "" {
		core.BuildImage(job)
	}

	// TODO split input
	job.Status = registry.StatusRunning
	core.RunTask(job)

	return nil
}

func (core *Core) BuildImage(job *registry.Job) error {
	// Build Images before run test task
	// TaskID: build-{JobID}-{randID}-{NumberOfScale}
	log.Infof("Create task for job(%v) to build image", job.ID)
	job.Image = "image-" + job.ID
	for index := 1; index <= job.BuildNodeNumber(); index++ {
		task := &registry.Task{
			Cpus:       BuildCPU,
			Mem:        BuildMem,
			ID:         "build-" + job.ID + "-" + strconv.Itoa(index),
			Name:       job.Name,
			Type:       registry.TaskTypeBuild,
			CreateTime: time.Now().UnixNano(),
			JobID:      job.ID,
			State:      "TASK_WAITING",
			SLA:        registry.SLAOnePerNode,
		}

		err := core.AddTask(task.ID, task)
		job.PushTask(task)
		if err != nil {
			log.Errorf("Error when add %d build image task: %v", index, err)
			task.State = "TASK_FAILED"
			job.PopLastTask()
			continue
		}

	}
	return nil
}

func (core *Core) RunTask(job *registry.Job) {
	// Run test task of specified job
	// TaskID: test-{JobID}-{randID}-{scaleNumber}
	var inputs []string
	if job.Splitter != nil && job.InputPath != "" {
		inputs, err := job.Splitter.Split(job.InputPath)
		log.Warn(inputs)
		if err != nil {
			log.Errorf("Error when split job: %v", err)
			return
		}

		jobScale := len(inputs)
		core.addTask(job, jobScale, inputs)
		return
	}

	core.addTask(job, 0, inputs)
}

func (core *Core) addTask(job *registry.Job, scale int, inputs []string) {
	for _, task := range job.Tasks {
		randID, err := utils.Encode(6)
		if err != nil {
			log.Errorf("Error when generate id to task %s of job %v", task.ID, job.ID)
			continue
		}

		if task.Scale <= 0 {
			task.Scale = 1
		}

		if scale == 0 {
			scale = task.Scale
		}

		job.TotalTaskLen = scale
		for index := 1; index <= scale; index++ {
			// To avoid use same pointer of ports
			// Instantiate a new array
			var ports []*registry.Port
			for _, port := range task.Ports {
				ports = append(ports, &registry.Port{
					ContainerPort: port.ContainerPort,
					HostPort:      port.HostPort,
				})
			}

			taskInstance := &registry.Task{
				JobID:       job.ID,
				ID:          "task-" + job.ID + "-" + randID + "-" + strconv.Itoa(index),
				Name:        job.Name,
				DockerImage: job.Image,
				Cpus:        task.Cpus,
				Mem:         task.Mem,
				Disk:        task.Disk,
				Ports:       ports,
				Command:     task.Command,
				Volumes:     task.Volumes,
				Resources:   task.Resources,
				Attributes:  task.Attributes,
				CreateTime:  time.Now().UnixNano(),
				Type:        registry.TaskTypeTest,
				State:       "TASK_WAITING",
			}

			// mount input path
			if len(inputs) > 0 {
				taskInstance.Volumes = append(taskInstance.Volumes, &registry.Volume{
					HostPath:      fs.NormalizePath(inputs[index-1]),
					ContainerPath: "/input",
				})
			}

			// mount work directory
			taskInstance.Volumes = append(taskInstance.Volumes, &registry.Volume{
				HostPath:      fs.NormalizePath(job.WorkDirectory),
				ContainerPath: "/workspace",
			})

			// mount output
			taskInstance.Volumes = append(taskInstance.Volumes, &registry.Volume{
				HostPath:      fs.NormalizePath(job.OutputPath),
				ContainerPath: "/output",
			})

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

			// TODO bugfix: task point to one pointer
			var taskArguments []string
			for _, arg := range task.Arguments {
				taskArguments = append(taskArguments, arg)
			}
			taskInstance.Arguments = taskArguments

			err = core.AddTask(taskInstance.ID, taskInstance)
			job.PushTask(taskInstance)
			if err != nil {
				task.State = "TASK_FAILED"
				log.Errorf("Error when running task %v: %v", task.ID, err)
				job.PopLastTask()
				continue
			}
		}
	}
}

// CollectResult collect result for task
func (core *Core) CollectResult(job *registry.Job, task *registry.Task) {
	taskInstance := &registry.Task{
		Cpus:       CollectCPU,
		Mem:        CollectMem,
		ID:         "collect-" + job.ID + "-" + task.ID,
		Name:       job.Name,
		Type:       registry.TaskTypeBuild,
		CreateTime: time.Now().UnixNano(),
		JobID:      job.ID,
		State:      "TASK_WAITING",
		SLA:        registry.SLAOnePerNode,
		Directory:  task.Directory,
	}
	err := core.AddTask(taskInstance.ID, taskInstance)
	job.PushTask(taskInstance)
	if err != nil {
		log.Errorf("Error when add %d result collector: %v", err)
		task.State = "TASK_FAILED"
		job.PopLastTask()
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
		if port.HostPort == 0 {
			port.HostPort = resource.GeneratePort(offer.GetResources())
		}

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
	executorServePath := "http://" + core.GetAddr() + "/executor/builder/image_builder"
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

func (core *Core) CreateCollectResultTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) (*mesosproto.TaskInfo, error) {
	log.Debugf("Collect result taskInfo of task(%v)", task.ID)
	job, err := core.GetJob(task.JobID)
	if err != nil {
		return nil, err
	}

	if !job.HasContextDir() {
		return nil, errors.New("Context directory needed.")
	}

	executorServePath := "http://" + core.GetAddr() + "/executor/collector/result_collector"
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
		Name: proto.String("Collect Result (APT-MESOS)"),
		Command: &mesosproto.CommandInfo{
			Uris:  executorUris,
			Value: proto.String("./result_collector " + task.ID + " " + task.Directory + " " + job.OutputPath),
		},
	}
	return &mesosproto.TaskInfo{
		Executor:  executorInfo,
		Name:      proto.String("collect-" + task.JobID),
		Resources: resources,
		SlaveId:   offer.SlaveId,
		TaskId: &mesosproto.TaskID{
			Value: proto.String(task.ID),
		},
	}, nil
}

func (core *Core) CreateTaskRunnerInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) (*mesosproto.TaskInfo, error) {
	log.Debugf("Create task-runner taskInfo of task(%v)", task.ID)

	executorServePath := "http://" + core.GetAddr() + "/executor/runner/task_runner"

	executorUris := []*mesosproto.CommandInfo_URI{
		{
			Value:      &executorServePath,
			Executable: proto.Bool(true),
		},
	}

	// Set the docker image if specified
	dockerInfo := &mesosproto.ContainerInfo_DockerInfo{
		Image: &task.DockerImage,
	}

	containerInfo := &mesosproto.ContainerInfo{
		Type:   mesosproto.ContainerInfo_DOCKER.Enum(),
		Docker: dockerInfo,
	}

	mode := mesosproto.Volume_RO
	containerInfo.Volumes = append(containerInfo.Volumes, &mesosproto.Volume{
		ContainerPath: proto.String("/task_runner"),
		HostPath:      proto.String("./task_runner"),
		Mode:          &mode,
	})

	commandInfo := &mesosproto.CommandInfo{
		Shell: proto.Bool(false),
	}
	if len(task.Arguments) > 0 {
		for _, argument := range task.Arguments {
			commandInfo.Arguments = append(commandInfo.Arguments, argument)
		}
	}

	executorInfo := &mesosproto.ExecutorInfo{
		ExecutorId: &mesosproto.ExecutorID{
			Value: proto.String(task.ID),
		},
		Name: proto.String("Run Task (APT-MESOS)"),
		Command: &mesosproto.CommandInfo{
			Uris:  executorUris,
			Value: proto.String("/task_runner"),
		},
		// Command:   commandInfo,
		Container: containerInfo,
	}

	return &mesosproto.TaskInfo{
		Executor:  executorInfo,
		Name:      proto.String("task-" + task.JobID),
		Resources: resources,
		SlaveId:   offer.SlaveId,
		TaskId: &mesosproto.TaskID{
			Value: proto.String(task.ID),
		},
		Container: containerInfo,
		Data:      []byte("file1"),
	}, nil
}
