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


// framework register to mesos master
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
