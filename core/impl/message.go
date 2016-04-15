package impl

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/icsnju/apt-mesos/mesosproto"
)

// Message is communication carrier between both side
type Message struct {
	UPID         *UPID
	Name         string
	ProtoMessage proto.Message
	Bytes        []byte
}

// NewMessage return a new message
func NewMessage(upid *UPID, protoMessage proto.Message, bytes []byte) *Message {
	return &Message{
		UPID:         upid,
		Name:         getMessageName(protoMessage),
		ProtoMessage: protoMessage,
		Bytes:        bytes,
	}
}

func (m *Message) RequestURI() string {
	return fmt.Sprintf("/%s/%s", m.UPID.ID, m.Name)
}

func getMessageName(msg proto.Message) string {
	return fmt.Sprintf("%v.%v", "mesos.internal", reflect.TypeOf(msg).Elem().Name())
}

// MesosMasterReachable test whether mesos-master is connectable or not
func MesosMasterReachable(masterAddr string) bool {
	_, err := http.Get("http://" + masterAddr + "/health")
	if err != nil {
		log.Errorf("Failed to connect to mesos %v Error: %v\n", masterAddr, err)
		return false
	}
	return true
}

// SendMessageToMesos is the api to send proto message to mesos-master
func SendMessageToMesos(sender *UPID, message *Message) error {
	log.Debugf("Sending message to %v from %v\n", message.UPID, sender)
	// marshal payload
	b, err := proto.Marshal(message.ProtoMessage)
	if err != nil {
		glog.Errorf("Failed to marshal message %v: %v\n", message, err)
		return err
	}
	message.Bytes = b
	// create request
	req, err := makeLibprocessRequest(sender, message)
	if err != nil {
		glog.Errorf("Failed to make libprocess request: %v\n", err)
		return err
	}
	// send it
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("Failed to POST: %v %v\n", err, resp)
		return err
	}
	resp.Body.Close()
	// ensure master acknowledgement.
	if (resp.StatusCode != http.StatusOK) &&
		(resp.StatusCode != http.StatusAccepted) {
		msg := fmt.Sprintf("Master %s rejected %s.  Returned status %s.", message.UPID, message.RequestURI(), resp.Status)
		glog.Errorln(msg)
		return fmt.Errorf(msg)
	}

	return nil
}

func makeLibprocessRequest(sender *UPID, msg *Message) (*http.Request, error) {
	hostport := net.JoinHostPort(msg.UPID.Host, msg.UPID.Port)
	targetURL := fmt.Sprintf("http://%s%s", hostport, msg.RequestURI())
	log.Debugf("Target URL %s", targetURL)
	req, err := http.NewRequest("POST", targetURL, bytes.NewReader(msg.Bytes))
	if err != nil {
		glog.Errorf("Failed to create request: %v\n", err)
		return nil, err
	}
	req.Header.Add("Libprocess-From", sender.String())
	req.Header.Add("Content-Type", "application/x-protobuf")
	req.Header.Add("Connection", "Keep-Alive")

	return req, nil
}

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
	messagePackage := NewMessage(core.masterUPID, message, nil)
	if err := SendMessageToMesos(core.coreUPID, messagePackage); err != nil {
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
