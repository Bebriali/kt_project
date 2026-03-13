package agent

import (
	"backend/internal/agent/sender"
	"backend/internal/models"
	"fmt"
	"os"
	"sync"
	"time"
)

func Run(exchanges []models.Exchange, coins []string) {
	var wg sync.WaitGroup

	var allowedPairs = map[string]bool{
		"BTCUSDT": true, "ETHUSDT": true, "DOGEUSDT": true,
		"ETHBTC": true, "SOLBTC": true, "BNBBTC": true,
		"SOLUSDT": true, "PEPEUSDT": true,
	}

	for i, basecoin := range coins {
		time.Sleep(500 * time.Millisecond)
		for _, quotecoin := range coins[i+1:] {

			symbol := basecoin + quotecoin
			if !allowedPairs[symbol] {
				// Если пары нет, пробуем перевернуть
				reverseSymbol := quotecoin + basecoin
				if !allowedPairs[reverseSymbol] {
					continue // Такой пары не существует вообще
				}
				symbol = reverseSymbol
			}

			for _, ex := range exchanges {
				wg.Add(1)
				go func(base string, quote string, e models.Exchange) {
					defer wg.Done()
					stat, err := ex.GetStat(base, quote)
					if err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: (%v)\n", err)
						return
					}
					if stat.AskPrice == 0 || stat.BidPrice == 0 {
						fmt.Fprintf(os.Stderr, "failed to get price for %s-%s\n", base, quote)
						return
					}
					err = sender.SendStat("http://127.0.0.1:8080/update", &stat)
					if err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: (%v)\n", err)
						return
					}
				}(basecoin, quotecoin, ex)
			}
		}
	}
	wg.Wait()
	fmt.Printf("end \n")
}
