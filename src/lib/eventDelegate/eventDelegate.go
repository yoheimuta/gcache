package eventDelegate

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"lib/config"

	"github.com/golang/groupcache"
	"github.com/hashicorp/memberlist"
)

func getPort() string {
	rootConf := config.Instance().Root()
	portConf := rootConf["port"].(map[interface{}]interface{})
	port := portConf["groupcache"].(string)
	return port
}

type eventDelegate struct {
	peers []string
	pool  *groupcache.HTTPPool
}

func (e *eventDelegate) NotifyJoin(node *memberlist.Node) {
	uri := e.groupcacheURI(node.Addr.String())
	e.removePeer(uri)
	e.peers = append(e.peers, uri)
	if e.pool != nil {
		e.pool.Set(e.peers...)
	}
	log.Print("Add peer: " + uri)
	log.Printf("Current peers: %v", e.peers)
}

func (e *eventDelegate) NotifyLeave(node *memberlist.Node) {
	uri := e.groupcacheURI(node.Addr.String())
	e.removePeer(uri)
	e.pool.Set(e.peers...)
	log.Print("Remove peer: " + uri)
	log.Printf("Current peers: %v", e.peers)
}

func (e *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	log.Print("Update the node: %+v\n", node)
}

func (e *eventDelegate) groupcacheURI(addr string) string {
	return fmt.Sprintf("http://%s:%s", addr, getPort())
}

func (e *eventDelegate) removePeer(uri string) {
	for i := 0; i < len(e.peers); i++ {
		if e.peers[i] == uri {
			e.peers = append(e.peers[:i], e.peers[i+1:]...)
			i--
		}
	}
}

func InitGroupCache() {
	eventHandler := &eventDelegate{}
	conf := memberlist.DefaultLANConfig()
	conf.Events = eventHandler
	if addr := os.Getenv("GROUPCACHE_ADDR"); addr != "" {
		conf.AdvertiseAddr = addr
	}

	list, err := memberlist.Create(conf)
	if err != nil {
		panic("Failed to created memberlist: " + err.Error())
	}

	self := list.Members()[0]
	addr := eventHandler.groupcacheURI(string(self.Addr))
	eventHandler.pool = groupcache.NewHTTPPool(addr)
	go http.ListenAndServe(addr, eventHandler.pool)

	if nodes := os.Getenv("JOIN_TO"); nodes != "" {
		if _, err := list.Join(strings.Split(nodes, ",")); err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}
}
