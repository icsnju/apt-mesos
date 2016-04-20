package main

import (
	"fmt"
	"bytes"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	var buf bytes.Buffer
	inputReader, err := os.Open("/vagrant/123.tar")
	opts := docker.BuildImageOptions {
		Name: "test",
		InputStream: inputReader,
		SuppressOutput: true,
		OutputStream: &buf,
	}
	err = client.BuildImage(opts)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
