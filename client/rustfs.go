package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type RustFSClient struct {
	client *s3.Client
}

var (
	rustFSInstance *RustFSClient
	rustFSOnce     sync.Once
	rustFSErr      error
)

func GetRustFSClient() (*RustFSClient, error) {
	rustFSOnce.Do(func() {
		rustFSInstance, rustFSErr = NewRustFSClient()
	})
	return rustFSInstance, rustFSErr
}

// Currently RustFSClient has solid s3.Client as its field,
// when logic get complex and needs for mocking arise, replace it with interface

func NewRustFSClient() (*RustFSClient, error) {
	endpoint := os.Getenv("RUSTFS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:9000"
	}
	accessKey := os.Getenv("RUSTFS_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("RUSTFS_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	//TODO: replace unstable accesskey override when testing with proper method

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &RustFSClient{client: client}, nil
}

func (c *RustFSClient) ListBuckets(ctx context.Context) ([]string, error) {
	out, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(out.Buckets))
	for i, b := range out.Buckets {
		names[i] = aws.ToString(b.Name)
	}
	return names, nil
}

func (c *RustFSClient) CreateBucket(ctx context.Context, bucket string) error {
	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

// Refer: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#NewPresignClient
// Direct listing, generating bucket and presigned URL is the role for RustFSClient(control)
// This Presigned url should be passed to Core client, and core should use this url to upload/download file.
func (c *RustFSClient) PresignPutObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)
	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (c *RustFSClient) PresignGetObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// HeadObject returns (true, nil) if the object exists, (false, nil) if only the object is missing,
// and (false, err) for all other errors including NoSuchBucket.
func (c *RustFSClient) HeadObject(ctx context.Context, bucket, key string) (bool, error) {
	_, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("HeadObject %s/%s: %w", bucket, key, err)
	}
	return true, nil
}

// ListObjects returns all object keys in the bucket matching the prefix.
// Returns an empty slice when no objects match. Returns an error if the bucket does not exist.
func (c *RustFSClient) ListObjects(ctx context.Context, bucket, prefix string) ([]string, error) {
	out, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("ListObjects %s: %w", bucket, err)
	}
	keys := make([]string, len(out.Contents))
	for i, obj := range out.Contents {
		keys[i] = aws.ToString(obj.Key)
	}
	return keys, nil
}

func (c *RustFSClient) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("DeleteObject %s/%s: %w", bucket, key, err)
	}
	return nil
}
