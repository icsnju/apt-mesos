package impl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	f = NewFileSplitter()
	ls = NewLineSplitter()
)

func TestFileSplitter(t *testing.T) {
	Convey("test file splitter", t, func() {
		// a correct path
		args, err := f.Split("./")
		So(err, ShouldBeNil)
		So(len(args), ShouldEqual, 0)

		// a fake path
		_, err = f.Split("/fakepath")
		So(err, ShouldNotBeNil)
	})
}

func TestLineSplitter(t *testing.T) {
	Convey("test line splitter", t, func() {
		// a correct path
		args, err := ls.Split("../../examples/testing/args")
		So(err, ShouldBeNil)
		So(len(args), ShouldEqual, 3)
		So(args[0], ShouldEqual, "1 3")

		// a fake path
		_, err = ls.Split("/fakepath")
		So(err, ShouldNotBeNil)
	})
}
