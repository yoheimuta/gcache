package gserver

import (
	"net"
	"os"
	"testing"
	"time"

	testSetup "test/setup"

	"github.com/hashicorp/memberlist"
	. "github.com/smartystreets/goconvey/convey"
)

type testDelegate struct {
	init   bool
	join   bool
	leave  bool
	update bool
}

func (e *testDelegate) InitGroupCachePool(node *memberlist.Node) {
	e.init = true
}

func (e *testDelegate) NotifyJoin(node *memberlist.Node) {
	e.join = true
}

func (e *testDelegate) NotifyLeave(node *memberlist.Node) {
	e.leave = true
}

func (e *testDelegate) NotifyUpdate(node *memberlist.Node) {
	e.update = true
}

func TestGserver_Start(t *testing.T) {
	testSetup.Start(nil)
	defer testSetup.Term()

	Convey("Invoking Start() should not panic", t, func() {

		So(func() { Start().Shutdown() }, ShouldNotPanic)
	})
}

func TestGserver_setup(t *testing.T) {
	testSetup.Start(nil)
	defer testSetup.Term()

	Convey("When setup() is invoked w/o JOIN_TO", t, func() {

		e := &testDelegate{}
		s := &server{handler: e}
		s.createListConf()

		So(func() { s.setup() }, ShouldNotPanic)
		defer s.Shutdown()

		Convey("InitGroupCachePool should be invoked", func() {
			So(e.init, ShouldBeTrue)
		})

		Convey("NotifyJoin should be invoked", func() {
			So(e.join, ShouldBeTrue)
		})

		Convey("NotifyLeave should not be invoked", func() {
			So(e.leave, ShouldBeFalse)
		})

		Convey("NotifyUpdate should not be invoked", func() {
			So(e.update, ShouldBeFalse)
		})
	})
}

func TestGserver_setup_join(t *testing.T) {
	testSetup.Start(nil)
	defer testSetup.Term()

	// Create a first node
	e := &testDelegate{}
	s := &server{handler: e}
	s.createListConf()

	addr := net.IPv4(127, 0, 0, 1)
	s.listConf.Name = addr.String()
	s.listConf.BindAddr = addr.String()
	s.setup()
	defer s.Shutdown()

	Convey("Cluster has a node", t, func() {

		So(len(s.list.Members()), ShouldEqual, 1)
		So(e.leave, ShouldBeFalse)
	})

	// Create a second node
	os.Setenv("JOIN_TO", s.listConf.BindAddr)
	e2 := &testDelegate{}
	s2 := &server{handler: e2}
	s2.createListConf()

	addr2 := net.IPv4(127, 0, 0, 2)
	s2.listConf.Name = addr2.String()
	s2.listConf.BindAddr = addr2.String()
	s2.setup()
	defer s2.Shutdown()

	Convey("Cluster has two node", t, func() {

		So(len(s.list.Members()), ShouldEqual, 2)
		So(len(s2.list.Members()), ShouldEqual, 2)
		So(e2.leave, ShouldBeFalse)
	})

	// Leave (u could also call Shutdown(), but waiting to sync gossip(suspected and marked dead) requires 10 sec, so instead use Leave())
	s2.list.Leave(time.Second)

	// Wait for leave
	time.Sleep(10 * time.Millisecond)

	Convey("Cluster has a node again", t, func() {

		So(len(s.list.Members()), ShouldEqual, 1)
		So(len(s2.list.Members()), ShouldEqual, 1)
		So(e.leave, ShouldBeTrue)
		So(e2.leave, ShouldBeTrue)
	})
}
