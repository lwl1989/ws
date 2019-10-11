package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/go-libraries/ws"
    "runtime"
    _ "net/http/pprof"
    c "ws/component"
    "github.com/go-libraries/timer"
    "time"
)
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var etcdServer = flag.String("etcdServer", "0.0.0.0:2380", "etcd service address")

func main()  {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()
    log.SetFlags(0)
    ws.Wsp.PLog = c.Logs
    go c.GetMessage(ws.Wsp)

    etcd
    sc := timer.GetTaskScheduler()
    sc.AddFuncSpace(int64(20*time.Second), time.Now().Unix() + int64(time.Hour*10000), func() {

    })
    log.Fatal(http.ListenAndServe(*addr, ws.Wsp))
}