package main

import (
	"context"
	"strings"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken

	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return generatePolicy("", "Deny", request.MethodArn), nil
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecretKey()), nil
	})

	if err != nil || !parsedToken.Valid {
		return generatePolicy("", "Deny", request.MethodArn), nil
	}

	principalID := claims.Subject
	if principalID == "" {
		principalID = "user"
	}

	return generatePolicy(principalID, "Allow", request.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
	}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}

func main() {
	awslambda.Start(handler)
}
