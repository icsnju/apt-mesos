package registry

import (
	"sort"
)

type FCFSQueue []*Task

func (queue FCFSQueue) Len() int {
	return len(queue)
}

func (queue FCFSQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue FCFSQueue) Less(i, j int) bool {
	return queue[i].CreatedTime < queue[j].CreatedTime
}

// NewFCFSQueue return a task queue in order of create time
func NewFCFSQueue(tasks []*Task) []*Task {
	sort.Sort(FCFSQueue(tasks))
	return tasks
}
