package fs

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/go-martini/martini"
	"github.com/icsnju/apt-mesos/fs"
	"github.com/prometheus/common/log"
)

type Handler struct {
	fileExplorer fs.FileExplorer
}

type Request struct {
	Action      string   `json:"action"`
	Path        string   `json:"path"`
	Item        string   `json:"item"`
	Content     string   `json:"content"`
	NewPath     string   `json:"newPath"`
	Items       []string `json:"items"`
	NewItemPath string   `json:"newItemPath"`
}

func NewHandler(fe fs.FileExplorer) *Handler {
	return &Handler{
		fileExplorer: fe,
	}
}

func (h *Handler) PostHandle() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// upload
		if r.FormValue("destination") != "" {
			path := r.FormValue("destination")
			for _, file := range r.MultipartForm.File {
				err := h.fileExplorer.Upload(path, file[0])
				if err != nil {
					writeError(w, err)
					return
				}
			}
			writeResponse(w, http.StatusOK, true)
			return
		}

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
				return
			}
			writeResponse(w, http.StatusOK, list)
		case "getContent":
			content, err := h.fileExplorer.Cat(req.Item)
			if err != nil {
				writeError(w, err)
				return
			}
			writeResponse(w, http.StatusOK, content)
		case "edit":
			err := h.fileExplorer.Write(req.Item, req.Content)
			if err != nil {
				writeError(w, err)
				return
			}
			writeResponse(w, http.StatusOK, true)
		case "createFolder":
			err := h.fileExplorer.Mkdir(req.NewPath)
			if err != nil {
				writeError(w, err)
				return
			}
			writeResponse(w, http.StatusOK, true)
		case "remove":
			for _, item := range req.Items {
				err := h.fileExplorer.Delete(item)
				if err != nil {
					writeError(w, err)
					return
				}
			}
			writeResponse(w, http.StatusOK, true)
		case "move":
			for _, item := range req.Items {
				_, file := path.Split(item)
				err := h.fileExplorer.Move(item, path.Join(req.NewPath, file))
				if err != nil {
					writeError(w, err)
					return
				}
				writeResponse(w, http.StatusOK, true)
			}
		case "rename":
			err := h.fileExplorer.Rename(req.Item, req.NewItemPath)
			if err != nil {
				writeError(w, err)
				return
			}
			writeResponse(w, http.StatusOK, true)
		}
	}
}

func (h *Handler) GetHandle() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		r.ParseForm()
		action := r.Form["action"][0]
		switch action {
		case "download":
			data, err := h.fileExplorer.Download(r.Form["path"][0])
			if err != nil {
				writeError(w, err)
				return
			}
			ctype := w.Header().Get("Content-Type")
			if ctype == "" {
				ctype = mime.TypeByExtension(filepath.Ext(r.Form["path"][0]))
				if ctype == "" {
					ctype = http.DetectContentType(data)
				}
				w.Header().Set("Content-Type", ctype)
			}
			w.Write(data)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, content interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		Code   int         `json:"code"`
		Result interface{} `json:"result"`
	}{
		code,
		content,
	}

	message, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
	}
	w.Write(message)
}

func writeError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	writeResponse(w, http.StatusInternalServerError, err.Error())
}
