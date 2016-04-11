package test

import (
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFCFSQueue(t *testing.T) {
	Convey("fcfs queue sort", t, func() {
		var arr []*registry.Task
		arr = append(arr, &registry.Task{
			ID:          "1",
			CreatedTime: 65,
		})
		arr = append(arr, &registry.Task{
			ID:          "2",
			CreatedTime: 75,
		})
		arr = append(arr, &registry.Task{
			ID:          "3",
			CreatedTime: 45,
		})
		queue := registry.NewFCFSQueue(arr)
		So(queue[0].ID, ShouldEqual, "3")
		So(queue[1].ID, ShouldEqual, "1")
		So(queue[2].ID, ShouldEqual, "2")
	})
}
