package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// Counsellor Content godoc
// @Tags Miscellaneous
// @Summary Update Notification status
// @Router /notification-status [put]
// @Param user_id query string true "ClientID/TherapistID/CounsellorID/ListernerID"
// @Param body body model.NotificationAllowSettingModel true "User type (counsellor:1/listener:2/client:3/therapists:4)"
// @Security JWTAuth
// @Produce json
// @Success 200
func NotificationInactiveORActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.NotificationRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	switch body["userType"] {
	case CONSTANT.CounsellorType:
		oK := DB.CheckIfExists(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("user_id")})
		if !oK {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		status, ok := DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("user_id")}, map[string]string{"notification_status": body["status"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	case CONSTANT.ListenerType:
		oK := DB.CheckIfExists(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("user_id")})
		if !oK {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		status, ok := DB.UpdateSQL(CONSTANT.ListenersTable, map[string]string{"listener_id": r.FormValue("user_id")}, map[string]string{"notification_status": body["status"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	case CONSTANT.ClientType:
		oK := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("user_id")})
		if !oK {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		status, ok := DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("user_id")}, map[string]string{"notification_status": body["status"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	case CONSTANT.TherapistType:
		oK := DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("user_id")})
		if !oK {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		status, ok := DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("user_id")}, map[string]string{"notification_status": body["status"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	default:
		UTIL.SetReponse(w, "400", "pass right user type", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
