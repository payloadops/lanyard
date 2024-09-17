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

//go:generate mockgen -package=mocks -destination=mocks/mock_actor_db_client.go "github.com/payloadops/lanyard/app/dal" ActorManager

// ActorManager defines the operations available for managing actors.
type ActorManager interface {
	CreateActor(ctx context.Context, orgID, serviceID string, actor *Actor) error
	GetActor(ctx context.Context, orgID, serviceID string, externalID string) (*Actor, error)
	UpdateActor(ctx context.Context, orgID, serviceID string, actor *Actor) error
	DeleteActor(ctx context.Context, orgID, serviceID string, externalID string) error
	ListActors(ctx context.Context, orgID, serviceID string) ([]Actor, error)
}

// Ensure ActorDBClient implements the ActorManager interface
var _ ActorManager = &ActorDBClient{}

// BillingInfo represents an actor's basic billing info in the system.
type BillingInfo struct {
	tier             string `json:"tier"`
	tierId           string `json:"tierId"`
	trialExpiry      string `json:"trialExpiry"`
	isTrialActive    bool   `json:"isTrialActive"`
	isTrialElgible   bool   `json:"isTrialActive"`
	stripeCustomerId string `json:"stripeCustomerId"`
}

// Actor represents a actor in the system.
type Actor struct {
	ActorID             string      `json:"actorId"`
	ExternalID          string      `json:"externalId"`
	MonthlyRequestLimit int         `json:"monthlyRequestLimit"`
	Deleted             bool        `json:"deleted"`
	BillingInfo         BillingInfo `json:"billingInfo"`
}

// ActorDBClient is a client for interacting with DynamoDB for actor-related operations.
type ActorDBClient struct {
	actor DynamoDBAPI
}

// NewActorDBClient creates a new ActorDBClient.
func NewActorDBClient(actor DynamoDBAPI) *ActorDBClient {
	return &ActorDBClient{
		actor: actor,
	}
}

// createActorCompositeKeys generates the partition key (pk) and sort key (sk) for a actor.
func createActorCompositeKeys(orgID, serviceID, actorExternalID string) (string, string) {
	return "Org#" + orgID + "Service#" + serviceID + "Actor", "Actor#" + actorExternalID
}

// createActorCompositeKeys generates the partition key (pk) and sort key (sk) for a actor.
func createActorGSI1(orgID, serviceID, actorID string) string {
	return "Org#" + orgID + "Service#" + serviceID + "Actor#" + actorID
}

// CreateActor creates a new actor in the DynamoDB table.
func (d *ActorDBClient) CreateActor(ctx context.Context, orgID, serviceID string, actor *Actor) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	actor.ActorID = ksuid
	pk, sk := createActorCompositeKeys(orgID, serviceID, actor.ExternalID)
	gsi1PK := createActorGSI1(orgID, serviceID, actor.ExternalID)

	av, err := attributevalue.MarshalMap(actor)
	if err != nil {
		return fmt.Errorf("failed to marshal actor: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk":     &types.AttributeValueMemberS{Value: pk},
		"sk":     &types.AttributeValueMemberS{Value: sk},
		"GSI1PK": &types.AttributeValueMemberS{Value: gsi1PK},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Services"),
		Item:      item,
	}

	_, err = d.actor.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetActor retrieves a actor by organization ID and actor ID from the DynamoDB table.
func (d *ActorDBClient) GetActor(ctx context.Context, orgID, serviceID, externalID string) (*Actor, error) {
	pk, sk := createActorCompositeKeys(orgID, serviceID, externalID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Services"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.actor.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var actor Actor
	err = attributevalue.UnmarshalMap(result.Item, &actor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	if actor.Deleted {
		return nil, nil
	}

	return &actor, nil
}

// GetActorByGSI retrieves an actor by using the GSI (Global Secondary Index) from the DynamoDB table.
func (d *ActorDBClient) GetActorByGSI(ctx context.Context, orgID, serviceID, actorID string) (*Actor, error) {
	// Create GSI1 partition key
	gsi1PK := createActorGSI1(orgID, serviceID, actorID)

	// Prepare the query input using the GSI
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Services"),
		IndexName:              aws.String("GSI1"), // Ensure GSI1 is set up in your DynamoDB table
		KeyConditionExpression: aws.String("GSI1PK = :gsi1pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1pk": &types.AttributeValueMemberS{
				Value: gsi1PK,
			},
		},
	}

	// Perform the query
	result, err := d.actor.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query GSI in DynamoDB: %v", err)
	}

	// Check if there are any results
	if len(result.Items) == 0 {
		return nil, nil // Return nil if no actor is found
	}

	// Unmarshal the result into an Actor
	var actor Actor
	err = attributevalue.UnmarshalMap(result.Items[0], &actor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal actor from DynamoDB: %v", err)
	}

	// Check if the actor is marked as deleted
	if actor.Deleted {
		return nil, nil // Return nil if the actor is marked as deleted
	}

	return &actor, nil
}

// UpdateActor updates the name, description, and updatedAt fields of an existing actor in the DynamoDB table.
func (d *ActorDBClient) UpdateActor(ctx context.Context, orgID, serviceID string, actor *Actor) error {
	pk, sk := createActorCompositeKeys(orgID, serviceID, actor.ExternalID)

	updateExpr := "SET #externalId = :externalId, #monthlyRequestLimit = :monthlyRequestLimit"
	exprAttrNames := map[string]string{
		"#externalId":          "ExternalId",
		"#monthlyRequestLimit": "MonthlyRequestLimit",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":externalId":          &types.AttributeValueMemberS{Value: actor.ExternalID},
		":monthlyRequestLimit": &types.AttributeValueMemberN{Value: strconv.Itoa(actor.MonthlyRequestLimit)},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("Services"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	}

	_, err := d.actor.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteActor marks a actor as deleted by organization ID and actor ID in the DynamoDB table.
func (d *ActorDBClient) DeleteActor(ctx context.Context, orgID, serviceID, externalID string) error {
	pk, sk := createActorCompositeKeys(orgID, serviceID, externalID)

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

	_, err := d.actor.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListActorsByOrganization retrieves all actors for a specific organization from the DynamoDB table.
func (d *ActorDBClient) ListActors(ctx context.Context, orgID, serviceID string) ([]Actor, error) {
	pk, _ := createActorCompositeKeys(orgID, serviceID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Services"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.actor.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var actors []Actor
	err = attributevalue.UnmarshalListOfMaps(result.Items, &actors)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []Actor{}
	for _, actor := range actors {
		if actor.Deleted {
			continue
		}
		results = append(results, actor)
	}

	return actors, nil
}
