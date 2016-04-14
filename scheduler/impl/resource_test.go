package impl

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRangeAdd(t *testing.T) {
	Convey("range add", t, func() {
		rangesA := &mesosproto.Value_Ranges{}
		range1 := &mesosproto.Value_Range{
			Begin: proto.Uint64(10),
			End:   proto.Uint64(15),
		}
		range2 := &mesosproto.Value_Range{
			Begin: proto.Uint64(21),
			End:   proto.Uint64(25),
		}
		rangesA.Range = append(rangesA.Range, range1)
		rangesA.Range = append(rangesA.Range, range2)

		rangesB := &mesosproto.Value_Ranges{}
		range3 := &mesosproto.Value_Range{
			Begin: proto.Uint64(16),
			End:   proto.Uint64(19),
		}
		range4 := &mesosproto.Value_Range{
			Begin: proto.Uint64(26),
			End:   proto.Uint64(28),
		}
		rangesA.Range = append(rangesA.Range, range3)
		rangesA.Range = append(rangesA.Range, range4)

		rangeC := RangeAdd(rangesA, rangesB)
		So(rangeC.Range[0].GetBegin(), ShouldEqual, 10)
		So(rangeC.Range[1].GetEnd(), ShouldEqual, 28)
	})
}

func TestRangeParse(t *testing.T) {
	Convey("range parse", t, func() {
		ranges, err := ParseRanges("[31000-31500, 31502-31504, 31505-31200] ")
		So(err, ShouldBeNil)
		So(ranges.Range[0].GetBegin(), ShouldEqual, 31000)
		So(ranges.Range[2].GetEnd(), ShouldEqual, 31200)
	})
}

func TestRangeInside(t *testing.T) {

}
