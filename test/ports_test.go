package test

import (
	"fmt"
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMergePort(t *testing.T) {
	Convey("merge port", t, func() {
		var ports []*registry.Port
		ports = append(ports, &registry.Port{
			ContainerPort: 8000,
			HostPort:      8000,
		})
		ports = append(ports, &registry.Port{
			ContainerPort: 8000,
			HostPort:      8001,
		})
		ports = append(ports, &registry.Port{
			ContainerPort: 8000,
			HostPort:      8003,
		})
		ports = append(ports, &registry.Port{
			ContainerPort: 8000,
			HostPort:      8004,
		})
		resource := c.MergePorts(ports)
		fmt.Print(resource)
	})
}
