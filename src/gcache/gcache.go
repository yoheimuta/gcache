package main

import (
	"encoding/json"
	"log"

	"lib/app"
	"lib/eventDelegate"
	"lib/index"

	"github.com/codegangsta/martini"
	"github.com/golang/groupcache"
)

const (
	groupName = "group"
	groupSize = 64 << 20
)

func main() {
	eventDelegate.InitGroupCache()
	gcache := groupcache.NewGroup(groupName, groupSize, groupcache.GetterFunc(app.Handle))

	idx := index.NewIndex()

	m := martini.Classic()
	m.Get("/_stats", func() []byte {
		v, err := json.Marshal(&gcache.Stats)
		if err != nil {
			log.Print(err)
		}
		return v
	})
	m.Get("/:key", func(params martini.Params) string {
		var result string
		if err := gcache.Get(idx, params["key"], groupcache.StringSink(&result)); err != nil {
			log.Print(err)
		}
		return result
	})
	m.Run()
}
