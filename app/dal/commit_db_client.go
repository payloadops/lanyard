package dal

/*
import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/api/cache"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_commit_db_client.go "github.com/payloadops/plato/api/dal" CommitManager

// CommitManager defines the operations available for managing commits.
type CommitManager interface {
	CreateCommit(ctx context.Context, commit Commit) error
	GetCommit(ctx context.Context, id string) (*Commit, error)
	ListCommits(ctx context.Context) ([]Commit, error)
	ListCommitsByBranch(ctx context.Context, branchID string) ([]Commit, error)
}

// Ensure CommitDBClient implements the CommitManager interface
var _ CommitManager = &CommitDBClient{}

// Commit represents a commit in the system.
type Commit struct {
	ID        string `json:"id"`
	BranchID  string `json:"branchId"`
	UserID    string `json:"userId"`
	Message   string `json:"message"`
	Content   string `json:"-"`
	Checksum  string `json:"checksum"`
	VersionID string `json:"versionId"`
	CreatedAt string `json:"createdAt"`
}

// CommitDBClient is a client for interacting with DynamoDB for commit-related operations.
type CommitDBClient struct {
	dynamoDb *dynamodb.Client
	s3       *s3.Client
	cache    cache.Cache
}

// NewCommitDBClient creates a new CommitDBClient with the AWS configuration.
func NewCommitDBClient(dynamoDb *dynamodb.Client, s3 *s3.Client, cache cache.Cache) *CommitDBClient {
	return &CommitDBClient{
		dynamoDb: dynamoDb,
		s3:       s3,
		cache:    cache,
	}
}

// / CreateCommit creates a new commit in the DynamoDB table.
func (d *CommitDBClient) CreateCommit(ctx context.Context, commit Commit) error {
	// Use the BranchID as the key, ensuring all commits on the same branch refer to the same object
	key := "commits/" + commit.BranchID + ".txt"
	obj, err := d.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String(key),
		Body:   strings.NewReader(commit.Content),
	})
	if err != nil {
		return fmt.Errorf("failed to upload commit content to S3: %v", err)
	}

	// Set version to commit so that we can retrieve it later, as well as the hash of the commit
	commit.VersionID = aws.ToString(obj.VersionId)
	commit.Checksum = aws.ToString(obj.ChecksumSHA256)

	// Remove content from the commit before saving to DynamoDB
	commit.Content = ""
	now := time.Now().UTC().Format(time.RFC3339)
	commit.CreatedAt = now

	av, err := attributevalue.MarshalMap(commit)
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Commits"),
		Item:      av,
	}

	_, err = d.dynamoDb.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	// Cache the latest commit content
	cacheKey := fmt.Sprintf("commit:%s", commit.BranchID)
	if err := d.cache.Set(ctx, cacheKey, commit.Content, 10*time.Minute); err != nil {
		return fmt.Errorf("failed to cache commit content: %v", err)
	}

	return nil
}

// GetCommit retrieves a commit by ID from the DynamoDB table.
func (d *CommitDBClient) GetCommit(ctx context.Context, id string) (*Commit, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Commits"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.dynamoDb.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var commit Commit
	err = attributevalue.UnmarshalMap(result.Item, &commit)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	// Try to get the content from the cache
	cacheKey := fmt.Sprintf("commit:%s", commit.BranchID)
	if err := d.cache.Get(ctx, cacheKey, &commit.Content); err == nil {
		return &commit, nil
	}

	// Retrieve the content from S3 using the BranchID and VersionID if not in cache
	key := "commits/" + commit.BranchID + ".txt"
	obj, err := d.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket:    aws.String("your-bucket-name"),
		Key:       aws.String(key),
		VersionId: aws.String(commit.VersionID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit content from S3: %v", err)
	}
	defer obj.Body.Close()

	content, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit content: %v", err)
	}

	commit.Content = string(content)

	// Cache the retrieved content
	if err := d.cache.Set(ctx, cacheKey, commit.Content, 10*time.Minute); err != nil {
		return nil, fmt.Errorf("failed to cache commit content: %v", err)
	}

	return &commit, nil
}

// ListCommits retrieves all commits from the DynamoDB table.
func (d *CommitDBClient) ListCommits(ctx context.Context) ([]Commit, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Commits"),
	}

	result, err := d.dynamoDb.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var commits []Commit
	err = attributevalue.UnmarshalListOfMaps(result.Items, &commits)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return commits, nil
}

// ListCommitsByBranch retrieves all commits belonging to a specific branch from the DynamoDB table.
func (d *CommitDBClient) ListCommitsByBranch(ctx context.Context, branchID string) ([]Commit, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Commits"),
		FilterExpression: aws.String("branchId = :branchId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":branchId": &types.AttributeValueMemberS{Value: branchID},
		},
	}

	result, err := d.dynamoDb.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var commits []Commit
	err = attributevalue.UnmarshalListOfMaps(result.Items, &commits)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return commits, nil
}
*/
