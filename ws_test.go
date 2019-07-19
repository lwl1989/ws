package ws

import (
    "testing"
    "fmt"
)


type tMessage struct {
    Msg chan []byte
}
func TestGetMessage(t *testing.T) {
    fmt.Println("sss")
    w := &tMessage{
        Msg:make(chan []byte),
    }
    go func() {
        w.Msg <- []byte("hello")
    }()

    for {
        select {
        case rMsg := <-w.Msg:
            fmt.Println(rMsg)
            TTT(rMsg)
            break
        }
    }

    //var f = make(chan bool)
    // <- f
}


func TTT(m []byte) {
   for i:=0;i<9999;i++ {
       fmt.Println(string(m[:]))
   }
}