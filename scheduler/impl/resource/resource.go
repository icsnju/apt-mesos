package resource

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

var (
	minCpus = 0.1
	minMem  = 16.0
)

// ResourcesMatch check if a offer fit task's resources
func ResourcesMatch(task *registry.Task, offer *mesosproto.Offer) bool {
	BuildBasicResources(task)
	for _, resource := range task.Resources {
		found := false
		for _, offerResource := range offer.GetResources() {
			if offerResource.GetName() == resource.GetName() {
				found = true
				if resource.GetType().String() == "SCALAR" {

					// check SCALAR type resource
					// resource.scalar.value should be greater than offerResource.scalar.value
					if offerResource.GetScalar().GetValue() < resource.GetScalar().GetValue() {
						return false
					}

				} else if resource.GetType().String() == "RANGES" {

					// check RANGES type resource
					// all ranges should be inside offerRanges
					// if not exist (exist == false) return false
					for _, taskRange := range resource.GetRanges().GetRange() {
						exists := false
						for _, offerRange := range offerResource.GetRanges().GetRange() {
							if RangeInside(offerRange, taskRange) {
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
		}

		// the node don't have this resource
		if !found {
			return false
		}

	}
	// all check pass
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

// BuildBasicResources build basic resources
// including cpus, mem, disk, ports
func BuildBasicResources(task *registry.Task) {
	if task.Build {
		return
	}
	if task.Cpus > 0 {
		task.Resources = append(task.Resources, &mesosproto.Resource{
			Name:   proto.String("cpus"),
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Cpus)},
		})
	}

	if task.Mem > 0 {
		task.Resources = append(task.Resources, &mesosproto.Resource{
			Name:   proto.String("mem"),
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Mem)},
		})
	}

	if task.Disk > 0 {
		task.Resources = append(task.Resources, &mesosproto.Resource{
			Name:   proto.String("disk"),
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: proto.Float64(task.Disk)},
		})
	}

	if len(task.Ports) > 0 {
		ranges := &mesosproto.Value_Ranges{}
		for _, port := range task.Ports {
			if port.HostPort == 0 {
				continue
			}
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
	task.Build = true
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

func GeneratePort(resources []*mesosproto.Resource) uint32 {
	for _, resource := range resources {
		if resource.GetName() == "ports" {
			for _, r := range resource.GetRanges().GetRange() {
				return uint32(GetPointOfRange(r))
			}
		}
	}
	return 0
}
