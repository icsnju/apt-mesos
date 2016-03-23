package registry

import (
	"errors"
	"sync"

	"github.com/icsnju/apt-mesos/mesosproto"
)

// Error definitions
var (
	ErrTaskNotExists = errors.New("Specific task not exist")
)

// Registry to manager jobs and tasks
type Registry struct {
	sync.RWMutex
	tasks map[string]*Task
}

// NewRegistry instantiate and return a new Registry
func NewRegistry() *Registry {
	return &Registry{
		tasks: make(map[string]*Task),
	}
}

// AddTask is called when user submit a task and add the task to the registry
func (registry *Registry) AddTask(id string, task *Task) error {
	registry.Lock()
	defer registry.Unlock()

	registry.tasks[id] = task
	return nil
}

// GetTask : Get the task that specified id
func (registry *Registry) GetTask(id string) (*Task, error) {
	registry.RLock()
	defer registry.RUnlock()

	task, exists := registry.tasks[id]
	if !exists {
		return nil, ErrTaskNotExists
	}

	return task, nil
}

// GetAllTasks Return all the tasks in registry
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

// DeleteTask Give an id and delete the task
func (registry *Registry) DeleteTask(id string) error {
	registry.Lock()
	defer registry.Unlock()

	delete(registry.tasks, id)
	return nil
}

// UpdateTask update a task which give the specfic string and a new struct
func (registry *Registry) UpdateTask(id string, task *Task) error {
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.tasks[id]
	if !exists {
		return ErrTaskNotExists
	}

	registry.tasks[id] = task
	return nil
}

// UpdateTaskState update task's state when a task is running
func (registry *Registry) UpdateTaskState(id string, state mesosproto.TaskState) error {
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.tasks[id]
	if !exists {
		return ErrTaskNotExists
	}

	registry.tasks[id].State = &state
	return nil
}
