package server

import (
    "net/http"


	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/api"
	"github.com/icsnju/apt-mesos/manager"
)

const (
	PORT = ":3030"
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

func ListenAndServe(addr string) {
	m := martini.Classic()
    m.Use(recovery())
    m.Use(martini.Static("static"))
    
	apis := api.NewAPI(manager.NewManager())

    m.Get("/api/handshake", apis.Handshake())
    m.Get("/api/tasks", apis.ListTasks())
    m.Post("/api/tasks", apis.AddTask())
    m.Delete("/api/tasks/:id", apis.DeleteTask())
    
    m.RunOnAddr(PORT)
}
