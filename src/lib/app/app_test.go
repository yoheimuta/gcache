package app

import (
	"testing"

	"lib/index"
	"test/setup"

	"github.com/golang/groupcache"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	command = "1417475105-4-str-HGET-ADINFO-1"
)

func TestApp(t *testing.T) {
	setup.Start(nil)
	defer setup.Term()

	idx := index.NewIndex()
	gcache := groupcache.NewGroup("group", 64<<20, groupcache.GetterFunc(Handle))

	Convey("When origin data is queried", t, func() {

		Convey("redis has no data", func() {
			var ret string
			err := gcache.Get(idx, command, groupcache.StringSink(&ret))

			So(err, ShouldNotBeNil)
			So(ret, ShouldEqual, "")
		})

		Convey("redis has data", func() {
			fixture(idx)

			var ret string
			err := gcache.Get(idx, command, groupcache.StringSink(&ret))

			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "test")

			Convey("data in redis is deleted, but the data in groupcache still exist", func() {
				teardown(idx)

				var ret2 string
				err2 := gcache.Get(idx, command, groupcache.StringSink(&ret2))

				So(err2, ShouldBeNil)
				So(ret2, ShouldEqual, "test")
			})
		})
	})

	Convey("private methods", t, func() {
		Convey("parseKeyString is invoked", func() {
			rettype, cmd, commandArgs, err := parseKeyString(command)
			So(err, ShouldBeNil)
			So(rettype, ShouldEqual, "str")
			So(cmd, ShouldEqual, "HGET")
			So(commandArgs, ShouldResemble, []interface{}{"ADINFO", "1"})
		})
		Convey("convertStrSliceToInterfaceSlice is invoked", func() {
			src := []string{"a", "b", "c"}
			dst := convertStrSliceToInterfaceSlice(src)
			So(dst, ShouldResemble, []interface{}{"a", "b", "c"})
		})
	})
}

func fixture(idx *index.Index) {
	if _, err := idx.Query("int", "hset", []interface{}{"ADINFO", "1", "test"}); err != nil {
		panic(err)
	}
}

func teardown(idx *index.Index) {
	if _, err := idx.Query("int", "del", []interface{}{"ADINFO"}); err != nil {
		panic(err)
	}
}
