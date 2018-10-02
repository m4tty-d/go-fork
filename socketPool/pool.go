package socketpool

import (
	"log"

	"gitlab.com/chess-fork/go-fork/types"

	"github.com/dchest/uniuri"
	"github.com/gorilla/websocket"
)

type client struct {
	id   string
	conn *websocket.Conn
}

var clients = make(map[string]client)

func Add(conn *websocket.Conn) {
	id := GenerateUniqueID()
	clients[id] = client{id: id, conn: conn}
}

func Remove(id string) {
	delete(clients, id)
}

func SendToAll(msg string) {
	for _, client := range clients {
		client.conn.WriteJSON(types.Response{Type: "message", Payload: msg})
	}
}

func Exists(id string) bool {
	if _, exists := clients[id]; exists {
		return true
	}

	return false
}

func GenerateUniqueID() string {
	id := uniuri.New()
	if Exists(id) {
		return GenerateUniqueID()
	}

	return id
}

func Print() {
	log.Println("{")
	for _, element := range clients {
		log.Println(" id:[" + element.id + "]" + " conn:[" + element.conn.RemoteAddr().String() + "]")
	}
	log.Println("}")
}
