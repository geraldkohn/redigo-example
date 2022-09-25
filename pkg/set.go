package pkg

import (
	"github.com/gomodule/redigo/redis"
)

// ----------------无序集合的操作----------------

// 向无序集合中添加元素
func (c *baseClient) SADD(key, item string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("sadd", key, item))
}

// 移除无序集合的元素
func (c *baseClient) SREM(key string, items ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}
	args.Add(key)
	for _, v := range items {
		args = args.Add(v)
	}
	return redis.Bool(conn.Do("srem", args...))
}

// 返回无序集合包含的所有元素
func (c *baseClient) SMEMBERS(key string) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("smembers", key))
}

// 检查给定元素是否在集合中
func (c *baseClient) SISMEMBER(key, item string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("sismember", key, item))
}
