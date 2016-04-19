package registry

import (
	"container/list"
	"path"

	"github.com/icsnju/apt-mesos/docker"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/utils"
)

type Job struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Image      string             `json:"image"`
	Dockerfile *docker.Dockerfile `json:"dockerfile"`
	ContextDir string             `json:"context_dir"`
	CreateTime int64              `json:"create_time"`
	Tasks      []*Task            `json:"tasks"`
	TaskQueue  list.List          `json:"task_queue"`
	Splitter   string             `json:"splitter"`
	Input      string             `json:"input"`

	UsedResources map[string]*mesosproto.Resource `json:"used_resource"`
	SLAOffers     map[string]string
}

func (job *Job) DockerfileExists() bool {
	if job.ContextDir != "" {
		dockerfilePath := path.Join(job.ContextDir, "Dockerfile")
		if !utils.Exists(dockerfilePath) {
			return false
		}
		return true
	}
	return false
}

func (job *Job) HasContextDir() bool {
	return job.ContextDir != ""
}

func (job *Job) FirstTask() *Task {
	return job.TaskQueue.Front().Value.(*Task)
}

func (job *Job) LastTask() *Task {
	return job.TaskQueue.Back().Value.(*Task)
}

func (job *Job) PushTask(task *Task) {
	job.TaskQueue.PushBack(task)
}

func (job *Job) PopFirstTask() *Task {
	task := job.TaskQueue.Front()
	job.TaskQueue.Remove(task)
	return task.Value.(*Task)
}

func (job *Job) PopLastTask() *Task {
	task := job.TaskQueue.Back()
	job.TaskQueue.Remove(task)
	return task.Value.(*Task)
}

func (job *Job) IsFinished() bool {
	return job.TaskQueue.Len() == 0
}
