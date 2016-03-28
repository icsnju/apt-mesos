package scheduler

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// Scheduler interface to schedule task to run
type Scheduler interface {
	Schedule(offers []*mesosproto.Offer, tasks []*registry.Task) (*mesosproto.Offer, *registry.Task, error)
}
