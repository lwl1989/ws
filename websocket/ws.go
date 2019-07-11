package websocket

import (
    "github.com/gorilla/websocket"
    "net/http"
    "github.com/lwl1989/ws/logger"
    "sync"
    "github.com/lwl1989/ws/message"
    "time"
    "fmt"
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

    Msg chan []byte
}


//one second read one message
func GetMessage() {
    timer := time.NewTicker(1 * time.Second)
    defer func() {
        // 如果程序异常, 停止当前定时任务,记录日志,重启任务
        if x := recover(); x != nil {
            timer.Stop()
            logger.Log.Println("update cache panic:UpdateCacheTickers panic :", x)
            go GetMessage()
        }
    }()
    for {
        // 监听IO
        select {
        // 如果时间通道数据读取成功,
        case <-timer.C:
            fmt.Println("one")
            go Wsp.Send()
        }
    }
}
func (w *WsProtocol) Send() {
    go Wsp.getMessage()
    //cs := w.All()
    for {
        select {
            case rMsg := <-w.Msg:
                go w.send(rMsg)
            default:
                //fmt.Println(3)
        }
    }
}

func (w *WsProtocol) send(b []byte) {
    all := w.All()
    for _,v := range all{
        v.send(b)
    }
}

func (wsc *WsConn) send(b []byte) {
    err := wsc.WriteMessage(websocket.BinaryMessage, b)
    if err != nil{
        logger.Log.Println("send message error message:"+ err.Error())
    }
}

func (w *WsProtocol) getMessage() {
    w.Msg <- []byte("hello")
    return
    bs,l,err := message.RMessage.GetMessage()
    if err != nil {
        logger.Log.Println("get message error:", err)
        return
    }

    if l == 0 {
        logger.Log.Println(time.Now().Unix(), "get content is null")
        return
    }

    w.Msg <- bs
}

func Handler(w http.ResponseWriter, r *http.Request)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    con, err := Wsp.Upgrade(w, r, nil)

    var wsConn  = &WsConn {
        UniqueKey:uniqueKey,
        Conn:con,
    }

    Wsp.Online(wsConn)
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
                logger.Log.Println("conn offline" +wsConn.GetUniqueKey())
                Wsp.OffLine(wsConn)
                return
            }
        }
    }
}

//conn connection,write lock
func (w *WsProtocol) Online(conn *WsConn) {
    w.rwm.Lock()
    w.Connections[conn.GetUniqueKey()] = conn
    w.rwm.Unlock()
}

//read lock
func (w *WsProtocol) All() map[string]*WsConn {
    w.rwm.RLock()
    defer w.rwm.RUnlock()
    return w.Connections
}

//write lock
func (w *WsProtocol) OffLine(conn *WsConn) {
    w.rwm.Lock()
    defer w.rwm.Unlock()

    delete(w.Connections, conn.GetUniqueKey())
    conn.Close()
}