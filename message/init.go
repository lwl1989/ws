package message

import (
    "github.com/go-redis/redis"
    "strconv"
    "github.com/lwl1989/ws/config"
)



var rs *RedisMessage

//数据落地
func init() {
    db,err := strconv.Atoi(config.Cf.Redis.Db)
    if err != nil {
        db = 0
    }
    rs.Rs = redis.NewClient(&redis.Options{
        Addr:    config.Cf.Redis.Host,
        Password: config.Cf.Redis.Pw, // no password set
        DB:      db,  // use default DB
    })
}

