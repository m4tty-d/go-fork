package handlers

import (
	"encoding/json"
	"log"

	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/rooms"
	"gitlab.com/chess-fork/go-fork/types"
)

func isTimeValid(base int, additional int) bool {
	return base >= 1 && base <= 120 && additional >= 0 && additional <= 120
}

func CreateRoom(conn *websocket.Conn, payload string) {
	var createGameReq types.CreateGameRequest
	err := json.Unmarshal([]byte(payload), &createGameReq)
	if err != nil {
		// conn.WriteJSON(types.Server{Type: "error", Payload: "Wrong json format!"})
		log.Println(err)
		return
	}

	log.Println(createGameReq)
	log.Println(createGameReq.BaseTime)
	log.Println(createGameReq.AdditionalTime)

	if !isTimeValid(createGameReq.BaseTime, createGameReq.AdditionalTime) {
		// conn.WriteJSON(types.Server{Type: "error", Payload: "Unsupported time!"})
		log.Println("timeNotValid")
		return
	}

	roomID := rooms.CreateRoom(createGameReq.BaseTime, createGameReq.AdditionalTime)
	conn.WriteJSON(types.Server{Type: "roomCreated", Payload: types.CreateGameResponse{RoomID: roomID, BaseTime: createGameReq.BaseTime, AdditionalTime: createGameReq.AdditionalTime}})

	var color piece.Color
	if createGameReq.Color == "white" {
		color = piece.White
	} else {
		color = piece.Black
	}

	playerID, _, err := rooms.AddPlayerToRoom(roomID, conn, color)
	if err != nil {
		conn.WriteJSON(types.Server{Type: "error", Payload: "Error while adding player to room!"})
	}
	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: types.CreatePlayerResponse{PlayerID: playerID, Color: createGameReq.Color}})
	rooms.PrintAll()
}

func JoinGame(conn *websocket.Conn, payload string) {
	log.Println("!!!!! JOINGAME")
	var joinGameReq types.JoinGameRequest
	err := json.Unmarshal([]byte(payload), &joinGameReq)
	if err != nil {
		// conn.WriteJSON(types.Server{Type: "error", Payload: "Wrong json format!"})
		log.Println(err)
		return
	}

	log.Println(joinGameReq)

	playerID, color, err := rooms.AddPlayerToRoom(joinGameReq.RoomID, conn, piece.BothColors)
	if err != nil {
		return
	}

	if color == "White" {
		color = "white"
	} else {
		color = "black"
	}

	conn.WriteJSON(types.Server{Type: "roomJoined", Payload: types.JoinGameResponse{BaseTime: rooms.GetRoom(joinGameReq.RoomID).BaseTime, AdditionalTime: rooms.GetRoom(joinGameReq.RoomID).AdditionalTime}})
	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: types.CreatePlayerResponse{PlayerID: playerID, Color: color}})
	rooms.NotifyPlayers(joinGameReq.RoomID, types.Server{Type: "gameCanStart"})
}

func Move(conn *websocket.Conn, payload string) {
	var moveReq types.MoveRequest
	err := json.Unmarshal([]byte(payload), &moveReq)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(moveReq)

	room := rooms.GetRoom(moveReq.RoomID)

	if !(*room.IsRunning) {
		*room.IsRunning = true
	}

	move, err := room.Game.Position().ParseMove(moveReq.Move)

	room.Game.MakeMove(move)

	// TODO: validálni, hogy helyes volt e a lépés

	fenStr, err := fen.Encode(room.Game.Position())

	player := rooms.GetPlayer(moveReq.RoomID, moveReq.PlayerID)
	playerTime := types.Stopper{H: player.Stopper.Hour(), M: player.Stopper.Minute(), S: player.Stopper.Second()}

	rooms.NotifyOtherPlayer(moveReq.RoomID, moveReq.PlayerID, types.Server{Type: "move", Payload: types.MoveResponse{Fen: fenStr, Move: moveReq.Move, Time: playerTime}})

	if room.Game.Status() != game.InProgress {
		rooms.NotifyPlayers(room.ID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: room.Game.Result()}})
		log.Println(room.Game.Result())
		*room.IsRunning = false
	}
}
