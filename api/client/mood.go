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

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

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
	moodResultID, _, ok := DB.InsertWithUniqueID(CONSTANT.MoodResultsTable, CONSTANT.MoodResultsDigits, map[string]string{
		"client_id":  body["client_id"],
		"name":       body["name"],
		"age":        body["age"],
		"gender":     body["gender"],
		"phone":      body["phone"],
		"mood_id":    body["mood_id"],
		"notes":      body["notes"],
		"date":       body["date"],
		"status":     CONSTANT.MoodResultActive,
		"created_at": UTIL.GetCurrentTime().UTC().String(),
	}, "mood_result_id")
	if !ok {
		status, ok := DB.UpdateSQL(CONSTANT.MoodResultsTable, map[string]string{"client_id": body["client_id"], "date": body["date"]}, map[string]string{
			"mood_id":     body["mood_id"],
			"notes":       body["notes"],
			"modified_at": UTIL.GetCurrentTime().UTC().String(),
		})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// UTIL.SetReponse(w, status, CONSTANT.MoodAlreadyAddedMessage, CONSTANT.ShowDialog, response)
		// return
	}

	if !(body["mood_id"] == "1") {

		moodTitle, status, ok := DB.SelectProcess("select title from "+CONSTANT.MoodsTable+" where id = ?", body["mood_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		UTIL.SendNotification(
			CONSTANT.ClientSelectSadMoodHeading,
			UTIL.ReplaceNotificationContentInString(
				CONSTANT.ClientSelectSadMoodContent,
				map[string]string{
					"###mood###": moodTitle[0]["title"],
				},
			),
			body["client_id"],
			CONSTANT.ClientType,
			UTIL.GetCurrentTime().String(),
			CONSTANT.NotificationSent,
			moodResultID,
		)

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

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	dates := strings.Split(r.FormValue("dates"), ",")
	if len(dates) < 2 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// get mood past results
	moodResults, status, ok := DB.SelectProcess("select mood_id, `date`, notes from "+CONSTANT.MoodResultsTable+" where client_id = ? and status = "+CONSTANT.MoodResultActive+" and `date` >= ? and `date` <= ? order by `date` asc", r.FormValue("client_id"), dates[0], dates[1])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["mood_results"] = moodResults
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// MoodContent godoc
// @Tags Client Mood
// @Summary Get mood content
// @Router /client/mood/content [get]
// @Param user_id query string true "Logged in user ID to get mood liked content"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListMoodContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	mood_id, status, ok := DB.SelectProcess("select * from "+CONSTANT.MoodResultsTable+" where client_id = ?  and status = 1 order by created_at desc limit 20", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	contents, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where mood_id = ? and training = 0 and status = 1 order by created_at desc limit 20", mood_id[0]["mood_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get liked content ids
	contentLiked, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["mood_content"] = contents
	response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(contentLiked, "content_id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
