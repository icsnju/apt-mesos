package communication

import (
	"net/http"

	"github.com/go-martini/martini"
	clientAPI "github.com/icsnju/apt-mesos/communication/api/client"
	mesosAPI "github.com/icsnju/apt-mesos/communication/api/mesos"
	"github.com/icsnju/apt-mesos/core"
	"github.com/martini-contrib/cors"
	"github.com/prometheus/common/log"
)

func logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, ctx martini.Context) {
		ctx.Next()
	}
}

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

func createRouter(core core.Core, clientHandlers *clientAPI.Handler, mesosHandlers *mesosAPI.Handler) martini.Router {
	router := martini.NewRouter()

	// create user endpoints
	router.Get("/api/handshake", clientHandlers.Handshake())
	router.Get("/api/tasks", clientHandlers.ListTasks())
	router.Post("/api/tasks", clientHandlers.AddTask())
	router.Delete("/api/tasks/:id", clientHandlers.DeleteTask())
	router.Put("/api/tasks/:id/kill", clientHandlers.KillTask())
	router.Get("/api/tasks/:id/file/:file", clientHandlers.GetFile())
	router.Get("/api/nodes", clientHandlers.GetNodes())

	// create monitor endpoints
	// router.Get("/api/system/metrics", clientHandlers.SystemMetrics())

	// create mesos endpoints
	for method, routes := range mesosHandlers.Endpoints {
		for route, function := range routes {
			switch method {
			case "POST":
				router.Post(route, function)
			case "GET":
				router.Get(route, function)
			case "DELETE":
				router.Delete(route, function)
			case "PUT":
				router.Put(route, function)
			}
		}
	}

	return router
}

// ListenAndServe start the server
func ListenAndServe(addr string, core core.Core) {
	log.Infof("REST listening at: http://%v", core.GetAddr())
	clientHandlers := clientAPI.NewHandler(core)
	mesosHandlers := mesosAPI.NewHandler(core)
	r := createRouter(core, clientHandlers, mesosHandlers)

	m := martini.New()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "x-requested-with", "Content-Type", "Content-Range", "Content-Disposition", "Content-Description"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	m.Use(logger())
	m.Use(recovery())
	m.Use(martini.Static("static"))
	m.Action(r.Handle)
	go m.RunOnAddr(addr)
}
