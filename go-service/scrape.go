package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// A struct that represents the data we want to scrape from a Walmart product page.
type WalmartProductData struct {
	URL          string  `json:"url"`
	Price        float64 `json:"final_price"`
	Currency     string  `json:"currency"`
	Specs        []Spec  `json:"specifications"`
	Brand        string  `json:"brand"`
	Retailer     string  `json:"retailer"`
	Availability bool    `json:"is_available"` // Just checking stock availability for now, but we can add more detailed availability information later if needed.
}

// A struct that represents a single specification coming from the product.
type Spec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WalmartSearchPageData struct {
}

// A function whose name is capitalized is exported and can be accessed from other packages.

// The purpose of this function is to take in a string that specifies the product to scrape and the use the Bright Data API to scrape the data,
// before parsing it and sending the parsed data to our Cloudflare R2 ORM for storage. The function returns an error if any step of the process fails.
func ScrapeWalmartProductData(productURL string, zipcode string) error {

	var apiKey = os.Getenv("BRIGHT_DATA_API_KEY")
	if apiKey == "" {
		panic("BRIGHT_DATA_API_KEY environment variable is not set")
	}

	// The datasetId tells Bright Data which scraper to use and how to parse it. The only way I know of finding this is to go to the scraper's
	//  page on the Bright Data website and look at the curl command. The datasetId will look like "dataset_id=gd_{x}" where {x} is a random string of letters and numbers.
	datasetId := "gd_l95fol7l1ru6rlo116"
	endpoint := "https://api.brightdata.com/datasets/v3/scrape?dataset_id=" + datasetId + "&format=json"
	payload := strings.NewReader(`[
    {"url": "` + productURL + `", "zipcode": "` + zipcode + `"}
	]`)

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}

	// Just adds the API key and content type to the request headers for authentication and content type specification.
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var product []WalmartProductData

	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return fmt.Errorf("Failed to decode JSON response: %v", err)
	}

	fmt.Printf("Scraped product data: %+v\n", product)

	return nil
}

// The purpose of this function is to take a search query and apply it to the walmart search url before scraping that specific
// search page result for all of the matching product urls.
func ScrapeWalmartSearchData(searchQuery string) error {
	return nil
}
