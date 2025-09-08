package routes

import (
	"AdvanceTradeEngine/handler"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/order", handler.PlaceOrder)
	r.HandleFunc("/api/getCandles", handler.GetCandles)
	r.HandleFunc("/api/getOrderBook", handler.GetOrderBook)
	r.HandleFunc("/api/getTrades", handler.GetTrades)
	return r
}
