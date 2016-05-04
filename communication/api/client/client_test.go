package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-martini/martini"
	core "github.com/icsnju/apt-mesos/core/impl"
	schedulerImpl "github.com/icsnju/apt-mesos/scheduler/impl"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s      = schedulerImpl.NewFCFSScheduler()
	c      = core.NewCore("192.168.33.1:3030", "192.168.33.10:5050", s)
	h      = NewHandler(c)
	m      = martini.Classic()
	reader io.Reader
)

func TestHandShake(t *testing.T) {
	Convey("handshake api", t, func() {
		m.Get("/handshake", h.Handshake())

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/handshake", nil)
		m.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusOK)
	})
}

func TestAddTask(t *testing.T) {
	Convey("add task api", t, func() {
		m.Post("/api/tasks", h.AddTask())

		taskJSON := `{"name": "busybox","docker_image": "busybox","cpu": 0.5,"mem":16,"cmd": "ls"}`
		reader = strings.NewReader(taskJSON)
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks", reader)
		m.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusOK)
	})
}

func TestListTask(t *testing.T) {
	Convey("list task api", t, func() {
		m.Get("/api/tasks", h.ListTasks())

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks", nil)
		m.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusOK)
	})
}

func TestDeleteNoExistTask(t *testing.T) {
	Convey("delete task api for not exist task", t, func() {
		m.Delete("/api/tasks/:id", h.DeleteTask())

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/tasks/1", nil)
		m.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusInternalServerError)
	})
}
