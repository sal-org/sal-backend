package functions

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/handlers"
	"github.com/sal-org/sal-backend/models"
)

// HandleAppointmentsRequest will handle all appointment gate api requests
func HandleAppointmentsRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	appointments, err := handlers.GetAllAppointments()

	if err != nil {
		return utils.APIResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	return utils.APIResponse(http.StatusOK, appointments)
}
