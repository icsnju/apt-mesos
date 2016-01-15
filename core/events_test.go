package core

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/mesos/mesos-go/mesosproto"
    "github.com/icsnju/apt-mesos/core"
)

var (
	c *core.Core
)

func init() {
	frameworkName := "api-mesos test" 
	user := "tester" 
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	log	:= logrus.New()
	log.Level = logrus.DebugLevel
	c = core.NewCore("127.0.0.1:3000", "192.168.33.10:5050", frameworkInfo, log)	
}

func TestAddEvent(t *testing.T) {
	err := c.AddEvent(mesosproto.Event_FAILURE, &mesosproto.Event{})
	assert.NoError(t, err)
	err = c.AddEvent(mesosproto.Event_SUBSCRIBED, &mesosproto.Event{})
	assert.Error(t, err)
}

func TestGetEvent(t *testing.T) {
	event := c.GetEvent(mesosproto.Event_FAILURE)
	assert.NotNil(t, event)
	event = c.GetEvent(mesosproto.Event_MESSAGE)
	assert.Nil(t, event)
}

