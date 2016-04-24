package client

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/core"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/utils"
)

// Handler bring out a interface of Restful API
type Handler struct {
	core core.Core
}

// NewHandler return a new API
func NewHandler(core core.Core) *Handler {
	return &Handler{
		core: core,
	}
}

// Handshake check the connection is connectable or not
// method:     GET
// path:       /api/handshake
func (h *Handler) Handshake() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, http.StatusOK, "Connection is OK")
	}
}

// ListTasks list all tasks
// method:		GET
// path:		/api/tasks
func (h *Handler) ListTasks() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := h.core.GetAllTasks()

		writeResponse(w, http.StatusOK, tasks)
	}
}

// AddTask submit a tasks
// method:		POST
// path:		/api/tasks
func (h *Handler) AddTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		task := &registry.Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Errorf("Cannot decode json: %v", err)
			writeError(w, err)
			return
		}

		// generate task id
		randID, err := utils.Encode(6)
		if err != nil {
			writeError(w, err)
			return
		}

		task.ID = "task-" + randID
		task.CreateTime = time.Now().UnixNano()
		task.State = "TASK_WAITING"
		log.Debugf("Receive task: %v", task)

		err = h.core.AddTask(task.ID, task)
		if err != nil {
			writeError(w, err)
			return
		}

		writeResponse(w, http.StatusOK, task.ID)
	}
}

// GetTask return a task
// method:		GET
// path:		/api/task/:id
func (h *Handler) GetTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		id := params["id"]

		task, err := h.core.GetTask(id)
		if err != nil {
			writeError(w, err)
		}

		writeResponse(w, http.StatusOK, task)
	}
}

// KillTask kill a task which is running
// method:		PUT
// path:		/api/tasks/:taskId/kill
func (h *Handler) KillTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		id := params["id"]
		task, err := h.core.GetTask(id)
		if err != nil {
			writeError(w, err)
			return
		}

		if task.State == "TASK_RUNNING" {
			if err := h.core.KillTask(id); err != nil {
				writeError(w, err)
				return
			}
		}
		writeResponse(w, http.StatusOK, "Successful kill the task")
	}
}

// DeleteTask delete and kill specific tasks
// method:		POST
// path:		/api/tasks/:taskId
func (h *Handler) DeleteTask() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		id := params["id"]

		task, err := h.core.GetTask(id)
		if err != nil {
			writeError(w, err)
			return
		}

		if task.State == "TASK_RUNNING" {
			if err := h.core.KillTask(id); err != nil {
				writeError(w, err)
				return
			}
		}

		if err := h.core.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}

		writeResponse(w, http.StatusOK, "Successful deleted task")
	}
}

// CreateJob is the endpoint to create a new job
// method: POST
// path:  /api/jobs
func (h *Handler) CreateJob() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		job := &registry.Job{}
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			log.Errorf("Cannot decode json: %v", err)
			writeError(w, err)
			return
		}

		err := job.InitBasicParams()
		if err != nil {
			writeError(w, err)
			return
		}

		log.WithField("Job", job).Infof("Receive job: %v", job.ID)
		err = h.core.AddJob(job.ID, job)
		if err != nil {
			writeError(w, err)
			return
		}

		err = h.core.StartJob(job)
		if err != nil {
			writeError(w, err)
			return
		}

		writeResponse(w, http.StatusOK, job.ID)
	}
}

// ListJobs list all jobs
// method:		GET
// path:		/api/tasks
func (h *Handler) ListJobs() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := h.core.GetAllJobs()

		writeResponse(w, http.StatusOK, jobs)
	}
}

// GetNodes return all node information
// method: 	Get
// path:   /api/nodes
func (h *Handler) GetNodes() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		node := h.core.GetAllNodes()

		writeResponse(w, http.StatusOK, node)
	}
}

// GetFile get the file content
func (h *Handler) GetFile() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		id := params["id"]
		file := params["file"]

		content, err := h.core.ReadFile(id, file)
		if err != nil {
			writeError(w, err)
			return
		}

		writeResponse(w, http.StatusOK, content)
	}
}

//
func (h *Handler) SystemUsage() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		metric := h.core.GetSystemUsage()

		writeResponse(w, http.StatusOK, metric)
	}
}

func (h *Handler) SystemMetric() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		metric := h.core.GetSystemMetric()

		writeResponse(w, http.StatusOK, metric)
	}
}
