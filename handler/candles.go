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

func GetCandles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	w.Header().Set("Content-Type", "application/json")

	log.Println("GetCandles(+)")

	var lResp models.CandleResp
	lResp.Status = common.Success

	lPair := r.URL.Query().Get("pair")
	lInterval := r.URL.Query().Get("interval")
	if lPair == "" || lInterval == "" {
		log.Println("Error: HGC01 Missing pair or interval")
		lResp.Status = common.Error
		lResp.Msg = "Missing pair or interval"
		json.NewEncoder(w).Encode(lResp)
		return
	}

	lKey := "candles:" + lInterval + ":" + lPair
	ctx := context.Background()

	lVals, lErr := redis.Rdb.LRange(ctx, lKey, 0, -1).Result()
	if lErr != nil {
		log.Println("Error: HGC02", lErr)
		lResp.Status = common.Error
		lResp.Msg = lErr.Error()
		json.NewEncoder(w).Encode(lResp)
		return
	}

	var lCandles []map[string]interface{}
	for _, lVal := range lVals {
		var lC map[string]interface{}
		if lErr := json.Unmarshal([]byte(lVal), &lC); lErr != nil {
			log.Println("Error: HGC03:", lErr)
			continue
		}
		lCandles = append(lCandles, lC)
	}

	lResp.Candles = make([]interface{}, len(lCandles))
	for i, lC := range lCandles {
		lResp.Candles[i] = lC
	}

	json.NewEncoder(w).Encode(lResp)
	log.Println("GetCandles(-)")
}
