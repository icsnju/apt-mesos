package fs

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	m = NewMfsFileExplorer()
)

func TestListDir(t *testing.T) {
	Convey("list dir", t, func() {
		list, _ := m.ListDir("/")
		fmt.Println(list)
	})
}
