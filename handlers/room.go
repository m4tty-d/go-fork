package handlers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"gitlab.com/chess-fork/go-fork/rooms"
	"gitlab.com/chess-fork/go-fork/types"
)

var log = logging.MustGetLogger("log")

func isTimeValid(base int, additional int) bool {
	return base >= 1 && base <= 120 && additional >= 0 && additional <= 120
}

// CreateRoom creates a room, insert a player in it, and
// sends back the room and player infos
func CreateRoom(conn *websocket.Conn, payload string) {
	var createGameReq types.CreateGameRequest
	err := json.Unmarshal([]byte(payload), &createGameReq)

	if err != nil {
		log.Error(err)
		return
	}

	log.Info(createGameReq)

	if !isTimeValid(createGameReq.BaseTime, createGameReq.AdditionalTime) {
		log.Error("Time is not valid")
		return
	}

	roomID := rooms.CreateRoom(createGameReq.BaseTime, createGameReq.AdditionalTime)
	conn.WriteJSON(types.Server{Type: "roomCreated", Payload: types.CreateGameResponse{RoomID: roomID, BaseTime: createGameReq.BaseTime, AdditionalTime: createGameReq.AdditionalTime}})

	color := piece.White
	if createGameReq.Color == "black" {
		color = piece.Black
	}

	playerID, _, err := rooms.AddAPlayerToRoom(roomID, conn, color)

	if err != nil {
		log.Error(err)
	}

	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: types.CreatePlayerResponse{PlayerID: playerID, Color: createGameReq.Color}})
}

// JoinGame insert a player into the given room, and
// sends back the room and the player infos
func JoinGame(conn *websocket.Conn, payload string) {
	var joinGameReq types.JoinGameRequest
	err := json.Unmarshal([]byte(payload), &joinGameReq)

	if err != nil {
		log.Error(err)
		return
	}

	log.Info(joinGameReq)

	playerID, color, err := rooms.AddAPlayerToRoom(joinGameReq.RoomID, conn, piece.NoColor)

	if err != nil {
		log.Error(err)
		return
	}

	color = strings.ToLower(color)

	room, _ := rooms.GetRoom(joinGameReq.RoomID)

	conn.WriteJSON(types.Server{Type: "roomJoined", Payload: types.JoinGameResponse{BaseTime: room.BaseTime, AdditionalTime: room.AdditionalTime}})
	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: types.CreatePlayerResponse{PlayerID: playerID, Color: color}})
	rooms.NotifyPlayers(joinGameReq.RoomID, types.Server{Type: "gameCanStart"})
}

// Move tries to make a move in a given room, and notifies the other player about it.
// It also checks if the game ended, and if so notifies the players.
func Move(payload string) {
	var moveReq types.MoveRequest
	err := json.Unmarshal([]byte(payload), &moveReq)

	if err != nil {
		log.Error(err)
		return
	}

	log.Info(moveReq)

	room, exists := rooms.GetRoom(moveReq.RoomID)

	if !exists {
		log.Error("No room exists with " + moveReq.RoomID + " ID")
		return
	}

	if !rooms.IsPlayerExistsInRoom(moveReq.RoomID, moveReq.PlayerID) {
		log.Error("No player exists in room " + moveReq.RoomID + " with ID " + moveReq.PlayerID)
		return
	}

	if !(*room.IsRunning) {
		*room.IsRunning = true
	}

	move, err := room.Game.Position().ParseMove(moveReq.Move)

	if err != nil {
		log.Error(err)
		return
	}

	room.Game.MakeMove(move)
	fenStr, _ := fen.Encode(room.Game.Position())

	player := rooms.GetPlayer(moveReq.RoomID, moveReq.PlayerID)
	otherPlayer := rooms.GetOtherPlayer(moveReq.RoomID, moveReq.PlayerID)

	if room.AdditionalTime != 0 {
		*player.Stopper = (*player.Stopper).Add(time.Duration(room.AdditionalTime) * time.Second)
	}

	opponentTime := player.Stopper.Hour()*60*60 + player.Stopper.Minute()*60 + player.Stopper.Second()
	playerTime := otherPlayer.Stopper.Hour()*60*60 + otherPlayer.Stopper.Minute()*60 + otherPlayer.Stopper.Second()

	if *otherPlayer.DrawOffered {
		*otherPlayer.DrawOffered = false
	}

	rooms.NotifyOtherPlayer(moveReq.RoomID, moveReq.PlayerID, types.Server{Type: "move", Payload: types.MoveResponse{Fen: fenStr, Move: moveReq.Move, PlayerSeconds: playerTime, OpponentSeconds: opponentTime}})

	if room.Game.Status() != game.InProgress {
		*room.IsRunning = false
		result := room.Game.Result()
		log.Info(result)

		rooms.IncreasePlayerScores(room.ID, result)

		rooms.NotifyPlayers(room.ID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: room.Game.Result()}})
	}
}

// Rematch resets the game, and notifies players
func Rematch(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Error(err)
		return
	}

	rooms.ResetGame(playerActionReq.RoomID)

	rooms.NotifyPlayers(playerActionReq.RoomID, types.Server{Type: "rematch"})
}

// Resign ends the game and notifies players with the result
func Resign(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Error(err)
		return
	}

	room, _ := rooms.GetRoom(playerActionReq.RoomID)
	player := rooms.GetPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	if &player == nil {
		return
	}

	result := "1-0"
	if player.Color == piece.White {
		result = "0-1"
	}

	rooms.IncreasePlayerScores(playerActionReq.RoomID, result)

	*room.IsRunning = false

	rooms.NotifyPlayers(playerActionReq.RoomID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: result}})
}

// OfferDraw notifies the other player with the offer
func OfferDraw(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Error(err)
		return
	}

	player := rooms.GetPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	*player.DrawOffered = true

	rooms.NotifyOtherPlayer(playerActionReq.RoomID, playerActionReq.PlayerID, types.Server{Type: "drawOffered"})
}

// AcceptDraw ends the game and notifies the players with the result
func AcceptDraw(payload string) {
	var playerActionReq types.PlayerActionRequest
	err := json.Unmarshal([]byte(payload), &playerActionReq)

	if err != nil {
		log.Error(err)
		return
	}

	room, _ := rooms.GetRoom(playerActionReq.RoomID)
	otherPlayer := rooms.GetOtherPlayer(playerActionReq.RoomID, playerActionReq.PlayerID)

	if !*otherPlayer.DrawOffered {
		return
	}

	*room.IsRunning = false

	result := "1/2-1/2"

	rooms.IncreasePlayerScores(playerActionReq.RoomID, result)

	rooms.NotifyPlayers(playerActionReq.RoomID, types.Server{Type: "gameover", Payload: types.GameOverResponse{Result: result}})
}
