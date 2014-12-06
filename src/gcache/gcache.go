package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/martini"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/groupcache"
	"github.com/hashicorp/memberlist"
)

const GroupcachePort = 8000

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
	return fmt.Sprintf("http://%s:%d", addr, GroupcachePort)
}

func (e *eventDelegate) removePeer(uri string) {
	for i := 0; i < len(e.peers); i++ {
		if e.peers[i] == uri {
			e.peers = append(e.peers[:i], e.peers[i+1:]...)
			i--
		}
	}
}

func initGroupCache() {
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
	addr := fmt.Sprintf("%s:%d", self.Addr, GroupcachePort)
	eventHandler.pool = groupcache.NewHTTPPool("http://" + addr)
	go http.ListenAndServe(addr, eventHandler.pool)

	if nodes := os.Getenv("JOIN_TO"); nodes != "" {
		if _, err := list.Join(strings.Split(nodes, ",")); err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}
}

func main() {
	initGroupCache()
	heavy := groupcache.NewGroup("ad", 64<<20, groupcache.GetterFunc(query))

	c := newCache()

	m := martini.Classic()
	m.Get("/_stats", func() []byte {
		v, err := json.Marshal(&heavy.Stats)
		if err != nil {
			log.Print(err)
		}
		return v
	})
	m.Get("/:key", func(params martini.Params) string {
		var result string
		if err := heavy.Get(c, params["key"], groupcache.StringSink(&result)); err != nil {
			log.Print(err)
		}
		return result
	})
	m.Run()
}

type cache struct {
	mu   sync.Mutex
	conn redis.Conn
}

func newCache() *cache {
	timeout := 10 * time.Second
	server := "localhost:6379"

	if conn, err := redis.DialTimeout(
		"tcp",
		server,
		timeout,
		timeout,
		timeout,
	); err != nil {
		panic(fmt.Sprintf("Failed to connect redis.Conn, because of %v\n", err))
	} else {
		log.Printf("Succeeded to connect redis.Conn:addr=%v\n", server)
		return &cache{conn: conn}
	}
}

func (this *cache) withConn(fn func(redis.Conn) (interface{}, error)) (interface{}, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return fn(this.conn)
}

func query(ctx groupcache.Context, key string, dst groupcache.Sink) error {
	if ctx == nil {
		return fmt.Errorf("nil Context is invalid")
	}

	// ex) [mtime]-[returntype]-[command]-[key]-[field] like 1417475105-str-HGET-ADINFO-1
	parts := strings.SplitN(key, "-", 5)
	if len(parts) < 4 {
		return fmt.Errorf("given key is invalid")
	}
	rettype := parts[1]
	command := parts[2]
	rkey := parts[3]
	var rfield string
	if len(parts) == 5 {
		rfield = parts[4]
	}

	c := ctx.(*cache)
	if got, err := c.withConn(func(conn redis.Conn) (interface{}, error) {
		return conn.Do(command, rkey, rfield)
	}); err != nil {
		log.Printf("Not Found origin data:err=%v, command=%v, rkey=%v, rfield=%v", err, command, rkey, rfield)
		return err
	} else {
		var ret string
		switch rettype {
		case "str":
			ret, err = redis.String(got, nil)
		case "int":
			var ret_int int
			if ret_int, err = redis.Int(got, nil); err == nil {
				ret = strconv.Itoa(ret_int)
			}
		default:
			return fmt.Errorf("Not Found rettype")
		}

		if err != nil {
			return err
		}

		dst.SetString(ret)
		return nil
	}
}
