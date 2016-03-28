package impl

import (
	"errors"

	comm "github.com/icsnju/apt-mesos/communication"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// ErrTaskNotExists defined errors
var (
	ErrTaskNotExists             = errors.New("Specific task not exist")
	ErrBasicResourceNotSatisfied = errors.New("Cpus and mem must be required in resources")
)

type filter func(task *registry.Task) bool

// AddTask is called when user submit a task and add the task to the registry
func (core *Core) AddTask(id string, task *registry.Task) error {
	if err := core.tasks.Add(id, task); err != nil {
		return err
	}
	return nil
}

// GetAllTasks Return all the tasks in registry
func (core *Core) GetAllTasks() []*registry.Task {
	rawList := core.tasks.List()
	tasks := make([]*registry.Task, len(rawList))

	for i, v := range rawList {
		tasks[i] = v.(*registry.Task)
	}
	return tasks
}

// GetTask : Get the task that specified id
func (core *Core) GetTask(id string) (*registry.Task, error) {
	if task := core.tasks.Get(id); task != nil {
		return task.(*registry.Task), nil
	}
	return nil, ErrTaskNotExists
}

// DeleteTask give an id and delete the task
func (core *Core) DeleteTask(id string) error {
	if err := core.tasks.Delete(id); err != nil {
		return err
	}
	return nil
}

// UpdateTask update task info
func (core *Core) UpdateTask(id string, task *registry.Task) error {
	return core.tasks.Update(id, task)
}

// KillTask kill the task
func (core *Core) KillTask(id string) error {
	if task := core.tasks.Get(id); task == nil {
		return ErrTaskNotExists
	}
	frameworkID := core.frameworkInfo.GetId().GetValue()
	message := &mesosproto.KillTaskMessage{
		FrameworkId: &mesosproto.FrameworkID{Value: &frameworkID},
		TaskId:      &mesosproto.TaskID{Value: &id},
	}
	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	return comm.SendMessageToMesos(core.coreUPID, messagePackage)
}

// FilterTask filter task by func
func (core *Core) FilterTask(choose filter) []*registry.Task {
	var result []*registry.Task
	for _, task := range core.tasks.List() {
		if choose(task.(*registry.Task)) {
			result = append(result, task.(*registry.Task))
		}
	}

	return result
}

// GetUnScheduledTask return all un schedulered task
func (core *Core) GetUnScheduledTask() []*registry.Task {
	return core.FilterTask(func(task *registry.Task) bool {
		return task.State == nil
	})
}
