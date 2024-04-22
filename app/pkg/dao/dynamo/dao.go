package dynamodao

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	_ "github.com/lib/pq"
)

type DAO struct {
	client dynamodbiface.DynamoDBAPI
}

func NewDAO(client dynamodbiface.DynamoDBAPI) *DAO {
	return &DAO{client: client}
}
