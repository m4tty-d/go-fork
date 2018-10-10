package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.com/chess-fork/go-fork/handlers"
)

func main() {
	handlers.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			//socketpool.VerifyTime()
		}
	}()
	log.Println("Server running on '8089' port!")
	http.HandleFunc("/ws", handlers.Handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8089", nil))
	//room := types.Room{Players: make(map[string]*types.Player)}
	/*rooms := make(map[string]types.Room)
	rooms["asd"] = types.Room{Players: make(map[string]*types.Player)}
	rooms["asd"].Players["dsa"] = nil
	room := types.Room{}
	room.ID = bson.NewObjectId().Hex()
	room.Game = game.New()
	room.IsRunning = false
	room.Players = make(map[string]*types.Player)
	room.Players["asd"] = nil*/
}
