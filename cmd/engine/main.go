package main

import (
	"AdvanceTradeEngine/common"
	"AdvanceTradeEngine/internal/kafka"
	"AdvanceTradeEngine/internal/redis"
	"AdvanceTradeEngine/models"
	"AdvanceTradeEngine/processor"
	"context"
	"log"

	jsoniter "github.com/json-iterator/go"

	"sync/atomic"
	"time"
)

var candleAggregator1s = processor.NewCandleAggregator(time.Second)
var candleAggregator1m = processor.NewCandleAggregator(time.Minute)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	ctx := context.Background()
	redis.InitRedis()
	lOrderBook := processor.NewOrderBook()
	var lErr error

	var lConsumer *kafka.Consumer
	var lProducer *kafka.Producer
	var processed uint64
	cfg := common.LoadConfig("toml/config.toml")
	lConsumer = kafka.NewConsumer([]string{cfg.Kafka.Broker}, cfg.Kafka.Topic1, "engine-group")
	lProducer = kafka.NewProducer([]string{cfg.Kafka.Broker}, cfg.Kafka.Topic2)

	lErr = lConsumer.Ping(ctx)
	if lErr != nil {
		log.Fatal("Error EM01:", lErr)
	}
	log.Println("Kafka ready, starting engine...")
	go StartTradelConsumer(ctx)

	defer lConsumer.Close()
	defer lProducer.Close()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		var last uint64
		for range ticker.C {
			cur := atomic.LoadUint64(&processed)
			log.Printf("[lConsumer] Orders/sec: %d", cur-last)
			last = cur
		}
	}()

	lErr = lConsumer.Consume(ctx, func(msg []byte) error {
		var order models.Order
		if lErr := json.Unmarshal(msg, &order); lErr != nil {
			log.Println("Error EM02:", lErr)
			return lErr
		}

		lTrades := lOrderBook.AddOrder(&order)

		for _, t := range lTrades {
			lBody, _ := json.Marshal(t)
			if lErr := lProducer.Publish(ctx, []byte(t.BuyOrderID), lBody); lErr != nil {
				log.Println("Error EM03:", lErr)
			}
		}
		atomic.AddUint64(&processed, 1)

		return nil
	})

	if lErr != nil {
		log.Fatal("lError EM04:", lErr)
	}
}

func StartTradelConsumer(ctx context.Context) {
	cfg := common.LoadConfig("toml/config.toml")
	lConsumer := kafka.NewConsumer([]string{cfg.Kafka.Broker}, cfg.Kafka.Topic2, "candle-group")
	defer lConsumer.Close()

	lConsumer.Consume(ctx, func(msg []byte) error {
		var lTrade models.Trade
		if lErr := json.Unmarshal(msg, &lTrade); lErr != nil {
			log.Println("lError EMS01:", lErr)
			return lErr
		}

		candleAggregator1s.AddTrade(lTrade)
		candleAggregator1m.AddTrade(lTrade)

		c1sJSON, _ := json.Marshal(candleAggregator1s.Candles[lTrade.Pair])
		redis.Rdb.RPush(ctx, "candles:1s:"+lTrade.Pair, c1sJSON)
		redis.Rdb.LTrim(ctx, "candles:1s:"+lTrade.Pair, -100, -1)

		c1mJSON, _ := json.Marshal(candleAggregator1m.Candles[lTrade.Pair])
		redis.Rdb.RPush(ctx, "candles:1m:"+lTrade.Pair, c1mJSON)
		redis.Rdb.LTrim(ctx, "candles:1m:"+lTrade.Pair, -100, -1)

		return nil
	})
}
