package client

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RustFSClient struct {
	client *s3.Client
}

// Currently RustFSClient has solid s3.Client as its field,
// when logic get complex and needs for mocking arise, replace it with interface

func NewRustFSClient() (*RustFSClient, error) {
	endpoint := os.Getenv("RUSTFS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:9001"
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

// / Refer:  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#NewPresignClient
// Direct listing, generating bucket and presigned URL is the role for RustFSClient(control)
// This Presigned url should be passed to Core client, and core should use this url to upload/download file.
// This Presigend related functions are yet not fully tested, so there might be some issues. Please refer to AWS SDK for Go v2 documentation
//
//	for more details and examples on how to use the PresignClient.
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

// newCLI := &RustFSClient{client: client}
// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
// defer cancel()
// err = newCLI.CreateBucket(ctx, "adsfasdfds")
// err = newCLI.CreateBucket(ctx, "asdfaifdsfn")
// err = newCLI.CreateBucket(ctx, "asdf")
// err = newCLI.CreateBucket(ctx, "zxvc")
// err = newCLI.CreateBucket(ctx, "asdfs")
// bucks, err := newCLI.ListBuckets(ctx)
// fmt.Println(bucks)
// if err != nil {
// 	fmt.Println(err)
// }
