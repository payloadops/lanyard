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

//go:generate mockgen -package=mocks -destination=mocks/mock_branch_db_client.go "github.com/payloadops/plato/api/dal" BranchManager

// BranchManager defines the operations available for managing branches.
type BranchManager interface {
	CreateBranch(ctx context.Context, branch Branch) error
	GetBranch(ctx context.Context, id string) (*Branch, error)
	DeleteBranch(ctx context.Context, id string) error
	ListBranches(ctx context.Context) ([]Branch, error)
	ListBranchesByPrompt(ctx context.Context, promptID string) ([]Branch, error)
}

// Ensure BranchDBClient implements the BranchManager interface
var _ BranchManager = &BranchDBClient{}

// Branch represents a branch in the system.
type Branch struct {
	ID        string `json:"id"`
	PromptID  string `json:"promptId"`
	CreatedAt string `json:"createdAt"`
}

// BranchDBClient is a client for interacting with DynamoDB for branch-related operations.
type BranchDBClient struct {
	service *dynamodb.Client
}

// NewBranchDBClient creates a new BranchDBClient with the AWS configuration.
func NewBranchDBClient() (*BranchDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &BranchDBClient{
		service: svc,
	}, nil
}

// CreateBranch creates a new branch in the DynamoDB table.
func (d *BranchDBClient) CreateBranch(ctx context.Context, branch Branch) error {
	now := time.Now().UTC().Format(time.RFC3339)
	branch.CreatedAt = now

	av, err := attributevalue.MarshalMap(branch)
	if err != nil {
		return fmt.Errorf("failed to marshal branch: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Branches"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetBranch retrieves a branch by ID from the DynamoDB table.
func (d *BranchDBClient) GetBranch(ctx context.Context, id string) (*Branch, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Branches"),
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

	var branch Branch
	err = attributevalue.UnmarshalMap(result.Item, &branch)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &branch, nil
}

// DeleteBranch deletes a branch by ID from the DynamoDB table.
func (d *BranchDBClient) DeleteBranch(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Branches"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := d.service.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %v", err)
	}

	return nil
}

// ListBranches retrieves all branches from the DynamoDB table.
func (d *BranchDBClient) ListBranches(ctx context.Context) ([]Branch, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Branches"),
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var branches []Branch
	err = attributevalue.UnmarshalListOfMaps(result.Items, &branches)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return branches, nil
}

// ListBranchesByPrompt retrieves all branches belonging to a specific prompt from the DynamoDB table.
func (d *BranchDBClient) ListBranchesByPrompt(ctx context.Context, promptID string) ([]Branch, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Branches"),
		FilterExpression: aws.String("promptId = :promptId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":promptId": &types.AttributeValueMemberS{Value: promptID},
		},
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var branches []Branch
	err = attributevalue.UnmarshalListOfMaps(result.Items, &branches)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return branches, nil
}
