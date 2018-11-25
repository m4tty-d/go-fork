package rooms

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/types"
	"gopkg.in/mgo.v2/bson"
)

var list = make(map[string]types.Room)

func getTimeFromBaseAndAdditional(base int) time.Time {
	m := base % 60
	h := base / 60

	return time.Date(0, 0, 0, h, m, 0, 0, time.UTC)
}

func CreateRoom(BaseTime int, AdditionalTime int) string {
	isRunning := false
	room := types.Room{ID: bson.NewObjectId().Hex(), BaseTime: BaseTime, AdditionalTime: AdditionalTime, Game: game.New(), IsRunning: &isRunning}
	list[room.ID] = room
	return room.ID
}

func AddPlayerToRoom(roomID string, conn *websocket.Conn, color piece.Color) (string, string, error) {
	room := list[roomID]

	log.Println("AddPlayerToRoom, room: ", room)

	if room.Player1 != nil && room.Player2 != nil {
		return "", "", errors.New("full")
	}

	playerStopper := getTimeFromBaseAndAdditional(list[roomID].BaseTime)

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
	list[roomID] = room
	return player.ID, player.Color.String(), nil

}

func PauseGame(conn *websocket.Conn) {
	for roomID, room := range list {
		if conn == room.Player1.Conn || conn == room.Player2.Conn {
			if *room.IsRunning {
				*room.IsRunning = false
				Print(roomID)
			}
			return
		}
	}
}

func PrintAll() {
	for roomID, _ := range list {
		Print(roomID)
	}
}

func Print(roomID string) {
	room := list[roomID]

	str := "\n{\n"
	str += " RoomID:[" + roomID + "]\n"
	str += " IsRunning:[" + strconv.FormatBool(*room.IsRunning) + "]\n"
	str += " BaseTime:[" + strconv.Itoa(room.BaseTime) + "]\n"
	str += " AdditonalTime:[" + strconv.Itoa(room.AdditionalTime) + "]\n"
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

func NotifyPlayers(roomID string, msg types.Server) {
	room := list[roomID]

	log.Println("NotifyPlayers:")
	Print(roomID)

	room.Player1.Conn.WriteJSON(msg)
	room.Player2.Conn.WriteJSON(msg)

	log.Println("majomkenyerfa")
}

func NotifyOtherPlayer(roomID string, PlayerID string, msg types.Server) {
	room := list[roomID]

	if room.Player1.ID == PlayerID {
		room.Player2.Conn.WriteJSON(msg)
	} else {
		room.Player1.Conn.WriteJSON(msg)
	}
}

func GetRoom(roomID string) types.Room {
	return list[roomID]
}

func GetPlayer(roomID string, playerID string) types.Player {
	room := GetRoom(roomID)
	if room.Player1.ID == playerID {
		return *room.Player1
	}

	return *room.Player2
}

func ActiveNonActivePlayer(room *types.Room) (*types.Player, *types.Player) {
	color := room.Game.ActiveColor()
	if room.Player1.Color == color {
		return room.Player1, room.Player2
	}
	return room.Player2, room.Player1
}

func VerifyTime() {
	for _, room := range list {
		if *room.IsRunning {
			activeP, _ := ActiveNonActivePlayer(&room)

			*activeP.Stopper = (*activeP.Stopper).Add(-time.Second)

			if activeP.Stopper.Hour() == 23 {
				// activeP.Conn.WriteJSON(types.Server{Type: "gameover", Payload: "gameover"})
				// nonActiveP.Conn.WriteJSON(types.Server{Type: "gameover", Payload: "gameover"})
				result := ""
				if activeP.Color == piece.White {
					result = "0-1"
				} else {
					result = "1-0"
				}
				NotifyPlayers(room.ID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: result}})
				*room.IsRunning = false
				continue
			}

			// actualTime := types.Stopper{H: activeP.Stopper.Hour(), M: activeP.Stopper.Minute(), S: activeP.Stopper.Second()}
			// activeP.Conn.WriteJSON(types.Server{Type: "clock", Payload: actualTime})
			// nonActiveP.Conn.WriteJSON(types.Server{Type: "enemyClock", Payload: actualTime})
		}
	}
}
