package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

	MODEL "salbackend/model"
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

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminProfileAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.AddUserProfileRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if admin already exists with specified username
	if DB.CheckIfExists(CONSTANT.AdminsTable, map[string]string{"username": body.Username}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add admin details
	admin := map[string]string{}
	admin["username"] = body.Username
	admin["password"] = UTIL.GetStringMD5Hash(body.Password)
	admin["type"] = body.Type
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

	adminID, ok := UTIL.Required(r.FormValue("admin_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, adminID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.UpdateUserProfileRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if admin already exists with specified username
	if DB.CheckIfExists(CONSTANT.AdminsTable, map[string]string{"admin_id": r.FormValue("admin_id")}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add admin details
	admin := map[string]string{}
	admin["username"] = body.Username
	admin["password"] = UTIL.GetStringMD5Hash(body.Password)
	admin["type"] = body.Type
	admin["status"] = body.Status
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

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminUserProfileRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.CreateProfileForRoleRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// add assessment result
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.RolesTable, CONSTANT.AssessmentResultsDigits, map[string]string{
		"profile_name":       body.ProfileName,
		"pc_add":             body.PCAdd,
		"pc_edit":            body.PCEdit,
		"pc_view":            body.PCView,
		"assessment_add":     body.AssessmentAdd,
		"assessment_edit":    body.AssessmentEdit,
		"assessment_view":    body.AssessmentView,
		"home_add":           body.HomeAdd,
		"home_edit":          body.HomeEdit,
		"home_view":          body.HomeView,
		"slot_view":          body.SlotView,
		"slot_add":           body.SlotAdd,
		"slot_edit":          body.SlotEdit,
		"inperson_cafe_add":  body.InpersonCafeAdd,
		"inperson_cafe_edit": body.InpersonCafeEdit,
		"inperson_cafe_view": body.InpersonCafeView,
		"link_add":           body.LinkAdd,
		"link_edit":          body.LinkEdit,
		"link_view":          body.LinkView,
		"noti_add":           body.NotiAdd,
		"noti_edit":          body.NotiEdit,
		"noti_view":          body.NotiView,
		"cont_add":           body.ContAdd,
		"cont_edit":          body.ContEdit,
		"cont_view":          body.ContView,
		"mq_view":            body.MQView,
		"mq_add":             body.MQAdd,
		"mq_edit":            body.MQEdit,
		"cent_add":           body.CentAdd,
		"cent_edit":          body.CentEdit,
		"cent_view":          body.CentView,
		"coun_view":          body.CounView,
		"coun_edit":          body.CounEdit,
		"coun_add":           body.CounAdd,
		"part_add":           body.PartAdd,
		"part_edit":          body.PartEdit,
		"part_view":          body.PartView,
		"part_loc_add":       body.PartLocAdd,
		"part_loc_edit":      body.PartLocEdit,
		"part_loc_view":      body.PartLocView,
		"list_add":           body.ListAdd,
		"list_edit":          body.ListEdit,
		"list_view":          body.ListView,
		"ther_view":          body.TherView,
		"ther_add":           body.TherAdd,
		"ther_edit":          body.TherEdit,
		"appoint_add":        body.AppointAdd,
		"appoint_edit":       body.AppointEdit,
		"appoint_view":       body.AppointView,
		"cafe_add":           body.CafeAdd,
		"cafe_view":          body.CafeView,
		"cafe_edit":          body.CafeEdit,
		"rept_view":          body.ReptView,
		"rept_edit":          body.ReptEdit,
		"rept_add":           body.ReptAdd,
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

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	orderID, ok := UTIL.Required(r.FormValue("id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, orderID, CONSTANT.ShowDialog, response)
		return
	}

	body := MODEL.CreateProfileForRoleRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if admin already exists with specified username
	if !DB.CheckIfExists(CONSTANT.RolesTable, map[string]string{"role_id": r.FormValue("id")}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AminUserProfileMessage, CONSTANT.ShowDialog, response)
		return
	}

	profile := map[string]string{}

	profile["profile_name"] = body.ProfileName
	profile["pc_add"] = body.PCAdd

	profile["pc_edit"] = body.PCEdit

	profile["pc_view"] = body.PCView

	profile["assessment_add"] = body.AssessmentAdd

	profile["assessment_edit"] = body.AssessmentEdit

	profile["assessment_view"] = body.AssessmentView

	profile["home_view"] = body.HomeView

	profile["home_edit"] = body.HomeEdit

	profile["home_add"] = body.HomeAdd

	profile["inperson_cafe_add"] = body.InpersonCafeAdd

	profile["inperson_cafe_edit"] = body.InpersonCafeEdit

	profile["inperson_cafe_view"] = body.InpersonCafeView

	profile["slot_view"] = body.SlotView

	profile["slot_add"] = body.SlotAdd

	profile["slot_edit"] = body.SlotEdit

	profile["link_add"] = body.LinkAdd

	profile["link_edit"] = body.LinkEdit

	profile["link_view"] = body.LinkView

	profile["noti_add"] = body.NotiAdd

	profile["noti_edit"] = body.NotiEdit

	profile["noti_view"] = body.NotiView

	profile["cont_add"] = body.ContAdd

	profile["cont_edit"] = body.ContEdit

	profile["cont_view"] = body.ContView

	profile["mq_view"] = body.MQView

	profile["mq_edit"] = body.MQEdit

	profile["mq_add"] = body.MQAdd

	profile["cent_add"] = body.CentAdd

	profile["cent_view"] = body.CentView

	profile["cent_edit"] = body.CentEdit

	profile["coun_view"] = body.CounView

	profile["coun_add"] = body.CounAdd

	profile["coun_edit"] = body.CounEdit

	profile["part_add"] = body.PartAdd

	profile["part_edit"] = body.PartEdit

	profile["part_view"] = body.PartView

	profile["part_loc_add"] = body.PartLocAdd

	profile["part_loc_edit"] = body.PartLocEdit

	profile["part_loc_view"] = body.PartLocView

	profile["list_view"] = body.ListView

	profile["list_add"] = body.ListAdd

	profile["list_edit"] = body.ListEdit

	profile["ther_view"] = body.TherView

	profile["ther_edit"] = body.TherEdit

	profile["ther_add"] = body.TherAdd

	profile["appoint_view"] = body.AppointView

	profile["appoint_add"] = body.AppointAdd

	profile["appoint_edit"] = body.AppointEdit

	profile["cafe_add"] = body.CafeAdd

	profile["cafe_view"] = body.CafeView

	profile["cafe_edit"] = body.CafeEdit

	profile["rept_view"] = body.ReptView

	profile["rept_edit"] = body.ReptEdit

	profile["rept_add"] = body.ReptAdd

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
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get quotes
	wheres := []string{}
	queryArgs := []any{}
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
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.AttachPermissionAddRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AdminUserRoleRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check if admin already exists with specified username
	if DB.CheckIfExists(CONSTANT.UsersPermissionTable, map[string]string{"user_name": body.Username}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AdminExistsMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add admin main table
	admin := map[string]string{}

	// add roleAttach details
	roleAttach := map[string]string{}
	roleAttach["role_id"] = body.RoleID
	roleAttach["user_name"] = body.Username
	roleAttach["password"] = UTIL.GetStringMD5Hash(body.Password)
	roleAttach["profile_name"] = body.ProfileName
	roleAttach["status"] = CONSTANT.AccessRoleActive
	roleAttach["created_at"] = UTIL.GetCurrentTime().String()

	userID, status, ok := DB.InsertWithUniqueID(CONSTANT.UsersPermissionTable, CONSTANT.CounsellorDigits, roleAttach, "user_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	admin["username"] = body.Username
	admin["user_id"] = userID
	admin["password"] = UTIL.GetStringMD5Hash(body.Password)
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
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	orderID, ok := UTIL.Required(r.FormValue("id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, orderID, CONSTANT.ShowDialog, response)
		return
	}

	body := MODEL.AttachPermissionUpdateRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if admin already exists with specified username
	if !DB.CheckIfExists(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("id")}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AminUserProfileMessage, CONSTANT.ShowDialog, response)
		return
	}

	// add roleAttach details

	roleAttach := map[string]string{}
	roleAttach["role_id"] = body.RoleID
	roleAttach["user_name"] = body.Username
	roleAttach["password"] = UTIL.GetStringMD5Hash(body.Password)
	roleAttach["profile_name"] = body.ProfileName
	roleAttach["status"] = body.Status
	roleAttach["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.UsersPermissionTable, map[string]string{"user_id": r.FormValue("id")}, roleAttach)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	admin := map[string]string{}
	admin["username"] = body.Username
	admin["password"] = UTIL.GetStringMD5Hash(body.Password)
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
