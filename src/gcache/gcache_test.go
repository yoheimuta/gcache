package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"lib/index"
	"test/setup"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	command = "1417475105-4-str-HGET-ADINFO-1"
)

func TestMain(t *testing.T) {
	setup.Start(nil)
	defer setup.Term()

	idx := index.NewIndex()
	port := getPort()

	go main()

	// Wait for listening
	time.Sleep(10 * time.Millisecond)

	Convey("When command request is arrived", t, func() {

		Convey("redis has no data", func() {
			resp, err := http.Get("http://localhost:" + port + "/" + command)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			byteArray, err2 := ioutil.ReadAll(resp.Body)

			So(err2, ShouldBeNil)
			So(string(byteArray), ShouldEqual, "")
		})

		Convey("redis has data", func() {
			fixture(idx)

			resp, err := http.Get("http://localhost:" + port + "/" + command)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			byteArray, err2 := ioutil.ReadAll(resp.Body)

			So(err2, ShouldBeNil)
			So(string(byteArray), ShouldEqual, "test")

			Convey("data in redis is deleted, but the data in groupcache still exist", func() {
				teardown(idx)

				resp, err := http.Get("http://localhost:" + port + "/" + command)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				byteArray, err2 := ioutil.ReadAll(resp.Body)

				So(err2, ShouldBeNil)
				So(string(byteArray), ShouldEqual, "test")
			})
		})
	})

	Convey("When stats request is arrived", t, func() {
		resp, err := http.Get("http://localhost:" + port + "/_stats")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		byteArray, err2 := ioutil.ReadAll(resp.Body)

		So(err2, ShouldBeNil)
		So(string(byteArray), ShouldNotEqual, "")
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
