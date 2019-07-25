package websocket

import (
    "time"
    "github.com/gorilla/websocket"
    "bytes"
    "github.com/lwl1989/ws/logger"
    "fmt"
)


type WsConn struct {
    *websocket.Conn
    UniqueKey string
    send chan []byte
}

func (wsc *WsConn) GetUniqueKey() string {
    return wsc.UniqueKey
}

func (wsc *WsConn) Send(b []byte) {
    wsc.send <- b
}


//read
func (wsc *WsConn) read() {
    defer func() {
        Wsp.OffLine(wsc)
    }()
    wsc.SetReadLimit(maxMessageSize)
    wsc.SetReadDeadline(time.Now().Add(pongWait))
    wsc.SetPongHandler(func(string) error { wsc.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        _, message, err := wsc.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                logger.Log.Println(fmt.Sprintf("error: %v", err))
            }
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        logger.Log.Println(fmt.Sprintf("error: %v", err))
        //Wsp.getMessageClient(message)
    }
}

//close and offline
func (wsc *WsConn) Close() {
    Wsp.OffLine(wsc)
    wsc.Conn.Close()
}

func (wsc *WsConn) write() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        wsc.Close()
    }()
    for {
        select {
        case message, ok := <- wsc.send:
            wsc.SetWriteDeadline(time.Now().Add(writeWait))
            //fmt.Println(message)
            if !ok {
                // The hub closed the channel.
                wsc.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := wsc.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued chat messages to the current websocket message.
            //n := len(wsc.send)
            //for i := 0; i < n; i++ {
            //    w.Write(newline)
            //    w.Write(<-wsc.send)
            //}

            if err := w.Close(); err != nil {
                return
            }
            //心跳
        case <-ticker.C:
            wsc.SetWriteDeadline(time.Now().Add(writeWait))
            if err := wsc.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}