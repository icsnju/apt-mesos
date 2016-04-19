package impl

import (
	"errors"

	"github.com/icsnju/apt-mesos/registry"
)

var (
	ErrJobNotExists = errors.New("Specific job not exist")
)

func (core *Core) AddJob(id string, job *registry.Job) error {
	if err := core.jobs.Add(id, job); err != nil {
		return err
	}
	core.scheduler.AddJob(job)
	return nil
}

func (core *Core) GetAllJobs() []*registry.Job {
	rawList := core.jobs.List()
	jobs := make([]*registry.Job, len(rawList))

	for i, v := range rawList {
		jobs[i] = v.(*registry.Job)
	}
	return jobs
}

func (core *Core) GetJob(id string) (*registry.Job, error) {
	if job := core.jobs.Get(id); job != nil {
		return job.(*registry.Job), nil
	}
	return nil, ErrJobNotExists
}

func (core *Core) DeleteJob(id string) error {
	if err := core.jobs.Delete(id); err != nil {
		return err
	}
	return nil
}

func (core *Core) UpdateJob(id string, job *registry.Job) error {
	return core.jobs.Update(id, job)
}

func (core *Core) updateJobByTask(id string, task *registry.Task) {
	if job, err := core.GetJob(id); err == nil {
		for _, resource := range task.Resources {
			// Only support scalar resource now
			if resource.GetType().String() == "SCALAR" {
				_, ok := job.UsedResources[resource.GetName()]
				if !ok {
					job.UsedResources[resource.GetName()] = resource

				} else {
					newScalar := job.UsedResources[resource.GetName()].GetScalar().GetValue() + resource.GetScalar().GetValue()
					job.UsedResources[resource.GetName()].Scalar.Value = &newScalar
				}
			}
		}
	}
}

func (core *Core) GetNotFinishedJobs() []*registry.Job {
	var result []*registry.Job
	for _, job := range core.GetAllJobs() {
		if !job.IsFinished() {
			result = append(result, job)
		}
	}
	return result
}
