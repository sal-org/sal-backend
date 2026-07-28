package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strings"
	"time"

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

		inPersonConnect, status, ok := DB.SelectSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, []string{"*"}, map[string]string{"counsellor_id": therapist[0]["therapist_id"], "status": "1"})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(inPersonConnect) > 0 {
			therapist[0]["in_person_connect"] = inPersonConnect[0]["status"]
		} else {
			therapist[0]["in_person_connect"] = "0"
		}

		// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, therapist[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
		// therapist[0]["photo"] = endPointURL

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken
		response["languages"] = languages
		response["topics"] = topics
		response["therapist"] = therapist[0]
		response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
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

	var response = make(map[string]any)

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TherapistProfileAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := Model.AddTherapistProfileRequest{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"phone": body.Phone}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if phone is verfied by OTP
	if !DB.CheckIfExists(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body.Phone}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.VerifyPhoneRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var typeOfService string

	switch body.CorporateTherapist {
	case "Individual Clients":
		typeOfService = "0"
	case "Corporate Clients":
		typeOfService = "2"
	case "Both":
		typeOfService = "1"
	}

	joinDate := body.StartDate

	currentTime := time.Now()

	nowDate := currentTime.Format("2006-01-02")

	gapYears := body.GapYears

	gapMonths := body.GapMonths

	experience := UTIL.CalculateExperience(joinDate, nowDate, gapYears, gapMonths)

	// add therapist details
	therapist := map[string]string{}
	therapist["first_name"] = body.FirstName
	therapist["last_name"] = body.LastName
	therapist["pronoun"] = body.Pronoun
	therapist["gender"] = body.Gender
	therapist["location"] = body.Location
	therapist["phone"] = body.Phone
	therapist["photo"] = body.Photo
	therapist["email"] = body.Email
	therapist["price"] = body.Price
	therapist["multiple_sessions"] = body.MultipleSessions
	therapist["price_3"] = body.Price3
	therapist["price_5"] = body.Price5
	therapist["education"] = body.Education
	therapist["experience"] = experience
	therapist["start_date"] = joinDate
	therapist["gap_years"] = gapYears
	therapist["gap_months"] = gapMonths
	therapist["therapeutic_approach"] = body.TherapeuticApproach
	therapist["about"] = body.About
	therapist["timezone"] = body.Timezone
	therapist["resume"] = body.Resume
	therapist["certificate"] = body.Certificate
	therapist["aadhar"] = body.Aadhar
	therapist["linkedin"] = body.Linkedin
	therapist["device_id"] = body.DeviceID
	therapist["payout_percentage"] = CONSTANT.CounsellorPayoutPercentageColumns
	therapist["payee_name"] = body.PayeeName
	therapist["bank_account_no"] = body.BankAccountNo
	therapist["ifsc"] = body.IFSC
	therapist["branch_name"] = body.BranchName
	therapist["bank_name"] = body.BankName
	therapist["bank_account_type"] = body.BankAccountType
	therapist["pan"] = body.PAN
	therapist["corporate_therpist"] = typeOfService
	therapist["status"] = CONSTANT.TherapistNotApproved
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
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body.Phone})

	// add languages, topics to therapist
	UTIL.AssociateLanguagesAndTopics(body.TopicIDs, body.LanguageIDs, therapistID)

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
	UTIL.SendNotification(CONSTANT.CounsellorAccountSignupCounsellorHeading, CONSTANT.CounsellorAccountSignupCounsellorContent, therapistID, CONSTANT.TherapistType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, therapistID, "")
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAccountSignupTextMessage,
			map[string]string{
				"###counsellor_name###": body.FirstName,
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		body.Phone,
		UTIL.GetCurrentTime().String(),
		therapistID,
		CONSTANT.InstantSendTextMessage,
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
		Media_URL:            CONFIG.MediaURLInCLOUDFRONT,
		First_Name:           therapist_details[0]["first_name"],
		Last_Name:            therapist_details[0]["last_name"],
		Pronoun:              therapist_details[0]["pronoun"],
		Gender:               therapist_details[0]["gender"],
		Location:             therapist_details[0]["location"],
		Type:                 "Therapist",
		Phone:                therapist_details[0]["phone"],
		Photo:                therapist_details[0]["photo"],
		Email:                therapist_details[0]["email"],
		Education:            therapist_details[0]["education"],
		CounsellingStartDate: UTIL.BuildOnlyDate(therapist_details[0]["start_date"]),
		CounsellingGap:       therapist_details[0]["gap_years"] + "Y" + " " + therapist_details[0]["gap_months"] + "M",
		Experience:           therapist_details[0]["experience"],
		TherapeuticApproach:  therapist_details[0]["therapeutic_approach"],
		About:                therapist_details[0]["about"],
		Resume:               therapist_details[0]["resume"],
		Certificate:          therapist_details[0]["certificate"],
		Aadhar:               therapist_details[0]["aadhar"],
		Linkedin:             therapist_details[0]["linkedin"],
		Status:               therapist_details[0]["status"],
	}

	filepath := "htmlfile/CounsellorProfile.html"

	emailbody := UTIL.GetHTMLTemplateForProfile(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorProfileWaitingForApprovalTitle,
		emailbody,
		CONFIG.OnboardingEmailID, // prod : CONSTANT.AkshayEmailID , dev : CONSTANT.ShivamEmailID
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

	// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, therapist_details[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// therapist_details[0]["photo"] = endPointURL

	response["therapist"] = therapist_details[0]
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

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

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := Model.AddTherapistProfileRequest{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if user already signed up with specified phone
	if DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"phone": body.Phone}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.PhoneExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update therapist details
	therapist := map[string]string{}
	if len(body.FirstName) > 0 {
		therapist["first_name"] = body.FirstName
	}
	if len(body.LastName) > 0 {
		therapist["last_name"] = body.LastName
	}
	if len(body.Pronoun) > 0 {
		therapist["pronoun"] = body.Pronoun
	}
	if len(body.Gender) > 0 {
		therapist["gender"] = body.Gender
	}
	if len(body.Location) > 0 {
		therapist["location"] = body.Location
	}
	if len(body.Photo) > 0 {
		therapist["photo"] = body.Photo
	}
	if len(body.Price) > 0 {
		therapist["price"] = body.Price
	}
	if len(body.MultipleSessions) > 0 {
		therapist["multiple_sessions"] = body.MultipleSessions
	}
	if len(body.Price3) > 0 {
		therapist["price_3"] = body.Price3
	}
	if len(body.Price5) > 0 {
		therapist["price_5"] = body.Price5
	}
	if len(body.Education) > 0 {
		therapist["education"] = body.Education
	}
	if len(body.Experience) > 0 {
		therapist["experience"] = body.Experience
	}
	if len(body.TherapeuticApproach) > 0 {
		therapist["therapeutic_approach"] = body.TherapeuticApproach
	}
	if len(body.About) > 0 {
		therapist["about"] = body.About
	}
	if len(body.Resume) > 0 {
		therapist["resume"] = body.Resume
	}
	if len(body.Certificate) > 0 {
		therapist["certificate"] = body.Certificate
	}
	if len(body.Aadhar) > 0 {
		therapist["aadhar"] = body.Aadhar
	}
	if len(body.Linkedin) > 0 {
		therapist["linkedin"] = body.Linkedin
	}
	if len(body.DeviceID) > 0 {
		therapist["device_id"] = body.DeviceID
	}
	if len(body.Timezone) > 0 {
		therapist["timezone"] = body.Timezone
	}
	if len(body.PayoutPercentage) > 0 {
		therapist["payout_percentage"] = body.PayoutPercentage
	}
	if len(body.BankAccountNo) > 0 {
		therapist["bank_account_no"] = body.BankAccountNo
	}
	if len(body.IFSC) > 0 {
		therapist["ifsc"] = body.IFSC
	}

	if len(body.PayeeName) > 0 {
		therapist["payee_name"] = body.PayeeName
	}

	if len(body.BranchName) > 0 {
		therapist["branch_name"] = body.BranchName
	}
	if len(body.BankAccountType) > 0 {
		therapist["bank_account_type"] = body.BankAccountType
	}
	if len(body.PAN) > 0 {
		therapist["pan"] = body.PAN
	}

	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid request", CONSTANT.ShowDialog, response)
		return
	}

	therapist["last_login_time"] = UTIL.GetCurrentTime().String()
	therapist["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("therapist_id")}, therapist)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update languages, topics to therapist
	UTIL.AssociateLanguagesAndTopics(body.TopicIDs, body.LanguageIDs, r.FormValue("therapist_id"))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
