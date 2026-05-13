package b2b

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strings"

	UTIL "salbackend/util"
)

// ProfileGet godoc
// @Tags Client Profile
// @Summary Get client profile with email, if signed up already
// @Router /client [get]
// @Param email query string true "Email of client - to get details, if signed up already"
// @Param device_id query string true "Device ID of client - to get details, if signed up already"
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": body["email"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(client) > 0 {
		// client already signed up
		// check if client is active
		if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		// generate access and refresh token
		// access token - jwt token with short expiry added in header for authorization
		// refresh token - jwt token with long expiry to get new access token if expired
		// if refresh token expired, need to login
		accessToken, ok := UTIL.CreateAccessTokenForWeb(client[0]["client_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		refreshToken, ok := UTIL.CreateRefreshTokenForWeb(client[0]["client_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}

		topics, status, ok := DB.SelectProcess("select topic from " + CONSTANT.TopicsTable + " where id in (" + client[0]["topic_ids"] + ")")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken
		response["topic"] = topics
		response["client"] = client[0]
	}

	encrypt, _ := UTIL.EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}

// ProfileAdd godoc
// @Tags Client Profile
// @Summary Add client profile after OTP verified to signup
// @Router /client [post]
// @Param body body model.ClientProfileAddRequest true "Request Body"
// @Produce json
// @Success 200
func ProfileAdd(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ClientProfileAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified email
	if DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": body["email"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EmailExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if phone is verfied by OTP
	if !DB.CheckIfExists(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.VerifyPhoneRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add client details
	client := map[string]string{}
	client["first_name"] = body["first_name"]
	client["last_name"] = body["last_name"]
	client["phone"] = body["phone"]
	client["email"] = body["email"]
	client["date_of_birth"] = body["date_of_birth"]
	client["photo"] = body["photo"]
	client["topic_ids"] = body["topic_ids"]
	client["gender"] = body["gender"]
	client["location"] = body["location"]
	client["timezone"] = body["timezone"]
	client["device_id"] = body["device_id"]
	client["platform"] = body["platform"]
	client["version"] = body["version"]
	client["status"] = CONSTANT.ClientActive
	client["last_login_time"] = UTIL.GetCurrentTime().String()
	client["created_at"] = UTIL.GetCurrentTime().String()
	clientID, status, ok := DB.InsertWithUniqueID(CONSTANT.ClientsTable, CONSTANT.ClientDigits, client, "client_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	clientD, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": clientID})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// using phone verified table to check if phone has been really verified by OTP
	// currently deleting if phone number is already present
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]})

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessTokenForWeb(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshTokenForWeb(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata := Model.EmailBodyMessageModel{
		Name:    body["first_name"],
		Message: CONSTANT.ClientSignupClientEmailBody,
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientSignupProfileTitle,
		emailBody,
		body["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientProfileTitleMessage,
			map[string]string{
				"###client_name###": body["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		body["phone"],
		UTIL.GetCurrentTime().String(),
		clientID,
		CONSTANT.InstantSendTextMessage,
	)

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	response["client"] = clientD[0]
	response["media_url"] = CONFIG.MediaURL

	encrypt, _ := UTIL.EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}

// ProfileUpdate godoc
// @Tags Client Profile
// @Summary Update client profile details
// @Router /client [put]
// @Param client_id query string true "Client ID to update details"
// @Param body body model.ClientProfileUpdateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func ProfileUpdate(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// update client details
	client := map[string]string{}
	if len(body["emp_id"]) > 0 {
		client["emp_id"] = body["emp_id"]
	}
	if len(body["first_name"]) > 0 {
		client["first_name"] = body["first_name"]
	}
	if len(body["last_name"]) > 0 {
		client["last_name"] = body["last_name"]
	}
	if len(body["location"]) > 0 {
		client["location"] = body["location"]
	}
	if len(body["department"]) > 0 {
		client["department"] = body["department"]
	}
	if len(body["date_of_birth"]) > 0 {
		client["date_of_birth"] = body["date_of_birth"]
	}
	if len(body["photo"]) > 0 {
		client["photo"] = body["photo"]
	}
	if len(body["topic_ids"]) > 0 {
		client["topic_ids"] = body["topic_ids"]
	}
	if len(body["gender"]) > 0 {
		client["gender"] = body["gender"]
	}
	if len(body["device_id"]) > 0 {
		client["device_id"] = body["device_id"]
	}
	if len(body["timezone"]) > 0 {
		client["timezone"] = body["timezone"]
	}
	client["last_login_time"] = UTIL.GetCurrentTime().String()
	client["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": body["client_id"]}, client)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AddProsculptStudentProfile(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	var isMobileExist = false
	var isEmailIDExist = false

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ProsculptStudentProfileRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	if DB.CheckIfExists(CONSTANT.B2B2CAppointmentTransitionsTable, map[string]string{"payment_id": body["payment_id"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ProsculptPaymentIDAlreadyMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"phone": body["mobile_no"]}) {
		isMobileExist = true
	}

	// check if user already signed up with specified email
	if DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": body["email_id"]}) {
		isEmailIDExist = true
	}

	if isEmailIDExist {

		if !isMobileExist {
			DB.UpdateSQL(CONSTANT.ClientsTable,
				map[string]string{
					"email": body["email_id"],
				},
				map[string]string{
					"company_name":  "prosculpt",
					"phone":         body["mobile_no"],
					"first_name":    body["first_name"],
					"last_name":     body["last_name"],
					"date_of_birth": body["dob"],
					"gender":        body["gender"],
					"location":      body["location"],
					"modified_at":   UTIL.GetCurrentTime().String(),
				},
			)

			isMobileExist = true
		}

	}

	if isMobileExist {

		if !isEmailIDExist {
			DB.UpdateSQL(CONSTANT.ClientsTable,
				map[string]string{
					"phone": body["mobile_no"],
				},
				map[string]string{
					"company_name":  "prosculpt",
					"email":         body["email_id"],
					"first_name":    body["first_name"],
					"last_name":     body["last_name"],
					"date_of_birth": body["dob"],
					"gender":        body["gender"],
					"location":      body["location"],
					"modified_at":   UTIL.GetCurrentTime().String(),
				},
			)

			isEmailIDExist = true
		}

	}

	if isEmailIDExist && isMobileExist {

		clientD, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": body["email_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		transition := map[string]string{}

		transition["client_id"] = clientD[0]["client_id"]
		transition["payment_id"] = body["payment_id"]
		transition["status"] = CONSTANT.ClientActive
		transition["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.B2B2CAppointmentTransitionsTable, CONSTANT.OrderTransitionsDigits, transition, "order_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	} else {

		// add client details
		client := map[string]string{}
		client["company_name"] = "prosculpt"
		client["first_name"] = body["first_name"]
		client["last_name"] = body["last_name"]
		client["phone"] = body["mobile_no"]
		client["email"] = body["email_id"]
		client["date_of_birth"] = body["dob"]
		client["gender"] = body["gender"]
		client["location"] = body["location"]
		client["timezone"] = "330"
		client["platform"] = "web"
		client["status"] = CONSTANT.ClientActive
		client["last_login_time"] = UTIL.GetCurrentTime().String()
		client["created_at"] = UTIL.GetCurrentTime().String()
		clientID, status, ok := DB.InsertWithUniqueID(CONSTANT.ClientsTable, CONSTANT.ClientDigits, client, "client_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		transition := map[string]string{}

		transition["client_id"] = clientID
		transition["payment_id"] = body["payment_id"]
		transition["status"] = CONSTANT.ClientActive
		transition["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.B2B2CAppointmentTransitionsTable, CONSTANT.OrderTransitionsDigits, transition, "order_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// send email to client
	// filepath_text := "htmlfile/emailmessagebody.html"

	// emaildata := Model.EmailBodyMessageModel{
	// 	Name:    body["first_name"],
	// 	Message: CONSTANT.ClientSignupClientEmailBody,
	// }

	// emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// // email for client
	// UTIL.SendEmail(
	// 	CONSTANT.ClientSignupProfileTitle,
	// 	emailBody,
	// 	body["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}
