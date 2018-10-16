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
	isRunning := false
	room := types.Room{ID: bson.NewObjectId().Hex(), GameTime: GameTime, Game: game.New(), IsRunning: &isRunning}
	list[room.ID] = room
	return room.ID
}

func AddPlayerToRoom(roomID string, conn *websocket.Conn, color piece.Color) (string, string, error) {
	room := list[roomID]

	if room.Player1 != nil && room.Player2 != nil {
		return "", "", errors.New("full")
	}

	playerStopper := list[roomID].GameTime

	if color != piece.BothColors {
		player := types.Player{ID: bson.NewObjectId().Hex(), Conn: conn, Color: color, Stopper: &playerStopper}
		room.Player1 = &player
		list[roomID] = room
		return player.ID, "", nil
	}

	player := types.Player{ID: bson.NewObjectId().Hex(), Conn: conn, Stopper: &playerStopper}
	if room.Player1.Color == piece.Black {
		player.Color = piece.White
	} else {
		player.Color = piece.Black
	}
	room.Player2 = &player
	return player.ID, player.Color.String(), nil

}

func PauseGame(conn *websocket.Conn) {
	for _, room := range list {
		if conn == room.Player1.Conn || conn == room.Player2.Conn {
			if *room.IsRunning {
				*room.IsRunning = false
				Print()
			}
			return
		}
	}
}

func Print() {
	for roomID, room := range list {
		str := "\n{\n"
		str += " RommID:[" + roomID + "]\n"
		str += " IsRunning:[" + strconv.FormatBool(*room.IsRunning) + "]\n"
		str += " GameTime:[" + room.GameTime.String() + "]\n"
		str += " {\n"
		str += "  PlayerID:[" + room.Player1.ID + "]\n"
		str += "  Connection:[" + room.Player1.Conn.RemoteAddr().String() + "]\n"
		str += "  Color:[" + room.Player1.Color.String() + "]\n"
		str += "  Stopper:[" + room.Player1.Stopper.String() + "]\n"
		str += " }\n"
		if room.Player2 != nil {
			str += " {\n"
			str += "  PlayerID:[" + room.Player2.ID + "]\n"
			str += "  Connection:[" + room.Player2.Conn.RemoteAddr().String() + "]\n"
			str += "  Color:[" + room.Player2.Color.String() + "]\n"
			str += "  Stopper:[" + room.Player2.Stopper.String() + "]\n"
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
