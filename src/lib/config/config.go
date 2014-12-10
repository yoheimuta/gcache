package config

import (
	"fmt"

	"bindata"

	"github.com/go-yaml/yaml"
)

const (
	path = "data/config.yaml"
)

type Config struct {
	data map[interface{}]interface{}
}

var instance *Config = nil

func newConfig() *Config {
	asset, err := bindata.Asset(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %v, because of %v\n", path, err))
	}
	if len(asset) == 0 {
		panic(fmt.Sprintf("Not Found asset %v\n", path))
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(asset, &data)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config yaml, because of %v\n", err))
	}

	return &Config{
		data: data,
	}
}

func Instance() *Config {
	if instance == nil {
		instance = newConfig()
	}
	return instance
}

func (this *Config) Root() map[interface{}]interface{} {
	return this.data
}
