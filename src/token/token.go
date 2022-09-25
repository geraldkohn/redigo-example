package token

import (
	"errors"
	"log"
	"time"

	"github.com/geraldkohn/redigo-example/pkg"
)

var (
	defaultSize = 1000000
)

type cancelFunc func()

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// 查询令牌
func getToken(userId string) (string, error) {
	client, err := pkg.GetRedisPool()
	checkError(err)

	s, err := client.HGET("toke-user", userId)
	checkError(err)

	return s, nil
}

// 更新令牌
func updateToken(userId, token string) error {
	client, err := pkg.GetRedisPool()
	checkError(err)

	ok, _ := client.HSET("token-user", userId, token)
	if ok == 0 {
		checkError(errors.New("设置失败"))
	}

	okk, _ := client.ZREM("userToken-time", userId)
	if !okk {
		checkError(errors.New("删除失败"))
	}
	okk, _ = client.ZADD("userToken-time", userId, time.Now().String())
	if !okk {
		checkError(errors.New("设置失败"))
	}

	return nil
}

// 清理令牌
// 调用cancel()就停止清理
func cleanToken() (cancelFunc, error) {
	ch := make(chan struct{}, 1)
	cancel := func() {
		ch <- struct{}{}
	}
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				client, err := pkg.GetRedisPool()
				checkError(err)
				size, _ := client.ZCARD("userToken-time")
				tokens, _ := client.ZRANGE("userToken-time", int64(defaultSize), -1, false)
				if size > defaultSize {
					ok, err := client.ZREM("userToken-time", tokens...)
					checkError(err)
					if !ok {
						checkError(errors.New("清理令牌失败"))
					}
				}
				time.Sleep(100 * time.Second)
			}
		}
	}()
	return cancel, nil
}
