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

	var response = make(map[string]any)

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

	var response = make(map[string]any)

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

	var response = make(map[string]any)

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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminUserProfileRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// add assessment result
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.RolesTable, CONSTANT.AssessmentResultsDigits, map[string]string{
		"profile_name":       body["profile_name"],
		"pc_add":             body["pc_add"],
		"pc_edit":            body["pc_edit"],
		"pc_view":            body["pc_view"],
		"assessment_add":     body["assessment_add"],
		"assessment_edit":    body["assessment_edit"],
		"assessment_view":    body["assessment_view"],
		"home_add":           body["home_add"],
		"home_edit":          body["home_edit"],
		"home_view":          body["home_view"],
		"slot_view":          body["slot_view"],
		"slot_add":           body["slot_add"],
		"slot_edit":          body["slot_edit"],
		"inperson_cafe_add":  body["inperson_cafe_add"],
		"inperson_cafe_edit": body["inperson_cafe_edit"],
		"inperson_cafe_view": body["inperson_cafe_view"],
		"link_add":           body["link_add"],
		"link_edit":          body["link_edit"],
		"link_view":          body["link_view"],
		"noti_add":           body["noti_add"],
		"noti_edit":          body["noti_edit"],
		"noti_view":          body["noti_view"],
		"cont_add":           body["cont_add"],
		"cont_edit":          body["cont_edit"],
		"cont_view":          body["cont_view"],
		"mq_view":            body["mq_view"],
		"mq_add":             body["mq_add"],
		"mq_edit":            body["mq_edit"],
		"cent_add":           body["cent_add"],
		"cent_edit":          body["cent_edit"],
		"cent_view":          body["cent_view"],
		"coun_view":          body["coun_view"],
		"coun_edit":          body["coun_edit"],
		"coun_add":           body["coun_add"],
		"part_add":           body["part_add"],
		"part_edit":          body["part_edit"],
		"part_view":          body["part_view"],
		"part_loc_add":       body["part_loc_add"],
		"part_loc_edit":      body["part_loc_edit"],
		"part_loc_view":      body["part_loc_view"],
		"list_add":           body["list_add"],
		"list_edit":          body["list_edit"],
		"list_view":          body["list_view"],
		"ther_view":          body["ther_view"],
		"ther_add":           body["ther_add"],
		"ther_edit":          body["ther_edit"],
		"appoint_add":        body["appoint_add"],
		"appoint_edit":       body["appoint_edit"],
		"appoint_view":       body["appoint_view"],
		"cafe_add":           body["cafe_add"],
		"cafe_view":          body["cafe_view"],
		"cafe_edit":          body["cafe_edit"],
		"rept_view":          body["rept_view"],
		"rept_edit":          body["rept_edit"],
		"rept_add":           body["rept_add"],
		"status":             CONSTANT.ClientActive,
		"created_at":         UTIL.GetCurrentTime().UTC().String(),
	}, "role_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UpdateProfileForUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

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

	profile["pc_add"] = body["pc_add"]

	profile["pc_edit"] = body["pc_edit"]

	profile["pc_view"] = body["pc_view"]

	profile["assessment_add"] = body["assessment_add"]

	profile["assessment_edit"] = body["assessment_edit"]

	profile["assessment_view"] = body["assessment_view"]

	profile["home_view"] = body["home_view"]

	profile["home_edit"] = body["home_edit"]

	profile["home_add"] = body["home_add"]

	profile["inperson_cafe_add"] = body["inperson_cafe_add"]

	profile["inperson_cafe_edit"] = body["inperson_cafe_edit"]

	profile["inperson_cafe_view"] = body["inperson_cafe_view"]

	profile["slot_view"] = body["slot_view"]

	profile["slot_add"] = body["slot_add"]

	profile["slot_edit"] = body["slot_edit"]

	profile["link_add"] = body["link_add"]

	profile["link_edit"] = body["link_edit"]

	profile["link_view"] = body["link_view"]

	profile["noti_add"] = body["noti_add"]

	profile["noti_edit"] = body["noti_edit"]

	profile["noti_view"] = body["noti_view"]

	profile["cont_add"] = body["cont_add"]

	profile["cont_edit"] = body["cont_edit"]

	profile["cont_view"] = body["cont_view"]

	profile["mq_view"] = body["mq_view"]

	profile["mq_edit"] = body["mq_edit"]

	profile["mq_add"] = body["mq_add"]

	profile["cent_add"] = body["cent_add"]

	profile["cent_view"] = body["cent_view"]

	profile["cent_edit"] = body["cent_edit"]

	profile["coun_view"] = body["coun_view"]

	profile["coun_add"] = body["coun_add"]

	profile["coun_edit"] = body["coun_edit"]

	profile["part_add"] = body["part_add"]

	profile["part_edit"] = body["part_edit"]

	profile["part_view"] = body["part_view"]

	profile["part_loc_add"] = body["part_loc_add"]

	profile["part_loc_edit"] = body["part_loc_edit"]

	profile["part_loc_view"] = body["part_loc_view"]

	profile["list_view"] = body["list_view"]

	profile["list_add"] = body["list_add"]

	profile["list_edit"] = body["list_edit"]

	profile["ther_view"] = body["ther_view"]

	profile["ther_edit"] = body["ther_edit"]

	profile["ther_add"] = body["ther_add"]

	profile["appoint_view"] = body["appoint_view"]

	profile["appoint_add"] = body["appoint_add"]

	profile["appoint_edit"] = body["appoint_edit"]

	profile["cafe_add"] = body["cafe_add"]

	profile["cafe_view"] = body["cafe_view"]

	profile["cafe_edit"] = body["cafe_edit"]

	profile["rept_view"] = body["rept_view"]

	profile["rept_edit"] = body["rept_edit"]

	profile["rept_add"] = body["rept_add"]

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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.DeleteSQL(CONSTANT.RolesTable, map[string]string{"role_id": r.FormValue("role_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AttachPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)


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

	var response = make(map[string]any)


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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.DeleteSQL(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("user_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
