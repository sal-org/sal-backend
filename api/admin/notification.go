package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"
	"time"

	CONFIG "salbackend/config"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

func NotificationGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

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
		case "notification_type":
			if len(val[0]) > 0 {
				wheres = append(wheres, " notification_type = ? ")
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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	body := MODEL.AddNotificationRequestInAdminPanel{}

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

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	notification := map[string]string{}
	notification["title"] = body.Title
	notification["body"] = body.Body
	notification["user_ids"] = body.UserIds
	notification["type"] = body.Type
	notification["notification_type"] = body.NotificationType
	notification["user_type"] = body.UserType
	notification["status"] = CONSTANT.NotificationActive
	notification["created_by"] = id
	notification["created_at"] = UTIL.GetCurrentTime().String()

	// get all deivce_ids and send notifications
	var (
		devices  []map[string]string
		status   string
		ok       bool
		userType string
	)
	if strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		if strings.EqualFold(notification["user_type"], CONSTANT.CompanyType) {

			// check if Partner Name exists
			if !DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"partner_name": body.PartnerName}) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid name", CONSTANT.ShowDialog, response)
				return
			}

			userType = "3"

			domain := DB.QueryRowSQL("select domain from "+CONSTANT.CorporatePartnersTable+" where partner_name = ? ", body.PartnerName)

			if body.PartnerLocation != "" {
				devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id, '3' as type from " + CONSTANT.ClientsTable + " where email like '%" + domain + "' and location = '" + body.PartnerLocation + "' and push_notification_status = '1'")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
			} else {
				devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id, '3' as type from " + CONSTANT.ClientsTable + " where email like '%" + domain + "' and push_notification_status = '1'")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
			}

		}
		// else {
		// 	if strings.EqualFold(notification["user_type"], CONSTANT.CounsellorType) {
		// 		devices, status, ok = DB.SelectProcess("select counsellor_id as user_id, device_id, '1' as type from " + CONSTANT.CounsellorsTable)
		// 		if !ok {
		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 			return
		// 		}
		// 	} else if strings.EqualFold(notification["user_type"], CONSTANT.ListenerType) {
		// 		devices, status, ok = DB.SelectProcess("select listener_id as user_id, device_id, '1' as type from " + CONSTANT.ListenersTable)
		// 		if !ok {
		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 			return
		// 		}
		// 	} else if strings.EqualFold(notification["user_type"], CONSTANT.ClientType) {
		// 		devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id, '3' as type from " + CONSTANT.ClientsTable)
		// 		if !ok {
		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 			return
		// 		}
		// 	} else {
		// 		devices, status, ok = DB.SelectProcess("select therapist_id as user_id, device_id, '1' as type from " + CONSTANT.TherapistsTable)
		// 		if !ok {
		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 			return
		// 		}
		// 	}
		// }

	} else {
		if strings.EqualFold(notification["user_type"], CONSTANT.CounsellorType) {

			userType = "4"
			devices, status, ok = DB.SelectProcess("select counsellor_id as user_id, device_id from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ListenerType) {

			userType = "4"
			devices, status, ok = DB.SelectProcess("select listener_id as user_id, device_id from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else if strings.EqualFold(notification["user_type"], CONSTANT.ClientType) {
			userType = "3"
			devices, status, ok = DB.SelectProcess("select client_id as user_id, device_id from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		} else {

			userType = "4"
			devices, status, ok = DB.SelectProcess("select therapist_id as user_id, device_id from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(strings.Split(notification["user_ids"], ","), "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		}
	}

	if strings.EqualFold(notification["user_type"], CONSTANT.ClientType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		UTIL.SendBulkNotification(body.Title, body.Body, CONSTANT.ClientType)
		// } else if strings.EqualFold(notification["user_type"], CONSTANT.CompanyType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		// 	for _, device := range devices {
		// 		UTIL.SendNotification(body["title"], body["body"], device["user_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, device["user_id"])
		// 	}
		// }
	} else if strings.EqualFold(notification["user_type"], CONSTANT.CounsellorType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		UTIL.SendBulkNotification(body.Title, body.Body, CONSTANT.CounsellorType)
		// } else if strings.EqualFold(notification["user_type"], CONSTANT.CompanyType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		// 	for _, device := range devices {
		// 		UTIL.SendNotification(body["title"], body["body"], device["user_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, device["user_id"])
		// 	}
		// }
	} else if strings.EqualFold(notification["user_type"], CONSTANT.ListenerType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		UTIL.SendBulkNotification(body.Title, body.Body, CONSTANT.ListenerType)
		// } else if strings.EqualFold(notification["user_type"], CONSTANT.CompanyType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		// 	for _, device := range devices {
		// 		UTIL.SendNotification(body["title"], body["body"], device["user_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, device["user_id"])
		// 	}
		// }
	} else if strings.EqualFold(notification["user_type"], CONSTANT.TherapistType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		UTIL.SendBulkNotification(body.Title, body.Body, CONSTANT.TherapistType)
		// } else if strings.EqualFold(notification["user_type"], CONSTANT.CompanyType) && strings.EqualFold(notification["type"], CONSTANT.BulkNotificationAll) {
		// 	for _, device := range devices {
		// 		UTIL.SendNotification(body["title"], body["body"], device["user_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, device["user_id"])
		// 	}
		// }
	} else {
		// send all notifications
		for _, device := range devices {
			// UTIL.SendNotification(body["title"], body["body"], device["user_id"], notification["user_type"], UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, device["user_id"])
			DB.InsertWithUniqueID(CONSTANT.NotificationsTable, CONSTANT.NotificationsDigits, map[string]string{
				"user_id":             device["user_id"],
				"onesignal_id":        device["device_id"],
				"title":               notification["title"],
				"body":                notification["body"],
				"type":                userType,
				"status":              CONSTANT.NotificationActive,
				"notification_status": CONSTANT.NotificationInProgress,
				"send_at":             UTIL.GetCurrentTime().Add(330 * time.Minute).UTC().String(),
				"created_at":          UTIL.GetCurrentTime().String(),
			}, "notification_id")
		}
	}

	// add notification
	status, ok = DB.InsertSQL(CONSTANT.NotificationsBulkTable, notification)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
