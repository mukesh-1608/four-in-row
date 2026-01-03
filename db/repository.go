package db

import (
        "connect4/game"
        "database/sql"
        "log"
        "os"
        "time"

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

        // Create table if not exists (logical schema)
        // In production, use migrations. Here for the assignment requirement "SCHEMA (logical, not migrations)",
        // we'll ensure it exists to make it functional if the user provides a DB.
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

        // Extract players
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

        // Winner is ID, need username or store ID.
        // The schema says "winner (text)".
        // Logic uses ID for winner. Let's find the username.
        var winnerUsername string
        if g.Winner == "draw" {
                winnerUsername = "draw"
        } else {
                if p, ok := g.Players[g.Winner]; ok {
                        winnerUsername = p.Username
                } else {
                        winnerUsername = g.Winner // Fallback
                }
        }

        // Basic query
        query := `
                INSERT INTO games (game_id, player1, player2, winner, created_at, finished_at)
                VALUES ($1, $2, $3, $4, $5, $6)
        `
        
        // Just use current time for timestamps if not tracked in game struct
        now := time.Now()
        
        _, err := r.db.Exec(query, g.ID, p1, p2, winnerUsername, now, now) // simplified timestamps
        if err != nil {
                log.Println("Failed to persist game:", err)
                // Requirement: Failures to persist must NOT crash the server
        } else {
                log.Println("Game persisted successfully:", g.ID)
        }
}

func (r *Repository) GetLeaderboard() ([]LeaderboardEntry, error) {
        if r == nil || r.db == nil {
                return []LeaderboardEntry{}, nil
        }

        // "Order by wins descending. Draws do NOT count as wins"
        query := `
                SELECT winner, COUNT(*) as total_wins
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
