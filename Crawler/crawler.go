package crawler

import (
	"io"
	// "github.com/gocolly/colly/v2"
	"encoding/json"
	"log"
	"net/http"
)

type APIresponse struct {
	Pagination struct {
		CurrentPage int `json:"currentPage"`
		NextPageNum int `json:"nextPageNum"`
		TotalPages  int `json:"totalPages"`
	} `json:"pagination"`

	PlpList plpList `json:"plpList"`
}

type plpList struct {
	ProductList []product `json:"productList"`
}

type price struct {
	PriceType      string `json:"priceType"`
	FormattedPrice string `json:"formattedPrice"`
	Price          float64 `json:"price"`
}

type color struct {
	ColorName  string `json:"-"`
	ColorShade string `json:"-"`
	ColorHex   string `json:"-"`
}

type img struct {
	Url string `json:"url"`
}

type swatch struct {
	Id           string `json:"articleId"`
	ColorName    string `json:"colorName"`
	ColorHex     string `json:"colorCode"`
	ProductImage string `json:"url"`
}

type product struct {
	Id           string   `json:"id"`
	ColorState   color    `json:"-"`
	PriceState   []price  `json:"prices"`
	Images       []img    `json:"images"`
	NewArrival   bool     `json:"newArrival"`
	ProductImage string   `json:"productImage"`
	ProductName  string   `json:"productName"`
	Swatches     []swatch `json:"swatches"`
}

func (p *product) UnmarshalJSON(data []byte) error {
	type Alias product
	aux := &struct {
		ColorHex   string `json:"colors"`
		ColorName  string `json:"colorName"`
		ColorShade string `json:"colourShades"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.ColorState = color{
		ColorName:  aux.ColorName,
		ColorShade: aux.ColorShade,
		ColorHex:   aux.ColorHex,
	}
	return nil
}

func CrawlHnM() []product {
	var response APIresponse
	resp, err := http.Get("https://api.hm.com/search-services/v1/en_US/listing/resultpage?pageSource=PLP&page=2&sort=RELEVANCE&pageId=/men/shop-by-product/trousers&page-size=36&categoryId=men_trousers&filters=sale:false||oldSale:false&touchPoint=DESKTOP&skipStockCheck=false")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	productList := response.PlpList.ProductList

	return productList
}
