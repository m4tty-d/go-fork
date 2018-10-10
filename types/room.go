package types

import (
	"time"

	"github.com/andrewbackes/chess/game"

	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID       string
	Players  map[string]Player
	Game     *game.Game
	GameTime time.Time
}

type Player struct {
	ID      string
	Conn    *websocket.Conn
	Color   piece.Color
	Stopper *time.Time
}

type Stopper struct {
	H int `json:"hour"`
	M int `json:"min"`
	S int `json:"sec"`
}
