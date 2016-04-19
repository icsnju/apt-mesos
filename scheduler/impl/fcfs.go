package impl

import (
	"container/list"

	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/scheduler/impl/resource"
)

// FCFSScheduler implements scheduler using FCFS algorithm
type FCFSScheduler struct {
	Queue *list.List
}

// NewFCFSScheduler create a new scheduler
func NewFCFSScheduler() *FCFSScheduler {
	return &FCFSScheduler{
		Queue: list.New(),
	}
}

func (scheduler *FCFSScheduler) AddJob(job *registry.Job) {
	scheduler.Queue.PushBack(job)
}

func (scheduler *FCFSScheduler) HasJob() bool {
	return scheduler.Queue.Len() > 0
}

func (scheduler *FCFSScheduler) CheckFinished() {
	for scheduler.Queue.Len() > 0 {
		element := scheduler.Queue.Front()
		job := element.Value.(*registry.Job)
		if job.IsFinished() {
			scheduler.Queue.Remove(element)
		} else {
			break
		}
	}
}

// Schedule implementation
func (scheduler *FCFSScheduler) Schedule(offers []*mesosproto.Offer) (*registry.Task, *mesosproto.Offer, bool) {
	log.Debugf("Schedule tasks, current registry len: %v", scheduler.Queue.Len())

	// get first task
	job := scheduler.Queue.Front().Value.(*registry.Job)
	task := job.FirstTask()

	for _, offer := range offers {
		if resource.ResourcesMatch(task, offer) && resource.ConstraintsMatch(task, offer) {
			job.PopFirstTask()
			return task, offer, true
		}
	}

	return nil, nil, false
}
