package main

import (
	"log"
	"net/http"

	"gitlab.com/chess-fork/go-fork/handlers"
)

func main() {
	handlers.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	http.HandleFunc("/ws", handlers.Handler)
	log.Fatal(http.ListenAndServe(":8089", nil))
}
