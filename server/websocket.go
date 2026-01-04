package server

import (
	"github.com/mukesh-1608/four-in-row/
analytics"
	"github.com/mukesh-1608/four-in-row/
db"
	"github.com/mukesh-1608/four-in-row/
game"
	"github.com/mukesh-1608/four-in-row/
game/bot"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for simplicity in this assignment
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const DisconnectTimeout = 30 * time.Second

// WebSocketHandler handles the websocket connection
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// NOTE: In a real app, we'd get gameID/playerID from query params or handshake.
	// For this step, we continue mocking the test game but add reconnect logic.
	
	// Query params for reconnect simulation: ?player=tester
	playerID := r.URL.Query().Get("player")
	if playerID == "" {
		playerID = "test-player-1"
	}
	gameID := "test-game"

	log.Printf("Client Connected: %s", playerID)

	g := game.Store.GetGame(gameID)
	
	// Initialize test game if not exists
	if g == nil {
		g = &game.Game{
			ID:          gameID,
			Players:     make(map[string]*game.Player),
			CurrentTurn: "test-player-1", // Fixed start
			Status:      "playing",
		}
		// Human
		g.Players["test-player-1"] = &game.Player{
			ID:          "test-player-1",
			Username:    "Tester",
			Color:       1,
			Conn:        conn,
			IsConnected: true,
		}
		// Bot
		g.Players["cpu"] = &game.Player{
			ID:          "cpu",
			Username:    "Bot",
			Color:       2,
			IsBot:       true,
			IsConnected: true,
		}
		game.Store.AddGame(g)

		// Analytics: Game Started
		analytics.Producer.Emit(analytics.GameEvent{
			Type:      "game_started",
			GameID:    g.ID,
			Timestamp: time.Now().Unix(),
		})
	} else {
		// Reconnect Logic
		if p, ok := g.Players[playerID]; ok {
			// Cancel disconnect timer if exists
			if p.DisconnectTimer != nil {
				p.DisconnectTimer.Stop()
				p.DisconnectTimer = nil
				log.Printf("Player %s reconnected, timer stopped", playerID)
			}
			p.Conn = conn
			p.IsConnected = true
			
			// Broadcast reconnect status? 
			// For simplicity, just send state to reconnected user
		} else {
			// New player trying to join existing game? (Not handled in this step, reject or ignore)
		}
	}

	// Send initial state
	sendState(conn, g)

	// Clean up on disconnect
	defer func() {
		log.Printf("Client Disconnected: %s", playerID)
		handleDisconnect(g, playerID)
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg game.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "move":
			handleMove(conn, g, playerID, msg.Payload)
		case "reset":
			handleReset(g)
		case "join":
			// Handled by connection setup
		}
	}
}

func handleReset(g *game.Game) {
	g.Board = [6][7]int{}
	g.Status = "playing"
	g.Winner = ""
	g.CurrentTurn = "test-player-1"
	broadcastState(g)
	log.Printf("Game %s reset", g.ID)
}

func handleDisconnect(g *game.Game, playerID string) {
	if g.Status == "finished" {
		return
	}

	p, ok := g.Players[playerID]
	if !ok || p.IsBot {
		return
	}

	p.IsConnected = false
	// Start forfeiture timer
	// Reconnect Window: 30 seconds.
	// If the player does not reconnect within this window, they forfeit the game.
	// This prevents games from hanging indefinitely due to connection issues.
	p.DisconnectTimer = time.AfterFunc(DisconnectTimeout, func() {
		// Check if still disconnected (double check lock)
		if !p.IsConnected {
			log.Printf("Player %s timed out. Forfeiting game.", playerID)
			g.Status = "finished"
			g.Winner = "cpu" // In 2p game, winner is opponent. Here it's CPU.
			
			// Notify other players (CPU doesn't care, but if p2 was human...)
			// Since we only have one connection in this scope, we can't easily broadcast 
			// without iterating players.
			broadcast(g, game.WSMessage{
				Type: "game_over",
				Payload: map[string]string{
					"winner": "cpu",
					"reason": "forfeit",
				},
			})
		}
	})
}

func handleMove(conn *websocket.Conn, g *game.Game, playerID string, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	var moveData struct {
		Column int `json:"column"`
	}
	if err := json.Unmarshal(payloadBytes, &moveData); err != nil {
		sendError(conn, "Invalid move payload")
		return
	}

	// Human Move
	if err := game.ApplyMove(g, playerID, moveData.Column); err != nil {
		sendError(conn, err.Error())
		return
	}
	
	// Analytics: Move Played
	analytics.Producer.Emit(analytics.GameEvent{
		Type:      "move_played",
		GameID:    g.ID,
		PlayerID:  playerID,
		Timestamp: time.Now().Unix(),
		Payload:   moveData,
	})

	broadcastState(g)

	if g.Status == "finished" {
		broadcastGameOver(g)
		return
	}

	// Check if next turn is BOT
	if g.CurrentTurn == "cpu" {
		// Trigger Bot Move
		// Small delay for realism (optional, strict rules say immediate, but "Do NOT use goroutines or delays yet" applied to "TRIGGERING")
		// "Bot should respond immediately after a human move" -> OK, synchronous call.
		
		botMoveCol, err := bot.GetBestMove(g, 2) // Bot is color 2
		if err != nil {
			log.Println("Bot error:", err)
			return // Should not happen with valid logic
		}

		// Apply Bot Move
		if err := game.ApplyMove(g, "cpu", botMoveCol); err != nil {
			log.Println("Bot invalid move:", err)
			return
		}

		// Analytics: Bot Move Played
		analytics.Producer.Emit(analytics.GameEvent{
			Type:      "move_played",
			GameID:    g.ID,
			PlayerID:  "cpu",
			Timestamp: time.Now().Unix(),
			Payload:   map[string]int{"column": botMoveCol},
		})

		broadcastState(g)

		if g.Status == "finished" {
			broadcastGameOver(g)
		}
	}
}

func broadcastState(g *game.Game) {
	msg := game.WSMessage{
		Type:    "update",
		Payload: g,
	}
	broadcast(g, msg)
}

func broadcastGameOver(g *game.Game) {
	msg := game.WSMessage{
		Type: "game_over",
		Payload: map[string]string{
			"winner": g.Winner,
		},
	}
	broadcast(g, msg)

	// Analytics: Game Completed
	analytics.Producer.Emit(analytics.GameEvent{
		Type:      "game_completed",
		GameID:    g.ID,
		Timestamp: time.Now().Unix(),
		Payload:   map[string]string{"winner": g.Winner},
	})

	// Persistence: Save Game
	db.Repo.SaveGame(g)
}

func broadcast(g *game.Game, msg game.WSMessage) {
	for _, p := range g.Players {
		if p.Conn != nil && p.IsConnected {
			p.Conn.WriteJSON(msg)
		}
	}
}

func sendState(conn *websocket.Conn, g *game.Game) {
	msg := game.WSMessage{
		Type:    "update",
		Payload: g,
	}
	conn.WriteJSON(msg)
}

func sendError(conn *websocket.Conn, errorMsg string) {
	msg := game.WSMessage{
		Type:    "error",
		Payload: errorMsg,
	}
	conn.WriteJSON(msg)
}
