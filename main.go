package main

import (
	"flag"
	"fmt"
	"github.com/go-libraries/ws"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	c "ws/component"
)

//var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
//var etcdServer = flag.String("etcdServer", "0.0.0.0:2380", "etcd service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	ws.Wsp.PLog = c.Logs
	//ws.Wsp.Config = c.Cf
	go c.GetMessage(ws.Wsp)
	ws.Wsp.Config = c.Cfg{}
	//etcd
	//sc := timer.GetTaskScheduler()
	//sc.AddFuncSpace(int64(20*time.Second), time.Now().Unix()+int64(time.Hour*10000), func() {
	//
	//})
	//s2 := s{}
	log.Fatal(http.ListenAndServe(":8080", ws.Wsp))
}

type s struct {

}
func (s1 s) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.Path, "/")
	l := len(res)
	fmt.Println(r.URL.String(), res[0], l)
}
