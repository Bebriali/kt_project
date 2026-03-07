package agent

import (
	kr_project/models
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getPrice(url string, data *models.Stat) { //Get information from market
	get_http, err := http.Get(url)
	defer get_http.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
	}
	fmt.Printf("STATUS: %v \n", get_http.Status)

	err = json.NewDecoder(get_http.Body).Decode(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: JSON was not decode | %v\n", err)
	}
}

func ConnectRedis() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Адрес, который мы пробросили в Docker
        Password: "",           // По умолчанию пароля нет
        DB:       0,            // Номер базы
    })
}
