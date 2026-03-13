package collector

import (
	"backend/internal/models"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Binance struct{}

func (b Binance) GetStat(baseCoin string, quoteCoin string) (models.Stat, error) { //Get information from market
	url := "https://api.binance.com/api/v3/ticker/bookTicker?symbol=" + baseCoin + quoteCoin

	var resp struct {
		Bidprice string `json:"bidPrice"`
		Askprice string `json:"askPrice"`
	}

	err := GetJSON(url, &resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not get BINANCE JSON (%v)\nwith api %s\n", err, url)
		return models.Stat{}, fmt.Errorf("binance api error: %w", err)
	}

	// checking if resp is empty
	if resp.Bidprice == "" || resp.Askprice == "" {
		return models.Stat{}, fmt.Errorf("no data for symbol %s%s", baseCoin, quoteCoin)
	}

	bidprice, _ := strconv.ParseFloat(resp.Bidprice, 64)
	askprice, _ := strconv.ParseFloat(resp.Askprice, 64)

	if bidprice == 0 || askprice == 0 {
		return models.Stat{}, fmt.Errorf("[BINANCE] could not resolve price for symbol %s-%s", baseCoin, quoteCoin)
	}

	return models.Stat{
		Base:     baseCoin,
		Quote:    quoteCoin,
		AskPrice: askprice,
		BidPrice: bidprice,
		Source:   "Binance",
		Timedump: time.Now(),
	}, nil
}
