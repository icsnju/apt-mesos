package structure

import (
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFCFSQueue(t *testing.T) {
	Convey("fcfs queue sort", t, func() {
		var arr []*registry.Job
		arr = append(arr, &registry.Job{
			ID:         "1",
			CreateTime: 65,
		})
		arr = append(arr, &registry.Job{
			ID:         "2",
			CreateTime: 75,
		})
		arr = append(arr, &registry.Job{
			ID:         "3",
			CreateTime: 45,
		})
		queue := NewFCFSQueue(arr)
		So(queue[0].ID, ShouldEqual, "3")
		So(queue[1].ID, ShouldEqual, "1")
		So(queue[2].ID, ShouldEqual, "2")
	})
}
