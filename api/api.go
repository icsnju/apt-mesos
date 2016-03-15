package api

import (
	"net/http"

	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/core"
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

func (api *API) WriteError(w http.ResponseWriter, err error) {
	var result Result
	result.Error = err
	result.Success = false
	result.Response(w)
}