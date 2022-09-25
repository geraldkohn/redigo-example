# 最受欢迎的文章

## 1 对文章进行投票

* 使用HASH存储文章

key: 文章ID, field: title, link, poster(userID), time, votes

* 使用ZSET针对不同的属性排序

分别按照评分排序

key: scores, item: 文章ID, score: 评分

* 防止重复投票, 使用SET类型记录每篇文章ID对应的投票过的userID的集合

key: 文章ID, item: userID

* 投票时: 先查看是否投过, 更新评分排序集合和投票集合.

## 2 发布文章

* HSET/HMSET写到HASH结构中
