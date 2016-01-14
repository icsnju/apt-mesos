package core

import (
	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

type MesosScheduler struct {
	executor      *mesos.ExecutorInfo
	tasksLaunched int
	tasksFinished int
	totalTasks    int
	commands      []string
	cpuPerTask    float64
	memPerTask    float64
}

func NewMesosScheduler(exec *mesos.ExecutorInfo, cpuPerTask float64, memPerTask float64) (*MesosScheduler, error) {
	return nil, nil
}

func (sched *MesosScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Registered with Master ", masterInfo)
}

func (sched *MesosScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Re-Registered with Master ", masterInfo)
}

func (sched *MesosScheduler) Disconnected(sched.SchedulerDriver) {
	log.Infoln("Scheduler Disconnected")
}


func (sched *MesosScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {

}

func (sched *MesosScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())

	// TODO 
	// extract function of taskManager
	if status.GetState() == mesos.TaskState_TASK_FINISHED {
		sched.tasksFinished++
		log.Infof("%v of %v tasks finished.", sched.tasksFinished, sched.totalTasks)
	}

	if sched.tasksFinished >= sched.totalTasks {
		log.Infoln("Total tasks completed, stopping framework.")
		driver.Stop(false)
	}

	if status.GetState() == mesos.TaskState_TASK_LOST ||
		status.GetState() == mesos.TaskState_TASK_KILLED ||
		status.GetState() == mesos.TaskState_TASK_FAILED {
		log.Infoln(
			"Aborting because task", status.TaskId.GetValue(),
			"is in unexpected state", status.State.String(),
			"with message", status.GetMessage(),
		)
		driver.Abort()
	}
}

func (sched *MesosScheduler) OfferRescinded(s sched.SchedulerDriver, id *mesos.OfferID) {
	log.Infof("Offer '%v' rescinded.\n", *id)
}

func (sched *MesosScheduler) FrameworkMessage(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, msg string) {
	log.Infof("Received framework message from executor '%v' on slave '%v': %s.\n", *exId, *slvId, msg)
}

func (sched *MesosScheduler) SlaveLost(s sched.SchedulerDriver, id *mesos.SlaveID) {
	log.Infof("Slave '%v' lost.\n", *id)
}

func (sched *MesosScheduler) ExecutorLost(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, i int) {
	log.Infof("Executor '%v' lost on slave '%v' with exit code: %v.\n", *exId, *slvId, i)
}

func (sched *MesosScheduler) Error(driver sched.SchedulerDriver, err string) {
	log.Infoln("Scheduler received error:", err)
}
