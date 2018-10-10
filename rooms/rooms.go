package rooms

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/andrewbackes/chess/piece"

	"github.com/gorilla/websocket"

	"github.com/andrewbackes/chess/game"

	"gopkg.in/mgo.v2/bson"

	"gitlab.com/chess-fork/go-fork/types"
)

var list = make(map[string]types.Room)

func CreateRoom(GameTime time.Time) string {
	room := types.Room{ID: bson.NewObjectId().Hex(), GameTime: GameTime, Game: game.New(), IsRunning: false, Players: make(map[string]*types.Player)}
	list[room.ID] = room
	return room.ID
}

func AddPlayerToRoom(roomID string, conn *websocket.Conn, color piece.Color) (string, error) {
	if len(list[roomID].Players) == 2 {
		return "", errors.New("full")
	}
	playerStopper := list[roomID].GameTime
	player := types.Player{ID: bson.NewObjectId().Hex(), Conn: conn, Color: color, Stopper: &playerStopper}
	list[roomID].Players[player.ID] = &player
	return player.ID, nil
}

func Print() {
	for roomID, room := range list {
		str := "\n{\n"
		str += " RommID:[" + roomID + "]\n"
		str += " IsRunning:[" + strconv.FormatBool(room.IsRunning) + "]\n"
		str += " GameTime:[" + room.GameTime.String() + "]\n"
		for playerID, player := range room.Players {
			str += " {\n"
			str += "  PlayerID:[" + playerID + "]\n"
			str += "  Connection:[" + player.Conn.RemoteAddr().String() + "]\n"
			str += "  Color:[" + player.Color.String() + "]\n"
			str += "  Stopper:[" + player.Stopper.String() + "]\n"
			str += " }\n"
		}
		str += "}"
		log.Println(str)
	}
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
