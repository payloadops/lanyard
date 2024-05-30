package dal_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/payloadops/plato/app/config"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cacheMocks "github.com/payloadops/plato/app/cache/mocks"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	mockS3 := mocks.NewMockS3API(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	client := dal.NewCommitDBClient(mockDynamoDB, mockS3, mockCache, &config.Config{
		PromptBucket: "test-bucket",
	})

	commit := &dal.Commit{
		UserID:  "user1",
		Message: "Initial commit",
		Content: "This is the commit content.",
	}

	mockS3.EXPECT().
		PutObject(gomock.Any(), gomock.Any()).
		Return(&s3.PutObjectOutput{VersionId: aws.String("1"), ChecksumSHA256: aws.String("abc123")}, nil)

	mockDynamoDB.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	mockCache.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := client.CreateCommit(context.Background(), "org1", "prompt1", "branch1", commit)
	assert.NoError(t, err)
	assert.NotEmpty(t, commit.CommitID)
	assert.Equal(t, "1", commit.CommitID)
}

func TestGetCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	mockS3 := mocks.NewMockS3API(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	client := dal.NewCommitDBClient(mockDynamoDB, mockS3, mockCache, &config.Config{
		PromptBucket: "test-bucket",
	})

	commit := dal.Commit{
		CommitID:  "commit1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(commit)
	mockDynamoDB.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	mockCache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return("", errors.New("cache miss"))

	mockS3.EXPECT().
		GetObject(gomock.Any(), gomock.Any()).
		Return(&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader("This is the commit content."))}, nil)

	mockCache.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	result, err := client.GetCommit(context.Background(), "org1", "prompt1", "branch1", "1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "commit1", result.CommitID)
	assert.Equal(t, "This is the commit content.", result.Content)
}

func TestListCommitsByBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	mockS3 := mocks.NewMockS3API(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	client := dal.NewCommitDBClient(mockDynamoDB, mockS3, mockCache, &config.Config{
		PromptBucket: "test-bucket",
	})

	commit := dal.Commit{
		CommitID:  "1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(commit)
	mockDynamoDB.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	results, err := client.ListCommitsByBranch(context.Background(), "org1", "prompt1", "branch1")
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 1)
	assert.Equal(t, "1", results[0].CommitID)
}
