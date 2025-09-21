package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func FetchPrice(coin string) (float64, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env.")
	}
	baseUrl := os.Getenv("BASE_URL_COIN_GECKO")
	resp, err := http.Get(baseUrl + "&ids=" + coin)

	if err != nil {
		fmt.Println("error: ", err)
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("error: ", err)
		return 0, err
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("error: ", err)
		return 0, err
	}

	return result[coin]["usd"], nil
}
