package dbdal

import (
	"context"
	"fmt"
	awsclient "plato/app/pkg/client/aws"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Team represents the structure for a team record in the database.
type Team struct {
	Id       string   `dynamodbav:"teamId" json:"id"`
	Name     string   `dynamodbav:"name" json:"name"`
	OrgId    string   `dynamodbav:"orgId" json:"org_id"`
	Users    []string `dynamodbav:"users" json:"users"`
	Projects []string `dynamodbav:"projects" json:"projects"`
}

// GetTeamById retrieves a team by its Id.
func GetTeamById(ctx context.Context, id string) (*Team, error) {
	team := &Team{
		Id: id,
	}
	pk := fmt.Sprintf("ORG#%s", id)
	sk := fmt.Sprintf("TEAM#%s", team.Id)

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &ORGS_TEAMS_USERS_TABLE,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &team)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling team: %w", err)
	}

	return team, nil
}

func AddTeam(
	ctx context.Context,
	name string,
	orgId string,
) (*Team, error) {
	team := &Team{
		Id:    util.GenUUIDString(),
		OrgId: orgId,
		Name:  name,
	}

	pk := fmt.Sprintf("ORG#%s", orgId)
	sk := fmt.Sprintf("TEAM#%s", team.Id)

	item, err := attributevalue.MarshalMap(team)
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

	return team, nil
}
