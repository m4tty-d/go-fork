package handlers

import (
	"encoding/json"
	"time"

	"github.com/andrewbackes/chess/piece"

	"gitlab.com/chess-fork/go-fork/rooms"
	"gitlab.com/chess-fork/go-fork/types"

	"github.com/gorilla/websocket"
)

type creategame struct {
	Color          string `json:"color"`
	BaseTime       int    `json:"baseTime"`
	AdditionalTime int    `json:"additionalTime"`
}

func CreateRoom(conn *websocket.Conn, payload string) {
	var creategame creategame
	err := json.Unmarshal([]byte(payload), &creategame)
	if err != nil {
		return
	}
	if creategame.BaseTime >= 1 && creategame.BaseTime <= 120 && creategame.AdditionalTime >= 0 && creategame.AdditionalTime <= 120 {
		s := creategame.AdditionalTime % 60
		m := (creategame.BaseTime + creategame.AdditionalTime/60) % 60
		h := (creategame.BaseTime + creategame.AdditionalTime/60) / 60
		time := time.Date(0, 0, 0, h, m, s, 0, time.UTC)
		roomID := rooms.CreateRoom(time)
		conn.WriteJSON(types.Server{Type: "roomID", Payload: roomID})
		var color piece.Color
		if creategame.Color == "white" {
			color = piece.White
		} else if creategame.Color == "black" {
			color = piece.Black
		} else {
			return
		}
		playerID, err := rooms.AddPlayerToRoom(roomID, conn, color)
		if err != nil {
			return
		}
		conn.WriteJSON(types.Server{Type: "playerID", Payload: playerID})
		rooms.Print()
	} else {
		return
	}
}
