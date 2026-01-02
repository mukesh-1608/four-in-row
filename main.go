package main

import (
	"connect4/server"
	"log"
	"net/http"
)

func main() {
	// Initialize the WebSocket handler
	http.HandleFunc("/ws", server.WebSocketHandler)

	// Serve static files for the minimal frontend
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	log.Println("Server starting on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
