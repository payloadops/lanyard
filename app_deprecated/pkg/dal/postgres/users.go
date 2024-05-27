package dbdal

import (
	"context"
	"fmt"
	awsclient "plato/app_deprecated/pkg/client/aws"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Id        string `dynamodbav:"userId" json:"id"`
	FirstName string `dynamodbav:"firstName" json:"first_name"`
	LastName  string `dynamodbav:"lastName" json:"last_name"`
	Email     string `dynamodbav:"email" json:"email"`
}

func GetUserById(ctx context.Context, id string) (*User, error) {
	user := &User{Id: id}
	pk := fmt.Sprintf("ORG#%s", id)
	sk := fmt.Sprintf("USER#%s", user.Id)

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

	err = attributevalue.UnmarshalMap(resp.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return user, nil
}

func AddUser(
	ctx context.Context,
	orgId string,
	firstName string,
	lastName string,
	email string,
) (*User, error) {
	user := &User{
		Id:        util.GenIDString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	pk := fmt.Sprintf("ORG#%s", orgId)
	sk := fmt.Sprintf("USER#%s", user.Id)

	item, err := attributevalue.MarshalMap(user)
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
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return user, nil
}
