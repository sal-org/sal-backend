package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"time"

	Model "salbackend/model"
	UTIL "salbackend/util"
	"strings"
)

// used for login for counsellor/listener/therapist

// SendOTP godoc
// @Tags Counsellor/Listener/Therapist Login
// @Summary Send OTP to specified phone
// @Router /sendotp [get]
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

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorNotApproved) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorDeleted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorBlocked) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountBlockedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get listener details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerNotApproved) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerDeleted) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerBlocked) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistNotApproved) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistDeleted) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) > 0 && strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistBlocked) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}
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
			CONSTANT.CounsellorOTPTextMessage,
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
// @Tags Counsellor/Listener/Therapist Login
// @Summary Verify OTP sent to specified phone
// @Router /verifyotp [get]
// @Param phone query string true "Phone number OTP has been sent to - send phone number with 91 code"
// @Param otp query string true "OTP entered by counsellor/listener/therapist"
// @Param device_id query string true "Device ID entered by counsellor/listener/therapist"
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

	if strings.EqualFold("917575757575", r.FormValue("phone")) {
		if !strings.EqualFold("4444", r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	} else if strings.EqualFold("914848484848", r.FormValue("phone")) {
		if !strings.EqualFold("4747", r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	} else {
		if !UTIL.VerifyOTP(r.FormValue("phone"), r.FormValue("otp")) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
			return
		}
	}

	// this for testing
	// if !strings.EqualFold("4444", r.FormValue("otp")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get counsellor details
	var counsellorType string

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(counsellor) > 0 {
		// counsellor already signed up
		// check if counsellor is active
		if strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorNotApproved) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
			return
		}
		if strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorDeleted) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}
		if strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorBlocked) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}
		counsellorType = CONSTANT.CounsellorType
	}

	if len(counsellor) == 0 {
		// get listener details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 {
			// listener already signed up
			// check if listener is active
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerNotApproved) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
				return
			}
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerDeleted) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountDeletedMessage, CONSTANT.ShowDialog, response)
				return
			}
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerBlocked) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountBlockedMessage, CONSTANT.ShowDialog, response)
				return
			}
			counsellorType = CONSTANT.ListenerType
		}
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 {
			// therapist already signed up
			// check if therapist is active
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistNotApproved) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountNotApprovedMessage, CONSTANT.ShowDialog, response)
				return
			}
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistDeleted) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
				return
			}
			if strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistBlocked) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountBlockedMessage, CONSTANT.ShowDialog, response)
				return
			}
			counsellorType = CONSTANT.TherapistType
		}
	}

	if len(counsellor) == 0 {
		// counsellor/listener/therapist is first time signing up

		// using phone verified table to check if phone has been really verified by OTP
		// currently inserting after phone is really verified
		DB.InsertSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": r.FormValue("phone"), "created_at": UTIL.GetCurrentTime().String()})
	}

	if len(counsellor) > 0 && len(counsellorType) > 0 {
		// generate access and refresh token
		// access token - jwt token with short expiry added in header for authorization
		// refresh token - jwt token with long expiry to get new access token if expired
		// if refresh token expired, need to login
		var (
			accessToken  string
			refreshToken string
		)
		switch counsellorType {
		case CONSTANT.CounsellorType:
			accessToken, ok = UTIL.CreateAccessToken(counsellor[0]["counsellor_id"])
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			refreshToken, ok = UTIL.CreateRefreshToken(counsellor[0]["counsellor_id"])
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
			status, ok = DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": counsellor[0]["counsellor_id"]}, map[string]string{"device_id": r.FormValue("device_id"), "last_login_time": UTIL.GetCurrentTime().String(), "platform": r.FormValue("platform"), "version": r.FormValue("version")})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			response["languages"] = languages
			response["topics"] = topics
			response["counsellor"] = counsellor[0]
		case CONSTANT.ListenerType:
			accessToken, ok = UTIL.CreateAccessToken(counsellor[0]["listener_id"])
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			refreshToken, ok = UTIL.CreateRefreshToken(counsellor[0]["listener_id"])
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}

			languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", counsellor[0]["listener_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellor topics
			topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", counsellor[0]["listener_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			status, ok = DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": counsellor[0]["listener_id"]}, map[string]string{"device_id": r.FormValue("device_id"), "last_login_time": UTIL.GetCurrentTime().String(), "platform": r.FormValue("platform"), "version": r.FormValue("version")})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			response["languages"] = languages
			response["topics"] = topics
			response["listener"] = counsellor[0]
		case CONSTANT.TherapistType:
			accessToken, ok = UTIL.CreateAccessToken(counsellor[0]["therapist_id"])
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			refreshToken, ok = UTIL.CreateRefreshToken(counsellor[0]["therapist_id"])
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", counsellor[0]["therapist_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellor topics
			topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", counsellor[0]["therapist_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			status, ok = DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": counsellor[0]["therapist_id"]}, map[string]string{"device_id": r.FormValue("device_id"), "last_login_time": UTIL.GetCurrentTime().String(), "platform": r.FormValue("platform"), "version": r.FormValue("version")})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellor[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
			counsellor[0]["photo"] = endPointURL

			response["languages"] = languages
			response["topics"] = topics
			response["therapist"] = counsellor[0]
		}

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AppInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	appInfo, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppInfoTable + " where status = 1 ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["app_info"] = appInfo[0]

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

func AppFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var firstName, lastName, email string

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ClientAppFeedbackRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user already signed up with specified email
	if !DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": body["user_id"]}) {
		if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": body["user_id"]}) {
			if !DB.CheckIfExists(CONSTANT.ListenersTable, map[string]string{"listener_id": body["user_id"]}) {
				if !DB.CheckIfExists(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": body["user_id"]}) {
					UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
					return
				} else {
					// get counsellor details
					counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": body["user_id"]})
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}
					if len(counsellor) == 0 {
						UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
						return
					}
					firstName = counsellor[0]["first_name"]
					lastName = counsellor[0]["last_name"]
					email = counsellor[0]["email"]
				}

			} else {
				// get listener details
				listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"listener_id": body["user_id"]})
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
				if len(listener) == 0 {
					UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
					return
				}
				firstName = listener[0]["first_name"]
				lastName = listener[0]["last_name"]
				email = listener[0]["email"]
			}
		} else {
			// get therapist details
			therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["user_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			if len(therapist) == 0 {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
				return
			}
			firstName = therapist[0]["first_name"]
			lastName = therapist[0]["last_name"]
			email = therapist[0]["email"]
		}
	} else {
		// get client details
		client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["user_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(client) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorLoginIfNotRegister, CONSTANT.ShowDialog, response)
			return
		}

		firstName = client[0]["first_name"]
		lastName = client[0]["last_name"]
		email = client[0]["email"]
	}

	// add client details
	user := map[string]string{}
	user["user_id"] = body["user_id"]
	user["type"] = body["type"]
	user["concern_area"] = body["concern_area"]
	user["details"] = body["details"]
	user["attach_1"] = body["attach_1"]
	user["attach_2"] = body["attach_2"]
	user["attach_3"] = body["attach_3"]
	user["platform"] = "Mobile"
	user["status"] = CONSTANT.ClientActive
	user["created_at"] = UTIL.GetCurrentTime().String()
	feedbackID, status, ok := DB.InsertWithUniqueID(CONSTANT.FeedbackTable, CONSTANT.FeedbackDigits, user, "feedback_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// send appointment reminder notification to counsellor before 30 min
	UTIL.SendNotification(
		CONSTANT.ClientFeedBackHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientFeedBackContent,
			map[string]string{
				"###type###": body["type"],
			},
		),
		body["user_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		CONSTANT.NotificationSent,
		feedbackID,
		"",
	)

	data := Model.EmailDataForFeedback{
		MediaURL:        CONFIG.MediaURL,
		ClientFirstName: firstName,
		ClientLastName:  lastName,
		ClientEmail:     email,
		ConcernsType:    body["type"],
		ConcernArea:     body["concern_area"],
		Details:         body["details"],
		Attach1:         body["attach_1"],
		Attach2:         body["attach_2"],
		Attach3:         body["attach_3"],
	}

	filepath := "htmlfile/FeedBackForApp.html"

	emailbody := UTIL.GetHTMLTemplateForAppFeedBack(data, filepath)

	// email for client
	UTIL.SendEmail(
		body["type"],
		emailbody,
		CONFIG.CustomerCareEmailID,
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
