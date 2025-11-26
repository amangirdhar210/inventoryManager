package lambda

import (
	"context"

	"github.com/amangirdhar210/inventory-manager/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client

func GetDynamoClient(ctx context.Context) (*dynamodb.Client, error) {
	if dynamoClient != nil {
		return dynamoClient, nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(config.GetAWSRegion()))
	if err != nil {
		return nil, err
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	return dynamoClient, nil
}
