package eventDelegate

import (
	"net"
	"testing"

	"test/setup"

	"github.com/hashicorp/memberlist"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEventDelegate(t *testing.T) {
	setup.Start(nil)
	defer setup.Term()

	Convey("InitGroupCachePool() is called", t, func() {
		e := &Delegate{}
		So(e.pool, ShouldBeNil)

		n := &memberlist.Node{Addr: net.IPv4(127, 0, 0, 1)}
		e.InitGroupCachePool(n)
		So(e.pool, ShouldNotBeNil)
	})

	Convey("NotifyJoin() is called", t, func() {
		e := &Delegate{}
		So(e.peers, ShouldBeNil)

		Convey("add first node", func() {
			n := &memberlist.Node{Addr: net.IPv4(127, 0, 0, 1)}
			e.NotifyJoin(n)
			So(e.peers, ShouldResemble, []string{"http://127.0.0.1:40000"})

			Convey("add second node", func() {
				n2 := &memberlist.Node{Addr: net.IPv4(127, 0, 0, 2)}
				e.NotifyJoin(n2)
				So(e.peers, ShouldResemble, []string{"http://127.0.0.1:40000", "http://127.0.0.2:40000"})
			})
		})
	})

	Convey("NotifyLeave() is called", t, func() {
		e := &Delegate{}
		e.peers = []string{"http://127.0.0.1:40000", "http://127.0.0.2:40000"}

		Convey("leave first node", func() {
			n := &memberlist.Node{Addr: net.IPv4(127, 0, 0, 2)}
			e.NotifyLeave(n)
			So(e.peers, ShouldResemble, []string{"http://127.0.0.1:40000"})

			Convey("leave second node", func() {
				n2 := &memberlist.Node{Addr: net.IPv4(127, 0, 0, 1)}
				e.NotifyLeave(n2)
				So(e.peers, ShouldResemble, []string{})
			})
		})
	})

	Convey("private methods", t, func() {
		Convey("removePeer is called", func() {
			e := &Delegate{}
			e.peers = []string{"a", "b", "c"}
			e.removePeer("b")
			So(e.peers, ShouldResemble, []string{"a", "c"})
		})
		Convey("groupcacheURI is called", func() {
			uri := groupcacheURI("127.0.0.1")
			So(uri, ShouldEqual, "http://127.0.0.1:40000")
		})
	})
}
