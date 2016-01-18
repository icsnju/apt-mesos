package core

import (
	"net/http"
	
	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

type Core struct {
	addr 			string
	master  		string
	frameworkInfo  	*mesosproto.FrameworkInfo
	log				*logrus.Logger
	registry		*registry.Registry
	events 			Events
	
	Endpoints		map[string]map[string]func(w http.ResponseWriter, r *http.Request) error
}

func NewCore(addr string, master string, frameworkInfo *mesosproto.FrameworkInfo, log *logrus.Logger) *Core{
	core := &Core{
		addr:			addr,
		master:			master,
		frameworkInfo: 	frameworkInfo,
		log:			log,
		Endpoints:		nil,
		events:			NewEvents(),
	}
	core.InitEndpoints()
	return core
}

// Framework register to mesos master
func (core *Core) RegisterFramework() error {
	core.log.WithFields(logrus.Fields{"master": core.master}).Info("Registering framework...")

	return core.SendMessageToMesos(&mesosproto.RegisterFrameworkMessage{
		Framework: core.frameworkInfo,
	}, "mesos.internal.RegisterFrameworkMessage")
}

// framework unregister from mesos master
func (core *Core) UnRegisterFramework() error {
	core.log.WithFields(logrus.Fields{"master": core.master}).Info("Unregistering framework...")

	return core.SendMessageToMesos(&mesosproto.UnregisterFrameworkMessage{
		FrameworkId: core.frameworkInfo.Id,
	}, "mesos.internal.UnRegisterFrameworkMessage")
}

// Send request to master for offers
func (core *Core) RequestOffers(resources []*mesosproto.Resource) ([]*mesosproto.Offer, error){
	core.log.Info("Request offers.")

	var event *mesosproto.Event
	select {
		case event = <-core.GetEvent(mesosproto.Event_OFFERS):
	}
	if event == nil {
		core.log.Info("send message")
		if err := core.SendMessageToMesos(&mesosproto.ResourceRequestMessage{
			FrameworkId: core.frameworkInfo.Id,
			Requests: []*mesosproto.Request{
				&mesosproto.Request{
					Resources: resources,
				},
			},
		}, "mesos.internal.ResourceRequestMessage"); err != nil {
			return nil, err
		}

		event = <-core.GetEvent(mesosproto.Event_OFFERS)
	}

	core.log.Infof("Received %d offer(s).", len(event.Offers.Offers))
	return event.Offers.Offers, nil	
}

// LauchTask with specific offer and resources
func (core *Core) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) error {
	core.log.WithFields(logrus.Fields{"ID": task.ID, "command": task.Command, "offerId": offer.Id, "dockerImage": task.DockerImage}).Info("Launching task...")

	taskInfo := createTaskInfo(offer, resources, task)

	return core.SendMessageToMesos(&mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		Tasks:       []*mesosproto.TaskInfo{taskInfo},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}

// Kill task with id
func (core *Core) KillTask(ID string) error {
	core.log.WithFields(logrus.Fields{"ID": ID}).Info("Killing task...")

	return core.SendMessageToMesos(&mesosproto.KillTaskMessage{
		FrameworkId: core.frameworkInfo.Id,
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
	}, "mesos.internal.KillTaskMessage")
}