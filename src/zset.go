package src

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

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

