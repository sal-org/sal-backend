package admin

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

	// get admin details
	admin, status, ok := DB.SelectSQL(CONSTANT.AdminsTable, []string{"*"}, map[string]string{"username": r.FormValue("username"), "password": UTIL.GetStringMD5Hash(r.FormValue("password"))})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(admin) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountWrongCredentialsMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(admin) > 0 && !strings.EqualFold(admin[0]["status"], CONSTANT.AdminActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, ok := UTIL.CreateAccessToken(admin[0]["admin_id"])
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, ok := UTIL.CreateRefreshToken(admin[0]["admin_id"])
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	response["admin"] = admin[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	fmt.Println("RefreshToken", r.Header, r.URL.Query())
	// check if refresh token is valid, not expired and token user id is same as user id given
	id, ok, access := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))
	if !ok || access || !strings.EqualFold(id, r.FormValue("admin_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if admin id is valid
	if !DB.CheckIfExists(CONSTANT.AdminsTable, map[string]string{"admin_id": r.FormValue("admin_id"), "status": CONSTANT.AdminActive}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate new access token
	accessToken, ok := UTIL.CreateAccessToken(r.FormValue("admin_id"))
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
