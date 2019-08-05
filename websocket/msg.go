package websocket


type RoomMsg struct {
    msg  []byte
    room string
}

func (rm *RoomMsg) SetMsg(m []byte) {
    rm.msg = m
}

func (rm *RoomMsg) SetRoom(r string) {
    rm.room = r
}

func (rm *RoomMsg) GetMsg() []byte {
    return rm.msg
}

func (rm *RoomMsg) GetRoom() string {
    return rm.room
}

