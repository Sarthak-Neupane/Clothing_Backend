package crawler

import (
	"encoding/json"
	"fmt"
	"io"

	// "log"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"

	"github.com/PuerkitoBio/goquery"
	"sync"
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
	ProductList []Product `json:"productList"`
}

type price struct {
	PriceType      string  `json:"priceType"`
	FormattedPrice string  `json:"formattedPrice"`
	Price          float64 `json:"price"`
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

type Product struct {
	Id           string   `json:"id"`
	// ColorState   color    `json:"-"`
	ColorName string `json:"colorName"`
	ColorCode string `json:"colors"`
	ColorShade string `json:"colourShades"`
	PriceState   []price  `json:"prices"`
	Images       []img    `json:"images"`
	NewArrival   bool     `json:"newArrival"`
	ProductImage string   `json:"productImage"`
	ProductName  string   `json:"productName"`
	Swatches     []swatch `json:"swatches"`	
	LinkToSite   string   `json:"url"`
	Details      details  `json:"details"`
}


var query = map[string]string{
	"pageSource":     "PLP",
	"page":           "1",
	"sort":           "RELEVANCE",
	"pageId":         "/men/shop-by-product/trousers",
	"page-size":      "36",
	"categoryId":     "men_trousers",
	"filters":        "sale:false||oldSale:false",
	"touchPoint":     "DESKTOP",
	"skipStockCheck": "false",
}

type Fit map[string]string

type Material map[string]string

type details struct {
	Description string   `json:"description"`
	Fit         Fit      `json:"fit"`
	Material    Material `json:"material"`
}

func buildQueryURL(baseUrl string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %v", err)
	}

	q := u.Query()
	for key, value := range queryParams {
		query[key] = value
	}
	for key, value := range query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func fetchAPI(apiUrl string) (APIresponse, error) {
	var response APIresponse
	resp, err := http.Get(apiUrl)
	if err != nil {
		return response, fmt.Errorf("failed to fetch API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response, nil
}

func crawlProductDetails(v *Product, url string, errChan chan<- error) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9")
	})

	c.OnError(func(r *colly.Response, err error) {
		errChan <- fmt.Errorf("request failed for URL %s: %v", r.Request.URL, err)
	})

	// Initialize details maps
	v.Details.Fit = make(Fit)
	v.Details.Material = make(Material)

	// Scrape description and fit details
	c.OnHTML(`#section-descriptionAccordion`, func(e *colly.HTMLElement) {
		v.Details.Description = e.ChildText(`p`)
		e.ForEach(`dl > div`, func(_ int, el *colly.HTMLElement) {
			v.Details.Fit[el.ChildText(`dt`)] = el.ChildText(`dd`)
		})
	})

	c.OnHTML(`#section-materialsAndSuppliersAccordion`, func(e *colly.HTMLElement) {
		ul := e.DOM.Find("ul").First()
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			labelText := li.Find("label").Text()
			materialText := li.Find("p").Text()

			if materialText == "" {
				return
			}

			if labelText != "" {
				v.Details.Material[labelText] = materialText
			} else {
				v.Details.Material[fmt.Sprintf("material_%d", i)] = materialText
			}
		})
	})

	if err := c.Visit("https://www2.hm.com" + url); err != nil {
		errChan <- fmt.Errorf("failed to visit URL %s: %v", url, err)
	}
}

func CrawlHnM(queryParams map[string]string) (products []Product, err error) {
	baseURL := "https://api.hm.com/search-services/v1/en_US/listing/resultpage?pageSource=PLP"

	apiURL, err := buildQueryURL(baseURL, queryParams)

	if err != nil {
		return nil, err
	}

	response, err := fetchAPI(apiURL)
	if err != nil {
		return nil, err
	}

	productsList := response.PlpList.ProductList
	var wg sync.WaitGroup
	errChan := make(chan error, len(productsList))


	for i := range productsList {
		wg.Add(1)
		go func(v *Product) {
			defer wg.Done()
			crawlProductDetails(v, v.LinkToSite, errChan)
		}(&productsList[i])
	}
	
	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return productsList, fmt.Errorf("encountered errors while crawling: %v", errs)
	}

	return productsList, nil
}
