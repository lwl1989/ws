package message

import (
    "github.com/go-redis/redis"
    "strconv"
    "github.com/lwl1989/ws/config"
    "fmt"
    "net"
)


var RMessage *RedisMessage

//数据落地
func init() {
    RMessage = &RedisMessage{}
    db,err := strconv.Atoi(config.Cf.Redis.Db)
    if err != nil {
        db = 0
    }
    fmt.Println(config.Cf.Redis)
    //&redis.Options{
    //    Addr: ":1234",
    //    Dialer: func() (net.Conn, error) {
    //        return net.Dial("tcp", redisAddr)
    //    },
    //}
    RMessage.Rs = redis.NewClient(&redis.Options{
        Addr:    config.Cf.Redis.Host,
        Dialer: func() (net.Conn, error) {
            return net.Dial("tcp", config.Cf.Redis.Host)
        },
        Password: config.Cf.Redis.Pw, // no password set
        DB:      db,  // use default DB
    })
}

