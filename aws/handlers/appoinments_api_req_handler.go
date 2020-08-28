package handlers

import (
	"net/http"

	"github.com/sal-org/sal-backend/counselor"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sal-org/sal-backend/aws/database"
	"github.com/sal-org/sal-backend/aws/repositories"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/models"
)

// HandleAppointmentsRequest will handle all appointment gate api requests
func HandleAppointmentsRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	db, err := database.CreateDB()

	if err != nil {
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	repository := repositories.AppointmentsRepository{DB: db}

	//TODO
	counselor := counselor.Counselor{
		Identifier: "random",
	}

	appointments, err := repository.FetchAll(&counselor)
	if err != nil {
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	return utils.GatewayResponse(http.StatusOK, appointments)
}
