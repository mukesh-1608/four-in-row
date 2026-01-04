package server

import (
	"encoding/json"
	"net/http"

	"github.com/mukesh-1608/four-in-row/db"
)

func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if db.Repo == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	leaderboard, err := db.Repo.GetLeaderboard()
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(leaderboard)
}
