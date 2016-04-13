package scheduler

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// Scheduler interface to schedule task to run
type Scheduler interface {
	Schedule(tasks []*registry.Task, offers []*mesosproto.Offer) (*mesosproto.Offer, *registry.Task, error)
}
