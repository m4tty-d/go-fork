package types

import (
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
}

type Player struct {
	Id   string
	Time time.Time
	Conn *websocket.Conn
}

type Watch struct {
	H int `json:"hour"`
	M int `json:"min"`
	S int `json:"sec"`
}
