package dal

import (
	"context"
	"fmt"

	"github.com/payloadops/lanyard/app/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_Org_db_client.go "github.com/payloadops/lanyard/app/dal" OrgManager

// OrgManager defines the operations available for managing Orgs.
type OrgManager interface {
	CreateOrg(ctx context.Context, orgID, serviceID string, Org *Org) error
	GetOrg(ctx context.Context, orgID, serviceID string, name string) (*Org, error)
	UpdateOrg(ctx context.Context, orgID, serviceID string, Org *Org) error
	DeleteOrg(ctx context.Context, orgID, serviceID string, name string) error
}

// Ensure OrgDBClient implements the OrgManager interface
var _ OrgManager = &OrgDBClient{}

type Org struct {
	OrgID           string `json:"orgId"`
	Name            string `json:"name"`
	StripeAccountId string `json:"stripeAccountId"`
	Domain          string `json:"domain"`
}

// OrgDBClient is a client for interacting with DynamoDB for Org-related operations.
type OrgDBClient struct {
	Org DynamoDBAPI
}

// NewOrgDBClient creates a new OrgDBClient.
func NewOrgDBClient(Org DynamoDBAPI) *OrgDBClient {
	return &OrgDBClient{
		Org: Org,
	}
}

// createOrgCompositeKeys generates the partition key (pk) and sort key (sk) for a Org.
func createOrgCompositeKeys(orgID string) (string, string) {
	return "Org#" + orgID, "Org#" + orgID
}

// CreateOrg creates a new Org in the DynamoDB table.
func (d *OrgDBClient) CreateOrg(ctx context.Context, orgID, serviceID string, Org *Org) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	Org.OrgID = ksuid
	pk, sk := createOrgCompositeKeys(orgID)

	av, err := attributevalue.MarshalMap(Org)
	if err != nil {
		return fmt.Errorf("failed to marshal Org: %v", err)
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

	_, err = d.Org.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetOrg retrieves a Org by organization ID and Org ID from the DynamoDB table.
func (d *OrgDBClient) GetOrg(ctx context.Context, orgID, serviceID, name string) (*Org, error) {
	pk, sk := createOrgCompositeKeys(orgID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Services"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.Org.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var Org Org
	err = attributevalue.UnmarshalMap(result.Item, &Org)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &Org, nil
}

// UpdateOrg updates the name, description, and updatedAt fields of an existing Org in the DynamoDB table.
func (d *OrgDBClient) UpdateOrg(ctx context.Context, orgID, serviceID string, Org *Org) error {
	pk, sk := createOrgCompositeKeys(orgID)

	updateExpr := "SET #name = :name, #domain = :domain, #stripeAccountId = :stripeAccountId"
	exprAttrNames := map[string]string{
		"#name":            "name",
		"#domain":          "domain",
		"#stripeAccountId": "stripeAccountId",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":name":            &types.AttributeValueMemberS{Value: Org.Name},
		":domain":          &types.AttributeValueMemberS{Value: Org.Domain},
		":stripeAccountId": &types.AttributeValueMemberS{Value: Org.StripeAccountId},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("Services"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	}

	_, err := d.Org.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteOrg marks a Org as deleted by organization ID and Org ID in the DynamoDB table.
func (d *OrgDBClient) DeleteOrg(ctx context.Context, orgID, serviceID, name string) error {
	pk, sk := createOrgCompositeKeys(orgID)

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

	_, err := d.Org.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListOrgsByOrganization retrieves all Orgs for a specific organization from the DynamoDB table.
func (d *OrgDBClient) ListOrgs(ctx context.Context, orgID, serviceID string) ([]Org, error) {
	pk, _ := createOrgCompositeKeys(orgID)
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Services"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.Org.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var Orgs []Org
	err = attributevalue.UnmarshalListOfMaps(result.Items, &Orgs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return Orgs, nil
}
