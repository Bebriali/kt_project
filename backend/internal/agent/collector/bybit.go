package collector

import (
	"fmt"
	"backend/internal/models"
	"os"
	"time"
	"strconv"
)

type Bybit struct {}


func (b Bybit) GetStat(coin string) models.Stat { //Get information from market
	url := "https://api.bybit.com/v5/market/tickers?category=spot&symbol=" + b.formatSymbolUSDT(coin)
	var resp struct {
    	Result struct {
    		List []struct {
    	        Price string `json:"lastPrice"` //JSON structure
    	    } `json:"list"`
    	} `json:"result"`
	}

	err := GetJSON(url, &resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Dont get JSON (%v)\n", err)
	}
	price, _:= strconv.ParseFloat(resp.Result.List[0].Price, 64)
	return models.Stat{
		Symbol:  coin,
		Price: price,
		Source: "Bybit",
		Timedump: time.Now(),
	}
}

func (b Bybit) formatSymbolUSDT(coin string) string {
    return coin + "USDT"
}
