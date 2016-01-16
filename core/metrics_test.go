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

func TestMetrics(t *testing.T) {
	_, err := c.Metrics()
	assert.NoError(t, err)
}

func TestMetricsData(t *testing.T) {
	_, err := c.GetMetricsData()
	assert.NoError(t, err)
}
