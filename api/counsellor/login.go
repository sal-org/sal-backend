package counsellor

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strings"
)

// RefreshToken godoc
// @Tags Counsellor Login
// @Summary Get new access token with refresh token
// @Router /counsellor/refresh-token [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("counsellor_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if counsellor id is valid
	if !DB.CheckIfExists(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "status": CONSTANT.CounsellorActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessToken(r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
