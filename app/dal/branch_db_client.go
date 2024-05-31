package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_branch_db_client.go "github.com/payloadops/plato/app/dal" BranchManager

// BranchManager defines the operations available for managing branches.
type BranchManager interface {
	CreateBranch(ctx context.Context, orgID, promptID string, branch *Branch) error
	GetBranch(ctx context.Context, orgID, promptID, branchName string) (*Branch, error)
	DeleteBranch(ctx context.Context, orgID, promptID, branchName string) error
	ListBranchesByPrompt(ctx context.Context, orgID, promptID string) ([]Branch, error)
}

// Ensure BranchDBClient implements the BranchManager interface
var _ BranchManager = &BranchDBClient{}

// Branch represents a branch in the system.
type Branch struct {
	Name      string `json:"name"`
	Deleted   bool   `json:"deleted"`
	CreatedAt string `json:"createdAt"`
}

// BranchDBClient is a client for interacting with DynamoDB for branch-related operations.
type BranchDBClient struct {
	service DynamoDBAPI
}

// NewBranchDBClient creates a new BranchDBClient.
func NewBranchDBClient(service DynamoDBAPI) *BranchDBClient {
	return &BranchDBClient{
		service: service,
	}
}

// createBranchCompositeKeys generates the partition key (pk) and sort key (sk) for a branch.
func createBranchCompositeKeys(orgID, promptID, branchName string) (string, string) {
	return "Org#" + orgID + "Prompt#" + promptID, "Branch#" + branchName
}

// CreateBranch creates a new branch in the DynamoDB table.
func (d *BranchDBClient) CreateBranch(ctx context.Context, orgID, promptID string, branch *Branch) error {
	pk, sk := createBranchCompositeKeys(orgID, promptID, branch.Name)

	now := time.Now().UTC().Format(time.RFC3339)
	branch.CreatedAt = now

	av, err := attributevalue.MarshalMap(branch)
	if err != nil {
		return fmt.Errorf("failed to marshal branch: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Branches"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetBranch retrieves a branch by orgID, prompt ID, and branch ID from the DynamoDB table.
func (d *BranchDBClient) GetBranch(ctx context.Context, orgID, promptID, branchName string) (*Branch, error) {
	pk, sk := createBranchCompositeKeys(orgID, promptID, branchName)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Branches"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
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

	if branch.Deleted {
		return nil, nil
	}

	return &branch, nil
}

// DeleteBranch marks a branch as deleted by org ID, prompt ID and branch ID in the DynamoDB table.
func (d *BranchDBClient) DeleteBranch(ctx context.Context, orgID, promptID, branchName string) error {
	pk, sk := createBranchCompositeKeys(orgID, promptID, branchName)
	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Branches"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListBranchesByPrompt retrieves all branches belonging to a specific prompt from the DynamoDB table.
func (d *BranchDBClient) ListBranchesByPrompt(ctx context.Context, orgID, promptID string) ([]Branch, error) {
	pk, _ := createBranchCompositeKeys(orgID, promptID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Branches"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.service.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var branches []Branch
	err = attributevalue.UnmarshalListOfMaps(result.Items, &branches)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []Branch{}
	for _, branch := range branches {
		if branch.Deleted {
			continue
		}
		results = append(results, branch)
	}

	return branches, nil
}
