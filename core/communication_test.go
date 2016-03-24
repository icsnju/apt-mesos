package core

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/core"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/stretchr/testify/assert"
)

var (
	c *core.Core
)

func init() {
	frameworkName := "api-mesos test"
	user := "tester"
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	log := logrus.New()
	c = NewCore("127.0.0.1:3000", "192.168.33.10:5050", frameworkInfo, log)
}

func TestMesosMasterReachable(t *testing.T) {
	result := c.mesosMasterReachable()
	assert.Equal(t, true, result)
}
