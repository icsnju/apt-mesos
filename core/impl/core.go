package impl

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	scheduler "github.com/icsnju/apt-mesos/scheduler/impl"
	"github.com/icsnju/apt-mesos/scheduler/impl/resource"

	log "github.com/Sirupsen/logrus"
)

// Core has the core function of sher
type Core struct {
	addr          string
	master        string
	frameworkInfo *mesosproto.FrameworkInfo
	masterUPID    *UPID
	coreUPID      *UPID
	events        Events
	tasks         registry.Registry
	nodes         registry.Registry
	jobs          registry.Registry
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
		jobs:          *registry.NewRegistry(),
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
	m, err := Parse("master@" + core.master)
	if err != nil {
		return err
	}
	core.masterUPID = m

	host, port, err := core.GetListenIPAndPort()
	if err != nil {
		return err
	}
	core.coreUPID = &UPID{
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
	messagePackage := NewMessage(core.masterUPID, message, nil)

	log.Debugf("Registering with master %s [%s] ", core.masterUPID, message)
	err = SendMessageToMesos(core.coreUPID, messagePackage)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	go core.monitor()
	go core.schedule()
	return nil
}

func (core *Core) schedule() {
	for {
		tasks := core.GetUnScheduledTask()
		if len(tasks) > 0 {
			offers, err := core.RequestOffers()
			if err != nil {
				log.Errorf("Request offers error: %v", err)
				continue
			}

			for _, offer := range offers {
				log.Debug(offer)
			}

			task, offer, success := core.scheduler.Schedule(tasks, offers)

			// if remained resource can run a task
			if success {
				log.Infof("Schedule task result: run task(%v) on %v", task.ID, offer.GetHostname())
				core.LaunchTask(task, offer, offers)
				core.updateNodeByTask(offer.GetSlaveId().GetValue(), task)
			} else {
				log.Infof("No enough resources remained, wait for other tasks finish")
				time.Sleep(3 * time.Second)
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

// GetAgentLisenPort return cadvisor listend port
func (core *Core) GetAgentLisenPort() string {
	return "18080"
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
					Resources: resource.BuildEmptyResources(),
				},
			},
		}
		messagePackage := NewMessage(core.masterUPID, message, nil)
		if err := SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
			return nil, err
		}

		event = <-core.GetEvent(mesosproto.Event_OFFERS)
		log.Debugf("Event: %v", event)
	}
	return event.Offers.Offers, nil
}

// AcceptOffer send message to mesos-master to accept a offer
func (core *Core) AcceptOffer(offer *mesosproto.Offer, resources []*mesosproto.Resource, taskInfo *mesosproto.TaskInfo) error {
	message := &mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		OfferIds:    []*mesosproto.OfferID{offer.Id},
		Tasks:       []*mesosproto.TaskInfo{taskInfo},
		Filters:     &mesosproto.Filters{},
	}

	messagePackage := NewMessage(core.masterUPID, message, nil)
	if err := SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send AcceptOffer message: %v\n", err)
	}
	return nil
}

// DeclineOffer decline a offer which is not suit for task
func (core *Core) DeclineOffer(offer *mesosproto.Offer, task *registry.Task) error {
	message := &mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		OfferIds:    []*mesosproto.OfferID{offer.Id},
		Tasks:       []*mesosproto.TaskInfo{},
		Filters:     &mesosproto.Filters{},
	}

	messagePackage := NewMessage(core.masterUPID, message, nil)
	if err := SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send DeclineOffer message: %v\n", err)
	}
	return nil
}

// LaunchTask with specific offer and resources
func (core *Core) LaunchTask(task *registry.Task, offer *mesosproto.Offer, offers []*mesosproto.Offer) error {
	core.generateResource(task)
	resources := resource.BuildResources(task)

	log.Infof("Launch task %v, on node %v", task.ID, offer.GetHostname())
	taskInfo := &mesosproto.TaskInfo{}
	var err error
	if task.Type == registry.TaskType_Test {
		taskInfo, err = core.CreateSingleTaskInfo(offer, resources, task)
	} else if task.Type == registry.TaskType_Build {
		taskInfo, err = core.CreateBuildImageTaskInfo(offer, resources, task)
	} else {
		return errors.New("Unknown task type received.")
	}

	if err != nil {
		return err
	}

	for _, value := range offers {
		if offer.GetSlaveId() == value.GetSlaveId() {
			if err := core.AcceptOffer(value, resources, taskInfo); err != nil {
				return err
			}
			task.State = "TASK_STAGING"
		} else {
			log.Debugf("Decline offer for node: %v", value.GetHostname())
			if err := core.DeclineOffer(value, task); err != nil {
				return err
			}
		}
	}
	return nil
}
