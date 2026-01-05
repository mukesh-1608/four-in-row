package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
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

// In-Memory Stats
var (
	gameStartTimes = make(map[string]int64) 
	durations      = []float64{}
	winCounts      = make(map[string]int)
	gamesPerHour   = make(map[string]int) 
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "game-events"

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create Consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	// Consume Partition 0
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to listen to topic: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("📊 Analytics Service Started. Listening for events...")

	// Handle Exit (Ctrl+C)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Loop forever
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			processMessage(msg.Value)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		case <-signals:
			log.Println("Shutting down analytics...")
			return
		}
	}
}

func processMessage(value []byte) {
	var event GameEvent
	if err := json.Unmarshal(value, &event); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	// 1. Track Games Per Hour
	t := time.Unix(event.Timestamp, 0)
	timeKey := t.Format("2006-01-02 15:00")
	gamesPerHour[timeKey]++

	switch event.Type {
	case "game_started":
		gameStartTimes[event.GameID] = event.Timestamp
		fmt.Printf("\n[EVENT] Game Started: %s (Type: %v)\n", event.GameID, event.Payload)

	case "game_finished":
		// Calculate Duration
		if startTime, exists := gameStartTimes[event.GameID]; exists {
			duration := float64(event.Timestamp - startTime)
			durations = append(durations, duration)
			delete(gameStartTimes, event.GameID) 
			
			var total float64
			for _, d := range durations { total += d }
			avg := total / float64(len(durations))
			fmt.Printf("⏱️  Game Over. Duration: %.0fs | Avg Duration: %.1fs\n", duration, avg)
		}

		// Track Wins
		winner := fmt.Sprintf("%v", event.Payload)
		if winner != "" && winner != "draw" {
			winCounts[winner]++
			fmt.Printf("🏆 Winner: %s | Total Wins: %d\n", winner, winCounts[winner])
		}
	}
}