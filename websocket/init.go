package websocket

import (
    "net/http"
    "sync"
)
var Wsp *WsProtocol

func init()  {
    Wsp = &WsProtocol{

    }
    Wsp.ReadBufferSize = 1024
    Wsp.WriteBufferSize = 1024
    Wsp.CheckOrigin = func(r *http.Request) bool {
        //todo: check origin is allowed
        return true
    }
    Wsp.rwm = &sync.RWMutex{}
    Wsp.Msg = make(chan []byte)
}