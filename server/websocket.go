package server

import (
	"connect4/game"
	"encoding/json"
	"log"
	"net/http"

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

// WebSocketHandler handles the websocket connection
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New Client Connected")

	// For step 2 testing, we'll create a dummy single-player game
	// so we can test the logic without full matchmaking.
	playerID := "test-player-1"
	gameID := "test-game"
	
	// Ensure game exists for testing
	g := game.Store.GetGame(gameID)
	if g == nil {
		g = &game.Game{
			ID:          gameID,
			Players:     make(map[string]*game.Player),
			CurrentTurn: playerID,
			Status:      "playing",
		}
		// Add test player
		g.Players[playerID] = &game.Player{
			ID:       playerID,
			Username: "Tester",
			Color:    1,
			Conn:     conn,
		}
		// Add dummy opponent so logic works
		g.Players["cpu"] = &game.Player{
			ID:    "cpu",
			Color: 2,
		}
		game.Store.AddGame(g)
	} else {
		// Update connection if re-connecting to test game
		if p, ok := g.Players[playerID]; ok {
			p.Conn = conn
		}
	}

	// Send initial state
	sendState(conn, g)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", message)

		var msg game.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "move":
			handleMove(conn, g, playerID, msg.Payload)
		case "join":
			// Handled by default setup above for Step 2
			log.Println("Join received (auto-handled for Step 2)")
		}
	}
}

func handleMove(conn *websocket.Conn, g *game.Game, playerID string, payload interface{}) {
	// Parse payload manually since it's interface{}
	payloadBytes, _ := json.Marshal(payload)
	var moveData struct {
		Column int `json:"column"`
	}
	if err := json.Unmarshal(payloadBytes, &moveData); err != nil {
		sendError(conn, "Invalid move payload")
		return
	}

	// Apply move logic
	err := game.ApplyMove(g, playerID, moveData.Column)
	if err != nil {
		sendError(conn, err.Error())
		return
	}

	// Broadcast update (to just this connection for now, simpler for Step 2)
	sendState(conn, g)

	// If game over, send game_over
	if g.Status == "finished" {
		gameOverMsg := game.WSMessage{
			Type: "game_over",
			Payload: map[string]string{
				"winner": g.Winner,
			},
		}
		conn.WriteJSON(gameOverMsg)
	} else if g.CurrentTurn != playerID {
		// Auto-switch turn back to player for testing purposes if it's single player test
		// In Step 2 we just want to test logic, so let's allow "self-play" by hacking the turn back?
		// Or better, just wait for the user to send another move as "cpu"?
		// The instructions say "Only the current player may move".
		// To test easily, let's just make the client send moves for both if needed, 
		// OR simpler: we updated the logic to switch turns. 
		// If we want to test solo, we might need to pretend the other player moved.
		// BUT: Requirements say "Reject moves made out of turn".
		// So we must respect the turn.
		// For Step 2 testing, let's just log that it's the other player's turn.
		log.Println("Turn switched to:", g.CurrentTurn)
	}
}

func sendState(conn *websocket.Conn, g *game.Game) {
	msg := game.WSMessage{
		Type:    "update",
		Payload: g,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
	}
}

func sendError(conn *websocket.Conn, errorMsg string) {
	msg := game.WSMessage{
		Type:    "error",
		Payload: errorMsg,
	}
	conn.WriteJSON(msg)
}
