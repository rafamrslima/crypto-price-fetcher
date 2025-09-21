package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rafamrslima/crypto-fetcher/internal/fileHelper"
	"github.com/rafamrslima/crypto-fetcher/internal/infra"
	"github.com/rafamrslima/crypto-fetcher/internal/models"
)

var allCoins = []string{"bitcoin", "ethereum", "cardano", "chainlink", "solana", "polkadot"}

func main() {
	coins := getCoinsUserInput()
	validCoins := validateCoinsInput(coins)
	data, err := getPrices(validCoins)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printPrices(data)
	saveToFile(data)
	fmt.Println("Exiting the program.")
}

func getCoinsUserInput() []string {
	fmt.Println("Available coins:")
	for _, coin := range allCoins {
		fmt.Println("-", coin)
	}
	fmt.Print("Enter coins to fetch (comma-separated, e.g., bitcoin,ethereum): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	coins := strings.Split(scanner.Text(), ",")
	return coins
}

func validateCoinsInput(coins []string) []string {
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

	return validCoins
}

func getPrices(coinsToFetch []string) ([]models.Coin, error) {
	var wg sync.WaitGroup
	coinsPrices := make([]models.Coin, len(coinsToFetch))
	wg.Add(len(coinsToFetch))

	for i, coin := range coinsToFetch {
		go func(i int) {
			price, err := infra.FetchPrice(coin)
			if err != nil {
				log.Printf("Error fetching %s: %v\n", coin, err)
				return
			}
			coinsPrices[i] = models.Coin{Coin: coin, Usd: price}
			wg.Done()
		}(i)
	}

	wg.Wait()
	return coinsPrices, nil
}

func printPrices(prices []models.Coin) {
	fmt.Println("+-------------+-----------+")
	fmt.Println("| Coin        | USD       |")
	fmt.Println("+-------------+-----------+")
	for _, p := range prices {
		fmt.Printf("| %-11s | %.2f |\n", p.Coin, p.Usd)
	}
	fmt.Println("+-------------+-----------+")
}

func saveToFile(data []models.Coin) {
	fmt.Print("If you want to save the prices, enter 'csv' or 'json'. Otherwise write 'exit' to finish the program. ")
	scannerExport := bufio.NewScanner(os.Stdin)
	scannerExport.Scan()
	input := scannerExport.Text()

	if strings.ToLower(input) == "csv" {
		fileHelper.GenerateCsv(data)
	} else if strings.ToLower(input) == "json" {
		fileHelper.GenerateJson(data)
	}
}
