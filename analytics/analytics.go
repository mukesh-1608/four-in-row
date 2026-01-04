package analytics

import "log"

type GameEvent struct {
	Type      string
	GameID    string
	PlayerID string
	Timestamp int64
	Payload   interface{}
}

type ProducerInterface interface {
	Emit(event GameEvent)
}

type StubProducer struct{}

func NewStubProducer() *StubProducer {
	return &StubProducer{}
}

func (s *StubProducer) Emit(event GameEvent) {
	log.Printf("[ANALYTICS] %+v\n", event)
}

var Producer ProducerInterface
