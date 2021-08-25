package listener

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strings"
)

// RefreshToken godoc
// @Tags Listener Login
// @Summary Get new access token with refresh token
// @Router /listener/refresh-token [get]
// @Param listener_id query string true "Logged in listener ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("listener_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if listener id is valid
	if !DB.CheckIfExists(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("listener_id"), "status": CONSTANT.ListenerActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessToken(r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
