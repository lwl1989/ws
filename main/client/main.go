package main

import (
    "flag"
    "log"
    "net/url"
    "os"
    "os/signal"

    "github.com/gorilla/websocket"
    "net/http"
    "github.com/google/uuid"
)
var addr = flag.String("addr", "localhost:8080", "http service address")

type connection struct {
    uuid string
    cs *websocket.Conn
}

func main()  {
    flag.Parse()
    log.SetFlags(0)

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
    log.Printf("connecting to %s", u.String())

    cons := make(map[string]*connection)


    for i := 0; i < 10000 ; i++ {
        header := http.Header{}
        header.Set("Sec-WebSocket-Key", uuid.New().String())
        c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
        if err != nil {
            log.Fatal("dial:", err)
        }
        go func() {
            //defer close(done)
            for {
                _, message, err := c.ReadMessage()
                if err != nil {
                    log.Println("read:", err)
                    return
                }
                log.Printf("recv: %s", message)
            }
        }()
        cons[header.Get("Sec-WebSocket-Key")] = &connection{
            header.Get("Sec-WebSocket-Key"),
            c,
        }
    }

    for _,v := range cons {
        v.cs.WriteMessage(websocket.TextMessage,[]byte("bye"))
        v.cs.Close()
    }
    //defer c.Close()

    //done := make(chan struct{})


}
