package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type WalmartProductData struct {
}

type WalmartSearchPageData struct {
}

func ScrapeWalmartProductData(productString string) error {
	var apiKey = os.Getenv("BRIGHT_DATA_API_KEY")
	if apiKey == "" {
		panic("BRIGHT_DATA_API_KEY environment variable is not set")
	}

	datasetId := "gd_l95fol7l1ru6rlo116"
	endpoint := "https://api.brightdata.com/datasets/v3/scrape?dataset_id=" + datasetId + "&format=json"
	payload := strings.NewReader(`[
    {"url": "https://www.walmart.com/ip/Hershey-s-Kit-Kat-And-Reese-s-Assorted-Flavored-Candy-Party-Pack-33-38-oz/466623588"}
	]`)

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
