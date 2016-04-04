package docker

import (
	"fmt"
	"testing"

	"github.com/icsnju/apt-mesos/docker"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("parse dockerfile", t, func() {
		dockerfile := docker.NewDockerfile("../examples/Dockerfile", "icsnju")
		out := dockerfile.Build()
		fmt.Println(out)
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.HasLocalSources(), ShouldBeTrue)
		dockerfile.BuildContext()
	})
}
