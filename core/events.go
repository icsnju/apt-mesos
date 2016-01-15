package core

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mesos/mesos-go/mesosproto"
)

type Events map[mesosproto.Event_Type]chan *mesosproto.Event

func NewEvents() Events {
	return Events{
		mesosproto.Event_FAILURE: 	 make(chan *mesosproto.Event, 64),
		mesosproto.Event_OFFERS:     make(chan *mesosproto.Event, 64),
		mesosproto.Event_UPDATE:     make(chan *mesosproto.Event, 64),
	}
}

func (core *Core) AddEvent(eventType mesosproto.Event_Type, event *mesosproto.Event) error {
	core.log.WithFields(logrus.Fields{"type": eventType}).Debug("Received event from master.")
	if c, ok := core.events[eventType]; ok {
		c <- event
		return nil
	}
	return fmt.Errorf("unknown event type: %v", eventType)
}

func (core *Core) GetEvent(eventType mesosproto.Event_Type) chan *mesosproto.Event {
	if c, ok := core.events[eventType]; ok {
		return c
	} else {
		return nil
	}
}
