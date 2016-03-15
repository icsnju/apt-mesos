package registry

import (
	"errors"
	"sync"
	"github.com/icsnju/apt-mesos/mesosproto"
)

var (
	TaskNotExistsErr = errors.New("Specific task not exist")
	JobNotExistsErr = errors.New("Specific job not exist")
)

type Registry struct {
	sync.RWMutex
	tasks 	map[string]*Task
	jobs	map[string]*Job
}

func NewRegistry() *Registry {
	return &Registry{
		tasks: 	make(map[string]*Task),
		jobs:	make(map[string]*Job),
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

func (registry *Registry) UpdateTaskState(id string, state mesosproto.TaskState) error{
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.tasks[id]
	if !exists {
		return TaskNotExistsErr
	}

	registry.tasks[id].State = &state	
	return nil		
}

func (registry *Registry) AddJob(id string, job *Job) error {
	registry.Lock()
	defer registry.Unlock()

	registry.jobs[id] = job
	return nil
}

func (registry *Registry) GetJob(id string) (*Job, error) {
	registry.RLock()
	defer registry.RUnlock()

	task, exists := registry.jobs[id]
	if !exists {
		return nil, JobNotExistsErr
	}

	return task, nil
}

func (registry *Registry) GetAllJobs() ([]*Job, error) {
	registry.RLock()
	defer registry.RUnlock()

	var i = 0
	result := make([]*Job, len(registry.jobs))

	for _, v := range registry.jobs {
		result[i] = v
		i++
	}

	return result, nil
}

func (registry *Registry) DeleteJob(id string) error {
	registry.Lock()
	defer registry.Unlock()

	delete(registry.jobs, id)
	return nil
}

func (registry *Registry) UpdateJob(id string, job *Job) error{
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.jobs[id]
	if !exists {
		return TaskNotExistsErr
	}

	registry.jobs[id] = job	
	return nil
}
