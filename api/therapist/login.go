package therapist

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strings"
)

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
