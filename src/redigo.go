package src

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	Addr = "192.168.100.26"
	Password = ""
)

type RedisConn redis.Conn

type BaseClient struct {
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
			option := redis.DialPassword(Password)
			c, err := redis.Dial("tcp", host, option)
			if err != nil {
				return nil, errors.New(err.Error() + host + Password)
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

func GetRedisPool() (*BaseClient, error) {
	pool := getRedisPool(Addr)
	return &BaseClient{
		redisPool: pool,
		address: Addr,
	}, nil
}

// 设置锁
func (c *BaseClient) Lock(key, value string, ttl int64) (bool, error) {
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

func (c *BaseClient) Unlock(key string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("del", key)
	return err
}

func (c *BaseClient) GetConn() redis.Conn {
	return c.redisPool.Get()
}

// 设置key的过期时间:秒
func (c *BaseClient) ExpireAt(key string, ttl int64) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("expire", key, ttl))
}

//--------------字符串的操作----------------
// 判断所在的key是否存在
func (c *BaseClient) EXIST(key string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("exist", key))
}

// 删除键值对
func (c *BaseClient) DEL(key string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("del", key)
	return err
}

// 设置键值对
func (c *BaseClient) SET(key, value string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, value)
	return err
}

// 获得键值对
func (c *BaseClient) GET(key string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

// 自增
func (c *BaseClient) INCR(key string) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("incr", key))
}

// 自减
func (c *BaseClient) DECR(key string) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("decr", key))
}

// 批量获取
func (c *BaseClient) MGET(keys... string) ([]string, error) {
	args := redis.Args{}
	for _, k := range keys {
		args = args.Add(k)
	}
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("mget", args...))
}

// 批量设置
func (c *BaseClient) MSET(keyAndValue... string) (string, error) {
	args := redis.Args{}
	//一项一项加上去 key, value, key, value ...
	for _, kv := range keyAndValue {
		args = args.Add(kv)
	}
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("mset", args...))
}


//-----------------有序集合的操作-------------------
// 向有序集合中添加元素
func (c *BaseClient) ZADD(key, item string, score interface{}) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("zadd", key, score, item))
}

// 移除集合元素
func (c *BaseClient) ZREM(key string, items ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}
	for _, v := range items {
		args = args.Add(v)
	}
	return redis.Bool(conn.Do("zrem", args...))
}

// 集合内元素数量
func (c *BaseClient) ZCARD(key string) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("zcard", key))
}

// 集合交集, 取最小值
// destination: 目的集合的名称
func (c *BaseClient) ZINTERSTOREMIN(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("min")
	return redis.Bool(conn.Do("zinterstore", args...))
}

// 集合交集, 取最大值
func (c *BaseClient) ZINTERSTOREMAX(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("max")
	return redis.Bool(conn.Do("zinterstore", args...))
}

// 集合交集, 取和
func (c *BaseClient) ZINTERSTORESUM(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("sum")
	return redis.Bool(conn.Do("zinterstore", args...))
}

// 集合并集, 取最小值
func (c *BaseClient) ZUNIONSTOREMIN(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("min")
	return redis.Bool(conn.Do("zunionstore", args...))
}

// 集合并集, 取最大值
func (c *BaseClient) ZUNIONSTOREMAX(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("max")
	return redis.Bool(conn.Do("zunionstore", args...))
}

// 集合并集, 取和
func (c *BaseClient) ZUNIONSTORESUM(destination string, keys ...string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(destination).Add(len(keys))
	for _, v := range keys {
		args = args.Add(v)
	}
	args = args.Add("aggregate").Add("sum")
	return redis.Bool(conn.Do("zunionstore", args...))
}

// 获取集合内元素
func (c *BaseClient) ZRANGE(key string, start, end int64, inverted bool) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	command := "zrange"
	if inverted {
		command = "zrevrange"
	}
	args := redis.Args{}.Add(key).Add(start).Add(end)
	return redis.Strings(conn.Do(command, args...))
}

// 获取集合内元素和它的分数
func (c *BaseClient) ZRANGEWITHSCORES(key string, start, end int64, inverted bool) (map[string]int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	command := "zrange"
	if inverted {
		command = "zrevrange"
	}
	args := redis.Args{}.Add(key).Add(start).Add(end).Add("withscores")
	return redis.Int64Map(conn.Do(command, args...))
}


// 获取集合内元素
func (c *BaseClient) ZRANGEBYSCORE(key string, min, max, offset, count int64, inverted bool) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	cmin := strconv.Itoa(int(min))
	cmax := strconv.Itoa(int(max))
	if cmin == "0" {
		cmin = "-inf"
	}
	if cmax == "0" {
		cmax = "+inf"
	}
	if inverted {
		cmin, cmax = cmax, cmin
	}
	command := "zrangebyscore"
	if inverted {
		command = "zrevrangebysocre"
	}
	args := redis.Args{}.Add(key).Add(cmin).Add(cmax).Add("limit").Add(offset).Add(count)
	return redis.Strings(conn.Do(command, args...))
}

// 获取集合内元素带分数
func (c *BaseClient) ZRANGEBYSCOREWITHSCORE(key string, min, max, offset, count int64, inverted bool) (map[string]int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	cmin := strconv.Itoa(int(min))
	cmax := strconv.Itoa(int(max))
	if cmin == "0" {
		cmin = "-inf"
	}
	if cmax == "0" {
		cmax = "+inf"
	}
	if inverted {
		cmin, cmax = cmax, cmin
	}
	command := "ZRANGEBYSCORE"
	if inverted {
		command = "ZREVRANGEBYSCORE"
	}
	args := redis.Args{}.Add(key).Add(cmin).Add(cmax).Add("limit").Add(offset).Add(count)
	return redis.Int64Map(conn.Do(command, args...))
}

// 获取集合内元素的分数值
func (c *BaseClient) ZSCORE(key, item string) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	v, err := redis.Int64(conn.Do("ZSCORE", key, item))
	if err == redis.ErrNil {
		return 0, nil
	}
	return v, err
}

// 移除有序集中，指定分数（score）区间内的所有成员
func (c *BaseClient) ZREMRANGEBYSCORE(key string, min, max int64) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	v, err := redis.Int64(conn.Do("ZREMRANGEBYSCORE", key, min, max))
	if err == redis.ErrNil {
		return 0, nil
	}
	return v, err
}

// 返回有序集合中指定成员的排名
func (c *BaseClient) ZRANK(key, item string, inverted bool) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	command := "ZRANK"
	if inverted {
		command = "ZREVRANK"
	}
	v, err := redis.Int64(conn.Do(command, key, item))
	if err == redis.ErrNil {
		return -1, nil
	}
	return v, err
}

//更新集合内元素的分数值，
func (c *BaseClient) ZINCRBY(key, item string, increment int64) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("ZINCRBY", key, increment, item))
}

