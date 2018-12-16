package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"gitlab.com/chess-fork/go-fork/types"
)

// Upgrader upgrades the HTTP connection to websocket
var Upgrader websocket.Upgrader

// Handler handles websocket messages
func Handler(w http.ResponseWriter, r *http.Request) {
	log := logging.MustGetLogger("log")
	conn, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(err)
		return
	}

	for {
		client := &types.Client{}
		err := conn.ReadJSON(client)
		if err != nil {
			log.Error(err)
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
			Move(client.Payload)
		case "rematch":
			Rematch(client.Payload)
		case "resign":
			Resign(client.Payload)
		case "offerDraw":
			OfferDraw(client.Payload)
		case "acceptDraw":
			AcceptDraw(client.Payload)
		case "spectateGame":
			Spectate(conn, client.Payload)
		default:
			return
		}
	}
}
