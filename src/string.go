package src

import "github.com/gomodule/redigo/redis"

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
func (c *BaseClient) SET(key string, value interface{}) error {
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
func (c *BaseClient) MSET(keyAndValue... interface{}) (string, error) {
	args := redis.Args{}
	//一项一项加上去 key, value, key, value ...
	for _, kv := range keyAndValue {
		args = args.Add(kv)
	}
	conn := c.redisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("mset", args...))
}
