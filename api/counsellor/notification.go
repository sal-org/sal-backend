package counsellor

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
)

// NotificationsGet godoc
// @Tags Counsellor Notifications
// @Summary Get notifications for counsellor
// @Router /counsellor/notification [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func NotificationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get notifications for counsellor
	notifications, status, ok := DB.SelectProcess("select * from "+CONSTANT.NotificationsTable+" where user_id = ? order by created_at desc limit "+strconv.Itoa(CONSTANT.NotificationsPerPage)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.NotificationsPerPage), r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number notifications for counsellor
	notificationsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.NotificationsTable+" where user_id = ?", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["notifications"] = notifications
	response["notifications_count"] = notificationsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(notificationsCount[0]["ctn"], CONSTANT.NotificationsPerPage))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
