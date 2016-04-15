package impl

import (
	"testing"

	"github.com/icsnju/apt-mesos/mesosproto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAddEvent(t *testing.T) {
	Convey("add events", t, func() {
		err := c.AddEvent(mesosproto.Event_FAILURE, &mesosproto.Event{})
		So(err, ShouldBeNil)
		err = c.AddEvent(mesosproto.Event_RESCIND, &mesosproto.Event{})
		So(err, ShouldNotBeNil)
	})
}

func TestGetEvent(t *testing.T) {
	Convey("get events", t, func() {
		event := c.GetEvent(mesosproto.Event_FAILURE)
		So(event, ShouldNotBeNil)
		event = c.GetEvent(mesosproto.Event_MESSAGE)
		So(event, ShouldBeNil)
	})
}
