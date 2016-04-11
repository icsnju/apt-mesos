package impl

import (
	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/registry"
)

// FCFSScheduler implements scheduler using FCFS algorithm
type FCFSScheduler struct{}

// NewScheduler create a new scheduler
func NewScheduler() *FCFSScheduler {
	return &FCFSScheduler{}
}

// Schedule implementation
func (scheduler *FCFSScheduler) Schedule(tasks []*registry.Task, nodes []*registry.Node) (*registry.Task, *registry.Node, bool) {
	log.Debugf("Schedule tasks, current registry len: %v", len(tasks))
	queue := registry.NewFCFSQueue(tasks)
	for _, task := range queue {
		for _, node := range nodes {
			if ResourcesMatch(task, node) && ConstraintsMatch(task, node) {
				return task, node, true
			}
		}
	}
	return nil, nil, false
}
