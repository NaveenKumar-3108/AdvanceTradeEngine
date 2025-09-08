package handler

import (
	"AdvanceTradeEngine/common"
	"AdvanceTradeEngine/internal/kafka"
	"AdvanceTradeEngine/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var cfg = common.LoadConfig("toml/config.toml")
var Producer = kafka.NewProducer([]string{cfg.Kafka.Broker}, cfg.Kafka.Topic1)

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	w.Header().Set("Content-Type", "application/json")

	log.Println("PlaceOrder(+)")

	if !strings.EqualFold("POST", r.Method) {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var lOrderReq models.Order
	var lOrderResp models.OrderResponse
	lOrderResp.Status = common.Success

	if lErr := json.NewDecoder(r.Body).Decode(&lOrderReq); lErr != nil {
		log.Println("Error :HAP01", lErr)
		lOrderResp.Status = common.Error
		lOrderResp.Msg = lErr.Error()
		_ = json.NewEncoder(w).Encode(lOrderResp)
		return
	}
	lErr := ValidatelOrderRequest(lOrderReq)
	if lErr != nil {
		log.Println("Error :HAP02", lErr)
		lOrderResp.Status = common.Error
		lOrderResp.Msg = lErr.Error()
		_ = json.NewEncoder(w).Encode(lOrderResp)
		return
	}

	lOrderReq.ID = uuid.New().String()
	lOrderReq.Timestamp = time.Now()

	lOrderBytes, lErr := json.Marshal(lOrderReq)
	if lErr != nil {
		log.Println("Error :HAP03", lErr)
		lOrderResp.Status = common.Error
		lOrderResp.Msg = lErr.Error()
		_ = json.NewEncoder(w).Encode(lOrderResp)
		return
	}

	if lErr := Producer.Publish(context.Background(), []byte(lOrderReq.ID), lOrderBytes); lErr != nil {
		log.Println("Error :HAP04", lErr)
		lOrderResp.Status = common.Error
		lOrderResp.Msg = lErr.Error()
		_ = json.NewEncoder(w).Encode(lOrderResp)
		return
	}

	lOrderResp.Msg = "Order placed"
	lOrderResp.OrderId = lOrderReq.ID

	_ = json.NewEncoder(w).Encode(lOrderResp)
	log.Println("PlaceOrder(-)")
}

func ValidatelOrderRequest(pReq models.Order) error {
	if strings.TrimSpace(pReq.Pair) == "" {
		return errors.New("pair is required")
	}

	if pReq.Side != "BUY" && pReq.Side != "SELL" {
		return errors.New("side must be 'BUY' or 'SELL'")
	}

	if pReq.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if pReq.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	return nil
}
