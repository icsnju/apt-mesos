package registry

type PortMapping struct {
	ContainerPort int64 	`json:"container_port,omitempty"`
	HostPort      int64 	`json:"host_port,omitempty"`
	Protocol      string 	`json:"protocol,omitempty"`	
}
