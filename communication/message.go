package communication

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
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
