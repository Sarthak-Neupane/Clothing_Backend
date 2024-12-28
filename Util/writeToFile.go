package Util

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/Sarthak-Neupane/Clothing_Backend.git/Crawler"

)

// type price struct {
// 	PriceType      string  `json:"priceType"`
// 	FormattedPrice string  `json:"formattedPrice"`
// 	Price          float64 `json:"price"`
// }
// type img struct {
// 	Url string `json:"url"`
// }

// type swatch struct {
// 	Id           string `json:"articleId"`
// 	ColorName    string `json:"colorName"`
// 	ColorHex     string `json:"colorCode"`
// 	ProductImage string `json:"url"`
// }

// type product struct {
// 	Id           string   `json:"id"`
// 	ColorName string `json:"colorName"`
// 	ColorCode string `json:"colors"`
// 	ColorShade string `json:"colourShades"`
// 	PriceState   []price  `json:"prices"`
// 	Images       []img    `json:"images"`
// 	NewArrival   bool     `json:"newArrival"`
// 	ProductImage string   `json:"productImage"`
// 	ProductName  string   `json:"productName"`
// 	Swatches     []swatch `json:"swatches"`	
// 	LinkToSite   string   `json:"url"`
// 	Details      details  `json:"details"`
// }

// type Fit map[string]string

// type Material map[string]string

// type details struct {
// 	Description string   `json:"description"`
// 	Fit         Fit      `json:"fit"`
// 	Material    Material `json:"material"`
// }

func WriteToFile(d []crawler.Product) {
	jsonData, err := json.MarshalIndent(d, "", "	")
	if err != nil {
		fmt.Printf("Error marshalling data: %v\n", err)
		return
	}

	file, err := os.Create("hnm.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	}

	defer file.Close()

	_, err = file.Write(jsonData)

	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return 
	}

	fmt.Println("Written successfully")

}