package handlers

import (
	"log"
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
		//TO-DO this is work in progress need to add api to get appointments per user/counselor/timeslot
		counselorId  := req.QueryStringParameters["counselor"]
		userId  := req.QueryStringParameters["user"]
		appointments, err := repository.FetchAll(counselorId)
		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}
		return utils.GatewayResponse(http.StatusOK, appointments)
		
	case "POST":
		db, err := database.CreateDB()

		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		repository := repositories.AppointmentsRepository{DB: db}
		appointment, err := repository.CreateAppointment(req)
		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		return utils.GatewayResponse(http.StatusOK, appointment)
	
	default:
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String("Method not allowed")})
	}
}
