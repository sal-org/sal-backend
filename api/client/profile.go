package client

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
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if len(r.FormValue("device_id")) < 0 {
		UTIL.SetReponse(w, "400", "device_id is required", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(client) > 0 && !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("email")}, map[string]string{"device_id": r.FormValue("device_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok = DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": r.FormValue("email")})
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
		accessToken, ok := UTIL.CreateAccessToken(client[0]["client_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		refreshToken, ok := UTIL.CreateRefreshToken(client[0]["client_id"])
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
		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileAdd godoc
// @Tags Client Profile
// @Summary Add client profile after OTP verified to signup
// @Router /client [post]
// @Param body body model.ClientProfileAddRequest true "Request Body"
// @Produce json
// @Success 200
func ProfileAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

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

	if len(body["notification_status"]) == 0 {
		body["notification_status"] = "1"
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
	client["notification_status"] = body["notification_status"]
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
	accessToken, ok := UTIL.CreateAccessToken(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshToken(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send notification to client
	UTIL.SendNotification(CONSTANT.ClientCompletedProfileHeading, CONSTANT.ClientCompletedProfileContent, clientID, CONSTANT.TherapistType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, clientID)

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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ProfileAddForCor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CorporateClientProfileAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	//check if otp is correct
	if !UTIL.VerifyOTPWithCorporateEmail(body["email"], body["otp"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check email format
	is_valid_email := UTIL.IsValidEmail(body["email"])
	if is_valid_email == "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct email id", CONSTANT.ShowDialog, response)
		return
	}

	// get domain
	domainName := strings.Split(body["email"], "@")

	// check domain exists or not
	ok = DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"domain": domainName[1]})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
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

	// // check if phone is verfied by OTP
	// if !DB.CheckIfExists(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]}) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.VerifyPhoneRequiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // default photo for client and listerner
	// if len(body["photo"]) == 0 || body["photo"] == "string" {
	// 	body["photo"] = CONSTANT.DefaultPhotoForClientAndListerner
	// }

	if len(body["notification_status"]) == 0 {
		body["notification_status"] = "1"
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
	client["department"] = body["cor_darpartment"]
	client["timezone"] = body["timezone"]
	client["device_id"] = body["device_id"]
	client["platform"] = body["platform"]
	client["version"] = body["version"]
	client["status"] = CONSTANT.ClientActive
	client["notification_status"] = body["notification_status"]
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
	// DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]})

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessToken(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshToken(clientID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send notification to client
	UTIL.SendNotification(CONSTANT.ClientCompletedProfileHeading, CONSTANT.ClientCompletedProfileContent, clientID, CONSTANT.TherapistType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, clientID)

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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
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
func ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update client details
	client := map[string]string{}
	if len(body["first_name"]) > 0 {
		client["first_name"] = body["first_name"]
	}
	if len(body["last_name"]) > 0 {
		client["last_name"] = body["last_name"]
	}
	if len(body["location"]) > 0 {
		client["location"] = body["location"]
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
	status, ok := DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("client_id")}, client)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
