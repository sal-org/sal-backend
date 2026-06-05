package client

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
)

// NotificationsGet godoc
// @Tags Client Notifications
// @Summary Get notifications for client
// @Router /client/notification [get]
// @Param client_id query string true "Logged in client ID"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func NotificationsGet(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var encryptedResponse = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get notifications for client
	notifications, status, ok := DB.SelectProcess("select * from "+CONSTANT.NotificationsTable+" where user_id = ? and notification_status='2' and created_at >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by created_at desc limit "+strconv.Itoa(CONSTANT.NotificationsPerPage)+" offset "+strconv.Itoa((UTIL.GetPageNumber(body["page"])-1)*CONSTANT.NotificationsPerPage), body["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number notifications for client
	notificationsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.NotificationsTable+" where user_id = ? and notification_status='2' ", body["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["notifications"] = notifications
	response["notifications_count"] = notificationsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(notificationsCount[0]["ctn"], CONSTANT.NotificationsPerPage))

	encrypt, _ := UTIL.EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}
