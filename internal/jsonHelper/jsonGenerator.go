package jsonHelper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rafamrslima/crypto-fetcher/internal/models"
)

func Generate(data []models.Coin) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error to convert data to json.", err)
		return
	}

	file, err := os.Create("output/data.json")
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
