package component

import (
    "github.com/go-redis/redis"
    "strconv"
    "log"
    "os"
)


var RMessage *RedisMessage
var Cf *Config
var Logs  *log.Logger
var EtcReg *ServiceReg

func init()   {
    Cf = &Config{
        LogConfig: &LogFileConfig{
            FilePath: "/tmp/ws_log",
        },
        Redis: &Redis{
            Host: "127.0.0.1:6379",
            Db:   "1",
            Pw:   "",
        },
        Etcd: &EtcdConfig{
            Addr: []string{"127.0.0.1:2379"},
            TimeOut: 5,
            LeaseTime: 3600,
        },
    }

    Logs = log.New(os.Stdout, "", 1)

    RMessage = &RedisMessage{}
    db,err := strconv.Atoi(Cf.Redis.Db)
    if err != nil {
        db = 0
    }

    RMessage.Rs = redis.NewClient(&redis.Options{
        Addr:    Cf.Redis.Host,
        //Dialer: func() (net.Conn, error) {
        //    return net.Dial("tcp", Cf.Redis.Host)
        //},
        Password: Cf.Redis.Pw, // no password set
        DB:      db,  // use default DB
    })

    if Cf.Etcd.Enabled {
        EtcReg,_ = NewServiceReg(Cf.Etcd.Addr, Cf.Etcd.LeaseTime, Cf.Etcd.TimeOut)
    } else{
        EtcReg = &ServiceReg{}
    }
}

