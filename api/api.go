package api

import (
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/mesos/mesos-go/mesosproto"
	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/core"
)

var (
	defaultState = mesosproto.TaskState_TASK_STAGING
)

type API struct {
	registry 	*registry.Registry
	core		*core.Core
}

func NewAPI(core *core.Core, registry *registry.Registry) *API{
	return &API{
		core:		core,
		registry: 	registry,
	}
}

/*
Check the connection 

method:     GET
path:       /api/handshake
*/
func (api *API) Handshake() martini.Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        var result Result
        result.Success = true
        result.Result = "OK"

        result.Response(w)
    }
}

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
			writeError(w, err)
			return
		}

		result.Success = true
		result.Result = tasks
		result.Response(w)
	}
}

/*
Submit a tasks

method:		POST
path:		/api/tasks
*/
func (api *API) AddTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result
		task := &registry.Task{State: &defaultState}

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, err)
			return
		}

		//generate task id
		id := make([]byte, 6)
		n, err := rand.Read(id)
		if n != len(id) || err != nil {
			writeError(w, err)
			return
		}
		task.ID = hex.EncodeToString(id)

		err = api.registry.AddTask(task.ID, task)
		if err != nil {
			writeError(w, err)
			return
		}

		result.Success = true
		result.Result = task.ID
		result.Response(w)
	}
}


/*
Delete and kill specific tasks

method:		POST
path:		/api/tasks
*/
func (api *API) DeleteTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		id := params["id"]

		// TODO killTask
		if err := api.registry.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}

		result.Success = true
		result.Result = "OK"
		result.Response(w)		
	}
}

func writeError(w http.ResponseWriter, err error) {
	var result Result
	result.Error = err
	result.Success = false
	result.Response(w)
}