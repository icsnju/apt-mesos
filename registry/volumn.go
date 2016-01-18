package registry

type Volume struct {
	ContainerPath string `json:"container_path,omitempty"`
	HostPath      string `json:"host_path,omitempty"`
	Mode          string `json:"mode,omitempty"`
}