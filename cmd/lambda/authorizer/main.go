package main

import (
	"context"
	"errors"
	"strings"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken

	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecretKey()), nil
	})

	if err != nil || !parsedToken.Valid {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	principalID := claims.Subject
	if principalID == "" {
		principalID = "user"
	}

	parts := strings.Split(request.MethodArn, ":")
	apiGatewayArnParts := strings.Split(parts[5], "/")
	resource := strings.Join(parts[:5], ":") + ":" + apiGatewayArnParts[0] + "/*/*"

	return generatePolicy(principalID, "Allow", resource), nil
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
