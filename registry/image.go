package registry

type Image struct {
	Name        string `json:"name"`
	DockerImage string `json:"docker_image"`
	Icon        string `json:"icon"`
}

func NewImage(name, dockerImage, icon string) *Image {
	return &Image{
		Name:        name,
		DockerImage: dockerImage,
		Icon:        icon,
	}
}
