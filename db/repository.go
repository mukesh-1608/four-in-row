package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"fourinrow/game"

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
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Println("DATABASE_URL not set, DB disabled")
		return
	}

	db, err := sql.Open("postgres", url)
	if err != nil || db.Ping() != nil {
		log.Println("DB connection failed")
		return
	}

	db.Exec(`
	CREATE TABLE IF NOT EXISTS games (
		game_id TEXT PRIMARY KEY,
		player1 TEXT,
		player2 TEXT,
		winner TEXT,
		created_at TIMESTAMP,
		finished_at TIMESTAMP
	)`)

	Repo = &Repository{db: db}
}

func (r *Repository) SaveGame(g *game.Game) {
	if r == nil {
		return
	}

	var p1, p2 string
	i := 0
	for _, p := range g.Players {
		if i == 0 {
			p1 = p.Username
		} else {
			p2 = p.Username
		}
		i++
	}

	winner := g.Winner
	if p, ok := g.Players[g.Winner]; ok {
		winner = p.Username
	}

	now := time.Now()
	r.db.Exec(`
	INSERT INTO games VALUES ($1,$2,$3,$4,$5,$6)
	ON CONFLICT (game_id) DO UPDATE SET winner=$4, finished_at=$6
	`, g.ID, p1, p2, winner, now, now)
}

func (r *Repository) GetLeaderboard() ([]LeaderboardEntry, error) {
	rows, err := r.db.Query(`
	SELECT winner, COUNT(*) FROM games
	WHERE winner != 'draw'
	GROUP BY winner
	ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		rows.Scan(&e.Username, &e.TotalWins)
		res = append(res, e)
	}
	return res, nil
}
