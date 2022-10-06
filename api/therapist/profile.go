package therapist

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
// @Tags Therapist Profile
// @Summary Get therapist profile with email, if signed up already
// @Router /therapist [get]
// @Param email query string false "Email of therapist - to get details, if signed up already"
// @Param therapist_id query string false "Therapist ID to update details"
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

	// get therapist details
	params := map[string]string{}
	if len(r.FormValue("email")) > 0 {
		params["email"] = r.FormValue("email")
	}
	if len(r.FormValue("therapist_id")) > 0 {
		params["therapist_id"] = r.FormValue("therapist_id")
	}
	if len(params) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, params)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) > 0 {
		// therapist already signed up
		// check if therapist is active
		if !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		// generate access and refresh token
		// access token - jwt token with short expiry added in header for authorization
		// refresh token - jwt token with long expiry to get new access token if expired
		// if refresh token expired, need to login
		accessToken, ok := UTIL.CreateAccessToken(therapist[0]["therapist_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		refreshToken, ok := UTIL.CreateRefreshToken(therapist[0]["therapist_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}

		languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", therapist[0]["therapist_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor topics
		topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", therapist[0]["therapist_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken
		response["languages"] = languages
		response["topics"] = topics
		response["therapist"] = therapist[0]
		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileAdd godoc
// @Tags Therapist Profile
// @Summary Add therapist profile after OTP verified to signup
// @Router /therapist [post]
// @Param body body model.TherapistProfileAddRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TherapistProfileAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if phone is verfied by OTP
	if !DB.CheckIfExists(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.VerifyPhoneRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add therapist details
	therapist := map[string]string{}
	therapist["first_name"] = body["first_name"]
	therapist["last_name"] = body["last_name"]
	therapist["gender"] = body["gender"]
	therapist["phone"] = body["phone"]
	therapist["photo"] = body["photo"]
	therapist["email"] = body["email"]
	therapist["price"] = body["price"]
	therapist["multiple_sessions"] = body["multiple_sessions"]
	therapist["price_3"] = body["price_3"]
	therapist["price_5"] = body["price_5"]
	therapist["education"] = body["education"]
	therapist["experience"] = body["experience"]
	therapist["about"] = body["about"]
	therapist["timezone"] = body["timezone"]
	therapist["resume"] = body["resume"]
	therapist["certificate"] = body["certificate"]
	therapist["aadhar"] = body["aadhar"]
	therapist["linkedin"] = body["linkedin"]
	therapist["device_id"] = body["device_id"]
	therapist["payout_percentage"] = CONSTANT.CounsellorPayoutPercentageColumns
	therapist["payee_name"] = body["payee_name"]
	therapist["bank_account_no"] = body["bank_account_no"]
	therapist["ifsc"] = body["ifsc"]
	therapist["branch_name"] = body["branch_name"]
	therapist["bank_name"] = body["bank_name"]
	therapist["bank_account_type"] = body["bank_account_type"]
	therapist["pan"] = body["pan"]
	therapist["status"] = CONSTANT.TherapistActive
	therapist["notification_status"] = CONSTANT.NotificationActive
	therapist["last_login_time"] = UTIL.GetCurrentTime().String()
	therapist["created_at"] = UTIL.GetCurrentTime().String()
	therapistID, status, ok := DB.InsertWithUniqueID(CONSTANT.TherapistsTable, CONSTANT.TherapistDigits, therapist, "therapist_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// using phone verified table to check if phone has been really verified by OTP
	// currently deleting if phone number is already present
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]})

	// add languages, topics to therapist
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], therapistID)

	// not available for next 30 days. change here when you change in add new slot cron
	for i := 0; i < 30; i++ {
		DB.InsertSQL(CONSTANT.SlotsTable, map[string]string{"counsellor_id": therapistID, "date": UTIL.GetCurrentTime().AddDate(0, 0, i).Format("2006-01-02")})
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessToken(therapistID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshToken(therapistID)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send account signup notification, message to therapist
	UTIL.SendNotification(CONSTANT.CounsellorAccountSignupCounsellorHeading, CONSTANT.CounsellorAccountSignupCounsellorContent, therapistID, CONSTANT.TherapistType, UTIL.GetCurrentTime().String(), therapistID)
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAccountSignupTextMessage,
			map[string]string{
				"###counsellor_name###": body["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		body["phone"],
		CONSTANT.LaterSendTextMessage,
	)

	/*orderdetails, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "gender", "phone", "photo", "email", "education", "experience", "about", "resume", "certificate", "aadhar", "linkedin", "status"}, map[string]string{"therapist_id": therapistID})
	counsellorbody := UTIL.GetHTMLTemplateForCounsellor(orderdetails)

	UTIL.SendEmail(
		CONSTANT.CounsellorProfileWaitingForApprovalTitle,
		counsellorbody,
		CONSTANT.AnandEmailID,
		CONSTANT.InstantSendEmailMessage,
	)*/

	therapist_details, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": therapistID})

	// counsellor_name := Model.CounsellorProfileSendEmailTextMessage{
	// 	First_Name: therapist_details[0]["first_name"],
	// }

	// filepath_text := "htmlfile/Counsellor_Profile_Text_Message.html"

	// emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(counsellor_name, filepath_text)

	// UTIL.SendEmail(
	// 	CONSTANT.CounsellorProfileWaitingForApprovalTitle,
	// 	emailBody,
	// 	therapist_details[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	data := Model.EmailDataForCounsellorProfile{
		First_Name:  therapist_details[0]["first_name"],
		Last_Name:   therapist_details[0]["last_name"],
		Gender:      therapist_details[0]["gender"],
		Type:        "Counsellor",
		Phone:       therapist_details[0]["phone"],
		Photo:       therapist_details[0]["photo"],
		Email:       therapist_details[0]["email"],
		Education:   therapist_details[0]["education"],
		Experience:  therapist_details[0]["experience"],
		About:       therapist_details[0]["about"],
		Resume:      therapist_details[0]["resume"],
		Certificate: therapist_details[0]["certificate"],
		Aadhar:      therapist_details[0]["aadhar"],
		Linkedin:    therapist_details[0]["linkedin"],
		Status:      therapist_details[0]["status"],
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
				"###First_Name###":  orderdetails[0]["first_name"],
				"###Last_Name###":   orderdetails[0]["last_name"],
				"###Gender###":      orderdetails[0]["gender"],
				"###Phone###":       orderdetails[0]["phone"],
				"###Email###":       orderdetails[0]["email"],
				"###Photo###":       orderdetails[0]["photo"],
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

	response["therapist"] = therapist_details[0]
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken
	response["media_url"] = CONFIG.MediaURL

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Therapist Profile
// @Summary Update therapist profile details
// @Router /therapist [put]
// @Param therapist_id query string true "Therapist ID to update details"
// @Param body body model.TherapistProfileUpdateRequest true "Request Body"
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

	// update therapist details
	therapist := map[string]string{}
	if len(body["first_name"]) > 0 {
		therapist["first_name"] = body["first_name"]
	}
	if len(body["last_name"]) > 0 {
		therapist["last_name"] = body["last_name"]
	}
	if len(body["gender"]) > 0 {
		therapist["gender"] = body["gender"]
	}
	if len(body["photo"]) > 0 {
		therapist["photo"] = body["photo"]
	}
	if len(body["price"]) > 0 {
		therapist["price"] = body["price"]
	}
	if len(body["multiple_sessions"]) > 0 {
		therapist["multiple_sessions"] = body["multiple_sessions"]
	}
	if len(body["price_3"]) > 0 {
		therapist["price_3"] = body["price_3"]
	}
	if len(body["price_5"]) > 0 {
		therapist["price_5"] = body["price_5"]
	}
	if len(body["education"]) > 0 {
		therapist["education"] = body["education"]
	}
	if len(body["experience"]) > 0 {
		therapist["experience"] = body["experience"]
	}
	if len(body["about"]) > 0 {
		therapist["about"] = body["about"]
	}
	if len(body["resume"]) > 0 {
		therapist["resume"] = body["resume"]
	}
	if len(body["certificate"]) > 0 {
		therapist["certificate"] = body["certificate"]
	}
	if len(body["aadhar"]) > 0 {
		therapist["aadhar"] = body["aadhar"]
	}
	if len(body["linkedin"]) > 0 {
		therapist["linkedin"] = body["linkedin"]
	}
	if len(body["device_id"]) > 0 {
		therapist["device_id"] = body["device_id"]
	}
	if len(body["timezone"]) > 0 {
		therapist["timezone"] = body["timezone"]
	}
	if len(body["timezone"]) > 0 {
		therapist["timezone"] = body["timezone"]
	}
	if len(body["payout_percentage"]) > 0 {
		therapist["payout_percentage"] = body["payout_percentage"]
	}
	if len(body["bank_account_no"]) > 0 {
		therapist["bank_account_no"] = body["bank_account_no"]
	}
	if len(body["ifsc"]) > 0 {
		therapist["ifsc"] = body["ifsc"]
	}
	if len(body["branch_name"]) > 0 {
		therapist["branch_name"] = body["branch_name"]
	}
	if len(body["bank_name"]) > 0 {
		therapist["bank_name"] = body["bank_name"]
	}
	if len(body["bank_account_type"]) > 0 {
		therapist["bank_account_type"] = body["bank_account_type"]
	}
	if len(body["pan"]) > 0 {
		therapist["pan"] = body["pan"]
	}

	therapist["last_login_time"] = UTIL.GetCurrentTime().String()
	therapist["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("therapist_id")}, therapist)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update languages, topics to therapist
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], r.FormValue("therapist_id"))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
