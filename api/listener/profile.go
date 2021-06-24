package listener

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

// ProfileGet godoc
// @Tags Listener Profile
// @Summary Get listener profile with email, if signed up already
// @Router /listener [get]
// @Param email query string true "Email of listener - to get details, if signed up already"
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get listener details
	listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"email": r.FormValue("email")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(listener) > 0 {
		// listener already signed up
		// check if listener is active
		if !strings.EqualFold(listener[0]["status"], CONSTANT.ListenerActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}

		// generate access and refresh token
		// access token - jwt token with short expiry added in header for authorization
		// refresh token - jwt token with long expiry to get new access token if expired
		// if refresh token expired, need to login
		accessToken, ok := UTIL.CreateAccessToken(listener[0]["listener_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		refreshToken, ok := UTIL.CreateRefreshToken(listener[0]["listener_id"])
		if !ok {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

		response["listener"] = listener[0]
		response["media_url"] = CONFIG.MediaURL
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileAdd godoc
// @Tags Listener Profile
// @Summary Add listener profile after OTP verified to signup
// @Router /listener [post]
// @Param body body model.ListenerProfileAddRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ListenerProfileAddRequiredFields)
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

	// add listener details
	listener := map[string]string{}
	listener["first_name"] = body["first_name"]
	listener["last_name"] = body["last_name"]
	listener["gender"] = body["gender"]
	listener["phone"] = body["phone"]
	listener["photo"] = body["photo"]
	listener["email"] = body["email"]
	listener["occupation"] = body["occupation"]
	listener["experience"] = body["experience"]
	listener["about"] = body["about"]
	listener["device_id"] = body["device_id"]
	listener["status"] = CONSTANT.ListenerNotApproved
	listener["created_at"] = UTIL.GetCurrentTime().String()
	listenerID, status, ok := DB.InsertWithUniqueID(CONSTANT.ListenersTable, CONSTANT.ListenerDigits, listener, "listener_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// using phone verified table to check if phone has been really verified by OTP
	// currently deleting if phone number is already present
	DB.DeleteSQL(CONSTANT.PhoneOTPVerifiedTable, map[string]string{"phone": body["phone"]})

	// add languages, topics to listener
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], listenerID)

	// add to availability - with 0 (not available)
	for i := 0; i < 7; i++ { // for 7 days of week
		DB.InsertSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": listenerID, "weekday": strconv.Itoa(i)})
	}

	response["listener_id"] = listenerID

	// send account signup notification to listener
	UTIL.SendNotification(CONSTANT.CounsellorAccountSignupCounsellorHeading, CONSTANT.CounsellorAccountSignupCounsellorContent, listenerID, CONSTANT.ListenerType)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Listener Profile
// @Summary Update listener profile details
// @Router /listener [put]
// @Param listener_id query string true "Listener ID to update details"
// @Param body body model.ListenerProfileUpdateRequest true "Request Body"
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

	// update listener details
	listener := map[string]string{}
	if len(body["first_name"]) > 0 {
		listener["first_name"] = body["first_name"]
	}
	if len(body["last_name"]) > 0 {
		listener["last_name"] = body["last_name"]
	}
	if len(body["gender"]) > 0 {
		listener["gender"] = body["gender"]
	}
	if len(body["photo"]) > 0 {
		listener["photo"] = body["photo"]
	}
	if len(body["occupation"]) > 0 {
		listener["occupation"] = body["occupation"]
	}
	if len(body["experience"]) > 0 {
		listener["experience"] = body["experience"]
	}
	if len(body["about"]) > 0 {
		listener["about"] = body["about"]
	}
	if len(body["device_id"]) > 0 {
		listener["device_id"] = body["device_id"]
	}
	listener["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("listener_id")}, listener)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update languages, topics to listener
	UTIL.AssociateLanguagesAndTopics(body["topic_ids"], body["language_ids"], r.FormValue("listener_id"))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
