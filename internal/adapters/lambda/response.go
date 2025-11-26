package lambda

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/aws/aws-lambda-go/events"
)

func RespondWithJSON(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	responseBody, err := json.Marshal(body)
	if err != nil {
		return RespondWithError(http.StatusInternalServerError, "Failed to marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(responseBody),
	}, nil
}

func RespondWithError(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: `{"error":"` + message + `"}`,
	}, nil
}

func HandleError(err error) (events.APIGatewayProxyResponse, error) {
	switch err {
	case domain.ErrProductNotFound:
		return RespondWithError(http.StatusNotFound, "Product not found")
	case domain.ErrProductInvalid:
		return RespondWithError(http.StatusBadRequest, "Invalid product data")
	case domain.ErrManagerNotFound:
		return RespondWithError(http.StatusUnauthorized, "Invalid credentials")
	case domain.ErrUnauthorized:
		return RespondWithError(http.StatusUnauthorized, "Unauthorized")
	case domain.ErrRepository:
		return RespondWithError(http.StatusInternalServerError, "Database error")
	default:
		return RespondWithError(http.StatusInternalServerError, err.Error())
	}
}
