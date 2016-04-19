package docker

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dockerfile = NewDockerfile("123", "../examples/nodejs")
)

func TestParse(t *testing.T) {
	Convey("parse dockerfile", t, func() {
		expectOut :=
			`FROM komljen/ubuntu
MAINTAINER Alen Komljen <alen.komljen@live.com>
RUN \
add-apt-repository -y ppa:chris-lea/node.js && \
apt-get update && \
apt-get -y install \
nodejs && \
rm -rf /var/lib/apt/lists/*
`
		out := dockerfile.Build()
		So(dockerfile, ShouldNotBeNil)
		So(out, ShouldEqual, expectOut)
	})
}

func TestBuildContext(t *testing.T) {
	Convey("build context", t, func() {
		err := dockerfile.BuildContext()
		So(err, ShouldBeNil)
		defer os.RemoveAll("./temp")
	})
}
