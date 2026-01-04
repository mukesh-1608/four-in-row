package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/mukesh-1608/four-in-row/game"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

type LeaderboardEntry struct {
	Username  string `json:"username"`
	TotalWins int    `json:"total_wins"`
}

var Repo *Repository

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("DATABASE_URL not set, persistence disabled")
		return
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to connect to DB:", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Println("Failed to ping DB:", err)
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS games (
		game_id TEXT PRIMARY KEY,
		player1 TEXT,
		player2 TEXT,
		winner TEXT,
		created_at TIMESTAMP,
		finished_at TIMESTAMP
	)`
	if _, err := db.Exec(query); err != nil {
		log.Println("Failed to create table:", err)
	}

	Repo = &Repository{db: db}
	log.Println("Database persistence enabled")
}

func (r *Repository) SaveGame(g *game.Game) {
	if r == nil || r.db == nil {
		return
	}

	var p1, p2 string
	i := 0
	for _, p := range g.Players {
		if i == 0 {
			p1 = p.Username
		} else if i == 1 {
			p2 = p.Username
		}
		i++
	}

	var winnerUsername string
	if g.Winner == "draw" {
		winnerUsername = "draw"
	} else if p, ok := g.Players[g.Winner]; ok {
		winnerUsername = p.Username
	} else {
		winnerUsername = g.Winner
	}

	query := `
	INSERT INTO games (game_id, player1, player2, winner, created_at, finished_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (game_id) DO UPDATE SET
		winner = EXCLUDED.winner,
		finished_at = EXCLUDED.finished_at
	`

	now := time.Now()
	_, err := r.db.Exec(query, g.ID, p1, p2, winnerUsername, now, now)
	if err != nil {
		log.Println("Failed to persist game:", err)
	} else {
		log.Println("Game persisted successfully:", g.ID)
	}
}

func (r *Repository) GetLeaderboard() ([]LeaderboardEntry, error) {
	if r == nil || r.db == nil {
		return []LeaderboardEntry{}, nil
	}

	query := `
	SELECT winner, COUNT(*) AS total_wins
	FROM games
	WHERE winner != 'draw' AND winner IS NOT NULL
	GROUP BY winner
	ORDER BY total_wins DESC
	LIMIT 10
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.Username, &e.TotalWins); err != nil {
			continue
		}
		leaderboard = append(leaderboard, e)
	}

	return leaderboard, nil
}
