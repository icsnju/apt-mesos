package structure

import (
	"sort"

	"github.com/icsnju/apt-mesos/registry"
)

type FCFSQueue []*registry.Job

func (queue FCFSQueue) Len() int {
	return len(queue)
}

func (queue FCFSQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue FCFSQueue) Less(i, j int) bool {
	return queue[i].CreateTime < queue[j].CreateTime
}

// NewFCFSQueue return a task queue in order of create time
func NewFCFSQueue(jobs []*registry.Job) []*registry.Job {
	sort.Sort(FCFSQueue(jobs))
	return jobs
}
