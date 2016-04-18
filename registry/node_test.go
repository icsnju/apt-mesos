package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	node = &Node{
		ID: "1",
		Tasks: []*Task{
			&Task{
				ID:  "task-1",
				SLA: SLAOnePerNode,
			},
			&Task{
				ID: "task-2",
			},
		},
	}
)

func TestGetAllTask(t *testing.T) {
	Convey("get all task of node", t, func() {
		tasks := node.GetTasks()
		So(len(tasks), ShouldEqual, 2)
	})
}

func TestGetSLATask(t *testing.T) {
	Convey("get sla task of node", t, func() {
		tasks := node.GetSLATasks()
		So(len(tasks), ShouldEqual, 1)
		So(tasks[0].SLA, ShouldEqual, SLAOnePerNode)
	})
}
