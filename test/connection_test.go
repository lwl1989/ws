package test

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"testing"
	"time"
    "fmt"
    "io/ioutil"
)

type connection struct {
	uuid string
	cs   *websocket.Conn
}
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func TestConnections(t *testing.T) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/test"}
	log.Printf("connecting to %s", u.String())
	u1 := url.URL{Scheme: "http", Host: *addr, Path: "/room/test"}
	cons := make(map[string]*connection)

	http.Get(u1.String())

	for i := 0; i < 10000; i++ {
		header := http.Header{}
		// header.Set("Sec-WebSocket-Key", GetRandomString(16))
		c, res, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
		    s,e := ioutil.ReadAll(res.Body)
		    fmt.Println(string(s[:]), e)
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
