package core

import (
	"fmt"

	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

/*
* FCFS scheduling algorithm
*/
func (core *Core) ScheduleTask(offers []*mesosproto.Offer, resources []*mesosproto.Resource, task *registry.Task) (*mesosproto.Offer, error) {
	for _, offer := range offers {
		if ScalarResource("cpus", offer) >= task.Cpus && ScalarResource("mem", offer) >= task.Mem && ScalarResource("disk", offer) >= task.Disk {
			core.log.WithField("offer-slave-id", offer.GetHostname()).Debug("Scheduler choose the offer")
			return offer, nil
		}
	}
	return nil, fmt.Errorf("No resource left")
}