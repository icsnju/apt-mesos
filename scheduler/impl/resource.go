package impl

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
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
	BuildBasicResources(task)
	for _, resource := range task.Resources {
		// the node don't have this resource
		if node.OfferedResources[resource.GetName()] == nil {
			return false
		}

		if resource.GetType().String() == "SCALAR" {
			// check SCALAR type resource
			if node.OfferedResources[resource.GetName()].GetScalar().GetValue() < resource.GetScalar().GetValue() {
				return false
			}
		} else if resource.GetType().String() == "RANGES" {
			// check RANGES type resource
			for _, taskRange := range resource.GetRanges().GetRange() {
				exists := false
				for _, nodeRange := range node.OfferedResources[resource.GetName()].GetRanges().GetRange() {
					if RangeInside(nodeRange, taskRange) {
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

// BuildBasicResources build basic resources
// including cpus, mem, disk, ports
func BuildBasicResources(task *registry.Task) {
	task.Resources = append(task.Resources, &mesosproto.Resource{
		Name:   proto.String("cpus"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Cpus)},
	})
	task.Resources = append(task.Resources, &mesosproto.Resource{
		Name:   proto.String("mem"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Mem)},
	})
	task.Resources = append(task.Resources, &mesosproto.Resource{
		Name:   proto.String("disk"),
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Disk)},
	})
	ranges := &mesosproto.Value_Ranges{}
	for _, port := range task.Ports {
		ranges.Range = append(ranges.Range, &mesosproto.Value_Range{
			Begin: proto.Uint64(uint64(port.HostPort)),
			End:   proto.Uint64(uint64(port.HostPort)),
		})
	}
	task.Resources = append(task.Resources, &mesosproto.Resource{
		Name:   proto.String("ports"),
		Type:   mesosproto.Value_RANGES.Enum(),
		Ranges: ranges,
	})
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

// BuildResourcesFromMap build resources from a map
func BuildResourcesFromMap(resourceMap map[string]interface{}) map[string]*mesosproto.Resource {
	resources := make(map[string]*mesosproto.Resource)
	for key, value := range resourceMap {
		// Now we only support scalar and ranges
		if reflect.TypeOf(value).Kind() == reflect.Float64 {
			resources[key] = &mesosproto.Resource{
				Name:   proto.String(key),
				Type:   mesosproto.Value_SCALAR.Enum(),
				Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(value.(float64))},
			}
		} else if reflect.TypeOf(value).Kind() == reflect.String {
			ranges, err := ParseRanges(value.(string))
			if err != nil {
				log.Error(err)
			}
			resources[key] = &mesosproto.Resource{
				Name:   proto.String(key),
				Type:   mesosproto.Value_RANGES.Enum(),
				Ranges: ranges,
			}
		}
	}
	return resources
}

func RangeUsedUpdate(taskRanges *mesosproto.Value_Ranges, nodeRanges *mesosproto.Value_Ranges) *mesosproto.Value_Ranges {
	var newNodeRanges []*mesosproto.Value_Range
	for _, taskRange := range taskRanges.GetRange() {
		for _, nodeRange := range nodeRanges.GetRange() {
			if RangeInside(nodeRange, taskRange) {
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

func RangeInside(large *mesosproto.Value_Range, small *mesosproto.Value_Range) bool {
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

func RangeAdd(rangeOne *mesosproto.Value_Ranges, rangeTwo *mesosproto.Value_Ranges) *mesosproto.Value_Ranges {
	line := make([]bool, 65536)
	for _, enumRange := range rangeOne.GetRange() {
		for index := enumRange.GetBegin(); index <= enumRange.GetEnd(); index++ {
			line[index] = true
		}
	}

	for _, enumRange := range rangeTwo.GetRange() {
		for index := enumRange.GetBegin(); index <= enumRange.GetEnd(); index++ {
			line[index] = true
		}
	}

	newRanges := &mesosproto.Value_Ranges{}
	for index := 0; index < len(line); {
		for index < len(line) && !line[index] {
			index++
		}
		if index < len(line) {
			uint64Begin := uint64(index)
			newRange := &mesosproto.Value_Range{
				Begin: &uint64Begin,
			}

			for index < len(line) && line[index] {
				index++
			}
			uint64End := uint64(index - 1)
			newRange.End = &uint64End
			newRanges.Range = append(newRanges.Range, newRange)
		}
	}

	return newRanges
}

func ParseRanges(text string) (*mesosproto.Value_Ranges, error) {
	reg := regexp.MustCompile(`([\d]+)`)
	ports := reg.FindAllString(text, -1)
	ranges := &mesosproto.Value_Ranges{}

	if len(ports)%2 != 0 {
		return nil, errors.New("Cannot parse range string: " + text)
	}

	for i := 0; i < len(ports)-1; i += 2 {
		i64Begin, err := strconv.Atoi(ports[i])
		if err != nil {
			return nil, errors.New("Cannot parse range string: " + text)
		}
		i64End, err := strconv.Atoi(ports[i+1])
		if err != nil {
			return nil, errors.New("Cannot parse range string: " + text)
		}

		ranges.Range = append(ranges.Range, &mesosproto.Value_Range{
			Begin: proto.Uint64(uint64(i64Begin)),
			End:   proto.Uint64(uint64(i64End)),
		})
	}
	return ranges, nil
}
