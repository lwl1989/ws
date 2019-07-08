package websocket

import (
    "net/http"
)
var ws *WsProtocol

func init()  {
    ws = &WsProtocol{

    }
    ws.ReadBufferSize = 1024
    ws.WriteBufferSize = 1024
    ws.CheckOrigin = func(r *http.Request) bool {
        //todo: check origin is allowed
        return true
    }
}