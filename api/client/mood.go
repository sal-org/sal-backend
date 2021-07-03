package client

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

	UTIL "salbackend/util"
)

// MoodAdd godoc
// @Tags Client Mood
// @Summary Add client mood
// @Router /client/mood [post]
// @Param body body model.MoodAddRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func MoodAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.MoodAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// add mood result
	moodResultID, status, ok := DB.InsertWithUniqueID(CONSTANT.MoodResultsTable, CONSTANT.MoodResultsDigits, map[string]string{
		"client_id":  body["client_id"],
		"name":       body["name"],
		"age":        body["age"],
		"gender":     body["gender"],
		"phone":      body["phone"],
		"mood_id":    body["mood_id"],
		"date":       body["date"],
		"status":     CONSTANT.MoodResultActive,
		"created_at": UTIL.GetCurrentTime().UTC().String(),
	}, "mood_result_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["mood_result_id"] = moodResultID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// MoodHistory godoc
// @Tags Client Mood
// @Summary Get mood history
// @Router /client/mood/history [get]
// @Param client_id query string true "Logged in client ID to get mood history"
// @Param dates query string true "Dates to get history - (start,end) - (2021-05-21,2021-06-10)"
// @Security JWTAuth
// @Produce json
// @Success 200
func MoodHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	dates := strings.Split(r.FormValue("dates"), ",")
	if len(dates) < 2 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// get mood past results
	moodResults, status, ok := DB.SelectProcess("select mood_id, `date` from "+CONSTANT.MoodResultsTable+" where client_id = ? and status = "+CONSTANT.MoodResultActive+" and `date` >= ? and `date` <= ? order by `date` asc", r.FormValue("client_id"), dates[0], dates[1])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["mood_results"] = moodResults
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}