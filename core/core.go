package core

import (
	"github.com/Sirupsen/logrus"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

type Core struct {
	addr 			string
	master  		string
	frameworkInfo  	*mesosproto.FrameworkInfo
	log				*logrus.Logger
	registry		*registry.Registry
	events 			Events
}

func NewCore(addr string, master string, frameworkInfo *mesosproto.FrameworkInfo, log *logrus.Logger) *Core{
	return &Core{
		addr:			addr,
		master:			master,
		frameworkInfo: 	frameworkInfo,
		log:			log,
		events:			NewEvents(),
	}
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
