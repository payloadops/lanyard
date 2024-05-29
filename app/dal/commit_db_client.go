package dal

/*
import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/app/cache"
	"github.com/payloadops/plato/app/utils"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_commit_db_client.go "github.com/payloadops/plato/app/dal" CommitManager

// CommitManager defines the operations available for managing commits.
type CommitManager interface {
	CreateCommit(ctx context.Context, commit *Commit) error
	GetCommit(ctx context.Context, branchID, commitID string) (*Commit, error)
	ListCommitsByBranch(ctx context.Context, branchID string) ([]Commit, error)
}

// Ensure CommitDBClient implements the CommitManager interface
var _ CommitManager = &CommitDBClient{}

// Commit represents a commit in the system.
type Commit struct {
	BranchID  string `json:"branchId"`
	CommitID  string `json:"commitId"`
	UserID    string `json:"userId"`
	Message   string `json:"message"`
	Content   string `json:"-"`
	Checksum  string `json:"checksum"`
	VersionID string `json:"versionId"`
	CreatedAt string `json:"createdAt"`
}

// CommitDBClient is a client for interacting with DynamoDB for commit-related operations.
type CommitDBClient struct {
	dynamoDb DynamoDBAPI
	s3       S3API
	cache    cache.Cache
}

// NewCommitDBClient creates a new CommitDBClient with the AWS configuration.
func NewCommitDBClient(dynamoDb DynamoDBAPI, s3 S3API, cache cache.Cache) *CommitDBClient {
	return &CommitDBClient{
		dynamoDb: dynamoDb,
		s3:       s3,
		cache:    cache,
	}
}

// createCommitCompositeKeys generates the partition key (PK) and sort key (SK) for a commit.
func createCommitCompositeKeys(branchID, commitID string) (string, string) {
	return "Branch#" + branchID, "Commit#" + commitID
}

// CreateCommit creates a new commit in the DynamoDB table.
func (d *CommitDBClient) CreateCommit(ctx context.Context, commit *Commit) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	commit.CommitID = ksuid
	pk, sk := createCommitCompositeKeys(commit.BranchID, commit.CommitID)

	now := time.Now().UTC().Format(time.RFC3339)
	commit.CreatedAt = now

	av, err := attributevalue.MarshalMap(commit)
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %v", err)
	}

	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Commits"),
		Item:      item,
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

// GetCommit retrieves a commit by branch ID and commit ID from the DynamoDB table.
func (d *CommitDBClient) GetCommit(ctx context.Context, branchID, commitID string) (*Commit, error) {
	pk, sk := createCommitCompositeKeys(branchID, commitID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Commits"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
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
	cacheKey := fmt.Sprintf("commit:%s", commit.CommitID)
	content, err := d.cache.Get(ctx, cacheKey)
	if err == nil {
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

// ListCommitsByBranch retrieves all commits belonging to a specific branch from the DynamoDB table.
func (d *CommitDBClient) ListCommitsByBranch(ctx context.Context, branchID string) ([]Commit, error) {
	pk := "Branch#" + branchID

	input := &dynamodb.QueryInput{
		TableName:              aws.String("Commits"),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.dynamoDb.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var commits []Commit
	err = attributevalue.UnmarshalListOfMaps(result.Items, &commits)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return commits, nil
}
*/
