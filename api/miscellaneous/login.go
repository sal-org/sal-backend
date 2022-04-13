package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

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
	if len(counsellor) > 0 && !strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get listener details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 && !strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
		if len(counsellor) > 0 && !strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
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

	// this for testing
	if !strings.EqualFold("4444", r.FormValue("otp")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

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
		if !strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
			if !strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerActive) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
			if !strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistActive) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
