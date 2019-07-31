package websocket

import (
    "net/http"
    "github.com/lwl1989/ws/logger"
    "sync"
    "github.com/lwl1989/ws/message"
    "time"
    "fmt"
)

type WsProtocol struct {
    // Register requests from the clients.
    register chan *WsConn

    // UnRegister requests from clients.
    unRegister chan *WsConn

    //all connections, It's mapping O(1)
    Connections map[string]*WsConn
    //todo: next splice connections

    //use rw mutex
    rwm *sync.RWMutex

    Msg chan []byte

    num int //count
}

func Handler(w http.ResponseWriter, r *http.Request)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    con, err := Up.Upgrade(w, r, nil)
    if err != nil {
        logger.Log.Println("handler err with message" + err.Error())
        panic("handler err with message" + err.Error())
    }

    var wsConn  = &WsConn {
        UniqueKey:uniqueKey,
        Conn:con,
        send: make(chan []byte, 256),
    }

    Wsp.Online(wsConn)

    go wsConn.read()
    go wsConn.write()
}

//one second read one message
func GetMessage() {
    timer := time.NewTicker(3 * time.Second)
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
            fmt.Println("now connections num is:",Wsp.num)
            if Wsp.num < 3000 {
                for k,v := range Wsp.All() {
                    fmt.Println(k,v)
                }
            }
            go Wsp.getMessage()
        }
    }
}

func (w *WsProtocol) send(b []byte) {
    all := w.All()
    for _,v := range all{
        v.Send(b)
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

func (w *WsProtocol) getMessageClient(msg []byte) {
    w.Msg <- msg
}
//conn connection,write lock
func (w *WsProtocol) Online(conn *WsConn) {
    w.register <- conn
}

//read lock
func (w *WsProtocol) All() map[string]*WsConn {
    w.rwm.RLock()
    defer w.rwm.RUnlock()
    return w.Connections
}

//write lock
func (w *WsProtocol) OffLine(conn *WsConn) {
    w.unRegister <- conn
}

func (w *WsProtocol) Run()  {
    for {
        select {
        case client := <-w.register:
            w.num = w.num +1
            w.Connections[client.GetUniqueKey()] = client
        case client := <-w.unRegister:
            w.num = w.num - 1
            delete(w.Connections, client.GetUniqueKey())
        case msg := <-w.Msg:
            w.send(msg)
        }
    }
}