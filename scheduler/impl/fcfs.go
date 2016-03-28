package impl

import "github.com/icsnju/apt-mesos/registry"

// FCFSScheduler implements scheduler using FCFS algorithm
type FCFSScheduler struct{}

// NewScheduler create a new scheduler
func NewScheduler() *FCFSScheduler {
	return &FCFSScheduler{}
}

// Schedule implementation
func (scheduler *FCFSScheduler) Schedule(tasks []*registry.Task, nodes []*registry.Node) (*registry.Task, *registry.Node, bool) {
	for _, task := range tasks {
		for _, node := range nodes {
			if ResourcesMatch(task, node) && ConstraintsMatch(task, node) {
				return task, node, true
			}
		}
	}
	return nil, nil, false
}
