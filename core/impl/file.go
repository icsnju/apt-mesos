package impl

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
)

// ErrReadFile error when read file
var (
	ErrReadFile = errors.New("Task's slavePID or directory is nil")
)

// ReadFile read file
func (core *Core) ReadFile(id string, filename string) (string, error) {
	task, err := core.GetTask(id)
	if err != nil {
		return "", err
	}

	if task.SlavePID == "" || task.Directory == "" {
		return "", ErrReadFile
	}
	v := url.Values{"path": []string{filepath.Join(task.Directory, filename)}, "offset": []string{"0"}}

	resp, err := http.Get("http://" + task.SlavePID + "/files/read.json?" + v.Encode())
	if err != nil {
		return "", err
	}

	data := struct {
		Data string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	return data.Data, nil
}
