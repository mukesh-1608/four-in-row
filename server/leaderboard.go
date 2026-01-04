package server

import (
	"encoding/json"
	"net/http"

	"fourinrow/db"
)

func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if db.Repo == nil {
		http.Error(w, "DB unavailable", 503)
		return
	}

	data, _ := db.Repo.GetLeaderboard()
	json.NewEncoder(w).Encode(data)
}
