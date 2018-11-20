package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.com/chess-fork/go-fork/types"
)

var Upgrader websocket.Upgrader

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		client := &types.Client{}
		err := conn.ReadJSON(client)
		if err != nil {
			log.Println(err)
			// rooms.PauseGame(conn)
			conn.Close()
			return
		}
		switch client.Type {
		case "createGame":
			CreateRoom(conn, client.Payload)
		case "joinGame":
			JoinGame(conn, client.Payload)
		case "move":
			Move(conn, client.Payload)
		case "spectateGame":
			Spectate(conn, client.Payload)
		default:
			return
		}
	}
}
