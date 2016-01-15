package core

import (
	"github.com/Sirupsen/logrus"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/mesos/mesos-go/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

type Core struct {
	addr 			string
	master  		string
	frameworkInfo  	*mesosproto.FrameworkInfo
	log				*logrus.Logger
	registry		*registry.Registry
}

func NewCore(addr string, master string, frameworkInfo *mesosproto.FrameworkInfo, log *logrus.Logger) *Core{
	return &Core{
		addr:			addr,
		master:			master,
		frameworkInfo: 	frameworkInfo,
		log:			log,
	}
}

func (core *Core) RegisterFramework() error {
	core.log.WithFields(logrus.Fields{"master": core.master}).Info("Registering framework...")

	return core.SendMessageToMesos(&mesosproto.RegisterFrameworkMessage{
		Framework: core.frameworkInfo,
	}, "mesos.internal.RegisterFrameworkMessage")
}

func (core *Core) ScheduleTasks(scheduler.SchedulerDriver, []*mesosproto.Offer) {
}