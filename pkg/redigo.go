package pkg

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	addr     = "192.168.100.26"
	password = ""
)

// type redisConn redis.Conn

type baseClient struct {
	redisPool *redis.Pool
	address   string
}

func getRedisPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,  //最大空闲连接数，没有redis操作进依然可以保持这个连接数量
		MaxActive:   400, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		Wait:        true,
		IdleTimeout: 120 * time.Second, //空闲连接关闭时间

		Dial: func() (redis.Conn, error) {
			option := redis.DialPassword(password)
			c, err := redis.Dial("tcp", host, option)
			if err != nil {
				return nil, errors.New(err.Error() + host + password)
			}
			if _, err := c.Do("ping"); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { //空闲连接状态检查
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func GetRedisPool() (*baseClient, error) {
	pool := getRedisPool(addr)
	return &baseClient{
		redisPool: pool,
		address:   addr,
	}, nil
}

// 设置锁
func (c *baseClient) Lock(key, value string, ttl int64) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("SET", key, value, "PX", ttl, "NX"))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *baseClient) Unlock(key string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("del", key)
	return err
}

// func (c *BaseClient) GetConn() redis.Conn {
// 	return c.redisPool.Get()
// }

// 设置key的过期时间:秒
func (c *baseClient) ExpireAt(key string, ttl int64) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("expire", key, ttl))
}

/**
The Do method converts command arguments to bulk strings for transmission to the server as follows:

Go Type                 Conversion
[]byte                  Sent as is
string                  Sent as is
int, int64              strconv.FormatInt(v)
float64                 strconv.FormatFloat(v, 'g', -1, 64)
bool                    true -> "1", false -> "0"
nil                     ""
all other types         fmt.Fprint(w, v)

Redis command reply types are represented using the following Go types:

Redis type              Go type
error                   redis.Error
integer                 int64
simple string           string
bulk string             []byte or nil if value not present.
array                   []interface{} or nil if value not present.

Use type assertions or the reply helper functions to convert from interface{} to the specific Go type for the command result.
*/
