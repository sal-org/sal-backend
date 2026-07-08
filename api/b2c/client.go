package b2c

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
)

func ClientBookDemo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ClientBookDemoInWebsiteAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.WebB2BBookDemoTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified email
	if DB.CheckIfExists(CONSTANT.WebB2BBookDemoTable, map[string]string{"email": body["email"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EmailExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add client details
	client := map[string]string{}
	client["name"] = body["name"]
	client["phone"] = body["phone"]
	client["email"] = body["email"]
	client["company_name"] = body["company_name"]
	client["company_location"] = body["location"]
	client["company_size"] = body["size"]
	client["message"] = body["message"]
	client["status"] = CONSTANT.ClientActive
	client["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.WebB2BBookDemoTable, CONSTANT.ClientDigits, client, "demo_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// send email to client
	filepath_text := "htmlfile/client_book_demo.html"

	data := Model.EmailDataForWebClientB2CProfile{ 
		Name: body["name"],
		Phone: body["phone"],
		Email: body["email"],
		CompanyName: body["company_name"],
		CompanyLocation: body["location"],
		CompanySize: body["size"],
		Message: body["message"],
	}

	emailbody := UTIL.GetHTMLTemplateForWebB2CClientProfile(data, filepath_text)


	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientB2BRegistrationForDemoProfileTitle,
		emailbody,
		CONFIG.EventEmailID,
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
