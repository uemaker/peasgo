package memcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"peasgo/cache"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Memcache struct {
	client  *memcache.Client
	servers []string
}

func NewMemCache() cache.Cache {
	return &Memcache{}
}

func (mc *Memcache) initClient() error {
	if len(mc.servers) == 0 {
		return errors.New("connection servers is empty!")
	}
	mc.client = memcache.New(mc.servers...)
	if mc.client == nil {
		return errors.New("memcache connection failed!")
	}
	return nil
}

func (mc *Memcache) Get(id string) interface{} {
	if mc.client == nil {
		err := mc.initClient()
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	if item, err := mc.client.Get(id); err == nil {
		return string(item.Value)
	}
	return nil
}

func (mc *Memcache) Mget(ids []string) []interface{} {
	count := len(ids)
	var vals []interface{}
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			for i := 0; i < count; i++ {
				vals = append(vals, nil)
			}
			fmt.Println(err)
			return vals
		}
	}
	mvals, err := mc.client.GetMulti(ids)
	if err == nil {
		for _, v := range mvals {
			vals = append(vals, string(v.Value))
		}
		return vals
	}
	for i := 0; i < count; i++ {
		vals = append(vals, nil)
	}
	fmt.Println(err)
	return vals
}

func (mc *Memcache) Set(id string, val interface{}, expire time.Duration) error {
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: id}
	if expire > 0 {
		item.Expiration = int32(expire / time.Second)
	}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only accept string and []byte")
	}
	return mc.client.Set(&item)
}

func (mc *Memcache) Add(id string, val interface{}, expire time.Duration) error {
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: id}
	if expire > 0 {
		item.Expiration = int32(expire / time.Second)
	}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only accept string and []byte")
	}
	return mc.client.Add(&item)
}

func (mc *Memcache) Delete(id string) error {
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			return err
		}
	}
	return mc.client.Delete(id)
}

func (mc *Memcache) Flush() error {
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			return err
		}
	}
	return mc.client.FlushAll()
}

func (mc *Memcache) Config(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["servers"]; !ok {
		return errors.New("config has no conn key")
	}
	mc.servers = strings.Split(cf["servers"], ";")
	if mc.client == nil {
		if err := mc.initClient(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
