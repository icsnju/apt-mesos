package registry

import (
	"testing"
	"encoding/json"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/stretchr/testify/assert"
)

var data = `
{
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

func TestJob(t *testing.T) {
    var job registry.Job;
	err := json.Unmarshal([]byte(data), &job)
	assert.NoError(t, err)
	assert.Equal(t, job.Environment[0].Name, "adb_server")
}
