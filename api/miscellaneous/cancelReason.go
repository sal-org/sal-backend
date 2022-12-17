package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
)

// Cancellation Reason godoc
// @Tags Counsellor Cancellation
// @Summary Cancellation Region Counsellor
// @Router /cancellationreason [put]
// @Param appointment_id query string true "Appointment ID to Cancellation Reason"
// @Param body body model.CancellationUpdateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func CancellationReason(w http.ResponseWriter, r *http.Request) {
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

	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CancellationUpdateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	exists := DB.CheckIfExists(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !exists {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.UpdateSQL(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")}, map[string]string{"cancel_reason_counsellor": body["cancellation_reason"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
