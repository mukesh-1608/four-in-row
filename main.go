package main

import (
	"log"
	"net/http"

	"fourinrow/analytics"
	"fourinrow/db"
	"fourinrow/server"
)

func main() {
	// Initialize database (safe if DATABASE_URL not set)
	db.InitDB()

	// Initialize analytics (stubbed)
	analytics.Producer = analytics.NewStubProducer()

	http.HandleFunc("/ws", server.WebSocketHandler)
	http.HandleFunc("/leaderboard", server.LeaderboardHandler)

	// Serve minimal frontend
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	log.Println("Server running on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
