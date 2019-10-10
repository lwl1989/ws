package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/go-libraries/ws"
    "runtime"
    _ "net/http/pprof"
    c "ws/component"
)
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main()  {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()
    log.SetFlags(0)
    ws.Wsp.PLog = c.Logs
    go c.GetMessage(ws.Wsp)
    log.Fatal(http.ListenAndServe(*addr, ws.Wsp))
}