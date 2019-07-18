package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/lwl1989/ws/websocket"
)
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main()  {
    flag.Parse()
    log.SetFlags(0)
    http.HandleFunc("/", websocket.Handler)
    go websocket.GetMessage()
    log.Fatal(http.ListenAndServe(*addr, nil))
}