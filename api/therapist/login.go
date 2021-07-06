package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strings"
)

// SendOTP godoc
// @Tags Therapist Login
// @Summary Send OTP to specified phone
// @Router /therapist/sendotp [get]
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

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"status"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) > 0 && !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
// @Tags Therapist Login
// @Summary Verify OTP sent to specified phone
// @Router /therapist/verifyotp [get]
// @Param phone query string true "Phone number OTP has been sent to - send phone number with 91 code"
// @Param otp query string true "OTP entered by therapist"
// @Produce json
// @Success 200
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if otp is correct
	if !UTIL.VerifyOTP(r.FormValue("phone"), r.FormValue("otp")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectOTPRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) == 0 {
		// therapist is first time signing up

		// using phone verified table to check if phone has been really verified by OTP
		// currently inserting after phone is really verified
		DB.InsertSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": r.FormValue("phone"), "created_at": UTIL.GetCurrentTime().String()})
	} else {
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

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

		response["therapist"] = therapist[0]
		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// RefreshToken godoc
// @Tags Therapist Login
// @Summary Get new access token with refresh token
// @Router /therapist/refresh-token [get]
// @Param therapist_id query string true "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("therapist_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if therapist id is valid
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("therapist_id"), "status": CONSTANT.TherapistActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessToken(r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
