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
	err = ProcessProductData("response.json")
	if err != nil {
		panic(err)
	}
}
