package src

import "github.com/gomodule/redigo/redis"

//-------------------list操作--------------------
// 返回list长度
func (c *BaseClient) LLEN(key string) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("llen", key))
}

// 遍历list, 按照start, end遍历list, 0表示第一个元素, -1表示最后一个
func (c *BaseClient) LRANGE(key string, start, end int64) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key).Add(start).Add(end)
	return redis.Strings(conn.Do("lrange", args...))
}

// 移除元素, count>0, 从表头向表尾搜索, count<0, 表尾向表头搜, count=0, 移除所有. 返回被移除的个数
func (c *BaseClient) LREM(key string, value interface{}, count int64) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key).Add(count).Add(value)
	return redis.Int64(conn.Do("lrem", args...))
}

// 在尾部插入多个元素. 返回插入后的元素个数
func (c *BaseClient) RPUSH(key string, value ...interface{}) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range value {
		args = args.Add(v)
	}
	return redis.Int64(conn.Do("rpush", args...))
}

// 在头部插入多个元素. 返回插入后的元素个数
func (c *BaseClient) LPUSH(key string, value ...interface{}) (int64, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range value {
		args = args.Add(v)
	}
	return redis.Int64(conn.Do("lpush", args...))
}

// 返回索引对应的元素, -1表示最后一个, -2表示倒数第二个, 从0开始
func (c *BaseClient) LINDEX(key string, index int64) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("LINDEX", key, index))
}
