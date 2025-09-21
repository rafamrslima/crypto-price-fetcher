package fileHelper

import (
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/rafamrslima/crypto-fetcher/internal/models"
)

func GenerateCsv(data []models.Coin) {
	file, err := os.Create("output/data.csv")
	if err != nil {
		log.Println("Error to generate csv file.", err)
		return
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&data, file); err != nil {
		log.Println("Error to generate csv file.", err)
		return
	}

	fmt.Println("CSV file saved successfully!")
}
