package socketPool

import (
	"log"

	"../types"

	"github.com/gorilla/websocket"
)

type browser struct {
	userid string
	conn   *websocket.Conn
}

var list []browser

func Add(userid string, conn *websocket.Conn) {
	brw := browser{userid: userid, conn: conn}
	list = append(list, brw)
}

func find(userid string) int {
	idx := -1
	for i := 0; i < len(list); i++ {
		if userid == list[i].userid {
			idx = i
			break
		}
	}
	return idx
}

func Status(userid string) bool {
	if find(userid) == -1 {
		return false
	} else {
		return true
	}
}

func Remove(userid string) {
	idx := find(userid)
	list[idx] = list[len(list)-1] // Copy last element to index i.
	list[len(list)-1] = browser{} // Erase last element (write zero value).
	list = list[:len(list)-1]     // Truncate slice.
}

func Print() {
	log.Println("{")
	for _, element := range list {
		log.Println(" userid:[" + element.userid + "]" + " conn:[" + element.conn.RemoteAddr().String() + "]")
	}
	log.Println("}")
}

func SendToAll(msg string) {
	for i := 0; i < len(list); i++ {
		list[i].conn.WriteJSON(types.Response{Type: "message", Payload: msg})
	}
}
