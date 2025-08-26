package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Coin struct {
	Coin string
	Usd  float64
}

const URL string = "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd"

var allCoins = []string{"bitcoin", "ethereum", "cardano", "chainlink", "solana", "polkadot"}

func main() {

	fmt.Println("Available coins:")
	for _, coin := range allCoins {
		fmt.Println("-", coin)
	}
	fmt.Print("Enter coins to fetch (comma-separated, e.g., bitcoin,ethereum): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	coins := strings.Split(scanner.Text(), ",")
	validCoins := []string{}

	for _, coin := range coins {
		coinIsAvailable := false
		coin = strings.ReplaceAll(coin, " ", "")
		for _, availableCoin := range allCoins {
			if strings.ToLower(coin) == availableCoin {
				coinIsAvailable = true
				validCoins = append(validCoins, coin)
				break
			}
		}

		if !coinIsAvailable {
			fmt.Println("coin not found, being ignored:", coin)
		}
	}

	data, err := getPrices(validCoins)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printPrices(data)
}

func getPrices(coinsToFetch []string) ([]Coin, error) {
	var wg sync.WaitGroup
	coinsPrices := make([]Coin, len(coinsToFetch))
	wg.Add(len(coinsToFetch))

	for i, coin := range coinsToFetch {
		go func() {
			price, err := fetchPrice(coin)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", coin, err)
				return
			}
			coinsPrices[i] = Coin{Coin: coin, Usd: price}
			wg.Done()
		}()
	}

	wg.Wait()
	return coinsPrices, nil
}

func fetchPrice(coin string) (float64, error) {
	resp, err := http.Get(URL + "&ids=" + coin)

	if err != nil {
		fmt.Println("error: ", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error: ", err)
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("error: ", err)
	}

	return result[coin]["usd"], nil
}

func printPrices(prices []Coin) {
	fmt.Println("+-------------+-----------+")
	fmt.Println("| Coin        | USD       |")
	fmt.Println("+-------------+-----------+")
	for _, p := range prices {
		fmt.Printf("| %-11s | %.2f |\n", p.Coin, p.Usd)
	}
	fmt.Println("+-------------+-----------+")
}
