package index

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"lib/config"

	"github.com/garyburd/redigo/redis"
)

type Index struct {
	mu   sync.Mutex
	conn redis.Conn
}

func getConfig() (timeout time.Duration, server string) {
	rootConf := config.Instance().Root()
	redisConf := rootConf["redis"].(map[interface{}]interface{})
	adIndexConf := redisConf["master"].(map[interface{}]interface{})

	timeout = time.Duration(adIndexConf["reconnect"].(int)) * time.Second
	server = adIndexConf["server"].(string)
	return timeout, server
}

func NewIndex() *Index {
	timeout, server := getConfig()

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
		return &Index{conn: conn}
	}
}

func (this *Index) Query(rettype, command string, commandArgs []interface{}) (string, error) {
	if got, err := this.withConn(func(conn redis.Conn) (interface{}, error) {
		return conn.Do(command, commandArgs...)
	}); err != nil {
		log.Printf("Not Found origin data:err=%v, command=%v, commandArgs=%v", err, command, commandArgs)
		return "", err
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
			return "", fmt.Errorf("Not Found rettype")
		}

		if err != nil {
			return "", err
		}

		return ret, nil
	}
}

func (this *Index) withConn(fn func(redis.Conn) (interface{}, error)) (interface{}, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return fn(this.conn)
}
