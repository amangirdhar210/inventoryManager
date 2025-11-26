package main

import (
	"context"
	"encoding/json"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/lambda"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/lambda/dto"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/repository"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/amangirdhar210/inventory-manager/utils/auth"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	client, err := lambda.GetDynamoClient(ctx)
	if err != nil {
		return lambda.RespondWithError(500, "Failed to initialize database client")
	}

	managerRepo := repository.NewManagerDynamoRepository(client, config.GetManagersTableName())
	tokenGenerator := auth.NewJWTGenerator(config.GetJWTSecretKey())
	authService := service.NewAuthService(managerRepo, tokenGenerator)

	var req dto.LoginRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return lambda.RespondWithError(400, "Invalid request body")
	}

	token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		return lambda.HandleError(err)
	}

	return lambda.RespondWithJSON(200, map[string]string{"token": token})
}

func main() {
	awslambda.Start(handler)
}
