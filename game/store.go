package game

import (
	"sync"
)

// GameStore holds all active games
type GameStore struct {
	Games map[string]*Game
	mu    sync.RWMutex
}

// Global store instance
var Store = &GameStore{
	Games: make(map[string]*Game),
}

// AddGame adds a new game to the store
func (gs *GameStore) AddGame(game *Game) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Games[game.ID] = game
}

// GetGame retrieves a game by ID
func (gs *GameStore) GetGame(id string) *Game {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.Games[id]
}

// RemoveGame removes a game from the store
func (gs *GameStore) RemoveGame(id string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	delete(gs.Games, id)
}
