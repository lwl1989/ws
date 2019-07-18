package ws

import (
    "flag"
    "log"
    "net/http"
    "github.com/lwl1989/ws/websocket"
)
var addr = flag.String("addr", "localhost:8080", "http service address")
func t()  {
    flag.Parse()
    log.SetFlags(0)
    http.HandleFunc("/", websocket.Handler)
    go websocket.GetMessage()
    log.Fatal(http.ListenAndServe(*addr, nil))
}