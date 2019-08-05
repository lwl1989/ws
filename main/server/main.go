package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/lwl1989/ws/websocket"
    "runtime"
    _ "net/http/pprof"
)
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main()  {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()
    log.SetFlags(0)
    //http.HandleFunc("/ws", websocket.Handler)
    go websocket.GetMessage()
    go func() {
        log.Println(http.ListenAndServe("localhost:10000", nil))
    }()
    log.Fatal(http.ListenAndServe(*addr, websocket.Wsp))
}