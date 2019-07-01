package websocket

var ws WsProtocol

func init()  {
    ws = WsProtocol{

    }
    ws.ReadBufferSize = 1024
    ws.WriteBufferSize = 1024

}