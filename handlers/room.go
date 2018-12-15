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

func Move(payload string) {
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

	if err != nil {
		return
	}

	room.Game.MakeMove(move)

	fenStr, err := fen.Encode(room.Game.Position())

	player := rooms.GetPlayer(moveReq.RoomID, moveReq.PlayerID)
	otherPlayer := rooms.GetOtherPlayer(moveReq.RoomID, moveReq.PlayerID)
	playerTime := types.Stopper{H: player.Stopper.Hour(), M: player.Stopper.Minute(), S: player.Stopper.Second()}

	if *otherPlayer.DrawOffered {
		*otherPlayer.DrawOffered = false
	}

	rooms.NotifyOtherPlayer(moveReq.RoomID, moveReq.PlayerID, types.Server{Type: "move", Payload: types.MoveResponse{Fen: fenStr, Move: moveReq.Move, Time: playerTime}})

	if room.Game.Status() != game.InProgress {
		rooms.NotifyPlayers(room.ID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: room.Game.Result()}})
		log.Println(room.Game.Result())
		*room.IsRunning = false
	}
}

func Rematch(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Println(err)
		return
	}

	room := rooms.GetRoom(playerActionReq.RoomID)

	room = types.Room{ID: room.ID, BaseTime: room.BaseTime, AdditionalTime: room.AdditionalTime, Game: game.New(), IsRunning: new(bool)}
	room.Game = game.New()
	*room.IsRunning = false

	if room.Player1.Color == piece.White {
		room.Player1.Color = piece.Black
		room.Player2.Color = piece.White
	} else {
		room.Player1.Color = piece.White
		room.Player2.Color = piece.Black
	}

	rooms.NotifyPlayers(room.ID, types.Server{Type: "rematch", Payload: types.GameOverResponse{Result: room.Game.Result()}})
}

func Resign(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Print(err)
		return
	}

	player := rooms.GetPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	if &player == nil {
		return
	}

	result := "1-0"

	if player.Color == piece.White {
		result = "0-1"
	}

	rooms.NotifyPlayers(playerActionReq.RoomID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: result}})
}

func OfferDraw(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Print(err)
		return
	}

	player := rooms.GetPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	*player.DrawOffered = true

	rooms.NotifyOtherPlayer(playerActionReq.RoomID, playerActionReq.PlayerID, types.Server{Type: "drawOffered"})
}

func AcceptDraw(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Print(err)
		return
	}

	otherPlayer := rooms.GetOtherPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	if !*otherPlayer.DrawOffered {
		return
	}

	rooms.NotifyPlayers(playerActionReq.RoomID, types.Server{Type: "gameover", Payload: "1/2-1/2"})
}
