# 管理令牌

## 查询令牌

将token和用户的映射关系存储在HASH结构

key: token-user, field: userID, value: userToken

更新ZSET结构, 重置时间戳

## 更新令牌

将用户令牌和最后一次更新时间映射到ZSET结构中, 并且更新HASH结构

key: userToken-time, item: userID, timestamp: 最后一次查询/更新时间

## 清理令牌

假设设置了一个最多缓存的令牌数量, 从zset中选择超出限制数量的令牌的userId

去HASH里删除userI对应记录, 并且在ZSET中删除记录.