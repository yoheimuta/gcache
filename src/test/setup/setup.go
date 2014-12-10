package setup

import (
	"bytes"
	"fmt"

	"lib/config"

	"github.com/stvp/tempredis"
)

const (
	REDIS_BIND = "127.0.0.1"
	REDIS_PORT = "6380"
	REDIS_CON  = REDIS_BIND + ":" + REDIS_PORT
)

var redisServer *tempredis.Server = nil

type Config struct {
	SkipMySQL       bool
	SkipMySQLTables bool
}

func newConfig() *Config {
	return &Config{
		SkipMySQL:       true,
		SkipMySQLTables: true,
	}
}

func Start(c *Config) {
	if c == nil {
		c = newConfig()
	}

	if err := setupRedis(); err != nil {
		panic(err)
	}
}

func Term() {
	teardownRedis()
}

func setupRedis() error {
	if r, err := tempredis.Start(
		tempredis.Config{
			"bind":      REDIS_BIND,
			"port":      REDIS_PORT,
			"databases": "8",
		},
	); err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Stderr)
		return fmt.Errorf("Failed to start redis-server w/ error=%v, stderr=%v\n", err, buf.String())
	} else {
		c := config.Instance().Root()
		c = c["redis"].(map[interface{}]interface{})
		ac := c["ad_index_master"].(map[interface{}]interface{})

		ac["server"] = REDIS_CON

		redisServer = r
		fmt.Println("Succeeded to start redis-server")
	}

	return nil
}

func teardownRedis() {
	if redisServer == nil {
		fmt.Println("Failed to term redis-server, because of redisServer is nil")
		return
	}
	if err := redisServer.Term(); err != nil {
		fmt.Sprintf("Failed to term redis-server, because of %v\n", err)
	}
	fmt.Println("Succeeded to term redis-server")
}
