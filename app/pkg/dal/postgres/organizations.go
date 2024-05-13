package dbdal

import (
	"context"
	"fmt"
	awsclient "plato/app/pkg/client/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ORGS_TEAMS_USERS_TABLE = "OrgsTeamsUsers"

// Organization represents the structure of an organization record in the database.
type Organization struct {
	Id   string `dynamodbav:"orgId" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
}

// GetOrganizationById retrieves an organization by its Id using Bun.
func GetOrganizationById(ctx context.Context, id string) (*Organization, error) {
	organization := &Organization{}
	pk := fmt.Sprintf("ORG#%s", organization.Id)
	sk := fmt.Sprintf("ORG#%s", organization.Id)

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &ORGS_TEAMS_USERS_TABLE,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &organization)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return organization, nil
}

// AddOrganization adds a new organization to the database with the provided name.
func AddOrganization(ctx context.Context, name string) (*Organization, error) {
	organization := &Organization{
		Name: name,
	}
	pk := fmt.Sprintf("ORG#%s", organization.Id)
	sk := fmt.Sprintf("ORG#%s", organization.Id)

	item, err := attributevalue.MarshalMap(organization)
	if err != nil {
		return nil, fmt.Errorf("error marshaling prompt: %w", err)
	}

	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}

	_, err = awsclient.GetDynamoClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &ORGS_TEAMS_USERS_TABLE,
		Item:      item,
	})

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling team: %w", err)
	}

	return organization, nil
}
