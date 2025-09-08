package redis

import (
	"AdvanceTradeEngine/common"
	"context"
	"log"

	redis "github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	cfg := common.LoadConfig("toml/config.toml")
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println(" Connected to Redis")
}
