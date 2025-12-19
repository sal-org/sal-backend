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

		status, ok = DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.CounsellorDeleted, "deletion_reason": r.FormValue("reason"), "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
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

		status, ok = DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.ListenerDeleted, "deletion_reason": r.FormValue("reason"), "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
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

		status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.ClientDeleted, "deletion_reason": r.FormValue("reason"), "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
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

		status, ok = DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("user_id")}, map[string]string{"status": CONSTANT.TherapistDeleted, "last_login_time": UTIL.GetCurrentTime().String(), "deletion_reason": r.FormValue("reason"), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func RestoreUserProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	userID := ""

	user, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	userType := "1"
	if len(user) == 0 {
		user, status, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		userType = "2"

		if len(user) == 0 {
			user, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"phone": r.FormValue("phone")})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			userType = "4"
		}
	}

	if userType == CONSTANT.CounsellorType {

		userID = user[0]["counsellor_id"]
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": user[0]["user_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		if counsellor[0]["status"] != CONSTANT.CounsellorDeleted {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": user[0]["counsellor_id"]}, map[string]string{"status": CONSTANT.CounsellorActive, "deletion_reason": "", "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	if userType == CONSTANT.ListenerType {

		userID = user[0]["listener_id"]
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"listener_id": user[0]["listener_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		if counsellor[0]["status"] != CONSTANT.ListenerDeleted {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": user[0]["listener_id"]}, map[string]string{"status": CONSTANT.ListenerActive, "deletion_reason": "", "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	if userType == CONSTANT.TherapistType {

		userID = user[0]["therapist_id"]
		// get counsellor details
		counsellor, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": user[0]["therapist_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		if counsellor[0]["status"] != CONSTANT.TherapistDeleted {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountBlockedMessage, CONSTANT.ShowDialog, response)
			return
		}

		status, ok = DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": user[0]["therapist_id"]}, map[string]string{"status": CONSTANT.TherapistActive, "deletion_reason": "", "last_login_time": UTIL.GetCurrentTime().String(), "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	UTIL.SendNotification(
		CONSTANT.ClientRestoreAccountClientHeading,
		CONSTANT.ClientRestoreAccountClientContent,
		userID,
		userType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationInProgress,
		userID,
		"",
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
