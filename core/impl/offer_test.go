package impl

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	offer1 = &mesosproto.Offer{
		Id: &mesosproto.OfferID{
			Value: proto.String("1"),
		},
	}
	offer2 = &mesosproto.Offer{
		Id: &mesosproto.OfferID{
			Value: proto.String("2"),
		},
	}
)

func TestDeleteOffer(t *testing.T) {
	Convey("delete offer", t, func() {
		c.offers = append(c.offers, offer1)
		c.offers = append(c.offers, offer2)
		So(len(c.offers), ShouldEqual, 2)
		c.deleteOffer(offer1)
		So(len(c.offers), ShouldEqual, 1)
		So(c.offers[0].GetId().GetValue(), ShouldEqual, "2")
	})
}
