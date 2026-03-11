package collector

import (
	"fmt"
	"kt_project/internal/models"
	"os"
	"time"
	"strconv"
)

type Binance struct {}

func (b Binance) GetStat(coin string) models.Stat { //Get information from market
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + b.formatSymbolUSDT(coin)

	var resp struct {
		Price string `json:"price"`
	}

	err := GetJSON(url, &resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Dont get JSON (%v)\n", err)
	}
	price, _:= strconv.ParseFloat(resp.Price, 64)
	return models.Stat{
		Name:  coin,
		Price: price,
		Source: "Binance",
		Timedump: time.Now(),
	}
}

func (b Binance) formatSymbolUSDT(coin string) string {
    return coin + "USDT"
}
