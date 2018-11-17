package handlers

import (
	"encoding/json"
	"log"

	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/rooms"
	"gitlab.com/chess-fork/go-fork/types"
)

type createGameRequest struct {
	Color          string `json:"color"`
	BaseTime       int    `json:"baseTime"`
	AdditionalTime int    `json:"AdditionalTime"`
}

type createGameResponse struct {
	RoomID         string `json:"roomId"`
	BaseTime       int    `json:"baseTime"`
	AdditionalTime int    `json:"additionalTime"`
}

type joinGameRequest struct {
	RoomID string `json:"roomId"`
}

type createPlayerResponse struct {
	PlayerID string `json:"playerId"`
	Color    string `json:"color"`
}

type joinGameResponse struct {
	BaseTime       int `json:"baseTime"`
	AdditionalTime int `json:"additionalTime"`
}

func isTimeValid(base int, additional int) bool {
	return base >= 1 && base <= 120 && additional >= 0 && additional <= 120
}

func CreateRoom(conn *websocket.Conn, payload string) {
	var createGameReq createGameRequest
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
	conn.WriteJSON(types.Server{Type: "roomCreated", Payload: createGameResponse{RoomID: roomID, BaseTime: createGameReq.BaseTime, AdditionalTime: createGameReq.AdditionalTime}})

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
	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: createPlayerResponse{PlayerID: playerID, Color: createGameReq.Color}})
	rooms.PrintAll()
}

func JoinGame(conn *websocket.Conn, payload string) {
	log.Println("!!!!! JOINGAME")
	var joinGameReq joinGameRequest
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

	conn.WriteJSON(types.Server{Type: "roomJoined", Payload: joinGameResponse{BaseTime: rooms.GetRoom(joinGameReq.RoomID).BaseTime, AdditionalTime: rooms.GetRoom(joinGameReq.RoomID).AdditionalTime}})
	conn.WriteJSON(types.Server{Type: "playerCreated", Payload: createPlayerResponse{PlayerID: playerID, Color: color}})
	rooms.NotifyPlayers(joinGameReq.RoomID, types.Server{Type: "gameCanStart"})
}
