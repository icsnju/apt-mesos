package api

import (
	"time"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

var (
	defaultState = mesosproto.TaskState_TASK_STAGING
)

/*
List all tasks

method:		GET
path:		/api/tasks
*/
func (api *API) ListTasks() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result

		tasks, err := api.registry.GetAllTasks()

		if err != nil {
			api.WriteError(w, err)
			return
		}

		result.Success = true
		result.Result = tasks
		result.Response(w)
	}
}

/*
Submit a task

method:		POST
path:		/api/tasks
*/
func (api *API) AddTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result
		task := &registry.Task{State: &defaultState}

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			api.WriteError(w, err)
			return
		}

		// generate task id
		id := make([]byte, 6)
		n, err := rand.Read(id)
		if n != len(id) || err != nil {
			api.WriteError(w, err)
			return
		}
		task.ID = hex.EncodeToString(id)

		err = api.registry.AddTask(task.ID, task)
		if err != nil {
			api.WriteError(w, err)
			return
		}

		// request for offers
		resources := api.core.BuildResources(task.Cpus, task.Mem, task.Disk)
		offers, err := api.core.RequestOffers(resources)
		if err != nil {
			api.WriteError(w, err)
			return
		}

		// schedule task
		offer, err := api.core.ScheduleTask(offers, resources, task)
		if err != nil {
			api.WriteError(w, err)
			return
		}

		// update task registry
		task.SlaveId = *offer.SlaveId.Value
		task.SlaveHostname, err = api.core.GetSlaveHostname(task.SlaveId)
		if err != nil {
			api.WriteError(w, err)
			return
		}
		task.CreatedTime = time.Now().Unix()*1000

		if err := api.registry.UpdateTask(task.ID, task); err != nil {
			api.WriteError(w, err)
			return 
		}

		// lauch task
		err = api.core.LaunchTask(offer, offers, resources, task)
		if err != nil {
			api.WriteError(w, err)
			return
		}
		
		result.Success = true
		result.Result = task.ID
		result.Response(w)
	}
}

/*
Kill a task which is running

method:		PUT
path:		/api/tasks/:taskId/kill
*/
func (api *API) KillTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		id := params["id"]

		if err := api.core.KillTask(id); err != nil {
			api.WriteError(w, err)
			return
		}

		result.Success = true
		result.Result = "OK"
		result.Response(w)		
	}
}

/*
Delete and kill specific tasks

method:		POST
path:		/api/tasks/:taskId
*/
func (api *API) DeleteTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		id := params["id"]

		if err := api.core.KillTask(id); err != nil {
			api.WriteError(w, err)
			return
		}

		if err := api.registry.DeleteTask(id); err != nil {
			api.WriteError(w, err)
			return
		}

		result.Success = true
		result.Result = "OK"
		result.Response(w)		
	}
}