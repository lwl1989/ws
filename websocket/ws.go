package websocket

import (
    "net/http"
    "github.com/lwl1989/ws/logger"
    "sync"
    //"github.com/lwl1989/ws/message"
    "time"
    "fmt"
    "strings"
)


type WsProtocol struct {
    // Register requests from the clients.
    register chan *WsConn

    // UnRegister requests from clients.
    unRegister chan *WsConn

    //all connections, It's mapping O(1)
    //Connections map[string]*WsConn
    //todo: map change to sync.Map
    //todo: next splice connections
    ConnectionsMap  map[string]*sync.Map

    //use rw mutex
    rwm *sync.RWMutex

    Msg chan RoomMsg
    num int //count
}

var m sync.Map

func (w *WsProtocol) ServeHTTP(rw http.ResponseWriter, r *http.Request)  {

    res := strings.Split(r.URL.Path, "/")
    l := len(res)
    if l < 2 || res[0] != "ws" {
        rw.WriteHeader(200)
        rw.Write([]byte("{}"))
        return
    }

    room := ""
    if len(res) == 2 {
        room = res[1]
    }

    if room == "" {
        rw.WriteHeader(200)
        rw.Write([]byte("{}"))
        return
    }

    if res[0] == "ws" {
        w.registerWs(rw, r, room)
    }

    if res[0] == "room" {
        w.registerRoom(rw, r, room)
    }

    rw.WriteHeader(200)
    rw.Write([]byte("{}"))
    return
}

//lock room
func (w *WsProtocol)  registerRoom(rw http.ResponseWriter, r *http.Request, room string) {
    w.rwm.Lock()
    defer func() {
        w.rwm.Unlock()
    }()

    if _,ok := w.ConnectionsMap[room]; !ok {
        w.ConnectionsMap[room] = new(sync.Map)
    }
}

func (w *WsProtocol)  registerWs(rw http.ResponseWriter, r *http.Request, room string)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    con, err := Up.Upgrade(rw, r, nil)
    if err != nil {
        logger.Log.Println("handler err with message" + err.Error())
        panic("handler err with message" + err.Error())
    }

    var wsConn  = &WsConn {
        UniqueKey:uniqueKey,
        Conn:con,
        send: make(chan []byte, 256),
        room:room,
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
            go Wsp.getMessage()
        }
    }
}

func (w *WsProtocol) send(msg RoomMsg) {
    all := w.All(msg.GetRoom())
    for _,v := range all{
        v.Send(msg.GetMsg())
    }
}

func (w *WsProtocol) getMessage() {
    w.Msg <- RoomMsg{
        msg:[]byte("hello"),
        room:"test",
    }
    return
    //bs,l,err := message.RMessage.GetMessage()
    //if err != nil {
    //    logger.Log.Println("get message error:", err)
    //    return
    //}
    //
    //if l == 0 {
    //    logger.Log.Println(time.Now().Unix(), "get content is null")
    //    return
    //}
    //
    //w.Msg <- bs
}

//func (w *WsProtocol) getMessageClient(msg []byte) {
//    w.Msg <- msg
//}
//conn connection,write lock
func (w *WsProtocol) Online(conn *WsConn) {
    w.register <- conn
}

//read lock
func (w *WsProtocol) All(room string) map[string]*WsConn {
    seen := make(map[string]*WsConn, w.num)

    m.Range(func(ki, vi interface{}) bool {
        k, v := ki.(string), vi.(*WsConn)
        seen[k] = v
        return true
    })

    return seen
}

//write lock
func (w *WsProtocol) OffLine(conn *WsConn) {
    w.unRegister <- conn
}

func (w *WsProtocol) Run()  {
    for {
        select {
        case client := <-w.register:
            room := client.GetRoom()
            w.num = w.num + 1
            w.ConnectionsMap[room].Store(client.GetUniqueKey(), client)
            //w.Connections[client.GetUniqueKey()] = client
        case client := <-w.unRegister:
            w.num = w.num - 1
            room := client.GetRoom()
            w.ConnectionsMap[room].Delete(client.GetUniqueKey())
            //delete(w.Connections, client.GetUniqueKey())
        case msg := <-w.Msg:
            w.send(msg)
        }
    }
}