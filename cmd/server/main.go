package main

import (
	"AdvanceTradeEngine/common"
	"AdvanceTradeEngine/handler"
	"AdvanceTradeEngine/internal/redis"
	"AdvanceTradeEngine/routes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	defer handler.Producer.Close()
	f, lErr := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatalf("Error opening file: %v", lErr)
	}
	defer f.Close()

	log.SetOutput(f)
	redis.InitRedis()
	Router := routes.SetupRoutes()
	cfg := common.LoadConfig("toml/config.toml")

	log.Println("API running on :", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Server.Port), Router))
}
