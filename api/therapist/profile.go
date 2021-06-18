package therapist

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
// @Tags Therapist Profile
// @Summary Get therapist profile with email, if signed up already
// @Router /therapist [get]
// @Param email query string true "Email of therapist - to get details, if signed up already"
// @Produce json
// @Success 200
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"email": r.FormValue("email")})
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

		response["access_token"] = accessToken
		response["refresh_token"] = refreshToken

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

	// add therapist details
	therapist := map[string]string{}
	therapist["first_name"] = body["first_name"]
	therapist["last_name"] = body["last_name"]
	therapist["gender"] = body["gender"]
	therapist["phone"] = body["phone"]
	therapist["photo"] = body["photo"]
	therapist["email"] = body["email"]
	therapist["price"] = body["price"]
	therapist["price_3"] = body["price_3"]
	therapist["price_5"] = body["price_5"]
	therapist["education"] = body["education"]
	therapist["experience"] = body["experience"]
	therapist["about"] = body["about"]
	therapist["resume"] = body["resume"]
	therapist["certificate"] = body["certificate"]
	therapist["aadhar"] = body["aadhar"]
	therapist["linkedin"] = body["linkedin"]
	therapist["device_id"] = body["device_id"]
	therapist["status"] = CONSTANT.TherapistNotApproved
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

	// add to availability - with 0 (not available)
	for i := 0; i < 7; i++ { // for 7 days of week
		DB.InsertSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": therapistID, "weekday": strconv.Itoa(i)})
	}

	response["therapist_id"] = therapistID

	// send account signup notification to therapist
	UTIL.SendNotification(CONSTANT.TherapistAccountSignupTherapistHeading, CONSTANT.TherapistAccountSignupTherapistContent, body["device_id"])

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Therapist Profile
// @Summary Update therapist profile details
// @Router /therapist [put]
// @Param therapist_id query string true "Therapist ID to update details"
// @Param body body model.TherapistProfileUpdateRequest true "Request Body"
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
