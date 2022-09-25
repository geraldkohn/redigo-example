package src

import (
	"errors"
	"log"
	"time"

	"github.com/geraldkohn/redigo-example/pkg"
	"github.com/google/uuid"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// 对文章进行投票
func vote(userId string, articleId string, score int) error {
	client, err := pkg.GetRedisPool()
	checkError(err)

	exist, err := client.SISMEMBER(articleId, userId)
	checkError(err)
	if !exist {
		client.SADD(articleId, userId)                    // 记录已经投票过了
		client.HINCRBY(articleId, "votes", 1)             // 增加文章的投票人数
		client.ZINCRBY("scores", articleId, int64(score)) // 更改文章的评分
		return nil
	}
	return errors.New("已经投票了")
}

// 发布文章
func postArticle(title, link, poster string) error {
	client, err := pkg.GetRedisPool()
	checkError(err)

	uid := uuid.New().String()
	_, err = client.HMSET(uid, []interface{}{"title", title, "link", link, "poster", poster, "time", time.Now().GoString(), "votes", 0})
	checkError(err)

	return nil
}
