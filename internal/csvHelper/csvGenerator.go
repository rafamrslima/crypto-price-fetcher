package csvHelper

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/rafamrslima/crypto-fetcher/internal/models"
)

func Generate(data []models.Coin) {
	file, err := os.Create("output/data.csv")
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
