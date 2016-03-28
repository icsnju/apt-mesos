package impl

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	comm "github.com/icsnju/apt-mesos/communication"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	scheduler "github.com/icsnju/apt-mesos/scheduler/impl"

	log "github.com/Sirupsen/logrus"
)

// Core has the core function of sher
type Core struct {
	addr          string
	master        string
	frameworkInfo *mesosproto.FrameworkInfo
	masterUPID    *comm.UPID
	coreUPID      *comm.UPID
	events        Events
	tasks         registry.Registry
	nodes         registry.Registry
	scheduler     scheduler.FCFSScheduler

	Endpoints map[string]map[string]func(w http.ResponseWriter, r *http.Request) error
}

// NewCore returns a new core by given params:
// 	addr: address of sher
//	master: address of mesos-master
// 	frameworkInfo: mesosproto, including framework name, host, id, etc.
// 	Log: logrus instance
// 	Endpoints: endpoints of all apis
//  events:	event channel
//  nodes: node registry
//  tasks: task registry
func NewCore(addr string, master string) *Core {
	core := &Core{
		addr:          addr,
		master:        master,
		frameworkInfo: &mesosproto.FrameworkInfo{},
		events:        NewEvents(),
		tasks:         *registry.NewRegistry(),
		nodes:         *registry.NewRegistry(),
		scheduler:     *scheduler.NewScheduler(),
		Endpoints:     nil,
	}
	return core
}

// Run register to mesos master
func (core *Core) Run() error {
	log.Info("Starting apt-mesos...")
	// get current hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// add masterUPID and coreUPID
	m, err := comm.Parse("master@" + core.master)
	if err != nil {
		return err
	}
	core.masterUPID = m

	host, port, err := core.GetListenIPAndPort()
	if err != nil {
		return err
	}
	core.coreUPID = &comm.UPID{
		ID:   "core",
		Host: host,
		Port: port,
	}

	// create frameworkInfo
	core.frameworkInfo = &mesosproto.FrameworkInfo{
		Name:       proto.String("apt-mesos"),
		User:       proto.String("root"),
		WebuiUrl:   proto.String("http://" + core.addr),
		Hostname:   proto.String(hostname),
		Checkpoint: proto.Bool(false),
	}

	// create message
	message := &mesosproto.RegisterFrameworkMessage{
		Framework: core.frameworkInfo,
	}
	messagePackage := comm.NewMessage(core.masterUPID, message, nil)

	log.Debugf("Registering with master %s [%s] ", core.masterUPID, message)
	err = comm.SendMessageToMesos(core.coreUPID, messagePackage)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	go core.schedule()
	return nil
}

func (core *Core) schedule() {
	for {
		tasks := core.GetUnScheduledTask()
		// log.Debugf("Task queue size: %v", len(tasks))
		if len(tasks) > 0 {
			offers, _ := core.RequestOffers()
			task, node, success := core.scheduler.Schedule(tasks, core.GetAllNodes())
			// if remained resource can run a task
			if success {
				log.Infof("Schedule task result: run task(%v) on %v", task.ID, node.Hostname)
				core.LaunchTask(task, node, offers)
				core.updateNodeByTask(node.ID, task)
			} else {
				log.Infof("No enough resources remained, wait for other tasks finish")
			}
		} else {
			// log.Debug("Task registry has no task, sleep for a while...")
			time.Sleep(3 * time.Second)
		}
	}
}

// UnRegister framework unregister from mesos master
func (core *Core) UnRegister() error {
	log.WithFields(log.Fields{"master": core.master}).Info("Unregistering framework...")

	return core.SendMessageToMesos(&mesosproto.UnregisterFrameworkMessage{
		FrameworkId: core.frameworkInfo.Id,
	}, "mesos.internal.UnRegisterFrameworkMessage")
}

// GetAddr return core's addr
func (core *Core) GetAddr() string {
	return core.addr
}

// GetListenIPAndPort return core listend ip
func (core *Core) GetListenIPAndPort() (string, string, error) {
	splits := strings.Split(core.addr, ":")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("Expect one `:' in the core addr")
	}
	return splits[0], splits[1], nil
}

// RequestOffers send request to master for offers
func (core *Core) RequestOffers() ([]*mesosproto.Offer, error) {
	var event *mesosproto.Event
	select {
	case event = <-core.GetEvent(mesosproto.Event_OFFERS):
	}
	if event == nil {
		log.Info("Offer channel is nil, send request message")
		message := &mesosproto.ResourceRequestMessage{
			FrameworkId: core.frameworkInfo.Id,
			Requests: []*mesosproto.Request{
				&mesosproto.Request{
					Resources: scheduler.BuildEmptyResources(),
				},
			},
		}
		messagePackage := comm.NewMessage(core.masterUPID, message, nil)
		if err := comm.SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
			return nil, err
		}

		event = <-core.GetEvent(mesosproto.Event_OFFERS)
		log.Debugf("Event: %v", event)
	}
	return event.Offers.Offers, nil
}

// AcceptOffer send message to mesos-master to accept a offer
func (core *Core) AcceptOffer(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) error {
	log.Infof("Lauch task %v, command(%v), docker_image(%v)", task.ID, task.Command, task.DockerImage)
	taskInfo := core.CreateTaskInfo(offer, resources, task)
	message := &mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		OfferIds:    []*mesosproto.OfferID{offer.Id},
		Tasks:       []*mesosproto.TaskInfo{taskInfo},
		Filters:     &mesosproto.Filters{},
	}

	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	if err := comm.SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send AcceptOffer message: %v\n", err)
	}
	return nil
}

// DeclineOffer decline a offer which is not suit for task
func (core *Core) DeclineOffer(offer *mesosproto.Offer, task *registry.Task) error {
	log.WithFields(log.Fields{"offerId": offer.Id, "slave": offer.GetHostname()}).Debug("Decline offer...")
	message := &mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		OfferIds:    []*mesosproto.OfferID{offer.Id},
		Tasks:       []*mesosproto.TaskInfo{},
		Filters:     &mesosproto.Filters{},
	}

	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	if err := comm.SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send DeclineOffer message: %v\n", err)
	}
	return nil
}

// CreateTaskInfo build taskInfo for task
func (core *Core) CreateTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) *mesosproto.TaskInfo {
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

// LaunchTask with specific offer and resources
func (core *Core) LaunchTask(task *registry.Task, node *registry.Node, offers []*mesosproto.Offer) error {
	resources := scheduler.BuildResources(task)
	for _, value := range offers {
		if node.ID == value.GetSlaveId().GetValue() {
			if err := core.AcceptOffer(value, resources, task); err != nil {
				return err
			}
			stagingState := mesosproto.TaskState_TASK_STAGING
			task.State = &stagingState
			core.UpdateTask(task.ID, task)
		} else {
			if err := core.DeclineOffer(value, task); err != nil {
				return err
			}
		}
	}
	return nil
}

// HandleFrameworkRegisteredMessage called when framework register to mesos-master
func (core *Core) HandleFrameworkRegisteredMessage(message *mesosproto.FrameworkRegisteredMessage) {
	log.WithField("frameworkId", message.FrameworkId).Debug("Receive framworkId from mesos master")
	core.frameworkInfo.Id = message.FrameworkId

	eventType := mesosproto.Event_REGISTERED
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	})
}

// HandleResourceOffersMessage called when framework receive offers from master
func (core *Core) HandleResourceOffersMessage(message *mesosproto.ResourceOffersMessage) {
	eventType := mesosproto.Event_OFFERS
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	})
}

// HandleStatusUpdateMessage called when slave's status updated
func (core *Core) HandleStatusUpdateMessage(statusMessage *mesosproto.StatusUpdateMessage) error {
	message := &mesosproto.StatusUpdateAcknowledgementMessage{
		FrameworkId: statusMessage.GetUpdate().FrameworkId,
		SlaveId:     statusMessage.GetUpdate().Status.SlaveId,
		TaskId:      statusMessage.GetUpdate().Status.TaskId,
		Uuid:        statusMessage.GetUpdate().Uuid,
	}
	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	if err := comm.SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send StatusAccept message: %v\n", err)
		return err
	}

	eventType := mesosproto.Event_UPDATE
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Update: &mesosproto.Event_Update{
			Uuid:   statusMessage.Update.Uuid,
			Status: statusMessage.Update.Status,
		},
	})

	return nil
}
