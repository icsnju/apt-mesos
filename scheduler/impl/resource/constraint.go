package resource

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ConstraintsMatch check if a offer fit task's constaints
func ConstraintsMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	if !SlAMatch(task, offer) {
		return false
	}

	// check task attributes
	for _, attribute := range task.Attributes {
		found := false
		for _, offerAttribute := range offer.GetAttributes() {
			if offerAttribute.GetName() == attribute.GetName() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
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
