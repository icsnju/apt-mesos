package server

import (
    "net/http"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/api"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/core"
)

func recovery() martini.Handler {
	return func(w http.ResponseWriter, ctx martini.Context) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}

func createRouter(apis *api.API) martini.Router {
	router := martini.NewRouter()

    router.Get("/api/handshake", apis.Handshake())
    router.Get("/api/tasks", apis.ListTasks())
    router.Post("/api/tasks", apis.AddTask())
    router.Delete("/api/tasks/:id", apis.DeleteTask())	

    return router
}

func ListenAndServe(addr string, registry *registry.Registry, core *core.Core) {
	apis := api.NewAPI(core, registry)
	r := createRouter(apis)

	m := martini.New()
    m.Use(recovery())
    m.Use(martini.Static("static"))
	m.Action(r.Handle)
    m.RunOnAddr(addr)
}
