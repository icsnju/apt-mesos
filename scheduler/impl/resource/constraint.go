package resource

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ConstraintsMatch check if a offer fit task's constaints
// TODO implementation
func ConstraintsMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	return SlAMatch(task, offer)
}

func SlAMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	for _, attribute := range offer.GetAttributes() {
		if attribute.GetName() == "SLATasks" {
			items := attribute.GetSet().GetItem()
			for _, item := range items {
				if item == task.Parse() {
					return false
				}
			}
			return true
		}
	}
	return true
}
