package main

import (
	"connect4/analytics"
	"connect4/db"
	"connect4/server"
	"log"
	"net/http"
)

func main() {
	// Initialize persistence and analytics
	db.InitDB()
	analytics.Producer = analytics.NewStubProducer()

	// Register handlers
	http.HandleFunc("/ws", server.WebSocketHandler)
	http.HandleFunc("/leaderboard", server.LeaderboardHandler)

	// Serve static files for the minimal frontend
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	log.Println("Server starting on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
