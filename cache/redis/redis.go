package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"peasgo/cache"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	client   *redis.Pool
	dbNum    int
	key      string
	password string
	server   string
}

func NEWRedis() cache.Cache {
	return &Redis{}
}

func (rds *Redis) initClient() error {
	if len(rds.server) == 0 {
		return errors.New("connection server is empty!")
	}
	rds.client = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rds.server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", rds.password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", rds.dbNum); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		}
	}

	if rds.client == nil {
		return errors.New("redis connection failed!")
	}
	return nil
}

func (rds *Redis) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rds.associate(args[0])
	c := rds.client.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (rds *Redis) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rds.key, originKey)
}

func (rds *Redis) Get(id string) interface{} {
	if v, err := rds.do("GET", key); err == nil {
		return v
	}
	return nil
}

func (rds *Redis) Mget(ids []string) []interface{} {
	c := rds.client.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

func (rds *Redis) Set(id string, val interface{}, expire time.Duration) error {
	_, err := rds.do("SETEX", key, int64(expire/time.Second), val)
	return err
}

func (rds *Redis) Add(id string, val interface{}, expire time.Duration) error {
	_, err := rds.do("APPEND", key, val)
	return err
}

func (rds *Redis) Delete(id string) error {
	if rds.client == nil {
		if err := rds.initClient(); err != nil {
			return err
		}
	}
	return rds.client.Delete(id)
}

func (rds *Redis) Flush() error {
	if rds.client == nil {
		if err := rds.initClient(); err != nil {
			return err
		}
	}
	return rds.client.FlushAll()
}

func (rds *Redis) Config(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["server"]; !ok {
		return errors.New("config has no server key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	rds.key = cf["key"]
	rds.server = cf["server"]
	rds.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rds.password = cf["password"]

	rds.initClient()

	c := rds.client.Get()
	defer c.Close()

	return c.Err()
}

func init() {
	cache.Register("redis", NEWRedis)
}
