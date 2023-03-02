package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

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

func AddProfileForUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// // read request body
	// body := Model.AddProfileUserRequest{}
	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// defer r.Body.Close()
	// err = json.Unmarshal(b, &body)
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminUserProfileRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// add assessment result
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.RolesTable, CONSTANT.AssessmentResultsDigits, map[string]string{
		"profile_name": body["profile_name"],
		"pc_add":       body["pc_add"],
		"pc_edit":      body["pc_edit"],
		"pc_view":      body["pc_view"],
		"noti_add":     body["noti_add"],
		"noti_edit":    body["noti_edit"],
		"noti_view":    body["noti_view"],
		"cont_add":     body["cont_add"],
		"cont_edit":    body["cont_edit"],
		"cont_view":    body["cont_view"],
		"mq_view":      body["mq_view"],
		"mq_add":       body["mq_add"],
		"mq_edit":      body["mq_edit"],
		"cent_add":     body["cent_add"],
		"cent_edit":    body["cent_edit"],
		"cent_view":    body["cent_view"],
		"coun_view":    body["coun_view"],
		"coun_edit":    body["coun_edit"],
		"coun_add":     body["coun_add"],
		"list_add":     body["list_add"],
		"list_edit":    body["list_edit"],
		"list_view":    body["list_view"],
		"ther_view":    body["ther_view"],
		"ther_add":     body["ther_add"],
		"ther_edit":    body["ther_edit"],
		"appoint_add":  body["appoint_add"],
		"appoint_edit": body["appoint_edit"],
		"appoint_view": body["appoint_view"],
		"cafe_add":     body["cafe_add"],
		"cafe_view":    body["cafe_view"],
		"cafe_edit":    body["cafe_edit"],
		"rept_view":    body["rept_view"],
		"rept_edit":    body["rept_edit"],
		"rept_add":     body["rept_add"],
		"status":       CONSTANT.ClientActive,
		"created_at":   UTIL.GetCurrentTime().UTC().String(),
	}, "role_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UpdateProfileForUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	// body := Model.AddProfileUserRequest{}
	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// defer r.Body.Close()
	// err = json.Unmarshal(b, &body)
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// // check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminProfileAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check if admin already exists with specified username
	if !DB.CheckIfExists(CONSTANT.RolesTable, map[string]string{"role_id": r.FormValue("id")}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AminUserProfileMessage, CONSTANT.ShowDialog, response)
		return
	}

	profile := map[string]string{}
	if len(body["profile_name"]) > 0 {
		profile["profile_name"] = body["profile_name"]
	}

	if len(body["pc_add"]) > 0 {
		profile["pc_add"] = body["pc_add"]
	}

	if len(body["pc_edit"]) > 0 {
		profile["pc_edit"] = body["pc_edit"]
	}

	if len(body["pc_view"]) > 0 {
		profile["pc_view"] = body["pc_view"]
	}

	if len(body["noti_add"]) > 0 {
		profile["noti_add"] = body["noti_add"]
	}

	if len(body["noti_edit"]) > 0 {
		profile["noti_edit"] = body["noti_edit"]
	}

	if len(body["noti_view"]) > 0 {
		profile["noti_view"] = body["noti_view"]
	}

	if len(body["cont_add"]) > 0 {
		profile["cont_add"] = body["cont_add"]
	}

	if len(body["cont_edit"]) > 0 {
		profile["cont_edit"] = body["cont_edit"]
	}

	if len(body["cont_view"]) > 0 {
		profile["cont_view"] = body["cont_view"]
	}

	if len(body["mq_view"]) > 0 {
		profile["mq_view"] = body["mq_view"]
	}

	if len(body["mq_edit"]) > 0 {
		profile["mq_edit"] = body["mq_edit"]
	}

	if len(body["mq_add"]) > 0 {
		profile["mq_add"] = body["mq_add"]
	}

	if len(body["cent_add"]) > 0 {
		profile["cent_add"] = body["cent_add"]
	}

	if len(body["cent_view"]) > 0 {
		profile["cent_view"] = body["cent_view"]
	}

	if len(body["cent_edit"]) > 0 {
		profile["cent_edit"] = body["cent_edit"]
	}

	if len(body["coun_view"]) > 0 {
		profile["coun_view"] = body["coun_view"]
	}

	if len(body["coun_add"]) > 0 {
		profile["coun_add"] = body["coun_add"]
	}

	if len(body["coun_edit"]) > 0 {
		profile["coun_edit"] = body["coun_edit"]
	}

	if len(body["list_view"]) > 0 {
		profile["list_view"] = body["list_view"]
	}

	if len(body["list_add"]) > 0 {
		profile["list_add"] = body["list_add"]
	}

	if len(body["list_edit"]) > 0 {
		profile["list_edit"] = body["list_edit"]
	}

	if len(body["ther_view"]) > 0 {
		profile["ther_view"] = body["ther_view"]
	}

	if len(body["ther_edit"]) > 0 {
		profile["ther_edit"] = body["ther_edit"]
	}

	if len(body["ther_add"]) > 0 {
		profile["ther_add"] = body["ther_add"]
	}

	if len(body["appoint_view"]) > 0 {
		profile["appoint_view"] = body["appoint_view"]
	}

	if len(body["appoint_add"]) > 0 {
		profile["appoint_add"] = body["appoint_add"]
	}

	if len(body["appoint_edit"]) > 0 {
		profile["appoint_edit"] = body["appoint_edit"]
	}

	if len(body["cafe_add"]) > 0 {
		profile["cafe_add"] = body["cafe_add"]
	}

	if len(body["cafe_view"]) > 0 {
		profile["cafe_view"] = body["cafe_view"]
	}

	if len(body["cafe_edit"]) > 0 {
		profile["cafe_edit"] = body["cafe_edit"]
	}

	if len(body["rept_view"]) > 0 {
		profile["rept_view"] = body["rept_view"]
	}

	if len(body["rept_edit"]) > 0 {
		profile["rept_edit"] = body["rept_edit"]
	}

	if len(body["rept_add"]) > 0 {
		profile["rept_add"] = body["rept_add"]
	}

	profile["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.RolesTable, map[string]string{"role_id": r.FormValue("id")}, profile)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// QuoteGet - get quote
func UserProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get quotes
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "role_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " role_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	// selectSQL := ""

	where := " where status = '1' "
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}

	// selectSQL += " where " + strings.Join(wheres, " and ")

	profile, status, ok := DB.SelectProcess("select * from "+CONSTANT.RolesTable+where+" order by id desc limit 20", queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["profiles"] = profile
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// UserProfileDelete - delete user profile
func UserProfileDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	status, ok := DB.DeleteSQL(CONSTANT.RolesTable, map[string]string{"role_id": r.FormValue("role_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AttachPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminUserRoleRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if admin already exists with specified username
	if DB.CheckIfExists(CONSTANT.UsersPermissionTable, map[string]string{"user_name": body["user_name"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add admin main table
	admin := map[string]string{}

	// add roleAttach details
	roleAttach := map[string]string{}
	roleAttach["role_id"] = body["role_id"]
	roleAttach["user_name"] = body["user_name"]
	roleAttach["password"] = UTIL.GetStringMD5Hash(body["password"])
	roleAttach["profile_name"] = body["profile_name"]
	roleAttach["status"] = "1"
	roleAttach["created_at"] = UTIL.GetCurrentTime().String()

	userID, status, ok := DB.InsertWithUniqueID(CONSTANT.UsersPermissionTable, CONSTANT.CounsellorDigits, roleAttach, "user_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	admin["username"] = body["user_name"]
	admin["user_id"] = userID
	admin["password"] = UTIL.GetStringMD5Hash(body["password"])
	admin["type"] = "1"
	admin["status"] = CONSTANT.AdminActive
	admin["created_at"] = UTIL.GetCurrentTime().String()

	_, status, ok = DB.InsertWithUniqueID(CONSTANT.AdminsTable, CONSTANT.CounsellorDigits, admin, "admin_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

// ProfileUpdate - update profile
func UpdateUserPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if admin already exists with specified username
	// if DB.CheckIfExists(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("id")}) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AminUserProfileMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// add roleAttach details
	roleAttach := map[string]string{}

	if len(body["role_id"]) > 0 {
		roleAttach["role_id"] = body["role_id"]
	}
	if len(body["user_name"]) > 0 {
		roleAttach["user_name"] = body["user_name"]
	}

	if len(body["password"]) > 0 {
		roleAttach["password"] = UTIL.GetStringMD5Hash(body["password"])
	}

	if len(body["profile_name"]) > 0 {
		roleAttach["profile_name"] = body["profile_name"]
	}

	if len(body["status"]) > 0 {
		roleAttach["status"] = body["status"]
	}
	roleAttach["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("id")}, roleAttach)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	admin := map[string]string{}
	if len(body["username"]) > 0 {
		admin["username"] = body["username"]
	}
	if len(body["password"]) > 0 {
		admin["password"] = UTIL.GetStringMD5Hash(body["password"])
	}

	admin["modified_at"] = UTIL.GetCurrentTime().String()

	status, ok = DB.UpdateSQL(CONSTANT.AdminsTable, map[string]string{"user_id": r.FormValue("id")}, admin)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// QuoteGet - get quote
func UserPermissionGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get quotes
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "user_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " user_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := " where status = '1' "
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	users, status, ok := DB.SelectProcess("select * from "+CONSTANT.UsersPermissionTable+where+" order by id desc limit 20", queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["users"] = users
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// UserProfileDelete - delete user profile
func UserPermissionDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	status, ok := DB.DeleteSQL(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("user_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
