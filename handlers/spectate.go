package handlers

import (
	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/spectate"
)

func Spectate(conn *websocket.Conn, roomID string) {
	spectate.AddConnToRoom(roomID, conn)
}
