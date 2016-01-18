package core

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"

	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/golang/protobuf/proto"
)

func (core *Core) MesosMasterReachable() bool {
	_, err := http.Get("http://" + core.master + "/health")
	if err != nil {
		core.log.Errorf("Failed to connect to mesos %v Error: %v\n", core.master, err)
		return false
	}
	return true
}

func (core *Core) SendMessageToMesos(msg proto.Message, path string) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/master/%s", core.master, path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-type", "application/octet-stream")
	req.Header.Add("Libprocess-From", fmt.Sprintf("core@%s", core.addr))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	
	if resp != nil && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code %d received while posting to: %s", resp.StatusCode, url)
	}
	return nil
}

func (core *Core) InitEndpoints() {
	core.Endpoints = map[string]map[string]func(w http.ResponseWriter, r *http.Request) error{
		"POST": {
			"/core/mesos.internal.FrameworkRegisteredMessage": core.FrameworkRegisteredMessage,
			"/core/mesos.internal.ResourceOffersMessage": core.ResourceOffersMessage,
		},
	}
}

func (core *Core) FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
    data, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	message := new(mesosproto.FrameworkRegisteredMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	core.log.WithField("frameworkId", message.FrameworkId).Debug("receive framworkId")
	core.frameworkInfo.Id = message.FrameworkId

	eventType := mesosproto.Event_REGISTERED
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

func (core *Core) ResourceOffersMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
    data, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	message := new(mesosproto.ResourceOffersMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}
	eventType := mesosproto.Event_OFFERS
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}