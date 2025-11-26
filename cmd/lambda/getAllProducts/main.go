package main

import (
	"context"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/lambda"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/notifier"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/repository"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	client, err := lambda.GetDynamoClient(ctx)
	if err != nil {
		return lambda.RespondWithError(500, "Failed to initialize database client")
	}

	productRepo := repository.NewProductDynamoRepository(client, config.GetProductsTableName())
	logNotifier := notifier.NewLogNotifier()
	inventoryService := service.NewInventoryService(productRepo, logNotifier)

	products, err := inventoryService.GetAllProducts()
	if err != nil {
		return lambda.HandleError(err)
	}

	return lambda.RespondWithJSON(200, products)
}

func main() {
	awslambda.Start(handler)
}
