package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(id string) interface{}
	Mget(ids []string) []interface{}
	Set(id string, val interface{}, expire time.Duration) error
	Add(id string, val interface{}, expire time.Duration) error
	Delete(id string) error
	Flush() error
	Config(config string) error
}

type Instance func() Cache

var instances = make(map[string]Instance)

func Register(name string, inst Instance) {
	if inst == nil {
		panic("cache instance register failed, instance is nil!")
	}
	if _, exist := instances[name]; !exist {
		instances[name] = inst
	}
}

func NewCache(name string, config string) (adapter Cache, err error) {
	instanceFunc, ok := instances[name]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", name)
		return
	}

	adapter = instanceFunc()
	err = adapter.Config(config)
	if err != nil {
		adapter = nil
	}
	return
}
