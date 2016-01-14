package core

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
)

type Core interface {
	
	Start(scheduler.SchedulerDriver, *mesosproto.MasterInfo) error

	ScheduleTasks(scheduler.SchedulerDriver, []*mesosproto.Offer)

	StatusUpdate(scheduler.SchedulerDriver, *mesosproto.TaskStatus)

	Stop()
}

