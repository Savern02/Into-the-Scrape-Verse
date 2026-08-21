package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

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
		return fmt.Errorf("Failed to make request: %v", err)
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

	// Now store this response.json file in the Cloudflare R2 bucket. The relevant info is stored in the environment variables.
	// First, initialize the R2 client using the AWS SDK for Go v2. The R2 endpoint is constructed using the account ID and the base URL for R2.
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	pathPrefix := "Snapshots/ItemSnapshots/" // This is the path in the R2 bucket where the file will be stored. Modify as needed.

	r2Endpoint := "https://" + accountID + ".r2.cloudflarestorage.com"

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return fmt.Errorf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	// Read the JSON data from the response.json file.
	file, err = os.Open("response.json")
	if err != nil {
		return fmt.Errorf("Failed to open response.json file: %v", err)
	}
	defer file.Close()

	// Generate a unique object key for the R2 bucket using the SHA-256 hash of the JSON data.
	objectKey := GenObjectKeyByUUID(pathPrefix)

	// Upload the JSON data to the R2 bucket.
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("Failed to upload JSON data to R2 bucket: %v", err)
	}

	return nil
}

// The purpose of this function is to generate unique object keys for the R2 bucket using the SHA-256 hash of the JSON data.
// Doing this ensures that duplicate JSON data that we try to store will cause a failure.
// The prefix is the path in the bucket where the file will be stored.
// This function works correctly, but the issue arises that Bright Data's scraper returns the JSON data with a timestamp, so even if the product data is the same, the JSON data will be different and thus the hash will be different.
// This means that we will end up with duplicate files in the R2 bucket for the same product data, which is not ideal.
// We need to find a way to remove the timestamp from the JSON data before hashing it, or find another way to generate unique object keys that doesn't rely on the JSON data itself.
func GenObjectKeyByHash(prefix string, filename string) string {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s: %v", filename, err))
	}

	hash := sha256.Sum256(jsonData)
	hashString := hex.EncodeToString(hash[:12]) // Use the first 12 bytes of the hash for a shorter key.
	return fmt.Sprintf("%s%s-output.json", prefix, hashString)
}

func GenObjectKeyByUUID(prefix string) string {
	timestamp := time.Now().Format("2006-01-02")
	uniqueId := uuid.New().String()[:8] // generates an 8 char uuid

	return fmt.Sprintf("%s%s-%s-output.json", prefix, timestamp, uniqueId)
}

// The purpose of this function is to take a search query and apply it to the walmart search url before scraping that specific
// search page result for all of the matching product urls.
func ScrapeWalmartSearchData(searchQuery string) error {
	return nil
}
