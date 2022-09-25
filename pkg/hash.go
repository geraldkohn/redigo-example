package pkg

import "github.com/gomodule/redigo/redis"

//----------------hash操作----------------

//字段赋值，旧值会被覆盖，设置成功返回1，被覆盖了返回0
func (c *baseClient) HSET(key, field string, value interface{}) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HSET", key, field, value))
}

//字段赋值，旧值会被覆盖，设置成功返回1，被覆盖了返回0
func (c *baseClient) HGET(key, field string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	v, err := redis.String(conn.Do("HGET", key, field))
	if err == redis.ErrNil {
		return "", nil
	}
	return v, err
}

// 删除字段, 删除成功返回1, 删除失败返回0
func (c *baseClient) HDEL(key, field string) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HDEL", key, field))
}

//命令用于查找所有符合给定模式 pattern 的 key 。。
func (c *baseClient) KEYS(pattern string) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("KEYS", pattern))
}

//获取给定多个字段的值
func (c *baseClient) HMGET(key string, field ...string) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range field {
		args = args.Add(v)
	}
	return redis.Strings(conn.Do("HMGET", args...))
}

//设置给定多个字段的值
func (c *baseClient) HMSET(key string, fieldAndValue ...interface{}) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range fieldAndValue {
		args = args.Add(v)
	}
	return redis.String(conn.Do("HMSET", args...))
}

//通过结构体设置哈希
func (c *baseClient) HMSETByStruct(key string, dest interface{}) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(dest)...))
}

//为哈希表 key 中的指定字段的整数值加上增量 increment 。
func (c *baseClient) HINCRBY(key, field string, increment int) (int, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HINCRBY", key, field, increment))
}

//返回key下所有的字段和值
func (c *baseClient) HGETALL(key string, dest interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	value, _ := redis.Values(conn.Do("HGETALL", key))
	return redis.ScanStruct(value, dest)
}

type Converter interface {
	Convert([]interface{}, interface{}) error
}

//批量返回hash
func (c *baseClient) HGETALLBatch(convert Converter, dest interface{}, key ...string) error {
	luaScript := `local rst={}; for i,v in pairs(KEYS) do rst[i]=redis.call('HGETALL', v) end;return rst`
	conn := c.redisPool.Get()
	defer conn.Close()
	args := redis.Args{}
	lua := redis.NewScript(len(key), luaScript)
	for _, v := range key {
		args = args.Add(v)
	}
	v, err := redis.Values(lua.Do(conn, args...))
	if err != nil {
		return err
	}
	return convert.Convert(v, dest)
}
