package impl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testJob = &registry.Job{
		ID:   "1",
		Name: "testJob",
	}
	testJob2 = &registry.Job{
		ID:   "2",
		Name: "testJob2",
	}
	task1 = &registry.Task{
		ID: "1",
	}
)

func TestCreateJob(t *testing.T) {
	Convey("create job", t, func() {
		err := c.AddJob("1", testJob)
		So(err, ShouldBeNil)
	})
}

func TestGetJob(t *testing.T) {
	Convey("get exist job", t, func() {
		job, err := c.GetJob("1")
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(job, testJob), ShouldBeTrue)
	})

	Convey("get not exist job", t, func() {
		_, err := c.GetNode("2")
		So(err, ShouldNotBeNil)
	})
}

func TestGetAllJobs(t *testing.T) {
	Convey("get all job", t, func() {
		jobs := c.GetAllJobs()
		So(len(jobs), ShouldEqual, 1)
	})
}

func TestUpdateJob(t *testing.T) {
	Convey("update job", t, func() {
		job := &registry.Job{
			ID:   "1",
			Name: "new test",
		}
		err := c.UpdateJob("1", job)
		So(err, ShouldBeNil)
		checkJob, err := c.GetJob("1")
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(job, checkJob), ShouldBeTrue)

	})
}

func TestDeleteJob(t *testing.T) {
	Convey("delete job", t, func() {
		err := c.DeleteJob("1")
		So(err, ShouldBeNil)
		jobs := c.GetAllJobs()
		So(len(jobs), ShouldEqual, 0)
	})
}

func TestGetNotFinishedJobs(t *testing.T) {
	Convey("get not finished jobs", t, func() {
		c.AddJob(testJob.ID, testJob)
		c.AddJob(testJob2.ID, testJob2)
		testJob.PushTask(task1)
		jobs := c.GetNotFinishedJobs()
		So(testJob.IsFinished(), ShouldBeFalse)
		So(len(jobs), ShouldEqual, 1)

		testJob.PopFirstTask()
		jobs = c.GetNotFinishedJobs()
		So(len(jobs), ShouldEqual, 0)
	})
}

func TestGeneratePortToTask(t *testing.T) {
	Convey("generate port to task", t, func() {
		port1 := &registry.Port{
			ContainerPort: 8080,
		}
		port2 := &registry.Port{
			ContainerPort: 8080,
		}
		task1 := &registry.Task{
			ID: "1",
		}
		task2 := &registry.Task{
			ID: "2",
		}
		task1.Ports = append(task.Ports, port1)
		task2.Ports = append(task.Ports, port2)
		testJob.TaskQueue.PushBack(task1)
		testJob.TaskQueue.PushBack(task2)
		task := testJob.PopFirstTask()
		fmt.Println(task.Ports[0].HostPort)
		task.Ports[0].HostPort = 2020
		task = testJob.PopFirstTask()
		So(task.Ports[0].HostPort, ShouldEqual, 0)
	})
}
