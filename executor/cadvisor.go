package main

import (
	"github.com/golang/glog"
	client "github.com/google/cadvisor/client/v2"
)

func main() {
	client, _ := client.NewClient("http://192.168.33.10:18080/")
	mInfo, _ := client.MachineInfo()
	glog.Info(mInfo)

	vInfo, _ := client.VersionInfo()
	glog.Info(vInfo)

	aInfo, _ := client.Attributes()
	glog.Info(aInfo)
}
