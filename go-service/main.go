package main

import (
	"github.com/joho/godotenv"
)

func main() {
	// Just loads all environment variabls from the .env file into the process's environment.
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Set up HTTP server.

	// Testing functions
	err = ScrapeProductData("walmart", "https://www.walmart.com/ip/Red-Bull-Amber-Edition-Sugar-Free-Energy-Drink-Strawberry-Apricot-12-fl-oz-Pack-of-4-Cans/5332753715", "74112")
	if err != nil {
		panic(err)
	}
}
