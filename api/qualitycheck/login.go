package qualitycheck

import (
	"fmt"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.QualityCheckLoginIDAndPassword)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get admin details
	qualitycheck, status, ok := DB.SelectSQL(CONSTANT.QualityCheckTable, []string{"*"}, map[string]string{"username": body["username"], "password": UTIL.GetStringMD5Hash(body["password"])})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(qualitycheck) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountWrongCredentialsMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(qualitycheck) > 0 && !strings.EqualFold(qualitycheck[0]["status"], CONSTANT.AdminActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessTokenForWeb(qualitycheck[0]["qualitycheck_id"])
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshTokenForWeb(qualitycheck[0]["qualitycheck_id"])
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	response["admin"] = qualitycheck[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	fmt.Println("RefreshToken", r.Header, r.URL.Query())
	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("qualitycheck_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if admin id is valid
	if !DB.CheckIfExists(CONSTANT.AdminsTable, map[string]string{"qualitycheck_id": r.FormValue("qualitycheck_id"), "status": CONSTANT.AdminActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessTokenForWeb(r.FormValue("qualitycheck_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
