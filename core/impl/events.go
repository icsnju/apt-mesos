package impl

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
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

		// update node
		core.updateNodesByOffer(event.Offers.Offers)
		core.updateNodesByMetrics()

		// update core.Offers
		offers := core.handleOffers(event.Offers.Offers)
		for _, offer := range offers {
			core.offers = append(core.offers, offer)
		}

	} else if eventType == mesosproto.Event_UPDATE {
		updateStatus := event.GetUpdate().GetStatus().GetState().String()
		if updateStatus == "TASK_FINISHED" {
			task, err := core.GetTask(event.GetUpdate().GetStatus().GetTaskId().GetValue())
			log.Debugf("Task %v finished", task.ID)

			// if task has jobID and is a build task
			// update it's node with new attributes images
			if err == nil && task.JobID != "" && task.Type == registry.TaskTypeBuild {
				job, err1 := core.GetJob(task.JobID)
				node, err2 := core.GetNode(event.GetUpdate().GetStatus().GetSlaveId().GetValue())
				if err1 == nil && err2 == nil {
					image := job.Image
					node.CustomAttributes = append(node.CustomAttributes, &mesosproto.Attribute{
						Name: proto.String("Image"),
						Text: &mesosproto.Value_Text{
							Value: proto.String(image),
						},
					})
				}
			}
		} else if updateStatus == "TASK_FAILED" || updateStatus == "TASK_KILLED" || updateStatus == "TASK_LOST" {
			task, err := core.GetTask(event.GetUpdate().GetStatus().GetTaskId().GetValue())
			log.Debugf("Task %v %s", task.ID, updateStatus)

			// if task has jobID and is failed || killed || lost
			// update its job's status and health
			if err == nil && task.JobID != "" {
				job, err := core.GetJob(task.JobID)
				if err == nil {
					job.Health = registry.UnHealthy
				}
			}
		}

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
