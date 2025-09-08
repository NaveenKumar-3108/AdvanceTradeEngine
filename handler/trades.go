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

func GetTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	w.Header().Set("Content-Type", "application/json")

	log.Println("GetTrades(+)")

	var lResp models.TradeResp
	lResp.Status = common.Success

	lPair := r.URL.Query().Get("pair")
	if lPair == "" {
		log.Println("Error HGT01: missing pair")
		lResp.Status = common.Error
		lResp.Msg = "Missing pair"
		json.NewEncoder(w).Encode(lResp)
		return
	}

	lKey := "trades:" + lPair
	lCtx := context.Background()

	lVals, lErr := redis.Rdb.LRange(lCtx, lKey, 0, 99).Result()
	if lErr != nil {
		log.Println("Error HGT02:", lErr)
		lResp.Status = common.Error
		lResp.Msg = lErr.Error()
		json.NewEncoder(w).Encode(lResp)
		return
	}

	var trades []interface{}
	for _, v := range lVals {
		var t map[string]interface{}
		if lErr := json.Unmarshal([]byte(v), &t); lErr != nil {
			log.Println("Error HGT03:", lErr)
			continue
		}
		trades = append(trades, t)
	}

	lResp.Trades = trades
	json.NewEncoder(w).Encode(lResp)
	log.Println("GetTrades(-)")
}
