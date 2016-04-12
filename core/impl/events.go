package impl

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/mesosproto"
)

// Events the channel array
type Events map[mesosproto.Event_Type]chan *mesosproto.Event

// NewEvents return a channel bus including events of registered, failure, offers, update.
func NewEvents() Events {
	return Events{
		mesosproto.Event_REGISTERED: make(chan *mesosproto.Event, 64),
		mesosproto.Event_FAILURE:    make(chan *mesosproto.Event, 64),
		mesosproto.Event_OFFERS:     make(chan *mesosproto.Event, 64),
		mesosproto.Event_UPDATE:     make(chan *mesosproto.Event, 64),
	}
}

// AddEvent add event to events bus
func (core *Core) AddEvent(eventType mesosproto.Event_Type, event *mesosproto.Event) error {
	log.WithFields(log.Fields{"type": eventType}).Debug("Received event from master.")
	if eventType == mesosproto.Event_OFFERS {
		log.Debugf("Received %d offer(s).", len(event.Offers.Offers))
		core.updateNodesByOffer(event.Offers.Offers)
		core.updateNodesByMetrics()
	}

	if c, ok := core.events[eventType]; ok {
		c <- event
		return nil
	}
	return fmt.Errorf("unknown event type: %v", eventType)
}

// GetEvent get an event from events bus
func (core *Core) GetEvent(eventType mesosproto.Event_Type) chan *mesosproto.Event {
	if c, ok := core.events[eventType]; ok {
		return c
	}
	return nil
}
