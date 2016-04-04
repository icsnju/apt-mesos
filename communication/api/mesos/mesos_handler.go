package mesos

import (
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/core"
	"github.com/icsnju/apt-mesos/mesosproto"

	log "github.com/Sirupsen/logrus"
)

// Handler has handlers to communicate with mesos
type Handler struct {
	core core.Core

	Endpoints map[string]map[string]func(w http.ResponseWriter, r *http.Request) error
}

// NewHandler returns a new mesos handler
func NewHandler(core core.Core) *Handler {
	h := &Handler{
		core: core,
	}
	h.Endpoints = map[string]map[string]func(w http.ResponseWriter, r *http.Request) error{
		"POST": {
			"/core/mesos.internal.FrameworkRegisteredMessage": h.FrameworkRegisteredMessage,
			"/core/mesos.internal.ResourceOffersMessage":      h.ResourceOffersMessage,
			"/core/mesos.internal.StatusUpdateMessage":        h.StatusUpdateMessage,
		},
	}
	return h
}

// FrameworkRegisteredMessage is the api that called when framework register to mesos-master
func (h *Handler) FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to unmarshal message %v\n", "mesos.internal.FrameworkRegisteredMessage")
		return err
	}

	message := new(mesosproto.FrameworkRegisteredMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	h.core.HandleFrameworkRegisteredMessage(message)
	w.WriteHeader(http.StatusOK)
	return nil
}

// ResourceOffersMessage is the api that called when framework receive offers from master
func (h *Handler) ResourceOffersMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	message := new(mesosproto.ResourceOffersMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		log.Errorf("Failed to unmarshal message %v\n", "mesos.internal.ResourceOffersMessage")
		return err
	}

	h.core.HandleResourceOffersMessage(message)
	w.WriteHeader(http.StatusOK)
	return nil
}

// StatusUpdateMessage called when slave's status updated
func (h *Handler) StatusUpdateMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	message := &mesosproto.StatusUpdateMessage{}
	if err := proto.Unmarshal(data, message); err != nil {
		log.Errorf("Failed to unmarshal message %v\n", "mesos.internal.StatusUpdateMessage")
		log.Println(string(data))
		return err
	}
	h.core.HandleStatusUpdateMessage(message)
	w.WriteHeader(http.StatusOK)
	return nil
}
