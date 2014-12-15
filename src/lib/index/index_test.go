package index

import (
	"testing"

	"test/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndex(t *testing.T) {
	setup.Start(nil)
	defer setup.Term()

	Convey("Invoking constructor should not panic", t, func() {

		So(func() { NewIndex() }, ShouldNotPanic)
	})

	idx := NewIndex()

	Convey("When hset command is sended", t, func() {
		ret, err := idx.Query("int", "hset", []interface{}{"INFO", "1", "test"})

		Convey("the command is valid", func() {
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "1")
		})

		Convey("When hget command is sended", func() {
			ret, err := idx.Query("str", "hget", []interface{}{"INFO", "1"})

			Convey("the command is valid", func() {
				So(err, ShouldBeNil)
				So(ret, ShouldEqual, "test")
			})
		})
	})
}
