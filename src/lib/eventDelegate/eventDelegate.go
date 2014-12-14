package eventDelegate

import (
	"fmt"
	"log"
	"net/http"

	"lib/config"

	"github.com/golang/groupcache"
	"github.com/hashicorp/memberlist"
)

type Delegate struct {
	peers []string
	pool  *groupcache.HTTPPool
}

func (e *Delegate) InitGroupCachePool(node *memberlist.Node) {
	uri := groupcacheURI(node.Addr.String())
	e.pool = groupcache.NewHTTPPool(uri)
	go http.ListenAndServe(uri, e.pool)
}

func (e *Delegate) NotifyJoin(node *memberlist.Node) {
	uri := groupcacheURI(node.Addr.String())
	e.removePeer(uri)
	e.peers = append(e.peers, uri)
	if e.pool != nil {
		e.pool.Set(e.peers...)
	}
	log.Print("Add peer: " + uri)
	log.Printf("Current peers: %v", e.peers)
}

func (e *Delegate) NotifyLeave(node *memberlist.Node) {
	uri := groupcacheURI(node.Addr.String())
	e.removePeer(uri)
	if e.pool != nil {
		e.pool.Set(e.peers...)
	}
	log.Print("Remove peer: " + uri)
	log.Printf("Current peers: %v", e.peers)
}

func (e *Delegate) NotifyUpdate(node *memberlist.Node) {
	log.Print("Update the node: %+v\n", node)
}

func (e *Delegate) removePeer(uri string) {
	for i := 0; i < len(e.peers); i++ {
		if e.peers[i] == uri {
			e.peers = append(e.peers[:i], e.peers[i+1:]...)
			i--
		}
	}
}

func groupcacheURI(addr string) string {
	rootConf := config.Instance().Root()
	portConf := rootConf["port"].(map[interface{}]interface{})
	port := portConf["groupcache"].(int)
	return fmt.Sprintf("http://%s:%d", addr, port)
}
