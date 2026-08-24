package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Structs to determine the data we want to scrape from the pages. Modify and correct as needed to get the relevant data.

// A struct that represents the data we want to scrape from a Walmart product page.
type WalmartProductData struct {
	ProductName  string   `json:"product_name"`
	URL          string   `json:"url"`
	Price        float64  `json:"final_price"`
	Currency     string   `json:"currency"`
	Specs        []Spec   `json:"specifications"`
	Images       []string `json:"image_urls"`
	Brand        string   `json:"brand"`
	Retailer     string   `json:"retailer"`
	Availability bool     `json:"is_available"` // Just checking stock availability for now, but we can add more detailed availability information later if needed.
	Zipcode      string   `jsdon:"zip_code"`
}

// A struct that represents a single specification coming from the product.
type Spec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WalmartSearchPageData struct {
	Products string `json:"url"`
}

func ProcessProductData(filename string) error {

	// This code block will be used to parse the JSON info from the snapshots stored in the Cloudflare R2 bucket.
	// Some adjustment will be needed to replace the response.json file with the actual snapshot file from the R2 bucket, but this is a good starting point for parsing the JSON data.
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open json file: %v", err)
	}
	defer reader.Close()

	var product []WalmartProductData
	err = json.NewDecoder(reader).Decode(&product)
	if err != nil {
		return fmt.Errorf("Failed to decode JSON raw: %v", err)
	}

	// Now read the parsed info into the response.json
	file, err := os.Create("response.json")
	if err != nil {
		return fmt.Errorf("Failed to create response.json file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ")
	err = encoder.Encode(product)
	if err != nil {
		return fmt.Errorf("Failed to write product data to file: %v", err)
	}

	//Store the info in Cloudflare
	err = StoreObjectInR2(file.Name(), "Items/", GenObjectKeyByHash)
	if err != nil {
		return fmt.Errorf("Failed to store product in R2: %v", err)
	}

	// Store the info in supabase.
	err = StoreInSupabase(product[0])
	if err != nil {
		return fmt.Errorf("Error storing product data: %v", err)
	}

	return nil
}

func ProcessPageData(filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open json file: %v", err)
	}
	defer reader.Close()

	var page []WalmartSearchPageData
	err = json.NewDecoder(reader).Decode(&page)
	if err != nil {
		return fmt.Errorf("Failed to decode JSON raw: %v", err)
	}

	// Now read the parsed info into the response.json
	file, err := os.Create("response.json")
	if err != nil {
		return fmt.Errorf("Failed to create response.json file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ")
	err = encoder.Encode(page)
	if err != nil {
		return fmt.Errorf("Failed to write page data to file: %v", err)
	}

	//Store the info in Cloudflare
	err = StoreObjectInR2("response.json", "SearchPages/", GenObjectKeyByHash)
	if err != nil {
		return fmt.Errorf("Failed to store search page urls in R2: %v", err)
	}

	fmt.Printf("Scraped product urls: %+v\n", page)

	return nil
}
