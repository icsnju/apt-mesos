package impl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

var (
	minCpus = 0.1
	minMem  = 16.0
)

// ConstraintsMatch check if a offer fit task's constaints
// TODO implementation
func ConstraintsMatch(task *registry.Task, node *registry.Node) bool {
	return true
}

// ResourcesMatch check if a offer fit task's resources
func ResourcesMatch(task *registry.Task, node *registry.Node) bool {
	for _, resource := range task.Resources {
		// the node don't have this resource
		if node.Resources[resource.GetName()] == nil {
			return false
		}

		if resource.GetType().String() == "SCALAR" {
			// check SCALAR type resource
			if node.Resources[resource.GetName()].GetScalar().GetValue() < resource.GetScalar().GetValue() {
				return false
			}
		} else if resource.GetType().String() == "RANGES" {
			// check RANGES type resource
			for _, taskRange := range resource.GetRanges().GetRange() {
				exists := false
				for _, nodeRange := range node.Resources[resource.GetName()].GetRanges().GetRange() {
					if rangeInside(nodeRange, taskRange) {
						exists = true
						break
					}
				}
				if !exists {
					return false
				}
			}
		}
	}
	return true
}

// BuildResources build Resource struct of given resource constraint
// TODO Check whether the resources is enough or not
func BuildResources(task *registry.Task) []*mesosproto.Resource {
	var resources = []*mesosproto.Resource{}
	for _, resource := range task.Resources {
		resources = append(resources, resource)
	}
	return resources
}

// BuildEmptyResources build empty resources
func BuildEmptyResources() []*mesosproto.Resource {
	var resources = []*mesosproto.Resource{}
	resources = append(resources, &mesosproto.Resource{
		Name:   proto.String("cpus"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(minCpus)},
	})
	resources = append(resources, &mesosproto.Resource{
		Name:   proto.String("mem"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(minMem)},
	})
	return resources
}

func RangeUsedUpdate(taskRanges *mesosproto.Value_Ranges, nodeRanges *mesosproto.Value_Ranges) *mesosproto.Value_Ranges {
	var newNodeRanges []*mesosproto.Value_Range
	for _, taskRange := range taskRanges.GetRange() {
		for _, nodeRange := range nodeRanges.GetRange() {
			if rangeInside(nodeRange, taskRange) {
				subedLeftRange, subedRightRange := rangeSub(nodeRange, taskRange)
				if subedLeftRange != nil {
					newNodeRanges = append(newNodeRanges, subedLeftRange)
				}
				if subedRightRange != nil {
					newNodeRanges = append(newNodeRanges, subedRightRange)
				}
			}
		}
	}
	return &mesosproto.Value_Ranges{
		Range: newNodeRanges,
	}
}

func rangeInside(large *mesosproto.Value_Range, small *mesosproto.Value_Range) bool {
	return large.GetBegin() <= small.GetBegin() && large.GetEnd() >= small.GetEnd()
}

func rangeSub(large *mesosproto.Value_Range, small *mesosproto.Value_Range) (*mesosproto.Value_Range, *mesosproto.Value_Range) {
	subedLeftBegin := large.GetBegin()
	subedLeftEnd := small.GetBegin() - 1
	subedLeftRange := &mesosproto.Value_Range{
		Begin: &subedLeftBegin,
		End:   &subedLeftEnd,
	}
	if subedLeftEnd < subedLeftBegin {
		subedLeftRange = nil
	}

	subedRightBegin := small.GetEnd() + 1
	subedRightEnd := large.GetEnd()
	subedRightRange := &mesosproto.Value_Range{
		Begin: &subedRightBegin,
		End:   &subedRightEnd,
	}
	if subedRightBegin > subedRightEnd {
		subedRightRange = nil
	}
	return subedLeftRange, subedRightRange
}
