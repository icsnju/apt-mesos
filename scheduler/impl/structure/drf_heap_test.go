package structure

import (
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	job1 = &registry.Job{
		ID: "1",
	}
	job2 = &registry.Job{
		ID: "2",
	}
	job3 = &registry.Job{
		ID: "2",
	}
)

func TestHeap(t *testing.T) {
	Convey("test heap", t, func() {
		// TODO fix bug
		// var h *structure.DRFHeap
		// ele1 := structure.DRFElement{
		// 	DominantResource: &structure.DominantResource{
		// 		Name:  "cpu",
		// 		Share: 0.3,
		// 	},
		// 	Job: job1,
		// }
		// ele2 := structure.DRFElement{
		// 	DominantResource: &structure.DominantResource{
		// 		Name:  "cpu",
		// 		Share: 0.5,
		// 	},
		// 	Job: job2,
		// }
		// ele3 := structure.DRFElement{
		// 	DominantResource: &structure.DominantResource{
		// 		Name:  "cpu",
		// 		Share: 0.1,
		// 	},
		// 	Job: job3,
		// }
		// heap.Init(h)
		// heap.Push(h, ele1)
		// heap.Push(h, ele2)
		// heap.Push(h, ele3)
		// item := heap.Pop(h).(*structure.DRFElement)
		// fmt.Printf("%.2d:%s ", item.Job.ID, item.DominantResource.Share)
	})
}
