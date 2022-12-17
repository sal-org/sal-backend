package cron

import (
	"net/http"

	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
)

func NoShowUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// events
	// get listener availability hours
	availability, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventTable + " where status = " + CONSTANT.EventToBeStarted)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for client no shows

	response["availability"] = availability
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
