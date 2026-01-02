package game

import (
	"github.com/gorilla/websocket"
	"time"
)

// Player represents a connected user
type Player struct {
	ID              string          `json:"id"`
	Username        string          `json:"username"`
	Color           int             `json:"color"` // 1 or 2
	Conn            *websocket.Conn `json:"-"`     // WebSocket connection
	IsBot           bool            `json:"isBot"`
	IsConnected     bool            `json:"isConnected"`
	DisconnectTimer *time.Timer     `json:"-"`
}

// Game represents the state of a Connect Four game
type Game struct {
	ID          string             `json:"id"`
	Board       [6][7]int          `json:"board"`
	Players     map[string]*Player `json:"players"`
	CurrentTurn string             `json:"currentTurn"` // Player ID
	Status      string             `json:"status"`      // "waiting", "playing", "finished"
	Winner      string             `json:"winner,omitempty"`
}

// Move represents a player's action
type Move struct {
	GameID   string `json:"gameId"`
	PlayerID string `json:"playerId"`
	Column   int    `json:"column"`
}

// WSMessage represents the standard message format
type WSMessage struct {
	Type    string      `json:"type"` // "join", "start", "move", "update", "game_over", "error", "opponent_disconnected"
	Payload interface{} `json:"payload"`
}
