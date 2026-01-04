package game

import "sync"

type GameStore struct {
	mu    sync.Mutex
	games map[string]*Game
}

var Store = &GameStore{
	games: make(map[string]*Game),
}

func (s *GameStore) AddGame(g *Game) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[g.ID] = g
}

func (s *GameStore) GetGame(id string) *Game {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.games[id]
}
