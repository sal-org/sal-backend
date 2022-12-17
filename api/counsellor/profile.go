package counsellor

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strings"
)

// ProfileGet godoc
// @Tags Counsellor Profile
// @Summary Get counsellor profile with email, if signed up already
// @Router /counsellor [get]
// @Param email query string false "Email of counsellor - to get details, if signed up already"
// @Param counsellor_id query string false "Counsellor ID to update details"
// @Security JWTAuth
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor details
	params := map[string]string{}
	if len(r.FormValue("email")) > 0 {
		params["email"] = r.FormValue("email")
	}
	if len(r.FormValue("counsellor_id")) > 0 {
		params["counsellor_id"] = r.FormValue("counsellor_id")
	}
	if len(params) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, params)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(counsellor) > 0 {
		// counsellor already signed up
		// check if counsellor is active
		if !strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		// generate access and refresh token
		// access token - jwt token with short expiry added in header for authorization
		// refresh token - jwt token with long expiry to get new access token if expired
		// if refresh token expired, need to login
		accessToken, ok := UTIL.CreateAccessToken(counsellor[0]["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		refreshToken, ok := UTIL.CreateRefreshToken(counsellor[0]["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}

		languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", counsellor[0]["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor topics
		topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", counsellor[0]["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken
		response["languages"] = languages
		response["topics"] = topics
		response["counsellor"] = counsellor[0]
		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileAdd godoc
// @Tags Counsellor Profile
// @Summary Add counsellor profile after OTP verified to signup
// @Router /counsellor [post]
// @Param body body model.CounsellorProfileAddRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorProfileAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.CounsellorsTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if phone is verfied by OTP
	if !DB.CheckIfExists(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.VerifyPhoneRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add counsellor details
	counsellor := map[string]string{}
	counsellor["first_name"] = body["first_name"]
	counsellor["last_name"] = body["last_name"]
	counsellor["gender"] = body["gender"]
	counsellor["phone"] = body["phone"]
	counsellor["photo"] = body["photo"]
	counsellor["email"] = body["email"]
	counsellor["price"] = body["price"]
	counsellor["multiple_sessions"] = body["multiple_sessions"]
	counsellor["price_3"] = body["price_3"]
	counsellor["price_5"] = body["price_5"]
	counsellor["education"] = body["education"]
	counsellor["experience"] = body["experience"]
	counsellor["about"] = body["about"]
	counsellor["timezone"] = body["timezone"]
	counsellor["resume"] = body["resume"]
	counsellor["certificate"] = body["certificate"]
	counsellor["aadhar"] = body["aadhar"]
	counsellor["linkedin"] = body["linkedin"]
	counsellor["payout_percentage"] = CONSTANT.CounsellorPayoutPercentageColumns
	counsellor["payee_name"] = body["payee_name"]
	counsellor["bank_account_no"] = body["bank_account_no"]
	counsellor["ifsc"] = body["ifsc"]
	counsellor["branch_name"] = body["branch_name"]
	counsellor["bank_name"] = body["bank_name"]
	counsellor["bank_account_type"] = body["bank_account_type"]
	counsellor["pan"] = body["pan"]
	counsellor["device_id"] = body["device_id"]
	counsellor["status"] = CONSTANT.CounsellorNotApproved
	counsellor["notification_status"] = CONSTANT.NotificationActive
	counsellor["last_login_time"] = UTIL.GetCurrentTime().String()
	counsellor["created_at"] = UTIL.GetCurrentTime().String()
	counsellorID, status, ok := DB.InsertWithUniqueID(CONSTANT.CounsellorsTable, CONSTANT.CounsellorDigits, counsellor, "counsellor_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// using phone verified table to check if phone has been really verified by OTP
	// currently deleting if phone number is already present
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]})

	// add languages, topics to counsellor
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], counsellorID)

	// not available for next 30 days. change here when you change in add new slot cron
	for i := 0; i < 30; i++ {
		DB.InsertSQL(CONSTANT.SlotsTable, map[string]string{"counsellor_id": counsellorID, "date": UTIL.GetCurrentTime().AddDate(0, 0, i).Format("2006-01-02")})
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessToken(counsellorID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshToken(counsellorID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send account signup notification, message to counsellor
	UTIL.SendNotification(CONSTANT.CounsellorAccountSignupCounsellorHeading, CONSTANT.CounsellorAccountSignupCounsellorContent, counsellorID, CONSTANT.CounsellorType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, counsellorID)
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAccountSignupTextMessage,
			map[string]string{
				"###counsellor_name###": body["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		body["phone"],
		UTIL.GetCurrentTime().String(),
		counsellorID,
		CONSTANT.InstantSendTextMessage,
	)

	// Counsellor details Send with SAL Team
	counsellor_details, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": counsellorID})

	// use this future
	// filepath_text := "htmlfile/emailmessagebody.html"

	// emaildata := Model.EmailBodyMessageModel{
	// 	Name:    counsellor_details[0]["first_name"],
	// 	Message: CONSTANT.CounsellorAccountSignupCounsellorEmailBody,
	// }

	// emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// // email for counsellor
	// UTIL.SendEmail(
	// 	CONSTANT.CounsellorProfileWaitingForApprovalTitle,
	// 	emailBody,
	// 	counsellor_details[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	data := Model.EmailDataForCounsellorProfile{
		First_Name:  counsellor_details[0]["first_name"],
		Last_Name:   counsellor_details[0]["last_name"],
		Gender:      counsellor_details[0]["gender"],
		Type:        "Counsellor",
		Phone:       counsellor_details[0]["phone"],
		Photo:       counsellor_details[0]["photo"],
		Email:       counsellor_details[0]["email"],
		Education:   counsellor_details[0]["education"],
		Experience:  counsellor_details[0]["experience"],
		About:       counsellor_details[0]["about"],
		Resume:      counsellor_details[0]["resume"],
		Certificate: counsellor_details[0]["certificate"],
		Aadhar:      counsellor_details[0]["aadhar"],
		Linkedin:    counsellor_details[0]["linkedin"],
		Status:      counsellor_details[0]["status"],
	}

	filepath := "htmlfile/CounsellorProfile.html"

	emailbody := UTIL.GetHTMLTemplateForProfile(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorProfileWaitingForApprovalTitle,
		emailbody,
		CONSTANT.AnandEmailID,
		CONSTANT.InstantSendEmailMessage,
	)

	/*UTIL.SendEmail(
		CONSTANT.CounsellorProfileWaitingForApprovalTitle,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorProfileHtml,
			map[string]string{
				"###First_Name###":  counsellor[0]["first_name"],
				"###Last_Name###":   orderdetails[0]["last_name"],
				"###Gender###":      orderdetails[0]["gender"],
				"###Phone###":       orderdetails[0]["phone"],
				"###Photo###":       orderdetails[0]["photo"],
				"###Email###":       orderdetails[0]["email"],
				"###Education###":   orderdetails[0]["education"],
				"###Experience###":  orderdetails[0]["experience"],
				"###About###":       orderdetails[0]["about"],
				"###Resume###":      orderdetails[0]["resume"],
				"###Certificate###": orderdetails[0]["certificate"],
				"###Aadhar###":      orderdetails[0]["aadhar"],
				"###Linkedin###":    orderdetails[0]["linkedin"],
				"###Status###":      orderdetails[0]["status"],
			},
		),
		CONSTANT.AnandEmailID,
		CONSTANT.InstantSendEmailMessage,
	)*/

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	response["counsellor"] = counsellor_details[0]
	response["media_url"] = CONFIG.MediaURL

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Counsellor Profile
// @Summary Update counsellor profile details
// @Router /counsellor [put]
// @Param counsellor_id query string true "Counsellor ID to update details"
// @Param body body model.CounsellorProfileUpdateRequest true "Request Body"
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

	// update counsellor details
	counsellor := map[string]string{}
	if len(body["first_name"]) > 0 {
		counsellor["first_name"] = body["first_name"]
	}
	if len(body["last_name"]) > 0 {
		counsellor["last_name"] = body["last_name"]
	}
	if len(body["gender"]) > 0 {
		counsellor["gender"] = body["gender"]
	}
	if len(body["photo"]) > 0 {
		counsellor["photo"] = body["photo"]
	}
	if len(body["price"]) > 0 {
		counsellor["price"] = body["price"]
	}
	if len(body["multiple_sessions"]) > 0 {
		counsellor["multiple_sessions"] = body["multiple_sessions"]
	}
	if len(body["price_3"]) > 0 {
		counsellor["price_3"] = body["price_3"]
	}
	if len(body["price_5"]) > 0 {
		counsellor["price_5"] = body["price_5"]
	}
	if len(body["education"]) > 0 {
		counsellor["education"] = body["education"]
	}
	if len(body["experience"]) > 0 {
		counsellor["experience"] = body["experience"]
	}
	if len(body["about"]) > 0 {
		counsellor["about"] = body["about"]
	}
	if len(body["resume"]) > 0 {
		counsellor["resume"] = body["resume"]
	}
	if len(body["certificate"]) > 0 {
		counsellor["certificate"] = body["certificate"]
	}
	if len(body["aadhar"]) > 0 {
		counsellor["aadhar"] = body["aadhar"]
	}
	if len(body["linkedin"]) > 0 {
		counsellor["linkedin"] = body["linkedin"]
	}
	if len(body["device_id"]) > 0 {
		counsellor["device_id"] = body["device_id"]
	}
	if len(body["timezone"]) > 0 {
		counsellor["timezone"] = body["timezone"]
	}
	if len(body["payout_percentage"]) > 0 {
		counsellor["payout_percentage"] = body["payout_percentage"]
	}
	if len(body["bank_account_no"]) > 0 {
		counsellor["bank_account_no"] = body["bank_account_no"]
	}
	if len(body["ifsc"]) > 0 {
		counsellor["ifsc"] = body["ifsc"]
	}
	if len(body["payee_name"]) > 0 {
		counsellor["payee_name"] = body["payee_name"]
	}
	if len(body["branch_name"]) > 0 {
		counsellor["branch_name"] = body["branch_name"]
	}
	if len(body["bank_name"]) > 0 {
		counsellor["bank_name"] = body["bank_name"]
	}
	if len(body["bank_account_type"]) > 0 {
		counsellor["bank_account_type"] = body["bank_account_type"]
	}
	if len(body["pan"]) > 0 {
		counsellor["pan"] = body["pan"]
	}

	counsellor["last_login_time"] = UTIL.GetCurrentTime().String()
	counsellor["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")}, counsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update languages, topics to counsellor
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], r.FormValue("counsellor_id"))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
