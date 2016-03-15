package registry

import (
	"testing"
	"fmt"
	"encoding/json"

	"github.com/icsnju/apt-mesos/registry"
	"github.com/stretchr/testify/assert"
)

func TestTaskManager(t *testing.T) {
	r := registry.NewRegistry()
	task := registry.TestTask("1")
	
	// test add task
	fmt.Println("Test case #1")
	err := r.AddTask(task.ID, task)
	assert.NoError(t, err)

	// test get non-exist task
	fmt.Println("Test case #2")
	task, err = r.GetTask("2")
	assert.Error(t, err)

	// test get exists task
	fmt.Println("Test case #3")
	task, err = r.GetTask("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", task.ID)

	// test update task
	fmt.Println("Test case #4")
	task.Command = "sleep"
	err = r.UpdateTask("1", task)
	assert.NoError(t, err)
	task, err = r.GetTask("1")
	assert.Equal(t, "sleep", task.Command)

	fmt.Println("Test case #5")
	err = r.DeleteTask("1")
	assert.NoError(t, err)
	task, err = r.GetTask("1")
	assert.Error(t, err)
}

func TestJobManager(t *testing.T) {
	r := registry.NewRegistry()
	var data = `
	{
		"id": "1",
	    "environment": [
	    	{
	    		"name": "adb_server",
	    		"image": "sorccu/adb",
	    		"network": "host",
	    		"port_mappings": [
	    			{
		    			"container_port": 80,
		    			"host_port": 8080,
		    			"protocol": "tcp"
		    		}
	    		],
			    "volumes": [
			        {
			            "container_path": "/data",
			            "host_path": "/vagrant",
			            "mode": "RW"
			        }
			    ],
			    "instance": 1	
	    	}, {
	    		"name": "adb_server",
	    		"image": "sorccu/adb",
	    		"network": "host",
	    		"port_mappings": [],
			    "volumes": [
			        {
			            "container_path": "/data",
			            "host_path": "/vagrant",
			            "mode": "RW"
			        }
			    ],
			    "instance": 4
	    	}
	    ]
	}
	`	
    var job registry.Job;
	err := json.Unmarshal([]byte(data), &job)
	assert.NoError(t, err)	
	job.ID = "1"

	err = r.AddJob(job.ID, &job)
	assert.NoError(t, err)

	newJob, err := r.GetJob(job.ID)
	assert.NoError(t, err)
	assert.Equal(t, newJob.ID, job.ID)

	jobs, err := r.GetAllJobs()
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)

	err = r.DeleteJob(job.ID)
	assert.NoError(t, err)

	jobs, err = r.GetAllJobs()
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 0)	
}
