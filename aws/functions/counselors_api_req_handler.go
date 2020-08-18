package functions

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sal-org/sal-backend/aws/database"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/models"
)

// HandleCounselorsRequest will handle all counselor gate api requests
func HandleCounselorsRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		counselors, err := database.GetAllCounselors()

		if err != nil {
			return utils.APIResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}
		return utils.APIResponse(http.StatusOK, counselors)
	case "POST":
		counselor, err := database.CreateCounselor(req)
		if err != nil {
			return utils.APIResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		return utils.APIResponse(http.StatusOK, counselor)
	default:
		return utils.APIResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String("Method not allowed")})
	}
}
