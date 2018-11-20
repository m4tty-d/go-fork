package spectate

import (
	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/types"
)

var list = make(map[string][]*websocket.Conn)

func AddRoom(ID string) {
	list[ID] = []*websocket.Conn{}
}

func AddConnToRoom(ID string, conn *websocket.Conn) {
	list[ID] = append(list[ID], conn)
}

func NotifyConns(ID, message string) {
	for _, conn := range list[ID] {
		conn.WriteJSON(types.Server{Type: "spectate", Payload: message})
	}
}
