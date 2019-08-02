package websocket

import (
    "net/http"
    "sync"
    "github.com/gorilla/websocket"
)
var Wsp *WsProtocol
var Up   *websocket.Upgrader


func init()  {
    Wsp = &WsProtocol{

    }
    Up = &websocket.Upgrader{

    }
    Up.ReadBufferSize = 1024
    Up.WriteBufferSize = 1024
    Up.CheckOrigin = func(r *http.Request) bool {
        //todo: check origin is allowed
        return true
    }
    //Wsp.Connections = make(map[string]*WsConn)
    Wsp.ConnectionsMap = new(sync.Map)
    //Wsp.rwm = new(sync.RWMutex)
    Wsp.Msg = make(chan []byte)
    Wsp.register = make(chan *WsConn)
    Wsp.unRegister = make(chan *WsConn)
    go Wsp.Run()
}