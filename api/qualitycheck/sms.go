package qualitycheck

import (
	"net/http"
	CONSTANT "salbackend/constant"
	UTIL "salbackend/util"
)

func SendSMS(w http.ResponseWriter, r *http.Request) {

	// set header application type
	w.Header().Set("Content-type", "application/json")

	// create an variable for response
	response := make(map[string]interface{})

	// Read all body request
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		// if request body not correct then send error message
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check all the required parameters
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.SMSServicesRequiredFields)
	if len(fieldCheck) > 0 {
		// if request body not correct then send error message
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// send sms to target user
	UTIL.SendMessage(body["body"], CONSTANT.TransactionalRouteTextMessage, body["phone"], UTIL.GetCurrentTime().String(), "2", true)

	// this is success response
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
