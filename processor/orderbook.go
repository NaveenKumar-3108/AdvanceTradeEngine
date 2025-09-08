package processor

import (
	"AdvanceTradeEngine/internal/redis"
	"AdvanceTradeEngine/models"
	"container/heap"
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

type OrderBook struct {
	Bids PriorityQueue
	Asks PriorityQueue
	mu   sync.RWMutex
}

func NewOrderBook() *OrderBook {
	lBids := make(PriorityQueue, 0)
	lAsks := make(PriorityQueue, 0)
	heap.Init(&lBids)
	heap.Init(&lAsks)
	return &OrderBook{Bids: lBids, Asks: lAsks}
}

func (orderbook *OrderBook) AddOrder(pOrder *models.Order) []models.Trade {
	orderbook.mu.Lock()
	defer orderbook.mu.Unlock()

	var lNewTrades []models.Trade
	var lUpdatedBids []*models.Order
	var lUpdatedAsks []*models.Order

	if pOrder.Side == "BUY" {
		for orderbook.Asks.Len() > 0 && pOrder.Quantity > 0 &&
			(pOrder.Type == "MARKET" || pOrder.Price >= orderbook.Asks[0].Price) {

			lBestAsk := heap.Pop(&orderbook.Asks).(*models.Order)
			lQty := min(pOrder.Quantity, lBestAsk.Quantity)

			lTrade := models.Trade{
				BuyOrderID:  pOrder.ID,
				SellOrderID: lBestAsk.ID,
				Price:       lBestAsk.Price,
				Pair:        pOrder.Pair,
				Quantity:    lQty,
				Timestamp:   time.Now(),
			}

			lTradeJSON, _ := json.Marshal(lTrade)
			ctx := context.Background()
			redis.Rdb.LPush(ctx, "trades:"+pOrder.Pair, lTradeJSON)
			redis.Rdb.LTrim(ctx, "trades:"+pOrder.Pair, 0, 99)

			lNewTrades = append(lNewTrades, lTrade)

			pOrder.Quantity -= lQty
			lBestAsk.Quantity -= lQty

			if lBestAsk.Quantity > 0 {
				heap.Push(&orderbook.Asks, lBestAsk)
			}
		}
		if pOrder.Quantity > 0 && pOrder.Type == "LIMIT" {
			heap.Push(&orderbook.Bids, pOrder)
		}
		lUpdatedBids = append(lUpdatedBids, pOrder)

	} else {
		for orderbook.Bids.Len() > 0 && pOrder.Quantity > 0 &&
			(pOrder.Type == "MARKET" || pOrder.Price <= orderbook.Bids[0].Price) {

			lBestBid := heap.Pop(&orderbook.Bids).(*models.Order)
			lQty := min(pOrder.Quantity, lBestBid.Quantity)

			lTrade := models.Trade{
				BuyOrderID:  lBestBid.ID,
				SellOrderID: pOrder.ID,
				Price:       lBestBid.Price,
				Pair:        pOrder.Pair,
				Quantity:    lQty,
				Timestamp:   time.Now(),
			}

			lTradeJSON, _ := json.Marshal(lTrade)
			ctx := context.Background()
			redis.Rdb.LPush(ctx, "trades:"+pOrder.Pair, lTradeJSON)
			redis.Rdb.LTrim(ctx, "trades:"+pOrder.Pair, 0, 99)
			lNewTrades = append(lNewTrades, lTrade)

			pOrder.Quantity -= lQty
			lBestBid.Quantity -= lQty

			if lBestBid.Quantity > 0 {
				heap.Push(&orderbook.Bids, lBestBid)
			}
		}
		if pOrder.Quantity > 0 && pOrder.Type == "LIMIT" {
			heap.Push(&orderbook.Asks, pOrder)
		}
		lUpdatedAsks = append(lUpdatedAsks, pOrder)
	}
	orderbook.persistToRedis(pOrder.Pair, lUpdatedBids, lUpdatedAsks)
	snapshot := map[string]interface{}{
		"bids": orderbook.Bids,
		"asks": orderbook.Asks,
	}

	data, _ := json.Marshal(snapshot)
	redis.Rdb.Set(context.Background(), "orderbook:"+pOrder.Pair, data, 0)

	return lNewTrades
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (orderbook *OrderBook) persistToRedis(pair string, lUpdatedBids, lUpdatedAsks []*models.Order) {
	ctx := context.Background()

	for _, bid := range lUpdatedBids {
		if bid.Quantity > 0 {
			redis.Rdb.HSet(ctx, "orderbook:"+pair+":bids", bid.Price, bid.Quantity)
		} else {
			lPrice := strconv.Itoa(int(bid.Price))
			redis.Rdb.HDel(ctx, "orderbook:"+pair+":bids", lPrice)
		}
	}

	for _, ask := range lUpdatedAsks {
		if ask.Quantity > 0 {
			redis.Rdb.HSet(ctx, "orderbook:"+pair+":asks", ask.Price, ask.Quantity)
		} else {
			lPrice := strconv.Itoa(int(ask.Price))
			redis.Rdb.HDel(ctx, "orderbook:"+pair+":asks", lPrice)
		}
	}

	diff := map[string]interface{}{
		"bids": lUpdatedBids,
		"asks": lUpdatedAsks,
	}
	diffJSON, _ := json.Marshal(diff)
	redis.Rdb.Publish(ctx, "orderbook:updates:"+pair, diffJSON)

	snapshot, _ := json.Marshal(orderbook)
	redis.Rdb.Set(ctx, "orderbook:"+pair+":snapshot", snapshot, 0)
}
