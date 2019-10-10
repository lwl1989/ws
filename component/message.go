package component

import (
    "time"
    "fmt"
    "github.com/go-libraries/ws"
)

//one second read one message
func GetMessage(wsp *ws.Protocol) {
    timer := time.NewTicker(3 * time.Second)
    defer func() {
        // 如果程序异常, 停止当前定时任务,记录日志,重启任务
        if x := recover(); x != nil {
            timer.Stop()
            wsp.PLog.Println("update cache panic:UpdateCacheTickers panic :", x)
            go GetMessage(wsp)
        }
    }()
    for {
        // 监听IO
        select {
        // 如果时间通道数据读取成功,
        case <-timer.C:
            fmt.Println("now connections num is:",wsp.GetNumber())
            go getMessage(wsp)
        }
    }
}

func getMessage(wsp *ws.Protocol) {
    msg := &ws.RoomMsg{}
    msg.SetMsg([]byte("hello"))
    msg.SetRoom("test")
    wsp.Send(msg)
}