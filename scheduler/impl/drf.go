package impl

import (
	"container/heap"

	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/scheduler/impl/resource"
	"github.com/icsnju/apt-mesos/scheduler/impl/structure"
)

// DRFScheduler implements scheduler using FCFS algorithm
type DRFScheduler struct {
	Heap          *structure.DRFHeap
	Queue         []*registry.Job
	init          bool
	totalResource map[string]float64
}

// NewDRFScheduler create a new scheduler
func NewDRFScheduler() *DRFScheduler {
	return &DRFScheduler{
		Heap: &structure.DRFHeap{},
	}
}

func (scheduler *DRFScheduler) AddJob(job *registry.Job) {
	// if heap has been initialed
	// push element into heap
	if scheduler.init {
		element := structure.NewDRFElement(job, scheduler.totalResource)
		heap.Push(scheduler.Heap, element)
	} else {
		// append element to list
		scheduler.Queue = append(scheduler.Queue, job)
	}
}

func (scheduler *DRFScheduler) HasJob() bool {
	if scheduler.init {
		return scheduler.Heap.Len() > 0
	} else {
		return len(scheduler.Queue) > 0
	}
}

func (scheduler *DRFScheduler) CheckFinished() {
	return
}

// Schedule implementation
func (scheduler *DRFScheduler) Schedule(offers []*mesosproto.Offer) (*registry.Task, *mesosproto.Offer, bool) {
	log.Debugf("Schedule tasks, current registry len: %v", scheduler.Heap.Len())
	if !scheduler.init {
		log.Debugf("Init DRF Heap, current waiting queue len: %d", len(scheduler.Queue))
		scheduler.totalResource = getTotalResource(offers)
		scheduler.init = true
		heap.Init(scheduler.Heap)
		for _, job := range scheduler.Queue {
			scheduler.AddJob(job)
		}
		return nil, nil, false
	}

	// get first task
	element := heap.Pop(scheduler.Heap).(*structure.DRFElement)
	job := element.Job
	log.Warn(element.DominantResource)
	task := job.FirstTask()

	// search suitable offer for first task of this job
	for _, offer := range offers {
		if resource.ResourcesMatch(task, offer) && resource.ConstraintsMatch(task, offer) {
			job.PopFirstTask()
			if !job.IsFinished() {
				scheduler.AddJob(job)
			}
			return task, offer, true
		}
	}

	if !job.IsFinished() {
		scheduler.AddJob(job)
	}

	return nil, nil, false
}

func getTotalResource(offers []*mesosproto.Offer) map[string]float64 {
	totalResource := make(map[string]float64)
	for _, offer := range offers {
		for _, resource := range offer.GetResources() {
			if resource.GetType().String() == "SCALAR" {
				_, exist := totalResource[resource.GetName()]
				if exist {
					totalResource[resource.GetName()] += float64(resource.GetScalar().GetValue())
				} else {
					totalResource[resource.GetName()] = float64(0.0)
				}
			}
		}
	}
	return totalResource
}
