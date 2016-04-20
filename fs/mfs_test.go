package fs

import (
	"fmt"
	"testing"

	"github.com/icsnju/apt-mesos/fs"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	m = fs.NewMfsFileExplorer()
)

func TestListDir(t *testing.T) {
	Convey("list dir", t, func() {
		list, _ := m.ListDir("/")
		fmt.Println(list)
	})
}
