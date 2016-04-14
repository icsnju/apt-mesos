package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
			ContextDir: "../examples/docker_context",
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
