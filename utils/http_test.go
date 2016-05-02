package utils

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDownload(t *testing.T) {
	Convey("download from remote uri", t, func() {
		_, err := Download("http://fake/path")
		So(err, ShouldBeNil)
		So(Exists("path"), ShouldBeTrue)
		defer os.Remove("path")
	})
}
