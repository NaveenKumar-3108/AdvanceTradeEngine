package main

import (
	"AdvanceTradeEngine/common"
	"context"

	"log"
	"math/rand"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Order struct {
	Pair      string    `json:"pair"`
	Side      string    `json:"side"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	numWorkers   = 16
	targetOrders = 220_000
)

func main() {
	ctx := context.Background()
	cfg := common.LoadConfig("toml/config.toml")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{cfg.Kafka.Broker},
		Topic:        cfg.Kafka.Topic1,
		Balancer:     &kafka.Hash{},
		BatchSize:    256 * 1024,
		BatchTimeout: 5 * time.Millisecond,
		Async:        true,
	})

	var lTotalSent uint64

	go func() {
		lTicker := time.NewTicker(1 * time.Second)
		defer lTicker.Stop()
		var lLast uint64
		for range lTicker.C {
			cur := atomic.LoadUint64(&lTotalSent)
			log.Printf("[Producer] Orders/sec: %d", cur-lLast)
			lLast = cur
		}
	}()

	lOrderCh := make(chan []byte, 500_000)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				lOrder := Order{
					Pair:      "BTC/USD",
					Side:      []string{"BUY", "SELL"}[rand.Intn(2)],
					Price:     20000 + rand.Float64()*1000,
					Quantity:  rand.Float64() * 5,
					Type:      "LIMIT",
					Timestamp: time.Now().UTC(),
				}
				data, _ := json.Marshal(lOrder)
				lOrderCh <- data
			}
		}()
	}

	for lMsg := range lOrderCh {
		err := writer.WriteMessages(ctx, kafka.Message{Value: lMsg})
		if err != nil {
			log.Println("Error:LM01", err)
			time.Sleep(1 * time.Millisecond)
			continue
		}

		atomic.AddUint64(&lTotalSent, 1)
	}
}
