package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/sal-org/sal-backend/aws/functions"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/constants"
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "counselors":
		return functions.HandleRequest(req)
	default:
		return utils.APIResponse(http.StatusMethodNotAllowed, constants.ErrorMethodNotAllowed)
	}
}
