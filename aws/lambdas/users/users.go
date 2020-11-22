package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/sal-org/sal-backend/aws/handlers"
)

func main() {
	lambda.Start(HandleRequest)
}

// HandleRequest handles the api gateway trigger for patients rest api
func HandleRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return handlers.HandleUsersRequest(req)
}
