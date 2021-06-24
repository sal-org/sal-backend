package listener

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
)

// NotificationsGet godoc
// @Tags Listener Notifications
// @Summary Get notifications for listener
// @Router /listener/notification [get]
// @Param listener_id query string true "Logged in listener ID"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func NotificationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get notifications for listener
	notifications, status, ok := DB.SelectProcess("select * from "+CONSTANT.NotificationsTable+" where person_id = ? order by created_at desc limit "+strconv.Itoa(CONSTANT.NotificationsPerPage)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.NotificationsPerPage), r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number notifications for listener
	notificationsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.PaymentsTable+" where person_id = ?", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["notifications"] = notifications
	response["notifications_count"] = notificationsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(notificationsCount[0]["ctn"], CONSTANT.NotificationsPerPage))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
