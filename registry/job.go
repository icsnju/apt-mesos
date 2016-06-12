package registry

import (
	"container/list"
	"path"
	"time"

	"github.com/icsnju/apt-mesos/docker"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/splitter"
	splitterImpl "github.com/icsnju/apt-mesos/splitter/impl"
	"github.com/icsnju/apt-mesos/utils"
)

type Job struct {
	// Basic
	ID   string `json:"id"`
	Name string `json:"name"`

	// Docker Settings
	Image      string             `json:"image"`
	Dockerfile *docker.Dockerfile `json:"dockerfile"`
	ContextDir string             `json:"context_dir"`

	// Task manager
	Tasks        []*Task   `json:"tasks"`
	TaskInstance []*Task   `json:"task_instance"`
	TotalTaskLen int       `json:"total_task_len"`
	TaskQueue    list.List `json:"task_queue"`
	TaskQueueLen int       `json:"task_queue_len"`

	// Timestamp
	CreateTime int64 `json:"create_time"`
	StartTime  int64 `json:"start_time"`
	FinishTime int64 `json:"finish_time"`

	// Splitter
	SplitterType  string `json:"splitter_type"`
	Splitter      splitter.Splitter
	InputPath     string `json:"input_path"`
	OutputPath    string `json:"output_path"`
	WorkDirectory string `json:"work_directory"`

	DominantShare float64 `json:"dominant_share"`

	// Job monitoring
	Health string `json:"health"`
	Status string `json:"status"`

	// Resources
	UsedResources map[string]*mesosproto.Resource `json:"used_resource"`
	SLAOffers     map[string]string
}

const (
	Healthy   = "Healthy"
	UnHealthy = "Unhealthy"
)

const (
	StatusRunning  = "Running"
	StatusFinished = "Finished"
	StatusFailed   = "Failed"
)

const (
	FileSplitterType = "file_splitter"
	LineSplitterType = "line_splitter"
)

func (job *Job) InitBasicParams() error {
	// generate task id
	randID, err := utils.Encode(6)
	if err != nil {
		return err
	}
	job.ID = randID
	job.CreateTime = time.Now().Unix()
	job.SLAOffers = make(map[string]string)
	job.UsedResources = make(map[string]*mesosproto.Resource)
	job.Health = Healthy
	job.Status = StatusRunning

	// init task
	if job.ContextDir != "" {
		job.TotalTaskLen = job.BuildNodeNumber()
	} else {
		job.TotalTaskLen = 0
	}
	for _, task := range job.Tasks {
		job.TotalTaskLen += task.Scale
	}

	// init splitter
	if job.SplitterType == FileSplitterType {
		job.Splitter = splitterImpl.NewFileSplitter()
	} else if job.SplitterType == LineSplitterType {
		job.Splitter = splitterImpl.NewLineSplitter()
	}
	return nil
}

func (job *Job) BuildNodeNumber() int {
	return len(job.Tasks)
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
	task := job.TaskQueue.Front().Value.(*Task)
	return task
}

func (job *Job) LastTask() *Task {
	return job.TaskQueue.Back().Value.(*Task)
}

func (job *Job) PushTask(task *Task) {
	job.TaskQueue.PushBack(task)
	job.TaskInstance = append(job.TaskInstance, task)
	job.TaskQueueLen = job.TaskQueue.Len()
}

func (job *Job) PopFirstTask() *Task {
	task := job.TaskQueue.Front()
	job.TaskQueue.Remove(task)
	job.TaskQueueLen = job.TaskQueue.Len()
	return task.Value.(*Task)
}

func (job *Job) PopLastTask() *Task {
	task := job.TaskQueue.Back()
	job.TaskQueue.Remove(task)
	job.TaskQueueLen = job.TaskQueue.Len()
	return task.Value.(*Task)
}

func (job *Job) IsFinished() bool {
	if job.TaskQueue.Len() == 0 {
		if job.Status == StatusRunning {
			if job.Health == UnHealthy {
				job.Status = StatusFailed
			} else {
				job.Status = StatusFinished
			}
		}
		return true
	}
	return false
}
