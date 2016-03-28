package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	comm "github.com/icsnju/apt-mesos/communication"
	core "github.com/icsnju/apt-mesos/core/impl"
)

var (
	addr   string
	master string
	debug  bool
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:3030", "Address to listen on <ip:port>")
	flag.StringVar(&master, "master", "127.0.0.1:5050", "Master to connect to <ip:port>")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.Parse()
}

func main() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// start a new core
	core := core.NewCore(addr, master)

	// Start HTTP server
	comm.ListenAndServe(addr, core)

	// try to register framework to master
	if err := core.Run(); err != nil {
		log.Fatal(err)
	}

	exit := make(chan bool)
	<-exit
}
