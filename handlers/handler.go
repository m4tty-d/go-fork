package handlers

import (
	"log"
	"net/http"

	"../socketPool"
	"../types"

	"github.com/gorilla/websocket"
)

var Upgrader websocket.Upgrader

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	socketPool.Add("randomid", conn)
	socketPool.Print()

	for {
		clientReq := &types.Request{}
		err := conn.ReadJSON(clientReq)
		if err != nil {
			log.Println(err)
			return
		}
		switch clientReq.Type {
		case "message":
			socketPool.SendToAll(clientReq.Payload)
		default:
			return
		}
	}
}
