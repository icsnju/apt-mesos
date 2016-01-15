package core

import (
	"fmt"
	"bytes"
	"net/http"

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