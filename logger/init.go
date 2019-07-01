package logger

import (
    log2 "github.com/lwl1989/logger"
    "github.com/lwl1989/ws/config"
)

var Log *log2.TTLog

func init()  {
    Log = log2.GetFileLogger(config.Cf.LogConfig)
}
