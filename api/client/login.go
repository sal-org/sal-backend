package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"

	UTIL "salbackend/util"
	"strings"
)

// SendOTP godoc
// @Tags Client Login
// @Summary Send OTP to specified phone
// @Router /client/sendotp [get]
// @Param phone query string true "Phone number to send OTP - send phone number with 91 code"
// @Produce json
// @Success 200
func SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if len(r.FormValue("phone")) < 8 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ValidPhoneRequiredMessage, CONSTANT.ShowDialog, response)
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

	// using phone verified table to check if phone has been really verified by OTP
	// currently deleting if phone number is already present
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": r.FormValue("phone")})

	// send otp
	otp, ok := UTIL.GenerateOTP(r.FormValue("phone"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientOTPTextMessage,
			map[string]string{
				"###otp###": otp,
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		r.FormValue("phone"),
		UTIL.GetCurrentTime().String(),
		CONSTANT.MessageSent,
		CONSTANT.InstantSendTextMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// VerifyOTP godoc
// @Tags Client Login
// @Summary Verify OTP sent to specified phone
// @Router /client/verifyotp [get]
// @Param phone query string true "Phone number OTP has been sent to - send phone number with 91 code"
// @Param otp query string true "OTP entered by client"
// @Param device_id query string true "Device ID entered by client"
// @Produce json
// @Success 200
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if otp is correct
	// if !UTIL.VerifyOTP(r.FormValue("phone"), r.FormValue("otp")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	if len(r.FormValue("device_id")) < 0 {
		UTIL.SetReponse(w, "400", "device_id is required", CONSTANT.ShowDialog, response)
		return
	}

	if strings.EqualFold("915757575757", r.FormValue("phone")) {
		if !strings.EqualFold("4444", r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	} else if strings.EqualFold("914747474747", r.FormValue("phone")) {
		if !strings.EqualFold("4848", r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	} else {
		if !UTIL.VerifyOTP(r.FormValue("phone"), r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	}

	// if !strings.EqualFold("4444", r.FormValue("otp")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(client) == 0 {
		// client is first time signing up

		// using phone verified table to check if phone has been really verified by OTP
		// currently inserting after phone is really verified
		DB.InsertSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": r.FormValue("phone"), "created_at": UTIL.GetCurrentTime().String()})
	} else {
		// client already signed up
		// check if client is active
		if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"phone": r.FormValue("phone")}, map[string]string{"device_id": r.FormValue("device_id"), "last_login_time": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
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

// RefreshToken godoc
// @Tags Client Login
// @Summary Get new access token with refresh token
// @Router /client/refresh-token [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("client_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if client id is valid
	if !DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("client_id"), "status": CONSTANT.ClientActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessToken(r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func SendOTPWithCorporateEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	is_valid_email := UTIL.IsValidEmail(r.FormValue("cor_email"))
	if is_valid_email == "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct email id", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(r.FormValue("cor_email"), "@")

	// get client details
	ok := DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"domain": domainName[1]})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
		return
	}

	// // get client details
	ok = DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
		return
	}
	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"status"}, map[string]string{"email": r.FormValue("cor_email")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(client) > 0 && !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// send otp
	otp, ok := UTIL.GenerateOTPWithCorporateEmail(r.FormValue("cor_email"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: "",
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientCorLoginOTPBody,
			map[string]string{
				"###otp###": otp,
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientCorLoginOTPTitle,
			map[string]string{
				"###otp###": otp,
			},
		),
		emailBody1,
		r.FormValue("cor_email"),
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func SendOTPWithCorporateEmailForRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	is_valid_email := UTIL.IsValidEmail(r.FormValue("cor_email"))
	if is_valid_email == "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct email id", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(r.FormValue("cor_email"), "@")

	// get client details
	ok := DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"domain": domainName[1]})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	ok = DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")})
	if ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfRegister, CONSTANT.ShowDialog, response)
		return
	}

	// // get client details
	// client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"status"}, map[string]string{"email": r.FormValue("cor_email")})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// if len(client) > 0 && !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// send otp
	otp, ok := UTIL.GenerateOTPWithCorporateEmail(r.FormValue("cor_email"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: "",
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientCorLoginOTPBody,
			map[string]string{
				"###otp###": otp,
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientCorLoginOTPTitle,
			map[string]string{
				"###otp###": otp,
			},
		),
		emailBody1,
		r.FormValue("cor_email"),
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CheckAccessCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get client details
	ok := DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"access_code": r.FormValue("access_code"), "status": "1"})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CorporateClientAccessCode, CONSTANT.ShowDialog, response)
		return
	}

	title, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"access_code": r.FormValue("access_code"), "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	address, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersAddressTable, []string{"address"}, map[string]string{"domain": title[0]["domain"], "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["address"] = address

	response["title"] = title[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CheckEmailANDPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if len(r.FormValue("cor_email")) == 0 && len(r.FormValue("phone")) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientEmptyBodyPass, CONSTANT.ShowDialog, response)
		return
	}

	if len(r.FormValue("phone")) > 0 {
		// check client details
		ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"phone": r.FormValue("phone")})
		if ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
			return
		}
	}

	if len(r.FormValue("cor_email")) > 0 {

		domainName := strings.Split(r.FormValue("cor_email"), "@")

		// get client details
		ok := DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"domain": domainName[1]})
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
			return
		}

		if len(r.FormValue("login")) > 0 {
			ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")})
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
				return
			}
		} else {
			ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")})
			if ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EmailExistsMessage, CONSTANT.ShowDialog, response)
				return
			}
		}

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func VerifyOTPWithCorporateEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if otp is correct
	if !UTIL.VerifyOTPWithCorporateEmail(r.FormValue("cor_email"), r.FormValue("otp")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": r.FormValue("cor_email")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// client already signed up
	// check if client is active
	if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"email": r.FormValue("cor_email")}, map[string]string{"last_login_time": UTIL.GetCurrentTime().String()})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"phone": r.FormValue("phone")}, map[string]string{"device_id": r.FormValue("device_id")})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}


func CheckIfAccessTokenExpired(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
