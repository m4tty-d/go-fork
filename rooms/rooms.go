package rooms

import (
	"time"

	"github.com/andrewbackes/chess/piece"

	"github.com/gorilla/websocket"

	"github.com/andrewbackes/chess/game"

	"gopkg.in/mgo.v2/bson"

	"gitlab.com/chess-fork/go-fork/types"
)

var rooms = make(map[string]types.Room)

func CreateRoom(GameTime time.Time) string {
	room := types.Room{}
	room.ID = bson.NewObjectId().Hex()
	room.GameTime = GameTime
	room.Game = game.New()
	rooms[room.ID] = room
	return room.ID
}

func AddPlayerToRoom(RoomID, PlayerID string, Conn *websocket.Conn, Color piece.Color) string {
	player := types.Player{}
	player.ID = bson.NewObjectId().Hex()
	player.Color = Color
	playerStopper := rooms[RoomID].GameTime
	player.Stopper = &playerStopper
	return player.ID
}

/*var players = make(map[string]types.Player)

func Add(conn *websocket.Conn, time *time.Time) {
	id := GenerateUniqueID()
	players[id] = types.Player{ID: id, Conn: conn, Time: time}
}

func RemoveById(id string) {
	delete(players, id)
}

func RemoveByConn(conn *websocket.Conn) {
	for _, player := range players {
		if player.Conn == conn {
			delete(players, player.ID)
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
	for _, element := range players {
		log.Println("\n{\n id:[" + element.ID + "]\n" + " conn:[" + element.Conn.RemoteAddr().String() + "]\n" + " time:[" + element.Time.String() + "]\n}")
	}
}

func VerifyTime() {
	for _, element := range players {
		(*element.Time) = (*element.Time).Add(-time.Second)
		if (*element.Time).Hour() == 23 {
			// time end
			element.Conn.WriteJSON(types.Server{Type: "clock", Payload: "end"})
		} else {
			h, m, s := element.Time.Clock()
			element.Conn.WriteJSON(types.Server{Type: "clock", Payload: types.Watch{H: h, M: m, S: s}})
		}
	}
}*/
