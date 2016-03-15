package api

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/core"
)

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
Endpoints to get system metrics data

method:		GET
path:		/api/system/metrics
*/
func (api *API) SystemMetrics() martini.Handler{
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		var metrics *core.Metrics
		metrics, states, err := api.core.Metrics()
		if err != nil {
			api.WriteError(w, err)
			return			
		}

		for id, state := range states {
			api.registry.UpdateTaskState(id, state)
		}
		result.Success = true
		result.Result = metrics
		result.Response(w)		
	}	
}

/*
Endpoints to get slave's metrics data

method:		GET
path:		/api/slave/metrics
*/
func (api *API) SlaveMetrics() martini.Handler{
	return func(w http.ResponseWriter, r *http.Request,params martini.Params) {
		var result Result
		var metrics *core.MetricsData
		metrics, err := api.core.GetMetricsData()
		if err != nil {
			api.WriteError(w, err)
			return			
		}

		result.Success = true
		result.Result = metrics.Slaves
		result.Response(w)		
	}	
}
