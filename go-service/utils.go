package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
	"uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// The purpose of this function is to take a specified file name and Cloudflare R2 bucket path, and then use that information to store the file in the
// designated location in the R2 storage bucket. The path parameter expects a path that ends in a /
func StoreObjectInR2(filename string, path string) error {
	// Initialize the R2 client using the AWS SDK for Go v2. The R2 endpoint is constructed using the account ID and the base URL for R2.
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")

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

	// Read the JSON data from the specified file.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open response.json file: %v", err)
	}
	defer file.Close()

	// Generate a unique object key for the R2 bucket using the SHA-256 hash of the JSON data.
	objectKey := GenObjectKeyByUUID(path)

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
