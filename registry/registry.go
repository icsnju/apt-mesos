package registry

import (
	"errors"
	"sync"
)

var (
	TaskNotExistsErr = errors.New("Specific task not exist")
)

type Registry struct {
	sync.RWMutex
	tasks 	map[string]*Task
}

func NewRegistry() *Registry {
	return &Registry{
		tasks: make(map[string]*Task),
	}
}

func (registry *Registry) AddTask(id string, task *Task) error {
	registry.Lock()
	defer registry.Unlock()

	registry.tasks[id] = task
	return nil
}

func (registry *Registry) GetTask(id string) (*Task, error) {
	registry.RLock()
	defer registry.RUnlock()

	task, exists := registry.tasks[id]
	if !exists {
		return nil, TaskNotExistsErr
	}

	return task, nil
}

func (registry *Registry) GetAllTasks() ([]*Task, error) {
	registry.RLock()
	defer registry.RUnlock()

	var i = 0
	result := make([]*Task, len(registry.tasks))

	for _, v := range registry.tasks {
		result[i] = v
		i++
	}

	return result, nil
}

func (registry *Registry) DeleteTask(id string) error {
	registry.Lock()
	defer registry.Unlock()

	delete(registry.tasks, id)
	return nil
}

func (registry *Registry) UpdateTask(id string, task *Task) error{
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.tasks[id]
	if !exists {
		return TaskNotExistsErr
	}

	registry.tasks[id] = task	
	return nil
}

