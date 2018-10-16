package handlers

import (
	"log"
	"net/http"

	"gitlab.com/chess-fork/go-fork/rooms"
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

	/*time, _ := time.Parse(time.Kitchen, "00:01AM")
	socketpool.Add(conn, &time)
	socketpool.Print()*/

	for {
		client := &types.Client{}
		err := conn.ReadJSON(client)
		if err != nil {
			log.Println(err)
			rooms.PauseGame(conn)
			conn.Close()
			return
		}
		switch client.Type {
		case "createGame":
			CreateRoom(conn, client.Payload)
		case "joinGame":
			JoinGame(conn, client.Payload)
		case "message":
			//socketpool.SendToAll(clientReq.Payload)
		default:
			return
		}
	}
}
