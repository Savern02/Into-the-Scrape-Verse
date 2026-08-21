package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// Just loads all environment variabls from the .env file into the process's environment.
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Set up HTTP server.

	// Testing functions.
	fmt.Println("Scraping Walmart product data...")
	err = ScrapeWalmartProductData("https://www.walmart.com/ip/Hershey-s-Kit-Kat-And-Reese-s-Assorted-Flavored-Candy-Party-Pack-33-38-oz/466623588", "74112")
	if err != nil {
		panic(err)
	}
}
