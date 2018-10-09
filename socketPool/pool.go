package socketpool

import (
	"log"
	"time"

	"gitlab.com/chess-fork/go-fork/types"

	"github.com/dchest/uniuri"
	"github.com/gorilla/websocket"
)

var players = make(map[string]types.Player)

func Add(conn *websocket.Conn, time time.Time) {
	id := GenerateUniqueID()
	players[id] = types.Player{Id: id, Conn: conn, Time: time}
}

func RemoveById(id string) {
	delete(players, id)
}

func RemoveByConn(conn *websocket.Conn) {
	for _, player := range players {
		if player.Conn == conn {
			delete(players, player.Id)
			break
		}
	}
}

func SendToAll(msg string) {
	for _, player := range players {
		player.Conn.WriteJSON(types.Server{Type: "message", Payload: msg})
	}
}

func Exists(id string) bool {
	if _, exists := players[id]; exists {
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
	for _, element := range players {
		log.Println(" id:[" + element.Id + "]\n" + " conn:[" + element.Conn.RemoteAddr().String() + "]\n" + " time:[" + element.Time.String() + "]")
	}
	log.Println("}")
}

func VerifyTime() {
	for _, element := range players {
		element.Time = element.Time.Add(-time.Second)
		if element.Time.Hour() == 23 {
			// time end
			element.Conn.WriteJSON(types.Server{Type: "clock", Payload: "end"})
			log.Println("time ended")
		} else {
			h, m, s := element.Time.Clock()
			element.Conn.WriteJSON(types.Server{Type: "clock", Payload: types.Watch{H: h, M: m, S: s}})
		}
	}
}
