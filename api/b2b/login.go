package b2b

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strings"
	"time"
)

func CheckAccessCode(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var encryptedResponse = make(map[string]interface{})

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get client details
	ok := DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"access_code": body["access_code"], "status": "1"})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CorporateClientAccessCode, CONSTANT.ShowDialog, response)
		return
	}

	title, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"access_code": body["access_code"], "status": "1"})
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

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}

func SendOTPWithCorporateEmail(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get client details
	ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": body["email"]})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"status", "company_name", "client_id"}, map[string]string{"email": body["email"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(client) > 0 && client[0]["status"] == CONSTANT.ClientBlocked {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(client) > 0 && !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAccountBlockedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if client[0]["company_name"] == "" || client[0]["company_name"] != "prosculpt" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	transitions, status, ok := DB.SelectProcess("select * from "+CONSTANT.B2B2CAppointmentTransitionsTable+" where client_id = ? and status != "+CONSTANT.AppointmentTransitionCompleted+" order by created_at asc", client[0]["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment transitions
	if len(transitions) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAppointmentPayment, CONSTANT.ShowDialog, response)
		return
	}

	if transitions[0]["status"] == CONSTANT.AppointmentTransitionCompleted {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAppointmentPayment, CONSTANT.ShowDialog, response)
		return
	}

	if transitions[0]["status"] == CONSTANT.AppointmentTransitionActive {

		expireCheck, _ := time.Parse("2006-01-02 15:04:05", transitions[0]["created_at"])

		expiry := expireCheck.Add(30 * 24 * time.Hour)

		if UTIL.GetCurrentTime().After(expiry) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAppointmentPayment, CONSTANT.ShowDialog, response)
			return
		}
	}

	// send otp
	otp, ok := UTIL.GenerateOTPWithCorporateEmail(body["email"])
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: "",
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientWebApplicationOTPBody,
			map[string]string{
				"###otp###": otp,
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientWebApplicationOTPTitle,
		emailBody1,
		body["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func VerifyOTPWithCorporateEmail(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	//check if otp is correct
	if !UTIL.VerifyOTPWithCorporateEmail(body["email"], body["otp"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	ok := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"email": body["email"]})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"email": body["email"]})
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

	// topics, status, ok := DB.SelectProcess("select topic from " + CONSTANT.TopicsTable + " where id in (" + client[0]["topic_ids"] + ")")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken
	// response["topic"] = topics
	response["client"] = client[0]
	response["media_url"] = CONFIG.MediaURL

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}
