package test

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

type connection struct {
	uuid string
	cs   *websocket.Conn
}
var addr = flag.String("addr", "localhost:8080", "http service address")

func TestConnections(t *testing.T) {
	//var addr = "0.0.0.0:8080"

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: "room_id=test"}
	log.Printf("connecting to %s", u.String())
	u1 := url.URL{Scheme: "http", Host: *addr, Path: "/room",RawQuery: "room_id=test"}
	fmt.Println(u1.String())
	cons := make(map[string]*connection)

	_ , _ = http.Get(u1.String())

	for i := 0; i < 1; i++ {
		header := http.Header{}
		// header.Set("Sec-WebSocket-Key", GetRandomString(16))
		c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		//fmt.Println(c,res ,err)
		time.Sleep(100*time.Millisecond)
		if err == nil {
			go func() {
				//defer close(done)
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						return
					}
					log.Printf("recv: %s from conn id ", message)
					c.WriteMessage(websocket.TextMessage , []byte("dddsdad"))
				}
			}()
			cons[header.Get("Sec-WebSocket-Key")] = &connection{
				header.Get("Sec-WebSocket-Key"),
				c,
			}
		}else {

		}
	}
	time.Sleep(1000 * time.Second)
	//for _,v := range cons {
	//    v.cs.WriteMessage(websocket.TextMessage,[]byte("bye"))
	//    v.cs.Close()
	//}
	//defer c.Close()

	done := make(chan struct{})
	<-done

}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
