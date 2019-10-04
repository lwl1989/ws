package ws

import (
    "net/http"
    "sync"
    "github.com/gorilla/websocket"
)
var Wsp *Protocol
var Up   *websocket.Upgrader


func init()  {
    Wsp = &Protocol{

    }
    Up = &websocket.Upgrader{

    }
    Up.ReadBufferSize = 1024
    Up.WriteBufferSize = 1024
    Up.CheckOrigin = func(r *http.Request) bool {
        //todo: check origin is allowed
        return true
    }
    //Wsp.Connections = make(map[string]*Connection)
    Wsp.ConnectionsMap = make(map[string]*sync.Map)
        //new(sync.Map)
    //Wsp.rwm = new(sync.RWMutex)
    Wsp.Msg = make(chan IMessage)
    Wsp.register = make(chan *Connection)
    Wsp.unRegister = make(chan *Connection)
    go Wsp.Run()
}