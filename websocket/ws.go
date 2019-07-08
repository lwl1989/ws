package websocket

import (
    "github.com/gorilla/websocket"
    "net/http"
    "github.com/lwl1989/ws/logger"
    "sync"
    "github.com/go-kit/kit/util/conn"
)

type WsConn struct {
    *websocket.Conn
    UniqueKey string
}

func (wsc *WsConn) GetUniqueKey() string {
    return wsc.UniqueKey
}

type WsProtocol struct {
    websocket.Upgrader

    //all connections, It's mapping O(1)
    Connections map[string]*WsConn

    //use rw mutex
    rwm *sync.RWMutex
}

func Handler(w http.ResponseWriter, r *http.Request)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    con, err := ws.Upgrade(w, r, nil)

    var wsConn  = &WsConn {
        UniqueKey:uniqueKey,
        Conn:con,
    }

    ws.Online(wsConn)
    if err != nil {
        logger.Log.Println("handler err with message" + err.Error())
        panic("handler err with message" + err.Error())
    }


    for {
        messageType, p, err := wsConn.ReadMessage()
        if err != nil {
            logger.Log.Println("read message error" +err.Error())
            return
        }

        if messageType == websocket.TextMessage {
            content := string(p)
            if content == "bye" {
                ws.OffLine(wsConn)
            }
        }
        if err := wsConn.WriteMessage(messageType, p); err != nil {
            logger.Log.Println("write message error" +err.Error())
            return
        }
    }
}

//conn connection,
func (w *WsProtocol) Online(conn *WsConn) {
    w.rwm.Lock()
    w.Connections[conn.GetUniqueKey()] = conn
    w.rwm.Unlock()
}

func (w *WsProtocol) All() map[string]*WsConn {
    w.rwm.RLock()
    defer w.rwm.RUnlock()
    return w.Connections
}

func (w *WsProtocol) OffLine(conn *WsConn) {
    w.rwm.Lock()
    defer w.rwm.Unlock()

    delete(w.Connections, conn.GetUniqueKey())
    conn.Close()
}