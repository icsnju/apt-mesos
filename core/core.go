package core

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// Core has the core function of sher
type Core struct {
	addr          string
	master        string
	frameworkInfo *mesosproto.FrameworkInfo
	events        Events

	Log       *logrus.Logger
	Endpoints map[string]map[string]func(w http.ResponseWriter, r *http.Request) error
}

// NewCore returns a new core by given params:
// 	addr: address of sher
//	master: address of mesos-master
// 	frameworkInfo: mesosproto, including framework name, host, id, etc.
// 	Log: logrus instance
// 	Endpoints: endpoints of all apis
//  events:	event channel
func NewCore(addr string, master string, frameworkInfo *mesosproto.FrameworkInfo, log *logrus.Logger) *Core {
	core := &Core{
		addr:          addr,
		master:        master,
		frameworkInfo: frameworkInfo,
		events:        NewEvents(),
		Log:           log,
		Endpoints:     nil,
	}
	core.InitEndpoints()
	return core
}

// RegisterFramework register to mesos master
func (core *Core) RegisterFramework() error {
	core.Log.WithFields(logrus.Fields{"master": core.master}).Info("Registering framework...")

	return core.SendMessageToMesos(&mesosproto.RegisterFrameworkMessage{
		Framework: core.frameworkInfo,
	}, "mesos.internal.RegisterFrameworkMessage")
}

// UnRegisterFramework framework unregister from mesos master
func (core *Core) UnRegisterFramework() error {
	core.Log.WithFields(logrus.Fields{"master": core.master}).Info("Unregistering framework...")

	return core.SendMessageToMesos(&mesosproto.UnregisterFrameworkMessage{
		FrameworkId: core.frameworkInfo.Id,
	}, "mesos.internal.UnRegisterFrameworkMessage")
}

// RequestOffers send request to master for offers
func (core *Core) RequestOffers(resources []*mesosproto.Resource) ([]*mesosproto.Offer, error) {
	core.Log.Info("Request offers.")

	var event *mesosproto.Event
	select {
	case event = <-core.GetEvent(mesosproto.Event_OFFERS):
	}
	if event == nil {
		core.Log.Info("send message")
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
	core.Log.Infof("Received %d offer(s).", len(event.Offers.Offers))
	return event.Offers.Offers, nil
}

// AcceptOffer send message to mesos-master to accept a offer
func (core *Core) AcceptOffer(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) error {
	core.Log.WithFields(logrus.Fields{"ID": task.ID, "command": task.Command, "offerId": offer.Id, "dockerImage": task.DockerImage}).Info("Launching task...")

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

// DeclineOffer decline a offer which is not suit for task
func (core *Core) DeclineOffer(offer *mesosproto.Offer, task *registry.Task) error {
	core.Log.WithFields(logrus.Fields{"offerId": offer.Id, "slave": offer.GetHostname()}).Debug("Decline offer...")

	return core.SendMessageToMesos(&mesosproto.LaunchTasksMessage{
		FrameworkId: core.frameworkInfo.Id,
		Tasks:       []*mesosproto.TaskInfo{},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}

// LaunchTask with specific offer and resources
func (core *Core) LaunchTask(offer *mesosproto.Offer, offers []*mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) error {
	for _, value := range offers {
		if offer.GetId().GetValue() == value.GetId().GetValue() {
			if err := core.AcceptOffer(value, resources, task); err != nil {
				return err
			}
		} else {
			if err := core.DeclineOffer(value, task); err != nil {
				return err
			}
		}
	}
	return nil
}

// KillTask kill task with id
func (core *Core) KillTask(ID string) error {
	core.Log.WithFields(logrus.Fields{"ID": ID}).Info("Killing task...")

	return core.SendMessageToMesos(&mesosproto.KillTaskMessage{
		FrameworkId: core.frameworkInfo.Id,
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
	}, "mesos.internal.KillTaskMessage")
}
