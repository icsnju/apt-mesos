package scheduler

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// Scheduler interface to schedule task to run
type Scheduler interface {
	Schedule(offers []*mesosproto.Offer) (*registry.Task, *mesosproto.Offer, bool)
	AddJob(job *registry.Job)
	HasJob() bool
	CheckFinished()
}

const (
	FCFS = "fcfs"
	DRF  = "drf"
)
