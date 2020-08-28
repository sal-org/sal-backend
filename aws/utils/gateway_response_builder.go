package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sal-org/sal-backend/constants"
)

// GatewayResponse function builds an APIGatewayProxyResponse object with the provided status and response body
func GatewayResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{Headers: map[string]string{constants.RequestContentType: constants.JSONResponse}}
	response.StatusCode = status

	stringBody, _ := json.Marshal(body)
	response.Body = string(stringBody)
	return &response, nil
}
