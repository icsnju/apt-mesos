package main

import (
	"flag"
	"reflect"

	log "github.com/Sirupsen/logrus"
	comm "github.com/icsnju/apt-mesos/communication"
	core "github.com/icsnju/apt-mesos/core/impl"
	"github.com/icsnju/apt-mesos/scheduler"
	schedulerImpl "github.com/icsnju/apt-mesos/scheduler/impl"
)

var (
	addr     string
	master   string
	strategy string
	debug    bool
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:3030", "Address to listen on <ip:port>")
	flag.StringVar(&master, "master", "127.0.0.1:5050", "Master to connect to <ip:port>")
	flag.StringVar(&strategy, "strategy", "fcfs", "Scheduling algorithm, FCFS default")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.Parse()
}

func main() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// create a new scheduler
	var schedule scheduler.Scheduler
	if strategy == scheduler.FCFS {
		schedule = schedulerImpl.NewFCFSScheduler()
	} else {
		schedule = schedulerImpl.NewDRFScheduler()
	}

	// start a new core
	core := core.NewCore(addr, master, schedule)

	// Start HTTP server
	comm.ListenAndServe(addr, core)
	log.Infof("Current scheduling stragy is %v", reflect.TypeOf(schedule))

	// try to register framework to master
	if err := core.Run(); err != nil {
		log.Fatal(err)
	}

	exit := make(chan bool)
	<-exit
}
