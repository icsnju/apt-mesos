package registry

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testItem struct {
	ID    string
	Value string
}

var (
	r     = NewRegistry()
	item1 = &testItem{
		ID:    "1",
		Value: "test",
	}
	item2 = &testItem{
		ID:    "2",
		Value: "test",
	}
)

func TestAdd(t *testing.T) {
	Convey("add item to registry", t, func() {
		err := r.Add(item1.ID, item1)
		So(err, ShouldBeNil)
		err = r.Add(item2.ID, item2)
		So(err, ShouldBeNil)
	})
}

func TestExists(t *testing.T) {
	Convey("registry exists item", t, func() {
		So(r.Exists(item1.ID), ShouldBeTrue)
		So(r.Exists("3"), ShouldBeFalse)
	})
}

func TestGetItem(t *testing.T) {
	Convey("get item", t, func() {
		item := r.Get(item1.ID)
		So(reflect.DeepEqual(item, item1), ShouldBeTrue)
		itemFake := r.Get("3")
		So(itemFake, ShouldBeNil)
	})
}

func TestUpdateItem(t *testing.T) {
	Convey("update item", t, func() {
		newItem := &testItem{
			ID:    "1",
			Value: "new_test",
		}
		err := r.Update("1", newItem)
		So(err, ShouldBeNil)
		itemUpdated := r.Get("1")
		So(reflect.DeepEqual(newItem, itemUpdated), ShouldBeTrue)
	})
}

func TestUpdateNoExistItem(t *testing.T) {
	Convey("update no exist item", t, func() {
		newItem := &testItem{
			ID:    "1",
			Value: "new_test",
		}
		err := r.Update("3", newItem)
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteNoExistItem(t *testing.T) {
	Convey("delete no exist item", t, func() {
		err := r.Delete("3")
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteItem(t *testing.T) {
	Convey("delete item", t, func() {
		err := r.Delete("2")
		So(err, ShouldBeNil)
	})
}

func TestListItems(t *testing.T) {
	Convey("list item", t, func() {
		list := r.List()
		So(len(list), ShouldEqual, 1)
	})
}
