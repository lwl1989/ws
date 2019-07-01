package websocket

import (
    "github.com/gorilla/websocket"
    "net/http"
    "github.com/lwl1989/ws/logger"
)

type WsProtocol struct {
    websocket.Upgrader
}

func Handler(w http.ResponseWriter, r *http.Request)  {
    conn, err := ws.Upgrade(w, r, nil)

    if err != nil {
        logger.Log.Println("handler err with message" + err.Error())
    }

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            logger.Log.Println(err)
            return
        }
        if err := conn.WriteMessage(messageType, p); err != nil {
            logger.Log.Println(err)
            return
        }
    }

}