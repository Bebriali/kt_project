package sender

import (
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//func SendStat(url string, data *models.Stat) error {
//					fmt.Printf("I HERE\n")
//	jsonData, err := json.Marshal(data)
//	if err != nil {
//		return fmt.Errorf("failed to marshal data stat: %w", err)
//	}
//
//	// making post request
//	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		return fmt.Errorf("failed to send request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// checking obj status for 201 (created)
//	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("error server response: %s", resp.Status)
//	}
//
//	fmt.Printf("[AGENT] successed sending data. status: %s\n", resp.Status)
//	return nil
//}

func SendStat(url string, data *models.Stat) error {
	if data.AskPrice == 0 || data.BidPrice == 0 {
		return fmt.Errorf("failed to parse price stat")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data stat: %w", err)
	}

	// 1. Создаем клиента, который ждет максимум 2 секунды
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// 2. Делаем запрос через этого клиента
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request (сервер не отвечает): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error server response: %s", resp.Status)
	}

	fmt.Printf("[AGENT] successed sending data. status: %s\n", resp.Status)
	return nil
}
