package types

import (
	"time"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID             string
	Player1        *Player
	Player2        *Player
	Game           *game.Game
	BaseTime       int
	AdditionalTime int
	IsRunning      *bool
}

type Player struct {
	ID          string
	Conn        *websocket.Conn
	Color       piece.Color
	Stopper     *time.Time
	DrawOffered *bool
	Score       float32
}

type Stopper struct {
	H int `json:"hour"`
	M int `json:"min"`
	S int `json:"sec"`
}

type CreateGameRequest struct {
	Color          string `json:"color"`
	BaseTime       int    `json:"baseTime"`
	AdditionalTime int    `json:"AdditionalTime"`
}

type CreateGameResponse struct {
	RoomID         string `json:"roomId"`
	BaseTime       int    `json:"baseTime"`
	AdditionalTime int    `json:"additionalTime"`
}

type JoinGameRequest struct {
	RoomID string `json:"roomId"`
}

type CreatePlayerResponse struct {
	PlayerID string `json:"playerId"`
	Color    string `json:"color"`
}

type JoinGameResponse struct {
	BaseTime       int `json:"baseTime"`
	AdditionalTime int `json:"additionalTime"`
}

type MoveRequest struct {
	PlayerID string `json:"playerId"`
	RoomID   string `json:"roomId"`
	Move     string `json:"move"`
}

type MoveResponse struct {
	Fen             string `json:"fen"`
	Move            string `json:"move"`
	PlayerSeconds   int    `json:"playerSeconds"`
	OpponentSeconds int    `json:"opponentSeconds"`
}

type GameOverResponse struct {
	Result string `json:"result"`
}

type PlayerActionRequest struct {
	PlayerID string `json:"playerId"`
	RoomID   string `json:"roomId"`
}

type RematchResponse struct {
	PlayerScore   float32 `json:"playerScore"`
	OpponentScore float32 `json:"opponentScore"`
}

type StateResponse struct {
	PlayerColor     string `json:"color"`
	BaseTime        int    `json:"baseTime"`
	AdditionalTime  int    `json:"additionalTime"`
	Fen             string `json:"fen"`
	PlayerSeconds   int    `json:"playerSeconds"`
	OpponentSeconds int    `json:"opponentSeconds"`
	Turn            string `json:"turn"`
}
