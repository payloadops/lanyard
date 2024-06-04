package dal

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/app/utils"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/app/cache"
	"github.com/payloadops/plato/app/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CommitTTL represents the duration that a prompt can be stored in the cache
const CommitTTL = 5 * time.Minute

//go:generate mockgen -package=mocks -destination=mocks/mock_commit_db_client.go "github.com/payloadops/plato/app/dal" CommitManager

// CommitManager defines the operations available for managing commits.
type CommitManager interface {
	CreateCommit(ctx context.Context, orgID, projectID, promptID, branchName string, commit *Commit) error
	GetCommit(ctx context.Context, orgID, projectID, promptID, branchName, commitID string) (*Commit, error)
	ListCommitsByBranch(ctx context.Context, orgID, projectID, promptID, branchName string) ([]Commit, error)
}

// Ensure CommitDBClient implements the CommitManager interface
var _ CommitManager = &CommitDBClient{}

// Commit represents a commit in the system.
type Commit struct {
	CommitID  string `json:"commitId"`
	UserID    string `json:"userId"`
	VersionID string `json:"versionId"`
	Message   string `json:"message"`
	Content   string `json:"-"`
	CreatedAt string `json:"createdAt"`
}

// CommitDBClient is a client for interacting with DynamoDB for commit-related operations.
type CommitDBClient struct {
	dynamoDb   DynamoDBAPI
	s3         S3API
	cache      cache.Cache
	bucketName string
}

// NewCommitDBClient creates a new CommitDBClient with the AWS configuration.
func NewCommitDBClient(dynamoDb DynamoDBAPI, s3 S3API, cache cache.Cache, config *config.Config) *CommitDBClient {
	return &CommitDBClient{
		dynamoDb:   dynamoDb,
		s3:         s3,
		cache:      cache,
		bucketName: config.PromptBucket,
	}
}

// createCommitCompositeKeys generates the partition key (pk) and sort key (sk) for a commit.
func createCommitCompositeKeys(orgID, promptID, branchName, commitID string) (string, string) {
	return fmt.Sprintf("Org#%sPrompt#%s#Branch#%s", orgID, promptID, branchName), fmt.Sprintf("Commit#%s", commitID)
}

// CreateCommit creates a new commit in the DynamoDB table.
func (d *CommitDBClient) CreateCommit(ctx context.Context, orgID, projectID, promptID, branchName string, commit *Commit) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	commit.CommitID = ksuid
	now := time.Now().UTC().Format(time.RFC3339)
	commit.CreatedAt = now
	pk, sk := createCommitCompositeKeys(orgID, promptID, branchName, commit.CommitID)

	// Upload the commit content to S3 and get the version ID
	s3Key := fmt.Sprintf("prompts/%s/%s/%s/%s.txt", orgID, projectID, promptID, branchName)
	putObjectOutput, err := d.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(s3Key),
		Body:   strings.NewReader(commit.Content),
	})
	if err != nil {
		return fmt.Errorf("failed to upload commit content to S3: %v", err)
	}

	commit.VersionID = aws.ToString(putObjectOutput.VersionId)
	content := commit.Content
	commit.Content = "" // Clear the content before saving to DynamoDB

	av, err := attributevalue.MarshalMap(commit)
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
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
	commit.Content = content
	cacheKey := fmt.Sprintf("commit:%s:%s:%s", promptID, branchName, commit.CommitID)
	if err := d.cache.Set(ctx, cacheKey, commit.Content, CommitTTL); err != nil {
		return fmt.Errorf("failed to cache commit content: %v", err)
	}

	return nil
}

// GetCommit retrieves a commit by prompt ID, branch ID, and commit ID from the DynamoDB table.
func (d *CommitDBClient) GetCommit(ctx context.Context, orgID, projectID, promptID, branchName, commitID string) (*Commit, error) {
	pk, sk := createCommitCompositeKeys(orgID, promptID, branchName, commitID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Commits"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
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
	cacheKey := fmt.Sprintf("commit:%s:%s:%s", promptID, branchName, commit.CommitID)
	if content, err := d.cache.Get(ctx, cacheKey, CommitTTL); content != "" && err == nil {
		commit.Content = content
		return &commit, nil
	}

	// Retrieve the content from S3 using the branchName and VersionID if not in cache
	s3Key := fmt.Sprintf("prompts/%s/%s/%s/%s.txt", orgID, projectID, promptID, branchName)
	obj, err := d.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket:    aws.String(d.bucketName),
		Key:       aws.String(s3Key),
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
	if err := d.cache.Set(ctx, cacheKey, commit.Content, CommitTTL); err != nil {
		return nil, fmt.Errorf("failed to cache commit content: %v", err)
	}

	return &commit, nil
}

// ListCommitsByBranch retrieves all commits belonging to a specific branch from the DynamoDB table.
func (d *CommitDBClient) ListCommitsByBranch(ctx context.Context, orgID, projectID, promptID, branchName string) ([]Commit, error) {
	pk, _ := createCommitCompositeKeys(orgID, promptID, branchName, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Commits"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
		ScanIndexForward: aws.Bool(false), // Retrieve results in descending order
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
