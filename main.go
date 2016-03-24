package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/core"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/server"
)

var (
	addr   string
	master string
	debug  bool

	frameworkName = "apt-mesos"
	user          = "vagrant"
	log           = logrus.New()
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:3030", "Address to listen on <ip:port>")
	flag.StringVar(&master, "master", "127.0.0.1:5050", "Master to connect to <ip:port>")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.Parse()
}

func main() {
	if debug {
		log.Level = logrus.DebugLevel
	}

	// get current hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	webuiURL := "http://" + addr

	// create frameworkInfo
	frameworkInfo := &mesosproto.FrameworkInfo{
		Name:     &frameworkName,
		User:     &user,
		WebuiUrl: &webuiURL,
		Hostname: &hostname,
	}

	// create registry
	registry := registry.NewRegistry()

	// start a new core
	core := core.NewCore(addr, master, frameworkInfo, log)

	// Start HTTP server
	log.Infof("HTTP Server run on %s", addr)
	server.ListenAndServe(addr, registry, core)

	// try to register framework to master
	if err := core.RegisterFramework(); err != nil {
		log.Fatal(err)
	}

	exit := make(chan bool)
	<-exit
}
