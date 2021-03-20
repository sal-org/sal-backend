package counsellor

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

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

	// add topics to counsellor
	for _, topicID := range strings.Split(body["topic_ids"], ",") {
		DB.InsertSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": counsellorID, "topic_id": topicID})
	}

	// add languages to counsellor
	for _, languageID := range strings.Split(body["language_ids"], ",") {
		DB.InsertSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": counsellorID, "language_id": languageID})
	}

	// add to availability - with 0 (not available)
	for i := 0; i < 7; i++ { // for 7 days of week
		DB.InsertSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": counsellorID, "weekday": strconv.Itoa(i)})
	}

	response["counsellor_id"] = counsellorID

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
	counsellor["updated_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")}, counsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(body["topic_ids"]) > 0 {
		// first delete all and add topics to counsellor - to update
		DB.DeleteSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")})
		for _, topicID := range strings.Split(body["topic_ids"], ",") {
			DB.InsertSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "topic_id": topicID})
		}
	}

	if len(body["language_ids"]) > 0 {
		// first delete all and add languages to counsellor - to update
		DB.DeleteSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")})
		for _, languageID := range strings.Split(body["language_ids"], ",") {
			DB.InsertSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "language_id": languageID})
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
