package impl

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchMetricData(t *testing.T) {
	Convey("fetch metric data", t, func() {
		data, err := c.FetchMetricData()
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		fmt.Println(data)
	})
}
