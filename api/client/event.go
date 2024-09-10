package client

import (
	"fmt"
	"math"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
	"time"
)

// EventsList godoc
// @Tags Client Event
// @Summary List available events
// @Router /client/events [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming events
	events, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventTable + " where status = " + CONSTANT.EventToBeStarted + " order by date desc, time desc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func InPersonEventsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var orderIDs []string
	var events []map[string]string

	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"email"}, map[string]string{"client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if client[0]["email"] == "anand.shah@clovemind.com" {

		eventsBooked, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and status = '1'", r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		orderIDs = UTIL.ExtractValuesFromArrayMap(eventsBooked, "order_id")

		// get upcoming events
		events, status, ok = DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventInPersonTable + " where status = " + CONSTANT.EventToBeStarted + " and date >= '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' order by date desc, time desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	} else {

		domainName := strings.Split(client[0]["email"], "@")

		partnerName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"*"}, map[string]string{"domain": domainName[1]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		eventsBooked, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and company_name = '"+partnerName[0]["partner_name"]+"' and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and status = '1'", r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		orderIDs = UTIL.ExtractValuesFromArrayMap(eventsBooked, "order_id")

		// get upcoming events
		events, status, ok = DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventInPersonTable + " where status = " + CONSTANT.EventToBeStarted + " and company_name = '" + partnerName[0]["partner_name"] + "' and date >= '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' order by date desc, time desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	}

	response["booked_event_id"] = orderIDs
	response["events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventDetail godoc
// @Tags Client Event
// @Summary Get event details
// @Router /client/event [get]
// @Param order_id query string true "Event order ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"*"}, map[string]string{"order_id": r.FormValue("order_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id = ?", event[0]["topic_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get event counsellor details
	var counsellor []map[string]string
	switch event[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"counsellor_id": event[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"therapist_id": event[0]["counsellor_id"]})
		break
	}
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["event"] = event[0]
	response["counsellor"] = counsellor[0]
	if len(topics) > 0 && len(topics[0]) > 0 {
		response["topic"] = topics[0]["topic"]
	}
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventInPersonDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventInPersonTable, []string{"*"}, map[string]string{"order_id": r.FormValue("order_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// // get event topics
	// topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id = ?", event[0]["topic_id"])
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get event counsellor details
	var counsellor []map[string]string
	switch event[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"counsellor_id": event[0]["counsellor_id"]})

	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"therapist_id": event[0]["counsellor_id"]})

	}
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", event[0]["counsellor_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", event[0]["counsellor_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["event"] = event[0]
	response["counsellor"] = counsellor[0]
	response["languages"] = languages
	response["topics"] = topics
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventsBooked godoc
// @Tags Client Event
// @Summary Get booked upcoming and past events
// @Router /client/event/booked [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsBooked(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") order by date asc, time asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventTable+" where user_id = ? and status > "+CONSTANT.OrderWaiting+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventsInPersonCancel(w http.ResponseWriter, r *http.Request) {
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

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and event_order_id = ? and status != '4' ", body["user_id"], body["order_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if client is valid
	if len(events) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["user_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventInPersonTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	DB.UpdateSQL(CONSTANT.OrderEventInPersonTable,
		map[string]string{
			"order_id": events[0]["order_id"],
		},
		map[string]string{
			"cancellation_reason": body["cancellation_reason"],
			"status":              CONSTANT.OrderCancel,
		},
	)

	remainingSeat, _ := strconv.Atoi(event[0]["remaining_seat"])

	remainingSeat = remainingSeat + 1

	remaining := strconv.Itoa(remainingSeat)

	DB.UpdateSQL(CONSTANT.OrderCounsellorEventInPersonTable,
		map[string]string{
			"order_id": event[0]["order_id"],
		},
		map[string]string{
			"remaining_seat": remaining,
		},
	)

	UTIL.RemoveNotification(events[0]["order_id"], events[0]["user_id"])

	// send to notification client cancellation
	UTIL.SendNotification(
		CONSTANT.ClientInPersonEventCancellationClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventCancellationClientContent,
			map[string]string{
				"###topic###": event[0]["title"],
			},
		),
		events[0]["user_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		events[0]["order_id"],
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventCancellationTextMessage,
			map[string]string{
				"###topic###": event[0]["title"],
				"###date###":  UTIL.BuildOnlyDate(event[0]["date"]),
				"###time###":  UTIL.GetTimeFromTimeSlotIN12Hour(event[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		events[0]["order_id"],
		CONSTANT.InstantSendTextMessage,
	)

	// event confirmation email
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventCancellationEmailBody,
			map[string]string{
				"###topic###": event[0]["title"],
			},
		),
	}

	filepath_text := "htmlfile/inpersonEventCancellation.html"

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientInPersonEventCancellationTitle,
		emailBody,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventsInPersonRequest(w http.ResponseWriter, r *http.Request) {
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

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventInPersonRequestRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	clients, status, ok := DB.SelectProcess("select first_name, last_name, phone from "+CONSTANT.ClientsTable+" where client_id = ?", body["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	getInPersonEventRequest, status, ok := DB.SelectSQL(CONSTANT.EventInPersonRequestTable, []string{"*"}, map[string]string{"client_id": body["client_id"], "order_id": body["order_id"], "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(getInPersonEventRequest) != 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAppointmentAlreadyExits, CONSTANT.ShowDialog, response)
		return
	}

	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventInPersonTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	eventInPersonRequest := map[string]string{}
	eventInPersonRequest["order_id"] = body["order_id"]
	eventInPersonRequest["client_id"] = body["client_id"]
	eventInPersonRequest["status"] = CONSTANT.AppointmentRequestProgress
	eventInPersonRequest["created_at"] = UTIL.GetCurrentTime().String()

	requestID, status, ok := DB.InsertWithUniqueID(CONSTANT.EventInPersonRequestTable, CONSTANT.AppointmentRequestDigits, eventInPersonRequest, "request_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// send to client
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientInPersonEventRequestTextMessage,
			map[string]string{
				"###topic###":    event[0]["title"],
				"###location###": event[0]["address"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		clients[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).UTC().String(),
		requestID,
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetEventInPersonRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	inpersonEventRequest, status, ok := DB.SelectSQL(CONSTANT.EventInPersonRequestTable, []string{"*"}, map[string]string{"client_id": r.FormValue("client_id"), "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["inperson_event_request"] = inpersonEventRequest

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func EventsInPersonRate(w http.ResponseWriter, r *http.Request) {
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

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and event_order_id = ? and status != '4' ", body["user_id"], body["order_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if client is valid
	if len(events) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	DB.UpdateSQL(CONSTANT.OrderEventInPersonTable,
		map[string]string{
			"order_id": events[0]["order_id"],
		},
		map[string]string{
			"question1": body["question1"],
			"question2": body["question2"],
			"question3": body["question3"],
			"question4": body["question4"],
			"question5": body["question5"],
		},
	)

	eventsOrder, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventInPersonTable+" where event_order_id = ? and status != '4' and question1 != ''", body["order_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	eventsOrderCount, status, ok := DB.SelectProcess("select count(*) as cnt from "+CONSTANT.OrderEventInPersonTable+" where event_order_id = ? and status != '4' and question1 != ''", body["order_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	avg := UTIL.AvgRatingFromula(eventsOrder, eventsOrderCount[0]["cnt"], "question1")

	DB.UpdateSQL(CONSTANT.OrderCounsellorEventInPersonTable,
		map[string]string{
			"order_id": events[0]["event_order_id"],
		},
		map[string]string{
			"total_rating": eventsOrderCount[0]["cnt"],
			"avg_rating":   avg,
		},
	)

	UTIL.SendNotification(
		CONSTANT.ClientEventRatingHeading,
		CONSTANT.ClientEventRatingContent,
		events[0]["user_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		events[0]["order_id"],
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetEventsInPersonRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select question1, question2, question3, question4, question5 from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and event_order_id = ? and status != '4' ", r.FormValue("user_id"), r.FormValue("order_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if client is valid
	if len(events) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["rating"] = events[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventsBookedInPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc, time asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// // get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	// events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status > "+CONSTANT.OrderWaiting+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("client_id"))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PastEventsInPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	eventsOrder, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+"", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["past_events"] = events
	response["past_order_event"] = eventsOrder

	// // get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	// events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status > "+CONSTANT.OrderWaiting+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("client_id"))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderCreate godoc
// @Tags Client Event
// @Summary Book a slot in an event
// @Router /client/event/order [post]
// @Param body body model.EventOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventOrderCreate(w http.ResponseWriter, r *http.Request) {
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

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["user_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if client is valid
	if len(client) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if client is active
	if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotAllowedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"*"}, map[string]string{"order_id": body["event_order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if event is valid
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if event is active
	if !strings.EqualFold(event[0]["status"], CONSTANT.EventToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["user_id"] = body["user_id"]
	order["event_order_id"] = body["event_order_id"]
	order["user_type"] = CONSTANT.ClientType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := event[0]["price"]

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and '"+UTIL.GetCurrentTime().String()+"' < end_by and (order_type = "+CONSTANT.OrderEventBookType+" or order_type = 0) and (client_id = ? or client_id = '') order by created_at desc limit 1", body["coupon_code"], body["user_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponInCorrectMessage, CONSTANT.ShowDialog, response)
			return
		}
		if !strings.EqualFold(coupon[0]["valid_for_order"], "0") { // coupon is valid for particular order
			// get total number of client appointment/event orders
			noOrders := DB.RowCount(CONSTANT.InvoicesTable, " user_id = ?", body["user_id"])
			// check if coupon applicable by order count and valid for order
			if !strings.EqualFold(coupon[0]["valid_for_order"], strconv.Itoa(noOrders+1)) { // add 1 to equal to valid for order value
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponNotApplicableMessage, CONSTANT.ShowDialog, response)
				return
			}
		}

		actualAmount, _ := strconv.ParseFloat(price, 64)
		minAmount, _ := strconv.ParseFloat(coupon[0]["minimum_order_value"], 64)
		if actualAmount >= minAmount { // coupon applicable only for minimum order value
			if strings.EqualFold(coupon[0]["type"], CONSTANT.CouponFlatType) { // flat
				order["discount"] = coupon[0]["discount"]
			} else if strings.EqualFold(coupon[0]["type"], CONSTANT.CouponPercentageType) { // percent
				discount, _ := strconv.ParseFloat(coupon[0]["discount"], 64) // percentage
				maxDiscount, _ := strconv.ParseFloat(coupon[0]["maximum_discount_value"], 64)
				discounted := actualAmount * discount / 100
				if discounted > maxDiscount {
					discounted = maxDiscount // maximum discount applied
				}
				order["discount"] = strconv.FormatFloat(discounted, 'f', 2, 64)
			}
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, strings.ReplaceAll(CONSTANT.CouponMinimumAmountRequiredMessage, "###amount###", coupon[0]["minimum_order_value"]), CONSTANT.ShowDialog, response)
			return
		}

		order["coupon_code"] = body["coupon_code"]
		order["coupon_id"] = coupon[0]["id"]
	}

	// calculate bill
	billing := UTIL.GetBillingDetails(price, order["discount"])
	order["paid_amount"] = billing["paid_amount"]
	order["discount"] = billing["discount"]
	order["tax"] = billing["tax"]
	order["actual_amount"] = billing["actual_amount"]
	order["cgst"] = billing["cgst"]
	order["sgst"] = billing["sgst"]

	amount, _ := strconv.ParseFloat(order["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.Itoa(int(math.Round(amount * 100)))
	response["paid_amount_razorpay"] = order["paid_amount_razorpay"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderEventTable, CONSTANT.OrderEventDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func EventOrderInPersonCreate(w http.ResponseWriter, r *http.Request) {
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

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["user_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if client is valid
	if len(client) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if client is active
	if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotAllowedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventInPersonTable, []string{"*"}, map[string]string{"order_id": body["event_order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if event is valid
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if event is active
	if !strings.EqualFold(event[0]["status"], CONSTANT.EventToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if event is active
	if strings.EqualFold(event[0]["remaining_seat"], "0") {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNoSeatMessage, CONSTANT.ShowDialog, response)
		return
	}

	ordertoCheck, status, ok := DB.SelectSQL(CONSTANT.OrderEventInPersonTable, []string{"*"}, map[string]string{"event_order_id": body["event_order_id"], "user_id": body["user_id"], "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(ordertoCheck) != 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderAlreadyExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["user_id"] = body["user_id"]
	order["event_order_id"] = body["event_order_id"]
	order["user_type"] = CONSTANT.ClientType
	order["status"] = CONSTANT.OrderInProgress
	order["created_at"] = UTIL.GetCurrentTime().String()

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderEventInPersonTable, CONSTANT.OrderEventDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	remainingSeat, _ := strconv.Atoi(event[0]["remaining_seat"])

	remainingSeat = remainingSeat - 1

	remaining := strconv.Itoa(remainingSeat)

	DB.UpdateSQL(CONSTANT.OrderCounsellorEventInPersonTable,
		map[string]string{
			"order_id": event[0]["order_id"],
		},
		map[string]string{
			"remaining_seat": remaining,
		},
	)

	// send to notification client
	UTIL.SendNotification(
		CONSTANT.ClientInPersonEventSucessClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventSucessClientContent,
			map[string]string{
				"###topic###": event[0]["title"],
			},
		),
		body["user_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		orderID,
	)

	// send appointment reminder notification to counsellor before 30 min
	UTIL.SendNotification(
		CONSTANT.ClientEventInPersonReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientEventInPersonRemiderClientContent,
			map[string]string{
				"###topic###": event[0]["title"],
				"###time###":  UTIL.GetTimeFromTimeSlotIN12Hour(event[0]["time"]),
			},
		),
		body["user_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(event[0]["date"], event[0]["time"]).Add(-30*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		orderID,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventConfirmationTextMessage,
			map[string]string{
				"###topic###": event[0]["title"],
				"###date###":  UTIL.BuildOnlyDate(event[0]["date"]),
				"###time###":  UTIL.GetTimeFromTimeSlotIN12Hour(event[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		orderID,
		CONSTANT.InstantSendTextMessage,
	)

	// event confirmation email
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonEventConfrimationEmailBody,
			map[string]string{
				"###topic###":    event[0]["title"],
				"###location###": event[0]["address"],
				"###date###":     UTIL.BuildOnlyDate(event[0]["date"]),
				"###time###":     UTIL.GetTimeFromTimeSlotIN12Hour(event[0]["time"]),
			},
		),
	}

	filepath_text := "htmlfile/inpersonEventConfirmation.html"

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientInPersonEventConfrimationTitle,
		emailBody,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderPaymentComplete godoc
// @Tags Client Event
// @Summary Call after payment is completed for event order
// @Router /client/event/paymentcomplete [post]
// @Param body body model.EventOrderPaymentCompleteRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
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

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderEventTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	razorPayTransaction := UTIL.GetRazorpayPayment(body["payment_id"])
	if !strings.EqualFold(razorPayTransaction.Description, body["order_id"]) { // check if razorpay payment id is associated with correct order id
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// capture razorpay payment
	amountRazorpay, _ := strconv.ParseFloat(order[0]["paid_amount_razorpay"], 64)
	UTIL.CaptureRazorpayPayment(body["payment_id"], amountRazorpay)

	// create invoice for the order
	invoice := map[string]string{}
	invoice["order_id"] = body["order_id"]
	invoice["payment_method"] = body["payment_method"]
	invoice["payment_id"] = body["payment_id"]
	invoice["user_id"] = order[0]["user_id"]
	invoice["user_type"] = CONSTANT.ClientType
	invoice["order_type"] = CONSTANT.OrderEventBookType
	invoice["actual_amount"] = order[0]["actual_amount"]
	invoice["tax"] = order[0]["tax"]
	invoice["cgst"] = order[0]["cgst"]
	invoice["sgst"] = order[0]["sgst"]
	invoice["discount"] = order[0]["discount"]
	invoice["coupon_code"] = order[0]["coupon_code"]
	invoice["coupon_id"] = order[0]["coupon_id"]
	invoice["paid_amount"] = order[0]["paid_amount"]
	invoice["status"] = CONSTANT.InvoiceInProgress
	invoice["created_at"] = UTIL.GetCurrentTime().String()

	invoiceID, status, ok := DB.InsertWithUniqueID(CONSTANT.InvoicesTable, CONSTANT.InvoiceDigits, invoice, "invoice_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update order with invoice id and change status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	orderUpdate["invoice_id"] = invoiceID
	status, ok = DB.UpdateSQL(CONSTANT.OrderEventTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	invoiceforemail, _, _ := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"id", "user_id", "discount", "paid_amount", "payment_id", "created_at"}, map[string]string{"invoice_id": invoiceID})
	orderdetails, _, _ := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"counsellor_id", "title", "date", "time", "price"}, map[string]string{"order_id": order[0]["event_order_id"]})
	counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "timezone"}, map[string]string{"counsellor_id": orderdetails[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"timezone", "email"}, map[string]string{"client_id": order[0]["user_id"]})

	// send event booking notification to client
	UTIL.SendNotification(
		CONSTANT.ClientEventPaymentSucessClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientEventPaymentSucessClientContent,
			map[string]string{
				"###cafe_name###":   orderdetails[0]["title"],
				"###paid_amount###": order[0]["paid_amount"],
				"###date_time###":   UTIL.ConvertTimezone(UTIL.BuildDateTime(orderdetails[0]["date"], orderdetails[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		order[0]["user_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		order[0]["event_order_id"],
	)

	// send event reminder notification to client before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientEventReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientEventRemiderClientContent,
			map[string]string{
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		order[0]["user_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(orderdetails[0]["date"], orderdetails[0]["time"]).Add(-15*time.Minute).String(),
		CONSTANT.NotificationInProgress,
		order[0]["event_order_id"],
	)

	receiptdata := UTIL.BuildDate(invoiceforemail[0]["created_at"])

	data := Model.EmailDataForPaymentReceipt{
		Date:        receiptdata,
		ReceiptNo:   invoiceforemail[0]["id"],
		ReferenceNo: invoiceforemail[0]["payment_id"],
		SPrice:      orderdetails[0]["price"],
		Qty:         CONSTANT.SalCafeQty,
		Total:       orderdetails[0]["price"],
		//SessionsType: CONSTANT.AppointmentSessionsTypeForReceipt,
		TPrice:   orderdetails[0]["price"],
		Discount: invoiceforemail[0]["discount"],
		TotalP:   invoiceforemail[0]["paid_amount"],
	}

	filepath := "htmlfile/index.html"

	emailbody, ok := UTIL.GetHTMLTemplateForReceipt(data, filepath)
	if !ok {
		fmt.Println("html body not create ")
	}

	UTIL.SendEmail(
		CONSTANT.ClientPaymentSucessClientTitle,
		emailbody,
		client[0]["email"], //
		CONSTANT.InstantSendEmailMessage,
	)

	/*UTIL.SendEmail(
		CONSTANT.ClientPaymentSucessClientTitle,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.SendAReceiptForClient,
			map[string]string{
				"###Date###":         receiptdata,
				"###ReceiptNo###":    invoiceforemail[0]["id"],
				"###ReferenceNo###":  invoiceforemail[0]["payment_id"],
				"###SPrice###":       orderdetails[0]["price"],
				"###Qty###":          CONSTANT.SalCafeQty,
				"###Total###":        orderdetails[0]["price"],
				"###SessionsType###": CONSTANT.SalCafeTypeForReceipt,
				"###TPrice###":       orderdetails[0]["price"],
				"###Discount###":     invoiceforemail[0]["discount"],
				"###TotalP###":       invoiceforemail[0]["paid_amount"],
			},
		),
		client[0]["email"], //client[0]["email"]
		CONSTANT.InstantSendEmailMessage,
	)*/

	response["invoice_id"] = invoiceID
	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)

}
