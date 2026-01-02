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

	// Simple read loop for the skeleton
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", message)

		// TODO: Handle different message types (join, move, etc.)
		var msg game.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Echo back for now to confirm receipt
		response := game.WSMessage{
			Type:    "echo",
			Payload: "Received: " + msg.Type,
		}
		
		if err := conn.WriteJSON(response); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
