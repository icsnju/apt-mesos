package main

import (
	"fmt"
	"flag"

	"github.com/icsnju/apt-mesos/server"
	"github.com/icsnju/apt-mesos/registry"
	global "github.com/icsnju/apt-mesos/global"
)

const (
	CPUS_PER_TASK = 1
	MEM_PER_TASK  = 128
	defaultArtifactPort = 8000
)

func init() {
	flag.Parse()
}

func main() {
	r := registry.NewRegistry()

	// Start HTTP server
	fmt.Printf("HTTP Server run on %v\n", global.Address)
	server.ListenAndServe(global.Address, r)
}

