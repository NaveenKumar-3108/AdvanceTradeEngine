package handler

import (
	"AdvanceTradeEngine/common"
	"AdvanceTradeEngine/internal/redis"
	"AdvanceTradeEngine/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	w.Header().Set("Content-Type", "application/json")

	log.Println("GetOrderBook(+)")

	var lResp models.OrderBookResponse
	lResp.Status = common.Success

	lPair := r.URL.Query().Get("pair")
	if lPair == "" {
		log.Println("Error:HGO01 missing pair")
		lResp.Status = common.Error
		lResp.Msg = "Missing pair"
		json.NewEncoder(w).Encode(lResp)
		return
	}

	lKey := "orderbook:" + lPair
	lCtx := context.Background()

	lVal, lErr := redis.Rdb.Get(lCtx, lKey).Result()
	if lErr != nil {
		log.Println("Error:HGO02:", lErr)
		lResp.Status = common.Error
		lResp.Msg = lErr.Error()
		json.NewEncoder(w).Encode(lResp)
		return
	}
	var lOrderBook interface{}
	if lErr := json.Unmarshal([]byte(lVal), &lOrderBook); lErr != nil {
		log.Println("Error:HGO03", lErr)
		lResp.Status = common.Error
		lResp.Msg = lErr.Error()
		json.NewEncoder(w).Encode(lResp)
		return
	}

	lResp.OrderBook = lOrderBook
	json.NewEncoder(w).Encode(lResp)
	log.Println("GetOrderBook(-)")
}
