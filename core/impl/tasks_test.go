package impl

import (
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	schedulerImpl "github.com/icsnju/apt-mesos/scheduler/impl"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s    = schedulerImpl.NewFCFSScheduler()
	c    = NewCore("192.168.33.1:3030", "192.168.33.10:5050", s)
	task = &registry.Task{
		ID:    "1",
		State: "TASK_WAITING",
	}
)

func TestAddTask(t *testing.T) {
	Convey("add task", t, func() {
		err := c.AddTask("1", task)
		So(err, ShouldBeNil)
	})
}

func TestGetTask(t *testing.T) {
	Convey("get exist task", t, func() {
		task, err := c.GetTask("1")
		So(err, ShouldBeNil)
		So(task, ShouldNotBeNil)
		So(task.ID, ShouldEqual, "1")
	})

	Convey("get not exist task", t, func() {
		_, err := c.GetTask("2")
		So(err, ShouldNotBeNil)
	})
}

func TestGetAllTask(t *testing.T) {
	Convey("get all tasks", t, func() {
		tasks := c.GetAllTasks()
		So(len(tasks), ShouldEqual, 1)
	})
}

func TestUpdateTask(t *testing.T) {
	Convey("update task", t, func() {
		task := &registry.Task{
			ID:       "1",
			DockerID: "2",
		}
		err := c.UpdateTask("1", task)
		So(err, ShouldBeNil)
		checkTask, err := c.GetTask("1")
		So(err, ShouldBeNil)
		So(checkTask.DockerID, ShouldEqual, "2")
	})
}

func TestDeleteTask(t *testing.T) {
	Convey("delete task", t, func() {
		err := c.DeleteTask("1")
		So(err, ShouldBeNil)
		tasks := c.GetAllTasks()
		So(len(tasks), ShouldEqual, 0)
	})
}

func TestGetUnscheduledTask(t *testing.T) {
	Convey("get unshceduled task", t, func() {
		c.AddTask("1", task)
		c.AddTask("2", &registry.Task{
			ID:    "2",
			State: "TASK_STAGING",
		})
		tasks := c.GetUnScheduledTask()
		So(len(tasks), ShouldEqual, 1)
		So(tasks[0].ID, ShouldEqual, "1")
	})
}
