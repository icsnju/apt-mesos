package utils

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExists(t *testing.T) {
	Convey("file exists", t, func() {
		exist := Exists("./fs.go")
		So(exist, ShouldBeTrue)
	})
}

func TestFileExists(t *testing.T) {
	Convey("file exists and check is a file", t, func() {
		ok := FileExists("./fs.go")
		So(ok, ShouldBeTrue)

		ok = FileExists("../utils")
		So(ok, ShouldBeFalse)
	})
}

func TestDirExists(t *testing.T) {
	Convey("file exists and check is a directory", t, func() {
		ok := DirExists("./fs.go")
		So(ok, ShouldBeFalse)

		ok = DirExists("../utils")
		So(ok, ShouldBeTrue)
	})
}

func TestCopyFile(t *testing.T) {
	Convey("copy file", t, func() {
		err := CopyFile("./fs.go", "./fs-copy.go")
		So(err, ShouldBeNil)

		ok := Exists("./fs-copy.go")
		So(ok, ShouldBeTrue)

		defer os.Remove("./fs-copy.go")
	})
}

func TestCopyDir(t *testing.T) {
	Convey("copy directory", t, func() {
		err := CopyDir("../utils", "../utils-copy")
		So(err, ShouldBeNil)

		ok := Exists("../utils-copy")
		So(ok, ShouldBeTrue)

		defer os.RemoveAll("../utils-copy")
	})
}

func TestListDir(t *testing.T) {
	ListDir("/mnt/mfsmount")
}
