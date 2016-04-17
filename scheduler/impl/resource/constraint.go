package resource

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ConstraintsMatch check if a offer fit task's constaints
// TODO implementation
func ConstraintsMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	return true
}

func SlAMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	// TODO implementation
	return true
}
