package registry

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	task1 = &Task{
		ID: "1",
	}
	task2 = &Task{
		ID: "2",
	}
)

func TestDockerfileNotExists(t *testing.T) {
	Convey("if dockerfile not exists", t, func() {
		job := &Job{
			ContextDir: "../docker",
		}
		So(job.DockerfileExists(), ShouldBeFalse)
	})
}

func TestDockerfileExists(t *testing.T) {
	Convey("if dockerfile exists", t, func() {
		job := &Job{
			ContextDir: "../examples/nodejs",
		}
		So(job.DockerfileExists(), ShouldBeTrue)
	})
}

func TestNoContextDir(t *testing.T) {
	Convey("if dockerfile has no context directory", t, func() {
		job := &Job{}
		So(job.DockerfileExists(), ShouldBeFalse)
	})
}

func TestJobQueue(t *testing.T) {
	Convey("job queue", t, func() {
		job := &Job{}
		job.PushTask(task1)
		job.PushTask(task2)
		So(job.IsFinished(), ShouldBeFalse)

		task3 := job.PopLastTask()
		So(reflect.DeepEqual(task3, task2), ShouldBeTrue)

		task4 := job.PopFirstTask()
		So(reflect.DeepEqual(task4, task1), ShouldBeTrue)

		So(job.IsFinished(), ShouldBeTrue)
	})
}
