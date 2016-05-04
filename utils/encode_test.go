package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEcode(t *testing.T) {
	Convey("enode with no error", t, func() {
		str, err := Encode(8)
		So(err, ShouldBeNil)
		So(len(str), ShouldBeGreaterThan, 0)
	})
}
