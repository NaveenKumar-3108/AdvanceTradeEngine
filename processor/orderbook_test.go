package processor

import (
	"AdvanceTradeEngine/internal/redis"
	"AdvanceTradeEngine/models"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkAddOrderTPS(b *testing.B) {
	ob := NewOrderBook()
	redis.InitRedis()
	order := &models.Order{
		ID:       "bench",
		Side:     "BUY",
		Type:     "LIMIT",
		Price:    100,
		Quantity: 1,
		Pair:     "BTC/USD",
	}

	start := time.Now()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ob.AddOrder(order)
	}

	b.StopTimer()
	elapsed := time.Since(start)

	tps := float64(b.N) / elapsed.Seconds()
	fmt.Printf("\nBenchmark TPS: %.2f orders/sec (N=%d, elapsed=%s)\n", tps, b.N, elapsed)
}

func TestOrderBook_AddOrder_TableDriven(t *testing.T) {
	redis.InitRedis()
	tests := []struct {
		name           string
		initialOrders  []*models.Order
		newOrder       *models.Order
		expectedTrades int
		expectedBids   int
		expectedAsks   int
	}{
		{
			name:          "Buy order no match",
			initialOrders: []*models.Order{},
			newOrder: &models.Order{
				ID:       "buy1",
				Side:     "BUY",
				Type:     "LIMIT",
				Price:    100,
				Quantity: 10,
				Pair:     "BTC/USD",
			},
			expectedTrades: 0,
			expectedBids:   1,
			expectedAsks:   0,
		},
		{
			name: "Sell matches existing buy",
			initialOrders: []*models.Order{
				{
					ID:       "buy1",
					Side:     "BUY",
					Type:     "LIMIT",
					Price:    100,
					Quantity: 10,
					Pair:     "BTC/USD",
				},
			},
			newOrder: &models.Order{
				ID:       "sell1",
				Side:     "SELL",
				Type:     "LIMIT",
				Price:    100,
				Quantity: 5,
				Pair:     "BTC/USD",
			},
			expectedTrades: 1,
			expectedBids:   1,
			expectedAsks:   0,
		},
		{
			name: "Sell exceeds buy quantity",
			initialOrders: []*models.Order{
				{
					ID:       "buy1",
					Side:     "BUY",
					Type:     "LIMIT",
					Price:    100,
					Quantity: 5,
					Pair:     "BTC/USD",
				},
			},
			newOrder: &models.Order{
				ID:       "sell1",
				Side:     "SELL",
				Type:     "LIMIT",
				Price:    100,
				Quantity: 10,
				Pair:     "BTC/USD",
			},
			expectedTrades: 1,
			expectedBids:   0,
			expectedAsks:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := NewOrderBook()

			for _, o := range tt.initialOrders {
				ob.AddOrder(o)
			}

			trades := ob.AddOrder(tt.newOrder)
			if len(trades) != tt.expectedTrades {
				t.Errorf("Expected %d trades, got %d", tt.expectedTrades, len(trades))
			}
			if len(ob.Bids) != tt.expectedBids {
				t.Errorf("Expected %d bids, got %d", tt.expectedBids, len(ob.Bids))
			}
			if len(ob.Asks) != tt.expectedAsks {
				t.Errorf("Expected %d asks, got %d", tt.expectedAsks, len(ob.Asks))
			}
		})
	}
}

func TestOrderBook_AddOrder_Concurrent(t *testing.T) {
	ob := NewOrderBook()

	tests := []struct {
		order *models.Order
	}{
		{order: &models.Order{ID: "buy1", Side: "BUY", Type: "LIMIT", Price: 100, Quantity: 1, Pair: "BTC/USD"}},
		{order: &models.Order{ID: "buy2", Side: "BUY", Type: "LIMIT", Price: 101, Quantity: 1, Pair: "BTC/USD"}},
		{order: &models.Order{ID: "sell1", Side: "SELL", Type: "LIMIT", Price: 100, Quantity: 1, Pair: "BTC/USD"}},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		go func(o *models.Order) {
			defer wg.Done()
			ob.AddOrder(o)
		}(tt.order)
	}
	wg.Wait()

	if len(ob.Bids)+len(ob.Asks) == 0 {
		t.Errorf("Expected orders in heap, got empty")
	}
}
