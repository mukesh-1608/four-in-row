package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"fourinrow/analytics"
	"fourinrow/db"
	"fourinrow/game"
	"fourinrow/game/bot"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const DisconnectTimeout = 30 * time.Second

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	playerID := r.URL.Query().Get("player")
	if playerID == "" {
		playerID = "test-player-1"
	}

	gameID := "test-game"
	log.Printf("Client Connected: %s", playerID)

	g := game.Store.GetGame(gameID)

	if g == nil {
		g = &game.Game{
			ID:          gameID,
			Players:     make(map[string]*game.Player),
			CurrentTurn: "test-player-1",
			Status:      "playing",
		}

		g.Players["test-player-1"] = &game.Player{
			ID:          "test-player-1",
			Username:    "Tester",
			Color:       1,
			Conn:        conn,
			IsConnected: true,
		}

		g.Players["cpu"] = &game.Player{
			ID:          "cpu",
			Username:    "Bot",
			Color:       2,
			IsBot:       true,
			IsConnected: true,
		}

		game.Store.AddGame(g)

		analytics.Producer.Emit(analytics.GameEvent{
			Type:      "game_started",
			GameID:    g.ID,
			Timestamp: time.Now().Unix(),
		})
	} else {
		if p, ok := g.Players[playerID]; ok {
			if p.DisconnectTimer != nil {
				p.DisconnectTimer.Stop()
				p.DisconnectTimer = nil
			}
			p.Conn = conn
			p.IsConnected = true
		}
	}

	sendState(conn, g)

	defer func() {
		log.Printf("Client Disconnected: %s", playerID)
		handleDisconnect(g, playerID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg game.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "move":
			handleMove(conn, g, playerID, msg.Payload)
		case "reset":
			handleReset(g)
		}
	}
}

func handleReset(g *game.Game) {
	g.Board = [6][7]int{}
	g.Status = "playing"
	g.Winner = ""
	g.CurrentTurn = "test-player-1"
	broadcastState(g)
}

func handleDisconnect(g *game.Game, playerID string) {
	p, ok := g.Players[playerID]
	if !ok || p.IsBot {
		return
	}

	p.IsConnected = false

	p.DisconnectTimer = time.AfterFunc(DisconnectTimeout, func() {
		if !p.IsConnected {
			g.Status = "finished"
			g.Winner = "cpu"
			broadcastGameOver(g)
		}
	})
}

func handleMove(conn *websocket.Conn, g *game.Game, playerID string, payload interface{}) {
	var data struct {
		Column int `json:"column"`
	}
	bytes, _ := json.Marshal(payload)
	json.Unmarshal(bytes, &data)

	if err := game.ApplyMove(g, playerID, data.Column); err != nil {
		sendError(conn, err.Error())
		return
	}

	analytics.Producer.Emit(analytics.GameEvent{
		Type:      "move_played",
		GameID:    g.ID,
		PlayerID:  playerID,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})

	broadcastState(g)

	if g.Status == "finished" {
		broadcastGameOver(g)
		return
	}

	if g.CurrentTurn == "cpu" {
		col, _ := bot.GetBestMove(g, 2)
		game.ApplyMove(g, "cpu", col)

		analytics.Producer.Emit(analytics.GameEvent{
			Type:      "move_played",
			GameID:    g.ID,
			PlayerID:  "cpu",
			Timestamp: time.Now().Unix(),
			Payload:   map[string]int{"column": col},
		})

		broadcastState(g)

		if g.Status == "finished" {
			broadcastGameOver(g)
		}
	}
}

func broadcastState(g *game.Game) {
	broadcast(g, game.WSMessage{
		Type:    "update",
		Payload: g,
	})
}

func broadcastGameOver(g *game.Game) {
	broadcast(g, game.WSMessage{
		Type: "game_over",
		Payload: map[string]string{
			"winner": g.Winner,
		},
	})

	analytics.Producer.Emit(analytics.GameEvent{
		Type:      "game_completed",
		GameID:    g.ID,
		Timestamp: time.Now().Unix(),
		Payload:   map[string]string{"winner": g.Winner},
	})

	if db.Repo != nil {
		db.Repo.SaveGame(g)
	}
}

func broadcast(g *game.Game, msg game.WSMessage) {
	for _, p := range g.Players {
		if p.Conn != nil && p.IsConnected {
			p.Conn.WriteJSON(msg)
		}
	}
}

func sendState(conn *websocket.Conn, g *game.Game) {
	conn.WriteJSON(game.WSMessage{
		Type:    "update",
		Payload: g,
	})
}

func sendError(conn *websocket.Conn, msg string) {
	conn.WriteJSON(game.WSMessage{
		Type:    "error",
		Payload: msg,
	})
}
