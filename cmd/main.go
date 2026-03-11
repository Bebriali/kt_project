package main

import (
	"kt_project/cmd/server"
	"kt_project/cmd/agent"
	"kt_project/internal/agent/collector"
	"kt_project/internal/models"
	"time"
)

// @title           Crypto Metrics API
// @version         1.0
// @description     Сервис для сбора и хранения котировок.
// @host            localhost:8080
// @BasePath        /
func main() {
	//
	// 	repo := repository.NewRedisStorage("localhost:6379", "", 0)
	// 	h := &handlers.Handler{Repo: repo}
	myExchanges := []models.Exchange{
		collector.Binance{},
		collector.Bybit{},
		// Kraken
		//CoinBase
	}
	targetCoins := []string{"BTC", "ETH", "DOGE"} //можно с командной строки

	go server.Run()

	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		agent.RunAgent(myExchanges, targetCoins)
	}
}
