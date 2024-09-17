package dal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/payloadops/lanyard/app/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_Tier_db_client.go "github.com/payloadops/lanyard/app/dal" TierManager

// TierManager defines the operations available for managing Tiers.
type TierManager interface {
	CreateTier(ctx context.Context, orgID, serviceID string, Tier *Tier) error
	GetTier(ctx context.Context, orgID, serviceID string, name string) (*Tier, error)
	UpdateTier(ctx context.Context, orgID, serviceID string, Tier *Tier) error
	DeleteTier(ctx context.Context, orgID, serviceID string, name string) error
	ListTiers(ctx context.Context, orgID, serviceID string) ([]Tier, error)
}

// Ensure TierDBClient implements the TierManager interface
var _ TierManager = &TierDBClient{}

// PricingTier represents an Tier's basic billing info in the system.
type Tier struct {
	TierID              string  `json:"tierId"`
	Name                string  `json:"tierId"`
	DefaultRequestLimit int     `json:"defaultRequestLimit"`
	Interval            int     `json:"interval"`
	OveragePrice        float32 `json:"overagePrice"`
}

// TierDBClient is a client for interacting with DynamoDB for Tier-related operations.
type TierDBClient struct {
	Tier DynamoDBAPI
}

// NewTierDBClient creates a new TierDBClient.
func NewTierDBClient(Tier DynamoDBAPI) *TierDBClient {
	return &TierDBClient{
		Tier: Tier,
	}
}

// createTierCompositeKeys generates the partition key (pk) and sort key (sk) for a Tier.
func createTierCompositeKeys(orgID, serviceID, TierName string) (string, string) {
	return "Org#" + orgID + "Service#" + serviceID + "Tier", "Tier#" + TierName
}

// CreateTier creates a new Tier in the DynamoDB table.
func (d *TierDBClient) CreateTier(ctx context.Context, orgID, serviceID string, Tier *Tier) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	Tier.TierID = ksuid
	pk, sk := createTierCompositeKeys(orgID, serviceID, Tier.Name)

	av, err := attributevalue.MarshalMap(Tier)
	if err != nil {
		return fmt.Errorf("failed to marshal Tier: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Services"),
		Item:      item,
	}

	_, err = d.Tier.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetTier retrieves a Tier by organization ID and Tier ID from the DynamoDB table.
func (d *TierDBClient) GetTier(ctx context.Context, orgID, serviceID, name string) (*Tier, error) {
	pk, sk := createTierCompositeKeys(orgID, serviceID, name)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Services"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.Tier.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var Tier Tier
	err = attributevalue.UnmarshalMap(result.Item, &Tier)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &Tier, nil
}

// UpdateTier updates the name, description, and updatedAt fields of an existing Tier in the DynamoDB table.
func (d *TierDBClient) UpdateTier(ctx context.Context, orgID, serviceID string, Tier *Tier) error {
	pk, sk := createTierCompositeKeys(orgID, serviceID, Tier.Name)

	updateExpr := "SET #name = :name, #defaultRequestLimit = :defaultRequestLimit, #overagePrice = :overagePrice"
	exprAttrNames := map[string]string{
		"#name":                "name",
		"#defaultRequestLimit": "defaultRequestLimit",
		"#overagePrice":        "overagePrice",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":name":                &types.AttributeValueMemberS{Value: Tier.Name},
		":defaultRequestLimit": &types.AttributeValueMemberN{Value: strconv.Itoa(Tier.DefaultRequestLimit)},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("Services"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	}

	_, err := d.Tier.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteTier marks a Tier as deleted by organization ID and Tier ID in the DynamoDB table.
func (d *TierDBClient) DeleteTier(ctx context.Context, orgID, serviceID, name string) error {
	pk, sk := createTierCompositeKeys(orgID, serviceID, name)

	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Services"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.Tier.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListTiersByOrganization retrieves all Tiers for a specific organization from the DynamoDB table.
func (d *TierDBClient) ListTiers(ctx context.Context, orgID, serviceID string) ([]Tier, error) {
	pk, _ := createTierCompositeKeys(orgID, serviceID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Services"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.Tier.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var Tiers []Tier
	err = attributevalue.UnmarshalListOfMaps(result.Items, &Tiers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return Tiers, nil
}
