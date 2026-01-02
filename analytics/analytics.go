package analytics

import (
	"log"
)

type GameEvent struct {
	Type      string      `json:"type"` // "game_started", "move_played", "game_completed"
	GameID    string      `json:"game_id"`
	PlayerID  string      `json:"player_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload,omitempty"`
}

type AnalyticsProducer interface {
	Emit(event GameEvent)
}

type StubProducer struct{}

func NewStubProducer() *StubProducer {
	return &StubProducer{}
}

func (p *StubProducer) Emit(event GameEvent) {
	// In a real implementation, this would send to Kafka.
	// For now, we just log to stdout as an abstraction.
	log.Printf("[ANALYTICS] Type: %s | Game: %s | Player: %s | Payload: %v",
		event.Type, event.GameID, event.PlayerID, event.Payload)
}

// Global instance
var Producer AnalyticsProducer = NewStubProducer()
