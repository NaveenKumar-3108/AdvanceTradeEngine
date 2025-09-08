package models

import "time"

type Config struct {
	Server struct {
		Port int
	}
	Kafka struct {
		Broker string
		Topic1 string
		Topic2 string
	}
	Redis struct {
		Address  string
		Password string
		DB       int
	}
}

type Order struct {
	ID        string    `json:"id"`
	Side      string    `json:"side"`
	Type      string    `json:"type"`
	Pair      string    `json:"pair"`
	Price     float64   `json:"price,omitempty"`
	Quantity  float64   `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}
type Trade struct {
	Pair        string    `json:"pair"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

type OrderResponse struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
	Msg     string `json:"msg"`
}

type OrderBookResponse struct {
	OrderBook interface{} `json:"orderbook"`
	Status    string      `json:"status"`
	Msg       string      `json:"msg"`
}

type Candle struct {
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
	Start  time.Time `json:"start"`
}

type CandleResp struct {
	Candles []interface{} `json:"candles"`
	Status  string        `json:"status"`
	Msg     string        `json:"msg"`
}

type TradeResp struct {
	Status string        `json:"status"`
	Msg    string        `json:"msg,omitempty"`
	Trades []interface{} `json:"trades,omitempty"`
}
