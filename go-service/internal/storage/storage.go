package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Store preserves raw scraper output before anything touches it.
//
// This is the load-bearing piece of "raw observations are preserved; everything
// else is derived." If the normalization logic is wrong on Friday, you reprocess
// from here instead of burning credits re-scraping.
type Store interface {
	Put(ctx context.Context, key string, body []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

// ---------------------------------------------------------------
// Cloudflare R2
// ---------------------------------------------------------------

type R2 struct {
	client *s3.Client
	bucket string
}

type R2Config struct {
	AccountID string // Cloudflare account ID, from the R2 overview page
	AccessKey string // R2 -> Manage API tokens -> Create token
	SecretKey string
	Bucket    string
}

func NewR2(ctx context.Context, c R2Config) (*R2, error) {
	// R2 is S3-compatible but has no regions. "auto" is required; any real
	// AWS region name will be rejected on signing.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("r2 config: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		// R2 does not support virtual-hosted-style addressing on the
		// cloudflarestorage.com endpoint. Omit this and every call 404s.
		o.UsePathStyle = true
	})

	return &R2{client: client, bucket: c.Bucket}, nil
}

func (r *R2) Put(ctx context.Context, key string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("r2 put %s: %w", key, err)
	}
	return nil
}

func (r *R2) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	out, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("r2 get %s: %w", key, err)
	}
	defer out.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(out.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Key is the layout from the architecture doc: raw/<retailer>/<date>/<snapshot>.json
// Partitioning by date is what makes "reprocess last Tuesday" a one-liner.
func Key(retailerSlug, snapshotID string, at time.Time) string {
	if snapshotID == "" {
		snapshotID = fmt.Sprintf("run-%d", at.UnixNano())
	}
	return fmt.Sprintf("raw/%s/%s/%s.json",
		retailerSlug, at.UTC().Format("2006-01-02"), snapshotID)
}

// ---------------------------------------------------------------
// Local fallback -- use this until R2 credentials exist so the
// pipeline is never blocked on someone else's dashboard access.
// ---------------------------------------------------------------

type Local struct{ Dir string }

func (l *Local) Put(_ context.Context, key string, body []byte) error {
	return writeLocal(l.Dir, key, body)
}

func (l *Local) Get(_ context.Context, key string) ([]byte, error) {
	return readLocal(l.Dir, key)
}
