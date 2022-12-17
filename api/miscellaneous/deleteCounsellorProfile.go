package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
)

// ListContentCategory godoc
// @Tags Miscellaneous
// @Summary delete User details
// @Router /delete-user [delete]
// @Param user_id query string true "User ID"
// @Param type query string true "1(counsellor)/2(listener)/3(Client)/4(therapist)"
// @Security JWTAuth
// @Produce json
// @Success 200
func DeleteUserProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	if r.FormValue("type") == CONSTANT.CounsellorType {
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("user_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.CounsellorBlocked, "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	if r.FormValue("type") == CONSTANT.ListenerType {
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"listener_id": r.FormValue("user_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.ListenerBlocked, "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	if r.FormValue("type") == CONSTANT.ClientType {
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": r.FormValue("user_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.ListenerBlocked, "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	if r.FormValue("type") == CONSTANT.TherapistType {
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": r.FormValue("user_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.TherapistBlocked, "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
