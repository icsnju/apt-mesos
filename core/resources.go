package core

import (
	// "github.com/mesos/mesos-go/mesosutil"
	"github.com/icsnju/apt-mesos/mesosproto"
)

func createScalarResource(name string, value float64) *mesosproto.Resource {
	return &mesosproto.Resource{
		Name:   &name,
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: &value},
	}
}

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

// func (core *Core) ScalarResource(name string, resources []*mesosproto.Resource) float64 {
// 	scalarResources := mesosutil.FilterResources(resources, func(res *mesosproto.Resource) bool {
// 		return res.GetType() == mesosproto.Value_SCALAR && res.GetName() == name
// 	})
// 	sum := 0.0
// 	for _, res := range scalarResources {
// 		sum += res.GetScalar().GetValue()
// 	}
// 	return sum
// }

