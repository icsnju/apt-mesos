package impl

import (
	"bytes"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
)

// MesosMasterReachable test whether mesos-master is connectable or not
func (core *Core) MesosMasterReachable() bool {
	_, err := http.Get("http://" + core.master + "/health")
	if err != nil {
		log.Errorf("Failed to connect to mesos %v Error: %v\n", core.master, err)
		return false
	}
	return true
}

// SendMessageToMesos is the api to send proto message to mesos-master
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
	log.Debugf("Sending message to %v", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code %d received while posting to: %s", resp.StatusCode, url)
	}
	return nil
}

// HandleFrameworkRegisteredMessage called when framework register to mesos-master
func (core *Core) HandleFrameworkRegisteredMessage(message *mesosproto.FrameworkRegisteredMessage) {
	log.WithField("frameworkId", message.FrameworkId).Debug("Receive framworkId from mesos master")
	core.frameworkInfo.Id = message.FrameworkId

	eventType := mesosproto.Event_REGISTERED
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	})
}

// HandleResourceOffersMessage called when framework receive offers from master
func (core *Core) HandleResourceOffersMessage(message *mesosproto.ResourceOffersMessage) {
	eventType := mesosproto.Event_OFFERS
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	})
}

// HandleStatusUpdateMessage called when slave's status updated
func (core *Core) HandleStatusUpdateMessage(statusMessage *mesosproto.StatusUpdateMessage) error {
	status := statusMessage.GetUpdate().GetStatus()
	if status.GetState() == mesosproto.TaskState_TASK_RUNNING {
		task, _ := core.GetTask(status.GetTaskId().GetValue())
		core.updateTaskByDockerInfo(task, status.GetData())
	}

	message := &mesosproto.StatusUpdateAcknowledgementMessage{
		FrameworkId: statusMessage.GetUpdate().FrameworkId,
		SlaveId:     statusMessage.GetUpdate().Status.SlaveId,
		TaskId:      statusMessage.GetUpdate().Status.TaskId,
		Uuid:        statusMessage.GetUpdate().Uuid,
	}
	messagePackage := comm.NewMessage(core.masterUPID, message, nil)
	if err := comm.SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
		log.Errorf("Failed to send StatusAccept message: %v\n", err)
		return err
	}

	eventType := mesosproto.Event_UPDATE
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Update: &mesosproto.Event_Update{
			Uuid:   statusMessage.Update.Uuid,
			Status: statusMessage.Update.Status,
		},
	})

	return nil
}
