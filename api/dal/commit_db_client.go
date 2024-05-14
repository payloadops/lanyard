package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// CommitDBClient is a client for interacting with DynamoDB for commit-related operations.
type CommitDBClient struct {
	service *dynamodb.Client
}

// NewCommitDBClient creates a new CommitDBClient with the AWS configuration.
func NewCommitDBClient() (*CommitDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &CommitDBClient{
		service: svc,
	}, nil
}

// CreateCommit creates a new commit in the DynamoDB table.
func (d *CommitDBClient) CreateCommit(ctx context.Context, commit Commit) error {
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

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
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

	result, err := d.service.GetItem(ctx, input)
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

	return &commit, nil
}

// ListCommits retrieves all commits from the DynamoDB table.
func (d *CommitDBClient) ListCommits(ctx context.Context) ([]Commit, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Commits"),
	}

	result, err := d.service.Scan(ctx, input)
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

	result, err := d.service.Scan(ctx, input)
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
