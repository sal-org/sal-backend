package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

	UTIL "salbackend/util"
)

// ProfileGet godoc
// @Tags Client Profile
// @Summary Get client profile with email, if signed up already
// @Router /client [get]
// @Param email query string true "Email of client - to get details, if signed up already"
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": r.FormValue("email")})
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

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

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

	// add client details
	client := map[string]string{}
	client["first_name"] = body["first_name"]
	client["last_name"] = body["last_name"]
	client["phone"] = body["phone"]
	client["email"] = body["email"]
	client["age"] = body["age"]
	client["gender"] = body["gender"]
	client["location"] = body["location"]
	client["timezone"] = body["timezone"]
	client["device_id"] = body["device_id"]
	client["status"] = CONSTANT.ClientActive
	client["last_login_time"] = UTIL.GetCurrentTime().String()
	client["created_at"] = UTIL.GetCurrentTime().String()
	clientID, status, ok := DB.InsertWithUniqueID(CONSTANT.ClientsTable, CONSTANT.ClientDigits, client, "client_id")
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

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	response["client_id"] = clientID
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
	if len(body["age"]) > 0 {
		client["age"] = body["age"]
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
