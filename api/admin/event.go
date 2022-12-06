package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"
	"time"

	_ "salbackend/model"
	UTIL "salbackend/util"
)

func EventGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events
	wheres := []string{}
	queryArgs := []interface{}{}
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
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add event
	event := map[string]string{}
	event["counsellor_id"] = body["counsellor_id"]
	event["title"] = body["title"]
	event["description"] = body["description"]
	event["photo"] = body["photo"]
	event["topic_id"] = body["topic_id"]
	event["date"] = body["date"]
	event["time"] = body["time"]
	event["duration"] = body["duration"]
	event["price"] = body["price"]
	event["status"] = body["status"]
	event["modified_by"] = body["modified_by"]
	event["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("order_id")}, event)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorType := CONSTANT.CounsellorType
	if len(DB.QueryRowSQL("select device_id from "+CONSTANT.TherapistsTable+" where therapist_id = ?", body["counsellor_id"])) > 0 {
		counsellorType = CONSTANT.TherapistType
	}

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("order_id"), body["counsellor_id"])

	// send event reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.CounsellorEventReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorEventReminderCounsellorContent,
			map[string]string{},
		),
		body["counsellor_id"],
		counsellorType,
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).String(),
		CONSTANT.NotificationInProgress,
		r.FormValue("order_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UploadEventFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventBookGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events booked
	wheres := []string{}
	queryArgs := []interface{}{}
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
