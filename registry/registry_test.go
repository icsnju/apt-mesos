package registry

import (
	"fmt"
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	r := registry.NewTaskRegistry()
	task := registry.TestTask("1")

	// test add task
	fmt.Println("Test case #1")
	err := r.AddTask(task.ID, task)
	assert.NoError(t, err)

	// test get non-exist task
	fmt.Println("Test case #2")
	task, err = r.GetTask("2")
	assert.Error(t, err)

	// test get exists task
	fmt.Println("Test case #3")
	task, err = r.GetTask("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", task.ID)

	// test update task
	fmt.Println("Test case #4")
	task.Command = "sleep"
	err = r.UpdateTask("1", task)
	assert.NoError(t, err)
	task, err = r.GetTask("1")
	assert.Equal(t, "sleep", task.Command)

	fmt.Println("Test case #5")
	err = r.DeleteTask("1")
	assert.NoError(t, err)
	task, err = r.GetTask("1")
	assert.Error(t, err)
}
