package utils

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTarFileNotExist(t *testing.T) {
	Convey("tar file not exists", t, func() {
		err := Tar("./fake", "./faketoo", true)
		So(err, ShouldNotBeNil)
	})
}

func TestTarFile(t *testing.T) {
	Convey("tar file src exists and dst isn't exist", t, func() {
		err := Tar("./tar.go", "./tar.tar", false)
		defer os.Remove("./tar.tar")
		So(err, ShouldBeNil)
		So(Exists("./tar.tar"), ShouldBeTrue)
	})
}

func TestTarDir(t *testing.T) {
	Convey("tar directory", t, func() {
		err := Tar("../registry", "./tar.tar", false)
		defer os.Remove("./tar.tar")
		So(err, ShouldBeNil)
		So(Exists("./tar.tar"), ShouldBeTrue)
	})
}

func TestTarFileDstExist(t *testing.T) {
	Convey("tar file exists, but do not overwrite it", t, func() {
		// tar a test file
		err := Tar("./tar.go", "./tar.tar", true)
		So(err, ShouldBeNil)
		defer os.Remove("./tar.tar")

		// overwrite it
		err = Tar("./tar.go", "./tar.tar", false)
		So(err, ShouldBeNil)

		// do not overwrite it
		err = Tar("./tar.go", "./tar.tar", true)
		So(err, ShouldNotBeNil)
	})
}

func TestUnTarFileSrcNotExist(t *testing.T) {
	Convey("untar file which isn't exist", t, func() {
		err := UnTar("fake.tar", "fakepath")
		So(err, ShouldNotBeNil)
	})
}

func TestUnTarFile(t *testing.T) {
	Convey("untar file which dst isn't exist", t, func() {
		err := UnTar("../examples/testing/dockerfile.tar", "fakepath")
		So(err, ShouldBeNil)
		So(Exists("fakepath/Dockerfile"), ShouldBeTrue)
		defer os.RemoveAll("fakepath")
	})
}
