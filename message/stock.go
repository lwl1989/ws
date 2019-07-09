package message

import (
    "github.com/go-redis/redis"
    "fmt"
    "time"
    "github.com/lwl1989/ws/logger"
    "encoding/json"
)

const RedisStockKey = "stock_key_%d"
type Stock interface {
    GetMessage() (bs []byte,len int64, err error)
}

type StockMessage struct{
    Time int64 `json:"time"`
    Contents []string `json:"contents"`
}

type RedisMessage struct {
    Rs *redis.Client
}

//impl Stock interface with redis
func (rMsg *RedisMessage) GetMessage() (bs []byte,len int64, err error)  {
    sm := rMsg.getMessage()
    cmd := rMsg.Rs.LRange(sm.getRedisKey(), 0, -1)

    if err := cmd.Err();  err != nil {
        logger.Log.Println("read message error"+err.Error())
        return bs,0,err
    }

    sm.Contents = make([]string, 0)
    for _,v := range cmd.Val() {
        sm.Contents = append(sm.Contents, v)
        len ++
    }

    bs,err = json.Marshal(sm)
    if err != nil {
        return bs,0,err
    }

    return bs,len,err
}

func (rMsg *RedisMessage) getMessage() StockMessage {
    return StockMessage{
        Time:time.Now().Unix(),
    }
}

func (sm StockMessage) getRedisKey() string {
    return fmt.Sprintf(RedisStockKey, time.Now().Unix())
}