package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/payloadops/plato/app/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_service_db_client.go "github.com/payloadops/plato/app/dal" ServiceManager

// ServiceManager defines the operations available for managing services.
type ServiceManager interface {
	CreateService(ctx context.Context, orgID string, service *Service) error
	GetService(ctx context.Context, orgID string, serviceID string) (*Service, error)
	UpdateService(ctx context.Context, orgID string, service *Service) error
	DeleteService(ctx context.Context, orgID string, serviceID string) error
	ListServicesByOrganization(ctx context.Context, orgID string) ([]Service, error)
}

// Ensure ServiceDBClient implements the ServiceManager interface
var _ ServiceManager = &ServiceDBClient{}

// Service represents a service in the system.
type Service struct {
	ServiceID   string `json:"serviceId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deleted     bool   `json:"deleted"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ServiceDBClient is a client for interacting with DynamoDB for service-related operations.
type ServiceDBClient struct {
	service DynamoDBAPI
}

// NewServiceDBClient creates a new ServiceDBClient.
func NewServiceDBClient(service DynamoDBAPI) *ServiceDBClient {
	return &ServiceDBClient{
		service: service,
	}
}

// createServiceCompositeKeys generates the partition key (pk) and sort key (sk) for a service.
func createServiceCompositeKeys(orgID, serviceID string) (string, string) {
	return "Org#" + orgID, "Service#" + serviceID
}

// CreateService creates a new service in the DynamoDB table.
func (d *ServiceDBClient) CreateService(ctx context.Context, orgID string, service *Service) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	service.ServiceID = ksuid
	pk, sk := createServiceCompositeKeys(orgID, service.ServiceID)

	now := time.Now().UTC().Format(time.RFC3339)
	service.CreatedAt = now
	service.UpdatedAt = now

	av, err := attributevalue.MarshalMap(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %v", err)
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

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetService retrieves a service by organization ID and service ID from the DynamoDB table.
func (d *ServiceDBClient) GetService(ctx context.Context, orgID, serviceID string) (*Service, error) {
	pk, sk := createServiceCompositeKeys(orgID, serviceID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Services"),
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

	var service Service
	err = attributevalue.UnmarshalMap(result.Item, &service)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	if service.Deleted {
		return nil, nil
	}

	return &service, nil
}

// UpdateService updates the name, description, and updatedAt fields of an existing service in the DynamoDB table.
func (d *ServiceDBClient) UpdateService(ctx context.Context, orgID string, service *Service) error {
	pk, sk := createServiceCompositeKeys(orgID, service.ServiceID)
	service.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	updateExpr := "SET #name = :name, #description = :description, #updatedAt = :updatedAt"
	exprAttrNames := map[string]string{
		"#name":        "Name",
		"#description": "Description",
		"#updatedAt":   "UpdatedAt",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":name":        &types.AttributeValueMemberS{Value: service.Name},
		":description": &types.AttributeValueMemberS{Value: service.Description},
		":updatedAt":   &types.AttributeValueMemberS{Value: service.UpdatedAt},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("Services"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteService marks a service as deleted by organization ID and service ID in the DynamoDB table.
func (d *ServiceDBClient) DeleteService(ctx context.Context, orgID, serviceID string) error {
	pk, sk := createServiceCompositeKeys(orgID, serviceID)
	now := time.Now().UTC().Format(time.RFC3339)

	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
		"UpdatedAt": {
			Value:  &types.AttributeValueMemberS{Value: now},
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

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListServicesByOrganization retrieves all services for a specific organization from the DynamoDB table.
func (d *ServiceDBClient) ListServicesByOrganization(ctx context.Context, orgID string) ([]Service, error) {
	pk, _ := createServiceCompositeKeys(orgID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Services"),
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

	var services []Service
	err = attributevalue.UnmarshalListOfMaps(result.Items, &services)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []Service{}
	for _, service := range services {
		if service.Deleted {
			continue
		}
		results = append(results, service)
	}

	return services, nil
}
