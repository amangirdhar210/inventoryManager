package repository

import (
	"context"
	"fmt"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type productDynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewProductDynamoRepository(client *dynamodb.Client, tableName string) *productDynamoRepository {
	return &productDynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (repo *productDynamoRepository) FindById(id string) (*domain.Product, error) {
	result, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, domain.ErrRepository
	}

	if result.Item == nil {
		return nil, domain.ErrProductNotFound
	}

	var product domain.Product
	err = attributevalue.UnmarshalMap(result.Item, &product)
	if err != nil {
		return nil, domain.ErrRepository
	}

	return &product, nil
}

func (repo *productDynamoRepository) Save(product *domain.Product) error {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return domain.ErrRepository
	}

	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      item,
	})
	if err != nil {
		return domain.ErrRepository
	}

	return nil
}

func (repo *productDynamoRepository) Update(product *domain.Product) error {
	_, err := repo.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: product.Id},
		},
		UpdateExpression: aws.String("SET #name = :name, price = :price, quantity = :quantity"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":     &types.AttributeValueMemberS{Value: product.Name},
			":price":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", product.Price)},
			":quantity": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.Quantity)},
		},
	})
	if err != nil {
		return domain.ErrRepository
	}

	return nil
}

func (repo *productDynamoRepository) DeleteById(id string) error {
	_, err := repo.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	})
	if err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if ok := err.(*types.ConditionalCheckFailedException); ok == conditionErr {
			return domain.ErrProductNotFound
		}
		return domain.ErrRepository
	}

	return nil
}

func (repo *productDynamoRepository) ListAll() ([]domain.Product, error) {
	result, err := repo.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	})
	if err != nil {
		return nil, domain.ErrRepository
	}

	var products []domain.Product
	err = attributevalue.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		return nil, domain.ErrRepository
	}

	if products == nil {
		products = []domain.Product{}
	}

	return products, nil
}
