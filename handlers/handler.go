package handlers

import (
	"log"
	"net/http"
	"time"

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

	time, _ := time.Parse(time.Kitchen, "00:20AM")
	socketpool.Add(conn, time)
	socketpool.Print()

	for {
		clientReq := &types.Client{}
		err := conn.ReadJSON(clientReq)
		if err != nil {
			log.Println(err)
			conn.Close()
			socketpool.RemoveByConn(conn)
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
