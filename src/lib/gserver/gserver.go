package gserver

import (
	"log"
	"os"
	"strings"

	"lib/config"
	"lib/eventDelegate"

	"github.com/hashicorp/memberlist"
)

type eventHandler interface {
	InitGroupCachePool(node *memberlist.Node)
	NotifyJoin(node *memberlist.Node)
	NotifyLeave(node *memberlist.Node)
	NotifyUpdate(node *memberlist.Node)
}

type server struct {
	handler  eventHandler
	list     *memberlist.Memberlist
	listConf *memberlist.Config
}

func Start() *server {
	s := &server{
		handler: &eventDelegate.Delegate{},
	}
	s.createListConf()
	s.setup()
	return s
}

func (s server) Shutdown() error {
	if err := s.list.Shutdown(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (s *server) createListConf() {
	conf := memberlist.DefaultLANConfig()
	conf.Events = s.handler
	conf.BindPort = getPort()
	if addr := os.Getenv("GROUPCACHE_ADDR"); addr != "" {
		conf.AdvertiseAddr = addr
	}
	s.listConf = conf
}

func (s *server) setup() {
	list, err := memberlist.Create(s.listConf)
	if err != nil {
		panic("Failed to created memberlist: " + err.Error())
	}
	s.list = list

	self := list.Members()[0]
	s.handler.InitGroupCachePool(self)

	if nodes := os.Getenv("JOIN_TO"); nodes != "" {
		if _, err := list.Join(strings.Split(nodes, ",")); err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}
}

func getPort() int {
	rootConf := config.Instance().Root()
	portConf := rootConf["port"].(map[interface{}]interface{})
	port := portConf["memberlist"].(int)
	return port
}
