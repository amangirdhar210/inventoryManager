package repository

import (
	"context"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type managerDynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewManagerDynamoRepository(client *dynamodb.Client, tableName string) *managerDynamoRepository {
	return &managerDynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (repo *managerDynamoRepository) FindByEmail(email string) (*domain.Manager, error) {
	result, err := repo.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(repo.tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, domain.ErrManagerNotFound
	}

	var manager domain.Manager
	err = attributevalue.UnmarshalMap(result.Items[0], &manager)
	if err != nil {
		return nil, err
	}

	return &manager, nil
}
