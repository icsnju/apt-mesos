package structure

import "github.com/icsnju/apt-mesos/registry"

type DRFElement struct {
	DominantResource *DominantResource
	Job              *registry.Job
}

type DominantResource struct {
	Name  string
	Value float64
	Share float64
}

type DRFHeap []*DRFElement

func (h DRFHeap) Len() int {
	return len(h)
}

func (h DRFHeap) Less(i, j int) bool {
	return h[i].DominantResource.Share < h[j].DominantResource.Share
}

func (h DRFHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *DRFHeap) Push(x interface{}) {
	element := x.(*DRFElement)
	*h = append(*h, element)
}

func (h *DRFHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewDRFElement(job *registry.Job, totalResource map[string]float64) *DRFElement {
	dr := &DominantResource{
		Name:  "None",
		Value: 0.0,
		Share: 0.0,
	}
	for name, resource := range job.UsedResources {
		if resource.GetType().String() == "SCALAR" {
			_, exist := totalResource[name]
			if exist {
				share := float64(resource.GetScalar().GetValue() / totalResource[name])
				if share > dr.Share {
					dr.Name = name
					dr.Value = resource.GetScalar().GetValue()
					dr.Share = share
				}
			}
		}
	}
	return &DRFElement{
		DominantResource: dr,
		Job:              job,
	}
}
