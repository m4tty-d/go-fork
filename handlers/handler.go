package handlers

import (
	"log"
	"net/http"

	"gitlab.com/chess-fork/go-fork/socketpool"
	"gitlab.com/chess-fork/go-fork/types"

	"github.com/gorilla/websocket"
)

var Upgrader websocket.Upgrader

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	socketpool.Add(conn)
	socketpool.Print()

	for {
		clientReq := &types.Request{}
		err := conn.ReadJSON(clientReq)
		if err != nil {
			log.Println(err)
			return
		}
		switch clientReq.Type {
		case "message":
			socketpool.SendToAll(clientReq.Payload)
		default:
			return
		}
	}
}
