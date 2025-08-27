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

	"github.com/gocarina/gocsv"
)

type Coin struct {
	Coin string
	Usd  float64
}

const URL string = "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd"

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
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error: ", err)
		return 0, err
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("error: ", err)
		return 0, err
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

func saveToFile(data []Coin) {
	fmt.Print("If you want to save the prices, enter 'csv' or 'json'. otherwise write 'exit' to finish the program. ")
	scannerExport := bufio.NewScanner(os.Stdin)
	scannerExport.Scan()
	input := scannerExport.Text()

	if strings.ToLower(input) == "csv" {
		generateCSV(data)
	} else if strings.ToLower(input) == "json" {
		generateJson(data)
	}
}

func generateCSV(data []Coin) {
	file, err := os.Create("data.csv")
	if err != nil {
		fmt.Println("Error to generate csv file.", err)
		return
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&data, file); err != nil {
		fmt.Println("Error to generate csv file.", err)
		return
	}

	fmt.Println("CSV file saved successfully!")
}

func generateJson(data []Coin) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error to convert data to json.", err)
		return
	}

	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Error creating json file.", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		fmt.Println("Error writing to json file.", err)
		return
	}

	fmt.Println("JSON file saved successfully!")
}
