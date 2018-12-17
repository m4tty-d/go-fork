package rooms

import (
	"errors"
	"time"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"gitlab.com/chess-fork/go-fork/types"
	"gopkg.in/mgo.v2/bson"
)

var log = logging.MustGetLogger("log")
var roomMap = make(map[string]types.Room)

func convertMinutesToTime(base int) time.Time {
	h := base / 60
	m := base % 60

	return time.Date(0, 0, 0, h, m, 0, 0, time.UTC)
}

// CreateRoom create a new room into the room map, and returns, if successful
// the ID of the room
func CreateRoom(BaseTime int, AdditionalTime int) string {
	room := types.Room{ID: bson.NewObjectId().Hex(), BaseTime: BaseTime, AdditionalTime: AdditionalTime, Game: game.New(), IsRunning: new(bool)}
	roomMap[room.ID] = room

	log.Info(room)
	log.Info("rooms length: %d", len(roomMap))

	return room.ID
}

// AddAPlayerToRoom create a player into the specified room, and returns
// the ID and color of the player and also an error if something went wrong
func AddAPlayerToRoom(roomID string, conn *websocket.Conn, color piece.Color) (string, string, error) {
	room := roomMap[roomID]

	if &room == nil {
		return "", "", errors.New("room not found")
	}

	if room.Player1 != nil && room.Player2 != nil {
		return "", "", errors.New("full")
	}

	playerStopper := convertMinutesToTime(roomMap[roomID].BaseTime)

	if color == piece.NoColor {
		otherPlayer := GetOtherPlayer(roomID, "")

		if otherPlayer.Color == piece.White {
			color = piece.Black
		} else {
			color = piece.White
		}
	}

	player := types.Player{ID: bson.NewObjectId().Hex(), Conn: conn, Color: color, Stopper: &playerStopper, DrawOffered: new(bool)}

	if room.Player1 == nil {
		room.Player1 = &player
	} else {
		room.Player2 = &player
	}

	roomMap[roomID] = room

	return player.ID, player.Color.String(), nil
}

func PauseGame(conn *websocket.Conn) {
	for _, room := range roomMap {
		if conn == room.Player1.Conn || conn == room.Player2.Conn {
			if *room.IsRunning {
				*room.IsRunning = false
			}
			return
		}
	}
}

// NotifyPlayers sends a message to all players in the specified room
func NotifyPlayers(roomID string, msg types.Server) {
	room := roomMap[roomID]

	room.Player1.Conn.WriteJSON(msg)
	room.Player2.Conn.WriteJSON(msg)

	log.Info(msg)
}

// NotifyOtherPlayer sends a message to the other player,
// the one who is not specified
func NotifyOtherPlayer(roomID string, PlayerID string, msg types.Server) {
	room := roomMap[roomID]

	if room.Player1.ID == PlayerID {
		room.Player2.Conn.WriteJSON(msg)
	} else {
		room.Player1.Conn.WriteJSON(msg)
	}

	log.Info(msg)
}

// GetRoom returns a room
func GetRoom(roomID string) (types.Room, bool) {
	room, exists := roomMap[roomID]

	return room, exists
}

// GetPlayer returns a player
func GetPlayer(roomID string, playerID string) types.Player {
	room, _ := GetRoom(roomID)

	if room.Player1.ID == playerID {
		return *room.Player1
	}

	return *room.Player2
}

// IsPlayerExistsInRoom checks whether a player exists with the given ID,
// in the given room
func IsPlayerExistsInRoom(roomID string, playerID string) bool {
	room, exists := GetRoom(roomID)
	if exists {
		return room.Player1.ID == playerID || room.Player2.ID == playerID
	}

	return false
}

// GetOtherPlayer returns the other player,
// the one who is not specified
func GetOtherPlayer(roomID string, playerID string) types.Player {
	room, _ := GetRoom(roomID)
	if room.Player1 != nil && room.Player1.ID == playerID {
		return *room.Player2
	}

	return *room.Player1
}

// ActiveNonActivePlayer returns the active, and the not active player (in that order)
func ActiveNonActivePlayer(room *types.Room) (*types.Player, *types.Player) {
	color := room.Game.ActiveColor()
	if room.Player1.Color == color {
		return room.Player1, room.Player2
	}
	return room.Player2, room.Player1
}

// VerifyTime checks the time in all rooms, and if it runs out
// ends the games as needed
func VerifyTime() {
	for _, room := range roomMap {
		if *room.IsRunning {
			activeP, _ := ActiveNonActivePlayer(&room)

			*activeP.Stopper = (*activeP.Stopper).Add(-time.Second)

			if activeP.Stopper.Hour() == 23 {
				result := ""

				if activeP.Color == piece.White {
					result = "0-1"
				} else {
					result = "1-0"
				}

				IncreasePlayerScores(room.ID, result)

				NotifyPlayers(room.ID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: result}})

				*room.IsRunning = false
				continue
			}
		}
	}
}

// ResetGame resets the game, and switches player colors
func ResetGame(roomID string) {
	room, exists := GetRoom(roomID)

	if exists {
		room.Game = game.New()
		*room.IsRunning = false

		room.Player1.Color, room.Player2.Color = room.Player2.Color, room.Player1.Color
		*room.Player1.Stopper = convertMinutesToTime(room.BaseTime)
		*room.Player2.Stopper = convertMinutesToTime(room.BaseTime)

		roomMap[roomID] = room
	}
}

func getWinnerColorByResult(result string) piece.Color {
	var winnerColor piece.Color

	if result == "1-0" {
		winnerColor = piece.White
	} else if result == "0-1" {
		winnerColor = piece.Black
	} else {
		winnerColor = piece.NoColor
	}

	return winnerColor
}

// IncreasePlayerScores increases the players scores according to a result string
func IncreasePlayerScores(roomID string, result string) {
	room, exists := GetRoom(roomID)

	if exists {
		winnerColor := getWinnerColorByResult(result)

		if winnerColor == room.Player1.Color {
			(*room.Player1).Score++
		} else if winnerColor == room.Player2.Color {
			(*room.Player2).Score++
		} else {
			(*room.Player1).Score += 0.5
			(*room.Player2).Score += 0.5
		}
	}

	roomMap[roomID] = room
}

func GetPlayerTimeInSeconds(player types.Player) int {
	return player.Stopper.Hour()*60*60 + player.Stopper.Minute()*60 + player.Stopper.Second()
}
