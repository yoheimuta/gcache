package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Invoking Instance() should not panic", t, func() {
		So(func() { Instance() }, ShouldNotPanic)
	})

	Convey("Given an Config with a starting value", t, func() {
		c := Instance()

		Convey("The yaml has value of key 'home'", func() {
			So(c.Root()["home"], ShouldEqual, "/gcache/gcache/")
		})
	})
}
