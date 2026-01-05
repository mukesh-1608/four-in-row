package analytics

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type GameEvent struct {
	Type      string      `json:"type"`
	GameID    string      `json:"game_id"`
	PlayerID  string      `json:"player_id"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type ProducerInterface interface {
	Emit(event GameEvent)
	Close()
}

// ---------------------------------------------------------
// 1. KAFKA PRODUCER (The Real Deal)
// ---------------------------------------------------------
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("[ANALYTICS] ⚠️ Failed to start Kafka producer: %v (Using Stub instead)", err)
		return nil
	}

	log.Println("[ANALYTICS] ✅ Connected to Kafka!")
	return &KafkaProducer{producer: p, topic: topic}
}

func (k *KafkaProducer) Emit(event GameEvent) {
	// Ensure timestamp is set
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	// Convert event to JSON bytes
	val, err := json.Marshal(event)
	if err != nil {
		log.Printf("[ANALYTICS] JSON Error: %v", err)
		return
	}

	// Send message
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(event.GameID), // Ensure events for same game go to same partition
		Value: sarama.ByteEncoder(val),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("[ANALYTICS] Failed to send message: %v", err)
	}
}

func (k *KafkaProducer) Close() {
	k.producer.Close()
}

// ---------------------------------------------------------
// 2. STUB PRODUCER (Fallback)
// ---------------------------------------------------------
type StubProducer struct{}

func NewStubProducer() *StubProducer { return &StubProducer{} }
func (s *StubProducer) Emit(event GameEvent) {
	log.Printf("[ANALYTICS STUB] %+v\n", event)
}
func (s *StubProducer) Close() {}

// Global Instance
var Producer ProducerInterface