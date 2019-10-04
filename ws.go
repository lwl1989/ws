package ws

import (
    "net/http"
    "sync"
    //"github.com/lwl1989/ws/message"
    "strings"
)


type Protocol struct {
    // Register requests from the clients.
    register chan *Connection

    // UnRegister requests from clients.
    unRegister chan *Connection

    //all connections, It's mapping O(1)
    //Connections map[string]*Connection
    ConnectionsMap  map[string]*sync.Map

    //use rw mutex
    rwm *sync.RWMutex

    Msg chan IMessage
    num int //count
    PLog ILog
}

var m sync.Map

func (w *Protocol) ServeHTTP(rw http.ResponseWriter, r *http.Request)  {

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
        return
    }

    if res[0] == "room" {
        w.registerRoom(rw, r, room)
        return
    }

    rw.WriteHeader(200)
    rw.Write([]byte("{}"))
    return
}

//lock room
func (w *Protocol)  registerRoom(rw http.ResponseWriter, r *http.Request, room string) {
    w.rwm.Lock()
    defer func() {
        w.rwm.Unlock()
    }()

    if _,ok := w.ConnectionsMap[room]; !ok {
        w.ConnectionsMap[room] = new(sync.Map)
    }
}

func (w *Protocol)  registerWs(rw http.ResponseWriter, r *http.Request, room string)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    con, err := Up.Upgrade(rw, r, nil)
    if err != nil {
        w.PLog.Println("handler err with message" + err.Error())
        panic("handler err with message" + err.Error())
    }

    var wsConn  = &Connection {
        UniqueKey:uniqueKey,
        Conn:con,
        send: make(chan []byte, 256),
        room:room,
    }

    Wsp.Online(wsConn)

    go wsConn.read()
    go wsConn.write()
}

func (w *Protocol) send(msg IMessage) {
    all := w.All(msg.GetRoom())
    bs,length,err := msg.GetMessage()

    if err != nil {
        w.PLog.Println(err)
        return
    }
    if length < 1 || len(bs) == 0 {
        w.PLog.Println("message is nil")
        return
    }
    for _,v := range all{
        v.Send(bs)
    }
}

//conn connection,write lock
func (w *Protocol) Online(conn *Connection) {
    w.register <- conn
}

//read lock
func (w *Protocol) All(room string) map[string]*Connection {
    seen := make(map[string]*Connection, w.num)

    m.Range(func(ki, vi interface{}) bool {
        k, v := ki.(string), vi.(*Connection)
        seen[k] = v
        return true
    })

    return seen
}

//write lock
func (w *Protocol) OffLine(conn *Connection) {
    w.unRegister <- conn
}

func (w *Protocol) Run()  {
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