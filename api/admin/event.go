package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	"strconv"
	"strings"
	"time"

	_ "salbackend/model"
	UTIL "salbackend/util"
)

func EventGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "topic_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " topic_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "order_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " order_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor, order ids to get details
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(events, "counsellor_id")
	orderIDs := UTIL.ExtractValuesFromArrayMap(events, "order_id")

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get event booked count
	eventBookedCount, status, ok := DB.SelectProcess("select event_order_id, count(*) as ctn from " + CONSTANT.OrderEventTable + " where event_order_id in ('" + strings.Join(orderIDs, "','") + "') and status > " + CONSTANT.OrderWaiting + " group by event_order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of events
	eventsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.OrderCounsellorEventTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get topics
	topics, status, ok := DB.SelectProcess("select id, topic from " + CONSTANT.TopicsTable)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["events"] = events
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["topics"] = UTIL.ConvertMapToKeyMap(topics, "id")
	response["events_booked_count"] = UTIL.ConvertMapToKeyMap(eventBookedCount, "event_order_id")
	response["events_count"] = eventsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	orderID, ok := UTIL.Required(r.FormValue("order_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, orderID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.EventUpdateRequestInAdminPanel{}

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

	// check if corporate_id exists
	if !DB.CheckIfExists(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": orderID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add event
	event := map[string]string{}
	event["counsellor_id"] = body.CounsellorID
	event["title"] = body.Title
	event["description"] = body.Description
	event["photo"] = body.Photo
	event["date"] = body.Date
	event["time"] = body.Time
	event["duration"] = body.Duration
	event["price"] = body.Price
	event["status"] = body.Status
	event["modified_by"] = id
	event["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("order_id")}, event)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorType := CONSTANT.CounsellorType
	if len(DB.QueryRowSQL("select device_id from "+CONSTANT.TherapistsTable+" where therapist_id = ?", body.CounsellorID)) > 0 {
		counsellorType = CONSTANT.TherapistType
	}

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("order_id"), body.CounsellorID)

	// send event reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.CounsellorEventReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorEventReminderCounsellorContent,
			map[string]string{},
		),
		body.CounsellorID,
		counsellorType,
		UTIL.BuildDateTime(body.Date, body.Time).Add(-15*time.Minute).String(),
		CONSTANT.NotificationInProgress,
		r.FormValue("order_id"),
		"",
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventInPersonGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "order_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " order_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor, order ids to get details
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(events, "counsellor_id")
	orderIDs := UTIL.ExtractValuesFromArrayMap(events, "order_id")

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get event booked count
	eventBookedCount, status, ok := DB.SelectProcess("select event_order_id, count(*) as ctn from " + CONSTANT.OrderEventInPersonTable + " where event_order_id in ('" + strings.Join(orderIDs, "','") + "') and status > " + CONSTANT.OrderWaiting + " group by event_order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of events
	eventsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.OrderCounsellorEventInPersonTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// for _, event := range events {
	// 	urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, event["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURLPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
	// 	event["photo"] = endPointURLPhoto

	// 	urlBackGroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, event["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURLBackGroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackGroundPhoto)
	// 	event["background_photo"] = endPointURLBackGroundPhoto
	// }

	response["events"] = events
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["events_booked_count"] = UTIL.ConvertMapToKeyMap(eventBookedCount, "event_order_id")
	response["events_count"] = eventsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventInPersonAdd(w http.ResponseWriter, r *http.Request) {
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

	body := MODEL.InPersonEventAddRequestInAdminPanel{}

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

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	type1 := "1"

	if len(body.CounsellorID) != 0 {
		// get client details
		counsellor, _, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name, email"}, map[string]string{"counsellor_id": body.CounsellorID})
		if !ok {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			type1 = "4"
		} else {
			type1 = "1"
		}
	}

	// add event
	event := map[string]string{}
	event["counsellor_id"] = body.CounsellorID
	event["title"] = body.Title
	event["description"] = body.Description
	event["total_seat"] = body.TotalSeat
	event["remaining_seat"] = body.TotalSeat
	event["document"] = body.Document
	event["carry_things"] = body.CarryThings
	event["company_name"] = body.CompanyName
	event["company_location"] = body.CompanyLocation
	event["address"] = body.Address
	event["photo"] = body.Photo
	event["background_photo"] = body.BackgroundPhoto
	event["date"] = body.Date
	event["time"] = body.Time
	event["type"] = type1
	event["duration"] = body.Duration
	event["status"] = body.Status
	event["created_by"] = id
	event["created_at"] = UTIL.GetCurrentTime().String()
	// status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("order_id")}, event)
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	_, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderCounsellorEventInPersonTable, CONSTANT.EventDigits, event, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventInPersonUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	orderID, ok := UTIL.Required(r.FormValue("order_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, orderID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.InPersonEventAddRequestInAdminPanel{}

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

	// check if order_id exists
	if !DB.CheckIfExists(CONSTANT.OrderCounsellorEventInPersonTable, map[string]string{"order_id": orderID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	type1 := "1"

	if len(body.CounsellorID) != 0 {
		// get client details
		counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name, email", "type"}, map[string]string{"counsellor_id": body.CounsellorID})
		// if !ok {
		// 	UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		if len(counsellor) == 0 {
			type1 = "4"
		} else {
			type1 = "1"
		}
	}

	// endPointURLPhoto := UTIL.GetEndpointFromURL(body["photo"])
	// body["photo"] = endPointURLPhoto

	// endPointURLBackgroundPhoto := UTIL.GetEndpointFromURL(body["background_photo"])
	// body["background_photo"] = endPointURLBackgroundPhoto

	// add event
	event := map[string]string{}
	event["counsellor_id"] = body.CounsellorID
	event["title"] = body.Title
	event["description"] = body.Description
	event["carry_things"] = body.CarryThings
	event["document"] = body.Document
	event["company_name"] = body.CompanyName
	event["company_location"] = body.CompanyLocation
	event["address"] = body.Address
	event["photo"] = body.Photo
	event["background_photo"] = body.BackgroundPhoto
	event["date"] = body.Date
	event["time"] = body.Time
	event["type"] = type1
	event["duration"] = body.Duration
	event["status"] = body.Status
	event["modified_by"] = id
	event["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventInPersonTable, map[string]string{"order_id": r.FormValue("order_id")}, event)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// cancellation process if user booked event and admin change status to cancelled
	if body.Status == "5" {
		// get event booked count
		eventBooked, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderEventInPersonTable + " where event_order_id  = '" + r.FormValue("order_id") + "' and status = " + CONSTANT.InPersonEventOrderCompleted + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, booking := range eventBooked {
			// update each booking status to cancelled
			status, ok := DB.UpdateSQL(CONSTANT.OrderEventInPersonTable, map[string]string{"order_id": booking["order_id"]}, map[string]string{
				"status":           CONSTANT.InPersonEventOrderCancelledByAdmin,
				"cancelled_reason": "Event cancelled by administrator",
				"modified_at":      UTIL.GetCurrentTime().String(),
			})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// send event reminder notification to counsellor before 15 min
			UTIL.SendNotification(
				CONSTANT.AdminCancelledInPersonCafeHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.AdminCancelledInPersonCafeContent,
					map[string]string{
						"###InPersonCafeName###": body.Title,
						"###date###":             UTIL.BuildOnlyDate(booking["date"]),
					},
				),
				booking["user_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body.Date, body.Time).String(),
				CONSTANT.NotificationSent,
				booking["order_id"],
				"",
			)
		}
	}

	counsellorType := CONSTANT.CounsellorType
	if len(DB.QueryRowSQL("select device_id from "+CONSTANT.TherapistsTable+" where therapist_id = ?", body.CounsellorID)) > 0 {
		counsellorType = CONSTANT.TherapistType
	}

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("order_id"), body.CounsellorID)

	// send event reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.CounsellorEventReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorEventReminderCounsellorContent,
			map[string]string{},
		),
		body.CounsellorID,
		counsellorType,
		UTIL.BuildDateTime(body.Date, body.Time).Add(-15*time.Minute).String(),
		CONSTANT.NotificationInProgress,
		r.FormValue("order_id"),
		"",
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UploadEventFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var fileName string
	// file upload
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("UploadEventFile", err)
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	if file != nil {
		defer file.Close()

		name, uploaded := UTIL.UploadToS3(CONFIG.S3Bucket, CONSTANT.EventS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(handler.Filename), CONSTANT.S3PublicRead, file)
		if !uploaded {
			fmt.Println("UploadEventFile", err)
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name
	}

	// urlFile := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, fileName, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlFile)
	// fileName = endPointURL

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventBookGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events booked
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "user_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " user_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "event_order_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " event_order_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	wheres = append(wheres, " status > "+CONSTANT.OrderWaiting+" ")

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	eventsBooked, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get user ids to get details
	userIDs := UTIL.ExtractValuesFromArrayMap(eventsBooked, "user_id")

	// get user details
	users, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(userIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(userIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(userIDs, "','") + "')) union (select client_id as id, first_name, last_name from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of events booked
	eventsBookedCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.OrderEventTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["events_booked"] = eventsBooked
	response["users"] = UTIL.ConvertMapToKeyMap(users, "id")
	response["events_booked_count"] = eventsBookedCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsBookedCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventBookInPersonGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events booked
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "user_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " user_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "event_order_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " event_order_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	wheres = append(wheres, " status > "+CONSTANT.OrderWaiting+" ")

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	eventsBooked, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventInPersonTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get user ids to get details
	userIDs := UTIL.ExtractValuesFromArrayMap(eventsBooked, "user_id")

	// get user details
	users, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(userIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(userIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(userIDs, "','") + "')) union (select client_id as id, first_name, last_name from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of events booked
	eventsBookedCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.OrderEventInPersonTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["events_booked"] = eventsBooked
	response["users"] = UTIL.ConvertMapToKeyMap(users, "id")
	response["events_booked_count"] = eventsBookedCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsBookedCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToUploadEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.EventS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	// urlFile := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, fileName, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlFile)
	// fileName = endPointURL

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
