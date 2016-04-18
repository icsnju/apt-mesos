package impl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
)

func (core *Core) deleteOffer(offer *mesosproto.Offer) {
	size := len(core.offers)
	for i := 0; i < size; i++ {
		if core.offers[i].GetId() == offer.GetId() {
			core.offers = append(core.offers[:i], core.offers[i+1:]...)
			break
		}
	}
}

func (core *Core) handleOffers(offers []*mesosproto.Offer) []*mesosproto.Offer {
	for _, offer := range offers {
		node, err := core.GetNode(offer.GetSlaveId().GetValue())
		if err != nil {
			continue
		}
		set := &mesosproto.Value_Set{}
		for _, task := range node.GetSLATasks() {
			taskID := task.Parse()
			set.Item = append(set.Item, taskID)
		}
		offer.Attributes = append(offer.Attributes, &mesosproto.Attribute{
			Name: proto.String("SLATasks"),
			Set:  set,
		})
	}
	return offers
}
