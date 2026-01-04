package main

import (
	"log"
	"net/http"

	"github.com/mukesh-1608/four-in-row/analytics"
	"github.com/mukesh-1608/four-in-row/db"
	"github.com/mukesh-1608/four-in-row/server"
)

func main() {
	// Initialize database
	db.InitDB()

	// Initialize analytics (stubbed Kafka producer)
	analytics.Producer = analytics.NewStubProducer()

	// WebSocket endpoint
	http.HandleFunc("/ws", server.WebSocketHandler)

	// Leaderboard endpoint
	http.HandleFunc("/leaderboard", server.LeaderboardHandler)

	// Serve static frontend files
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	log.Println("Connect Four server running on http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
