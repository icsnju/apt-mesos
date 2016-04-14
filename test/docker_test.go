package test

import (
	"testing"

	"github.com/icsnju/apt-mesos/docker"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	dockerfile = docker.NewDockerfile("123", "../examples/docker_context", "")
)

func TestParse(t *testing.T) {
	Convey("parse dockerfile", t, func() {
		expectOut :=
			`FROM scratch
COPY hello /
COPY test_folder/file_in_folder /2
CMD ["/hello"]
`
		out := dockerfile.Build()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.HasLocalSources(), ShouldBeTrue)
		So(out, ShouldEqual, expectOut)
	})
}

func TestBuildContext(t *testing.T) {
	Convey("build context", t, func() {
		err := dockerfile.BuildContext()
		So(err, ShouldBeNil)
	})
}
