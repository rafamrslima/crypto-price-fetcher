package fileHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rafamrslima/crypto-fetcher/internal/models"
)

func GenerateJson(data []models.Coin) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error to convert data to json.", err)
		return
	}

	file, err := os.Create("output/data.json")
	if err != nil {
		log.Println("Error creating json file.", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		log.Println("Error writing to json file.", err)
		return
	}

	fmt.Println("JSON file saved successfully!")
}
