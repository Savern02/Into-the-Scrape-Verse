package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ScrapeInput map[string]string

// The purpose of this function is to take in a string that specifies the product to scrape and the use the Bright Data API to scrape the data,
// before parsing it and sending the parsed data to our Cloudflare R2 ORM for storage. The function returns an error if any step of the process fails.
func ScrapeProductData(retailer string, productURL string, zipcode string) error {

	retailerDatasetMap := map[string]string{
		"walmart": "gd_m693oc1r1gebnayxq",
		"target":  "gd_ltppk5mx2lp0v1k0vo",
	}

	var apiKey = strings.TrimSpace(os.Getenv("BRIGHT_DATA_API_KEY"))
	if apiKey == "" {
		panic("BRIGHT_DATA_API_KEY environment variable is not set")
	}

	// The datasetId tells Bright Data which scraper to use and how to parse it. The only way I know of finding this is to go to the scraper's
	//  page on the Bright Data website and look at the curl command. The datasetId will look like "dataset_id=gd_{x}" where {x} is a random string of letters and numbers.
	datasetId, retailerExists := retailerDatasetMap[retailer]
	if retailerExists == false {
		return fmt.Errorf("Retailer does not exist in dataset.")
	}
	endpoint := "https://api.brightdata.com/datasets/v3/scrape?dataset_id=" + datasetId + "&format=json"

	var input []ScrapeInput = nil
	switch retailer {
	case "walmart":
		input = []ScrapeInput{
			{"url": productURL, "zip_code": zipcode, "store_id": ""},
		}
	case "target":
		input = []ScrapeInput{
			{"url": productURL, "zipcode": zipcode},
		}
	}
	payloadBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("Could not marshal payload json for request: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("Failed to make request: %v", err)
	}

	// Just adds the API key and content type to the request headers for authentication and content type specification.
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	// Create a new file to store the response body as JSON.
	file, err := os.Create("response.json")
	if err != nil {
		return fmt.Errorf("Failed to create response.json file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to write response body to file: %v", err)
	}

	// Now store this response.json file in the Cloudflare R2 bucket.
	err = StoreObjectInR2("response.json", "Snapshots/ItemSnapshots/", GenObjectKeyByUUID)
	if err != nil {
		return fmt.Errorf("Failed to store object in R2 bucket: %v", err)
	}

	return nil
}

// The purpose of this function is to take a search query and apply it to the walmart search url before scraping that specific
// search page result for all of the matching product urls.
func ScrapeSearchData(retailer string, searchQuery string, zipcode string) error {

	retailerDatasetMap := map[string]string{
		"walmart": "gd_m7khey0wb7wviejgj",
		"target":  "gd_ltppk5mx2lp0v1k0vo",
	}

	var apiKey = strings.TrimSpace(os.Getenv("BRIGHT_DATA_API_KEY"))
	if apiKey == "" {
		panic("BRIGHT_DATA_API_KEY environment variable is not set")
	}

	// The datasetId tells Bright Data which scraper to use and how to parse it. The only way I know of finding this is to go to the scraper's
	//  page on the Bright Data website and look at the curl command. The datasetId will look like "dataset_id=gd_{x}" where {x} is a random string of letters and numbers.
	datasetId, retailerExists := retailerDatasetMap[retailer]
	if retailerExists == false {
		return fmt.Errorf("Retailer does not exist in dataset.")
	}
	endpoint := "https://api.brightdata.com/datasets/v3/scrape?dataset_id=" + datasetId + "&format=json"

	var input []ScrapeInput = nil
	switch retailer {
	case "walmart":
		input = []ScrapeInput{
			{"url": "https://www.walmart.com/search?q=" + searchQuery, "zip_code": zipcode, "store_id": ""},
		}
	case "target":
		input = []ScrapeInput{
			{"keywords": searchQuery, "zipcode": zipcode},
		}
	}
	payloadBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("Could not marshal payload json for request: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("Failed to make request: %v", err)
	}

	// Just adds the API key and content type to the request headers for authentication and content type specification.
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Store the response into a json file that we can then submit to the R2 object storage.
	file, err := os.Create("response.json")
	if err != nil {
		return fmt.Errorf("Failed to create response.json file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to write response body to file: %v", err)
	}

	// Now store this response.json file in the Cloudflare R2 bucket.
	err = StoreObjectInR2("response.json", "Snapshots/SearchSnapshots/", GenObjectKeyByUUID)
	if err != nil {
		return fmt.Errorf("Failed to store object in R2 bucket: %v", err)
	}

	return nil
}
