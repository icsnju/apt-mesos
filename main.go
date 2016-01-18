package main

import (
	"flag"
	
	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/server"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/core"
)

var (
	addr 		string
	master		string

	frameworkName 	= "apt-mesos"
	user			= ""
	log				= logrus.New()
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:3030", "Address to listen on <ip:port>")
	flag.StringVar(&master, "master", "127.0.0.1:5050", "Master to connect to <ip:port>")
	flag.Parse()
}

func main() {
	// create frameworkInfo
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	
	// create registry
	registry := registry.NewRegistry()
	
	// start a new core
	core := core.NewCore(addr, master, frameworkInfo, log)

	// Start HTTP server
	log.Infof("HTTP Server run on %s", addr)
	server.ListenAndServe(addr, registry, core)

	exit := make(chan bool)
	<-exit
}

