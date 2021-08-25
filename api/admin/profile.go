package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
)

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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminProfileAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if admin already exists with specified username
	if DB.CheckIfExists(CONSTANT.AdminsTable, map[string]string{"username": body["username"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add admin details
	admin := map[string]string{}
	admin["username"] = body["username"]
	admin["password"] = UTIL.GetStringMD5Hash(body["password"])
	admin["type"] = body["type"]
	admin["status"] = CONSTANT.AdminActive
	admin["created_at"] = UTIL.GetCurrentTime().String()
	adminID, status, ok := DB.InsertWithUniqueID(CONSTANT.AdminsTable, CONSTANT.AdminDigits, admin, "admin_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["admin_id"] = adminID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update admin details
	admin := map[string]string{}
	if len(body["username"]) > 0 {
		admin["username"] = body["username"]
	}
	if len(body["password"]) > 0 {
		admin["password"] = UTIL.GetStringMD5Hash(body["password"])
	}
	if len(body["type"]) > 0 {
		admin["type"] = body["type"]
	}
	if len(body["status"]) > 0 {
		admin["status"] = body["status"]
	}
	admin["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.AdminsTable, map[string]string{"admin_id": r.FormValue("admin_id")}, admin)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
