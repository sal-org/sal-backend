package counsellor

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
// @Tags Counsellor Profile
// @Summary Get counsellor profile with email, if signed up already
// @Router /counsellor [get]
// @Param email query string true "Email of counsellor - to get details, if signed up already"
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"email": r.FormValue("email")})
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

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

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

	// add counsellor details
	counsellor := map[string]string{}
	counsellor["first_name"] = body["first_name"]
	counsellor["last_name"] = body["last_name"]
	counsellor["gender"] = body["gender"]
	counsellor["phone"] = body["phone"]
	counsellor["photo"] = body["photo"]
	counsellor["email"] = body["email"]
	counsellor["price"] = body["price"]
	counsellor["price_3"] = body["price_3"]
	counsellor["price_5"] = body["price_5"]
	counsellor["education"] = body["education"]
	counsellor["experience"] = body["experience"]
	counsellor["about"] = body["about"]
	counsellor["resume"] = body["resume"]
	counsellor["certificate"] = body["certificate"]
	counsellor["aadhar"] = body["aadhar"]
	counsellor["linkedin"] = body["linkedin"]
	counsellor["device_id"] = body["device_id"]
	counsellor["status"] = CONSTANT.CounsellorNotApproved
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

	// add to availability - with 0 (not available)
	for i := 0; i < 7; i++ { // for 7 days of week
		DB.InsertSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": counsellorID, "weekday": strconv.Itoa(i)})
	}

	response["counsellor_id"] = counsellorID

	// send notification
	UTIL.SendNotification(CONSTANT.CounsellorAccountSignupHeading, CONSTANT.CounsellorAccountSignupContent, body["device_id"])

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Counsellor Profile
// @Summary Update counsellor profile details
// @Router /counsellor [put]
// @Param counsellor_id query string true "Counsellor ID to update details"
// @Param body body model.CounsellorProfileUpdateRequest true "Request Body"
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
