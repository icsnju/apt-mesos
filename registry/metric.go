package registry

// Metrics provide some useful metrics to frontend
type Metrics struct {
	TotalCpus    float64 `json:"total_cpus"`
	TotalMem     float64 `json:"total_mem"`
	TotalDisk    float64 `json:"total_disk"`
	UsedCpus     float64 `json:"used_cpus"`
	UsedMem      float64 `json:"used_mem"`
	UsedDisk     float64 `json:"used_disk"`
	TaskRunning  int64   `json:"task_running"`
	TaskStaging  int64   `json:"task_staging"`
	TaskFinished int64   `json:"task_finished"`
	TaskKilled   int64   `json:"task_killed"`
	//TODO add customed resources
}

// MetricsData is a struct suit for json data from mesos-master
type MetricsData struct {
	Frameworks []struct {
		Tasks []struct {
			ExecutorID string `json:"executor_id"`
			ID         string
			SlaveID    string `json:"slave_id"`
			Resources  struct {
				Cpus float64
				Mem  float64
				Disk float64
			}
			State string `json:"state"`
		}
		CompletedTasks []struct {
			ExecutorID string `json:"executor_id"`
			ID         string
			SlaveID    string `json:"slave_id"`
			State      string `json:"state"`
		} `json:"completed_tasks"`
		ID string
	}
	Slaves []struct {
		ID        string
		PID       string
		Hostname  string
		Resources struct {
			Cpus float64
			Mem  float64
			Disk float64
		}
	}
}
