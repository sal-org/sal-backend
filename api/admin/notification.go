package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

func NotificationGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get notifications
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "type":
			if len(val[0]) > 0 {
				wheres = append(wheres, " type = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "user_type":
			if len(val[0]) > 0 {
				wheres = append(wheres, " user_type = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	notifications, status, ok := DB.SelectProcess("select * from "+CONSTANT.NotificationsBulkTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of notifications
	notificationsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.NotificationsBulkTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["notifications"] = notifications
	response["notifications_count"] = notificationsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(notificationsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func NotificationAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	notification := map[string]string{}
	notification["title"] = body["title"]
	notification["body"] = body["body"]
	notification["user_ids"] = body["user_ids"]
	notification["type"] = body["type"]
	notification["user_type"] = body["user_type"]
	notification["status"] = "1"
	notification["created_by"] = body["created_by"]
	notification["created_at"] = UTIL.GetCurrentTime().String()

	// get all deivce_ids and send notifications
	var (
		devices []map[string]string
		status  string
	)
	if strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		if strings.EqualFold(notification["user_type"], CONSTANT.CounsellorType) {
			devices, status, ok = DB.SelectProcess("select counsellor_id as user_id, device_id from " + CONSTANT.CounsellorsTable)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ListenerType) {
			devices, status, ok = DB.SelectProcess("select listener_id as user_id, device_id from " + CONSTANT.ListenersTable)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ClientType) {
			devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id from " + CONSTANT.ClientsTable)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else {
			devices, status, ok = DB.SelectProcess("select therapist_id as user_id, device_id from " + CONSTANT.TherapistsTable)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		}
	} else {
		if strings.EqualFold(notification["user_type"], CONSTANT.CounsellorType) {
			devices, status, ok = DB.SelectProcess("select counsellor_id as user_id, device_id from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ListenerType) {
			devices, status, ok = DB.SelectProcess("select listener_id as user_id, device_id from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ClientType) {
			devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else {
			devices, status, ok = DB.SelectProcess("select therapist_id as user_id, device_id from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		}
	}

	// send all notifications
	for _, device := range devices {
		DB.InsertWithUniqueID(CONSTANT.NotificationsTable, CONSTANT.NotificationsDigits, map[string]string{
			"user_id":             device["user_id"],
			"onesignal_id":        device["device_id"],
			"title":               notification["title"],
			"body":                notification["body"],
			"status":              CONSTANT.NotificationActive,
			"notification_status": CONSTANT.NotificationInProgress,
			"created_at":          UTIL.GetCurrentTime().String(),
		}, "notification_id")
	}

	// add notification
	status, ok = DB.InsertSQL(CONSTANT.NotificationsBulkTable, notification)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
