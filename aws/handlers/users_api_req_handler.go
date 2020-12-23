package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sal-org/sal-backend/aws/database"
	"github.com/sal-org/sal-backend/aws/repositories"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/models"
)

// HandleUsersRequest will handle all users gate api requests
func HandleUsersRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		db, err := database.CreateDB()

		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		repository := repositories.UsersRepository{DB: db}
		userId  := req.QueryStringParameters["user"]
		users, err := repository.FetchUser(userId)
		return utils.GatewayResponse(http.StatusOK, users)

	default:
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String("Method not allowed")})
	}
}
