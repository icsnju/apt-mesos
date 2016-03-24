package core

import (
	"github.com/icsnju/apt-mesos/mesosproto"
)

func createScalarResource(name string, value float64) *mesosproto.Resource {
	return &mesosproto.Resource{
		Name:   &name,
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: &value},
	}
}

// ScalarResource count scalar of specific type offer
func ScalarResource(name string, offer *mesosproto.Offer) float64 {
	for _, resource := range offer.GetResources() {
		if resource.GetName() == name {
			return resource.Scalar.GetValue()
		}
	}
	return 0
}

// BuildResources build Resource struct of given resource constraint
// TODO Check whether the resources is enough or not
func (core *Core) BuildResources(cpus, mem, disk float64) []*mesosproto.Resource {
	var resources = []*mesosproto.Resource{}

	if cpus > 0 {
		resources = append(resources, createScalarResource("cpus", cpus))
	}

	if mem > 0 {
		resources = append(resources, createScalarResource("mem", mem))
	}

	if disk > 0 {
		resources = append(resources, createScalarResource("disk", disk))
	}

	return resources
}
