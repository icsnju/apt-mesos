package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/icsnju/apt-mesos/docker"
)

func TestParse(t *testing.T) {
    fmt.Println(1)
	dockerfile := docker.NewDockerfile("../examples/Dockerfile", "icsnju")
    out := dockerfile.Build()
    fmt.Println(out)
    assert.NotNil(t, dockerfile)
}
