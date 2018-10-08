package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.com/chess-fork/go-fork/handlers"
	"gitlab.com/chess-fork/go-fork/socketpool"
)

func main() {
	handlers.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			socketpool.VerifyTime()
		}
	}()

	http.HandleFunc("/ws", handlers.Handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8089", nil))
}
