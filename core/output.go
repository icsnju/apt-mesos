package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
)

// getTaskDirectory return the directory of a task
func (c *Core) getTaskDirectory(slavePID, executorID string) (string, error) {
	resp, err := http.Get("http://" + slavePID + "/state.json")
	if err != nil {
		return "", err
	}

	data := struct {
		Frameworks []struct {
			Executors []struct {
				ID        string
				Directory string
			}
			CompletedExecutors []struct {
				ID        string
				Directory string
			} `json:"completed_executors"`
			ID string
		}
		CompletedFrameworks []struct {
			CompletedExecutors []struct {
				ID        string
				Directory string
			} `json:"completed_executors"`
			ID string
		} `json:"completed_frameworks"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	for _, framework := range data.Frameworks {
		if framework.ID != *c.frameworkInfo.Id.Value {
			continue
		}
		for _, executor := range framework.Executors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
	}

	for _, framework := range data.CompletedFrameworks {
		if framework.ID != *c.frameworkInfo.Id.Value {
			continue
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.ID == executorID {
				return executor.Directory, nil
			}
		}
	}

	return "", nil
}

func (c *Core) readFile(slavePid, directory, filename string) (string, error) {
	v := url.Values{"path": []string{filepath.Join(directory, filename)}, "offset": []string{"0"}}

	resp, err := http.Get("http://" + slavePid + "/files/read.json?" + v.Encode())
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

// ReadFile read all given files' content and return a map
func (c *Core) ReadFile(taskID string, filenames ...string) (map[string]string, error) {
	slavePID, executorID, err := c.getSlavePIDAndExecutorID(taskID)
	if err != nil {
		return nil, err
	}
	if slavePID == "" {
		return nil, fmt.Errorf("cannot get slave PID")
	}
	if executorID == "" {
		executorID = taskID
	}
	directory, err := c.getTaskDirectory(slavePID, executorID)
	if err != nil {
		return nil, err
	}

	var files = make(map[string]string)

	for _, filename := range filenames {
		file, err := c.readFile(slavePID, directory, filename)
		if err != nil {
			return nil, err
		}
		files[filename] = file
	}
	return files, nil
}
