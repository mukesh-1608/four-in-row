package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"fourinrow/analytics"
	"fourinrow/db"
	"fourinrow/server"
)

// spaHandler serves the index.html for any unknown route to support React Router (SPA)
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// join the static path with the requested URL path
	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if the file exists on disk
	_, err := os.Stat(path)

	// If the file does NOT exist (like /game, /leaderboard), serve index.html
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the file DOES exist, serve it normally
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	// 1. Initialize Database
	db.InitDB()

	// 2. Initialize Analytics
	// Check for environment variable first (for Cloud), otherwise default to localhost
	kafkaUrl := os.Getenv("KAFKA_BROKER")
	var kafkaBrokers []string
	if kafkaUrl != "" {
		kafkaBrokers = []string{kafkaUrl}
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	analytics.Producer = analytics.NewKafkaProducer(kafkaBrokers, "game-events")
	
	// If Kafka fails (or isn't configured on Cloud), fallback to Stub so app doesn't crash
	if analytics.Producer == nil {
		analytics.Producer = analytics.NewStubProducer()
	}
    defer analytics.Producer.Close()

	// 3. Setup Routes
	http.HandleFunc("/ws", server.WebSocketHandler)
	http.HandleFunc("/leaderboard", server.LeaderboardHandler)

	// 4. Serve Frontend
	spa := spaHandler{staticPath: "./client/dist", indexPath: "index.html"}
	http.Handle("/", spa)

	// 5. Start Server (Cloud Compatible)
	// Render/Heroku provide the PORT variable. We must use it.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default for local development
	}
	
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}