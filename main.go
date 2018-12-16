package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gitlab.com/chess-fork/go-fork/handlers"
	"gitlab.com/chess-fork/go-fork/rooms"
	"gitlab.com/chess-fork/go-fork/util"
)

func main() {
	log := util.InitLogger()

	handlers.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			rooms.VerifyTime()
		}
	}()

	port := os.Getenv("PORT")
	log.Debug("Server running on port " + port)
	http.HandleFunc("/ws", handlers.Handler)
	log.Critical(http.ListenAndServe("127.0.0.1:"+port, nil))
}
