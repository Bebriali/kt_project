package main

import (
	"backend/cmd/agent"
	"backend/cmd/server"
	"backend/internal/agent/collector"
	"backend/internal/models"
	"log"
	"time"
)

// @title           Crypto Monitoring API
// @version         1.0
// @description     This is a sample crypto metrics server.
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
	targetCoins := []string{
		"PEPE", "SHIB", // Мемы для волатильности
		"SOL", "BNB", "XRP", // Топ по капитализации
		"ETH", "DOGE", "BTC", // Твои текущие
		"USDT", "USDC", // Стейблы для пар
	}

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Критическая ошибка сервера: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		agent.Run(myExchanges, targetCoins)
	}
}
