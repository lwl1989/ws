package ws

type IMessage interface {
    GetMessage() (bs []byte,length int64, err error)
    GetRoom() string
}

type ILog interface {
    Println(v... interface{})
}