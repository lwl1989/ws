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
    close bool
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
        wsc.Close()
    }()
    wsc.SetReadLimit(maxMessageSize)
    wsc.SetReadDeadline(time.Now().Add(pongWait))
    wsc.SetPongHandler(func(string) error { wsc.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        msgType, message, err := wsc.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                logger.Log.Println(fmt.Sprintf("error: %v", err))
            }
            break
        }

        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        logger.Log.Println(fmt.Sprintf("receive message: %s", string(message[:])))

        switch msgType {
            case websocket.CloseMessage:
                logger.Log.Println("client send close:"+wsc.GetUniqueKey())
                return
            case websocket.TextMessage:
                //todo:预留
            case websocket.BinaryMessage:
                //todo:预留
            //default:
                //todo:忽略,select 移除default减少性能损失
        }
    }
}

//close and offline
func (wsc *WsConn) Close() {
    //panic: close of closed channel
    if !wsc.close {
        wsc.close = true
        Wsp.OffLine(wsc)
        wsc.Conn.Close()
        close(wsc.send)
    }
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

            err := wsc.WriteBytes(message)
            if err != nil {
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

func (wsc *WsConn) WriteString(message string) {
    wsc.WriteBytes([]byte(message))
}

func (wsc *WsConn) WriteBytes(message []byte) error {
    wsc.SetWriteDeadline(time.Now().Add(writeWait))

    w, err := wsc.NextWriter(websocket.TextMessage)
    if err != nil {
        return err
    }

    _,err = w.Write(message)
    if err != nil {
        return err
    }

    if err := w.Close(); err != nil {
        return err
    }

    return nil
}