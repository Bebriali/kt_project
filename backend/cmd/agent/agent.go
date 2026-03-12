package agent

import(
	"sync"
	"backend/internal/models"
	"backend/internal/agent/sender"
)

func RunAgent(exchanges []models.Exchange, coins []string) {
	var wg sync.WaitGroup

	for _, coin := range coins {
		for _, ex := range exchanges {
			wg.Add(1)
			go func(c string, e models.Exchange) {
				defer wg.Done()
				stat := ex.GetStat(coin)
				sender.SendStat("http://localhost:8080/update", &stat)
			}(coin, ex)
		}
	}
	wg.Wait()
}
