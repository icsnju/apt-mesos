package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	SingleTask = &Task{
		ID: "test-abcde",
	}
	JobTaskScale = &Task{
		ID:    "test-abcde-defgh-1",
		JobID: "test-abcde",
	}
)

func TestParseSingleTask(t *testing.T) {
	Convey("parse single task", t, func() {
		So(SingleTask.Parse(), ShouldEqual, "test-abcde")
	})
}
func TestParseJobTaskScale(t *testing.T) {
	Convey("parse task of job with scaling", t, func() {
		So(JobTaskScale.Parse(), ShouldEqual, "test-abcde-defgh")
	})
}
