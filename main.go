package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.com/chess-fork/go-fork/handlers"
	"gitlab.com/chess-fork/go-fork/rooms"
)

import _ "github.com/joho/godotenv/autoload"

func main() {
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

	log.Println("Server running on " + port + " port!")
	http.HandleFunc("/ws", handlers.Handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
