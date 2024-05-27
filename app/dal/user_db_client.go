package dal

/*
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

//go:generate mockgen -package=mocks -destination=mocks/mock_user_db_client.go "github.com/payloadops/plato/api/dal" UserManager

// UserManager defines the operations available for managing users.
type UserManager interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]User, error)
}

// Ensure UserDBClient implements the UserManager interface
var _ UserManager = &UserDBClient{}

// User represents a user in the system.
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UserDBClient is a client for interacting with DynamoDB for user-related operations.
type UserDBClient struct {
	service *dynamodb.Client
}

// NewUserDBClient creates a new UserDBClient with AWS configuration.
func NewUserDBClient() (*UserDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &UserDBClient{
		service: svc,
	}, nil
}

// CreateUser creates a new user in the DynamoDB table.
func (d *UserDBClient) CreateUser(ctx context.Context, user User) error {
	now := time.Now().UTC().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now

	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetUser retrieves a user by ID from the DynamoDB table.
func (d *UserDBClient) GetUser(ctx context.Context, id string) (*User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
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

	var user User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user in the DynamoDB table.
func (d *UserDBClient) UpdateUser(ctx context.Context, user User) error {
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteUser deletes a user by ID from the DynamoDB table.
func (d *UserDBClient) DeleteUser(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Users"),
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

// ListUsers retrieves all users from the DynamoDB table.
func (d *UserDBClient) ListUsers(ctx context.Context) ([]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var users []User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return users, nil
}
*/
