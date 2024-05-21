package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/plato/api/utils"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_branch_db_client.go "github.com/payloadops/plato/api/dal" BranchManager

// BranchManager defines the operations available for managing branches.
type BranchManager interface {
	CreateBranch(ctx context.Context, branch *Branch) error
	GetBranch(ctx context.Context, promptID, branchID string) (*Branch, error)
	DeleteBranch(ctx context.Context, promptID, branchID string) error
	ListBranchesByPrompt(ctx context.Context, promptID string) ([]Branch, error)
}

// Ensure BranchDBClient implements the BranchManager interface
var _ BranchManager = &BranchDBClient{}

// Branch represents a branch in the system.
type Branch struct {
	PromptID  string `json:"promptId"`
	BranchID  string `json:"branchId"`
	Deleted   bool   `json:"deleted"`
	CreatedAt string `json:"createdAt"`
}

// BranchDBClient is a client for interacting with DynamoDB for branch-related operations.
type BranchDBClient struct {
	service *dynamodb.Client
}

// NewBranchDBClient creates a new BranchDBClient.
func NewBranchDBClient(service *dynamodb.Client) *BranchDBClient {
	return &BranchDBClient{
		service: service,
	}
}

// createBranchCompositeKeys generates the partition key (PK) and sort key (SK) for a branch.
func createBranchCompositeKeys(promptID, branchID string) (string, string) {
	return "Prompt#" + promptID, "Branch#" + branchID
}

// CreateBranch creates a new branch in the DynamoDB table.
func (d *BranchDBClient) CreateBranch(ctx context.Context, branch *Branch) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	branch.BranchID = ksuid
	pk, sk := createBranchCompositeKeys(branch.PromptID, branch.BranchID)

	now := time.Now().UTC().Format(time.RFC3339)
	branch.CreatedAt = now

	av, err := attributevalue.MarshalMap(branch)
	if err != nil {
		return fmt.Errorf("failed to marshal branch: %v", err)
	}

	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
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

// GetBranch retrieves a branch by prompt ID and branch ID from the DynamoDB table.
func (d *BranchDBClient) GetBranch(ctx context.Context, promptID, branchID string) (*Branch, error) {
	pk, sk := createBranchCompositeKeys(promptID, branchID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Branches"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
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

// DeleteBranch marks a branch as deleted by prompt ID and branch ID in the DynamoDB table.
func (d *BranchDBClient) DeleteBranch(ctx context.Context, promptID, branchID string) error {
	pk, sk := createBranchCompositeKeys(promptID, branchID)
	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Branches"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListBranchesByPrompt retrieves all branches belonging to a specific prompt from the DynamoDB table.
func (d *BranchDBClient) ListBranchesByPrompt(ctx context.Context, promptID string) ([]Branch, error) {
	pk, _ := createBranchCompositeKeys(promptID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Branches"),
		KeyConditionExpression: aws.String("PK = :pk"),
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

	return branches, nil
}
