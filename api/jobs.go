package api

import (
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/registry"
)

/*
Get a job

method: 	GET
path:		/api/job/:jobId
*/

func (api *API) GetJob() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		id := params["id"]

		job, err := api.registry.GetJob(id)
		if err != nil {
			api.WriteError(w, err)
			return 
		}

		result.Success = true
		result.Result = job
		result.Response(w)
	}	
}
/*
List all jobs

method:		GET
path:		/api/jobs
*/
func (api *API) ListJobs() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result

		jobs, err := api.registry.GetAllJobs()

		if err != nil {
			api.WriteError(w, err)
			return
		}

		result.Success = true
		result.Result = jobs
		result.Response(w)
	}
}

/*
Submit a job, describe environment of testing.

method:		POST
path:		/api/jobs
*/
func (api *API) AddJob() martini.Handler{
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result
		job := &registry.Job{}

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
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
		job.ID = hex.EncodeToString(id)

		err = api.registry.AddJob(job.ID, job)
		if err != nil {
			api.WriteError(w, err)
			return
		}		

		result.Success = true
		result.Result = job.ID
		result.Response(w)		
	}	
}

/*
Delete a job

method:		POST
path:		/api/jobs/:jobId
*/
func (api *API) DeleteJob() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		id := params["id"]

		if err := api.registry.DeleteJob(id); err != nil {
			api.WriteError(w, err)
			return
		}

		result.Success = true
		result.Result = "OK"
		result.Response(w)		
	}
}