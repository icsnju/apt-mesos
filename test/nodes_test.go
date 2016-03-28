package test

import (
	"testing"

	"github.com/icsnju/apt-mesos/registry"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	node = &registry.Node{
		ID:       "1",
		Hostname: "node1",
		IP:       "192.168.33.1",
	}
)

func TestRegisterNode(t *testing.T) {
	Convey("register node", t, func() {
		err := c.RegisterNode("1", node)
		So(err, ShouldBeNil)
	})
}

func TestExistNode(t *testing.T) {
	Convey("node exist", t, func() {
		exist := c.ExistsNode("1")
		So(exist, ShouldEqual, true)
	})
}

func TestGetNode(t *testing.T) {
	Convey("get exist node", t, func() {
		node, err := c.GetNode("1")
		So(err, ShouldBeNil)
		So(node.Hostname, ShouldEqual, "node1")
	})

	Convey("get not exist node", t, func() {
		_, err := c.GetNode("2")
		So(err, ShouldNotBeNil)
	})
}

func TestGetAllNodes(t *testing.T) {
	Convey("get all nodes", t, func() {
		nodes := c.GetAllNodes()
		So(len(nodes), ShouldEqual, 1)
	})
}

func TestUpdateNode(t *testing.T) {
	Convey("update node", t, func() {
		node := &registry.Node{
			ID:       "1",
			Hostname: "node2",
			IP:       "192.168.33.2",
		}
		err := c.UpdateNode("1", node)
		So(err, ShouldBeNil)
		checkNode, err := c.GetNode("1")
		So(err, ShouldBeNil)
		So(checkNode.Hostname, ShouldEqual, "node2")
		So(checkNode.LastUpdateTime, ShouldBeGreaterThan, 0)
	})
}

func TestDeleteNode(t *testing.T) {
	Convey("delete node", t, func() {
		err := c.DeleteNode("1")
		So(err, ShouldBeNil)
		nodes := c.GetAllNodes()
		So(len(nodes), ShouldEqual, 0)
	})
}
