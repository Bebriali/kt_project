package collector

import (
	"backend/internal/models"
	"fmt"
	"strconv"
	"time"
)

type Bybit struct{}

func (b Bybit) GetStat(baseCoin string, quoteCoin string) (models.Stat, error) {
	// Формируем URL. Bybit чувствителен к регистру,
	// если в coins затесались строчные буквы — лучше добавить strings.ToUpper
	url := fmt.Sprintf("https://api.bybit.com/v5/market/tickers?category=spot&symbol=%s%s", baseCoin, quoteCoin)

	var resp struct {
		Result struct {
			List []struct {
				BidPrice string `json:"bid1Price"`
				AskPrice string `json:"ask1Price"`
			} `json:"list"`
		} `json:"result"`
	}

	err := GetJSON(url, &resp)
	if err != nil {
		// Не пишем в stderr здесь, отдаем ошибку наверх в Run
		return models.Stat{}, fmt.Errorf("bybit network error: %w", err)
	}

	// 1. Проверяем, что список не пуст (Bybit возвращает пустой list для неверных пар)
	if len(resp.Result.List) == 0 {
		return models.Stat{}, fmt.Errorf("[BYBIT] symbol %s%s not found", baseCoin, quoteCoin)
	}

	// Берём первый элемент для удобства
	data := resp.Result.List[0]

	// 2. Проверяем на пустые строки в JSON (бывает при технических работах на бирже)
	if data.BidPrice == "" || data.AskPrice == "" {
		return models.Stat{}, fmt.Errorf("[BYBIT] empty price strings for %s%s", baseCoin, quoteCoin)
	}

	// 3. Парсим с проверкой ошибок
	bidprice, errBid := strconv.ParseFloat(data.BidPrice, 64)
	askprice, errAsk := strconv.ParseFloat(data.AskPrice, 64)

	if errBid != nil || errAsk != nil || bidprice == 0 || askprice == 0 {
		return models.Stat{}, fmt.Errorf("[BYBIT] invalid price values for %s%s", baseCoin, quoteCoin)
	}

	return models.Stat{
		Base:     baseCoin,
		Quote:    quoteCoin,
		AskPrice: askprice,
		BidPrice: bidprice,
		Source:   "Bybit",
		Timedump: time.Now(),
	}, nil
}
