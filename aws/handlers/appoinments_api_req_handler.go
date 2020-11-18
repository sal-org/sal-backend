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

// HandleAppointmentsRequest will handle all appointment gate api requests
func HandleAppointmentsRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		db, err := database.CreateDB()

		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		repository := repositories.AppointmentsRepository{DB: db}
		counselors, err := repository.FetchAll(nil)
		return utils.GatewayResponse(http.StatusOK, counselors)
	case "POST":
		db, err := database.CreateDB()

		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		repository := repositories.AppointmentsRepository{DB: db}
		counselor, err := repository.CreateAppointment(req)
		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		return utils.GatewayResponse(http.StatusOK, counselor)
	default:
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String("Method not allowed")})
	}
}
