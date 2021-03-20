package listener

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

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

	// add topics to listener
	for _, topicID := range strings.Split(body["topic_ids"], ",") {
		DB.InsertSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": listenerID, "topic_id": topicID})
	}

	// add languages to listener
	for _, languageID := range strings.Split(body["language_ids"], ",") {
		DB.InsertSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": listenerID, "language_id": languageID})
	}

	// add to availability - with 0 (not available)
	for i := 0; i < 7; i++ { // for 7 days of week
		DB.InsertSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": listenerID, "weekday": strconv.Itoa(i)})
	}

	response["listener_id"] = listenerID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ProfileUpdate godoc
// @Tags Listener Profile
// @Summary Update listener profile details
// @Router /listener [put]
// @Param listener_id query string true "Listener ID to update details"
// @Param body body model.ListenerProfileUpdateRequest true "Request Body"
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
	if len(body["price"]) > 0 {
		listener["price"] = body["price"]
	}
	if len(body["price_3"]) > 0 {
		listener["price_3"] = body["price_3"]
	}
	if len(body["price_5"]) > 0 {
		listener["price_5"] = body["price_5"]
	}
	if len(body["education"]) > 0 {
		listener["education"] = body["education"]
	}
	if len(body["experience"]) > 0 {
		listener["experience"] = body["experience"]
	}
	if len(body["about"]) > 0 {
		listener["about"] = body["about"]
	}
	listener["updated_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("listener_id")}, listener)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(body["topic_ids"]) > 0 {
		// first delete all and add topics to listener - to update
		DB.DeleteSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": r.FormValue("listener_id")})
		for _, topicID := range strings.Split(body["topic_ids"], ",") {
			DB.InsertSQL(CONSTANT.CounsellorTopicsTable, map[string]string{"counsellor_id": r.FormValue("listener_id"), "topic_id": topicID})
		}
	}

	if len(body["language_ids"]) > 0 {
		// first delete all and add languages to listener - to update
		DB.DeleteSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": r.FormValue("listener_id")})
		for _, languageID := range strings.Split(body["language_ids"], ",") {
			DB.InsertSQL(CONSTANT.CounsellorLanguagesTable, map[string]string{"counsellor_id": r.FormValue("listener_id"), "language_id": languageID})
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
