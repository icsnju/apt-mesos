package resource

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ConstraintsMatch check if a offer fit task's constaints
func ConstraintsMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	for _, attribute := range offer.GetAttributes() {
		// SLA Match
		if attribute.GetName() == "SLATasks" {
			items := attribute.GetSet().GetItem()
			for _, item := range items {
				if item == task.Parse() {
					return false
				}
			}
		}
	}
	return true
}
