package functions

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/handlers"
	"github.com/sal-org/sal-backend/models"
)

// HandleRequest will handle all counselor gate api requests
func HandleRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	counselors, err := handlers.GetAllCounselors()

	if err != nil {
		return utils.APIResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	return utils.APIResponse(http.StatusOK, counselors)
}
