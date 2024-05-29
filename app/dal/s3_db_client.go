package dal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_s3_db_client.go "github.com/payloadops/plato/app/dal" S3API

// Ensure DynamoDBClient implements the DynamoDBAPI interface
var _ S3API = &s3.Client{}

// S3API is the interface for the AWS S3 client.
type S3API interface {
	GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}
