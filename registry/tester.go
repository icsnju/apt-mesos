package registry

type Tester struct {
	CPU            []float64            `json:"cpu"`
	Mem            []float64            `json:"mem"`
	DominantShare  map[string][]float64 `json:"dominant_share"`
	ResponseTime   int64                `json:"response_time"`
	JobLen         int64                `json:"job_size"`
	RunTime        int64                `json:"run_time"`
	MaxTaskRunTime int64                `json:"max_task_run_time"`
	TaskLen        int64                `json:"task_size"`
	TaskRunTime    int64                `json:"task_run_time"`
}

func NewTester() *Tester {
	return &Tester{
		ResponseTime: 0,
		JobLen:       0,
		RunTime:      0,
	}
}
func (t *Tester) AddMetric(m *Metrics) {
	t.CPU = append(t.CPU, m.UsedCpus/(m.FreeCpus+m.UsedCpus))
	t.Mem = append(t.Mem, m.UsedMem/(m.FreeMem+m.UsedMem))
}

func (t *Tester) AddJob(j *Job) {
	t.ResponseTime += (j.StartTime - j.CreateTime)
	t.JobLen++
}

func (t *Tester) AddTask(ta *Task) {
	t.TaskLen++
	runTime := (ta.FinishTime - ta.StartTime)
	t.TaskRunTime += runTime
	if runTime > t.MaxTaskRunTime {
		t.MaxTaskRunTime = runTime
	}
}
