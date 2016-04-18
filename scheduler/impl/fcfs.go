package impl

import (
	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/scheduler/impl/resource"
)

// FCFSScheduler implements scheduler using FCFS algorithm
type FCFSScheduler struct{}

// NewScheduler create a new scheduler
func NewScheduler() *FCFSScheduler {
	return &FCFSScheduler{}
}

// Schedule implementation
func (scheduler *FCFSScheduler) Schedule(tasks []*registry.Task, offers []*mesosproto.Offer) (*registry.Task, *mesosproto.Offer, bool) {
	log.Debugf("Schedule tasks, current registry len: %v", len(tasks))
	queue := registry.NewFCFSQueue(tasks)
	for _, task := range queue {
		for _, offer := range offers {
			log.Debug(task.ID)
			log.Debug(offer.GetHostname())
			log.Debug(offer.GetAttributes())
			log.Debug(resource.ResourcesMatch(task, offer))
			log.Debug(resource.ConstraintsMatch(task, offer))
			if resource.ResourcesMatch(task, offer) && resource.ConstraintsMatch(task, offer) {
				return task, offer, true
			}
		}
	}
	return nil, nil, false
}
