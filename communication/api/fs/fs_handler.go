package fs

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/fs"
	"github.com/prometheus/common/log"
)

type Handler struct {
	fileExplorer fs.FileExplorer
}

type Request struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

func NewHandler(fe fs.FileExplorer) *Handler {
	return &Handler{
		fileExplorer: fe,
	}
}

func (h *Handler) Handle() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("Cannot decode json: %v", err)
			writeError(w, err)
			return
		}
		switch req.Action {
		case "list":
			list, err := h.fileExplorer.ListDir(req.Path)
			if err != nil {
				writeError(w, err)
			}
			writeResponse(w, http.StatusOK, list)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, content interface{}) {
	data := struct {
		Code   int         `json:"code"`
		Result interface{} `json:"result"`
	}{
		code,
		content,
	}

	message, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)
}

func writeError(w http.ResponseWriter, err error) {
	writeResponse(w, http.StatusInternalServerError, err.Error())
}
