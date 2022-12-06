package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"time"

	UTIL "salbackend/util"
	"strings"
)

// ListenerProfile godoc
// @Tags Client Listener
// @Summary Get listener details
// @Router /client/listener [get]
// @Param listener_id query string true "Listener ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListenerProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get listener details
	listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "slot_type", "age_group"}, map[string]string{"listener_id": r.FormValue("listener_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(listener) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get listener languages
	languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get listener topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get last 10 listener apppointment reviews
	reviews, status, ok := DB.SelectProcess("select a.comment, a.rating, a.modified_at, c.first_name, c.last_name from "+CONSTANT.AppointmentsTable+" a, "+CONSTANT.ClientsTable+" c where a.client_id = c.client_id and a.counsellor_id = ? and a.status = "+CONSTANT.AppointmentCompleted+" and a.comment != '' order by a.modified_at desc limit 10 ", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor latest content
	contents, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where counsellor_id = ? and training = 0 and status = 1 order by created_at desc limit 20", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["listener"] = listener[0]
	response["languages"] = languages
	response["topics"] = topics
	response["reviews"] = reviews
	response["contents"] = contents
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerSlots godoc
// @Tags Client Listener
// @Summary Get listener slots
// @Router /client/listener/slots [get]
// @Param listener_id query string true "Listener ID to get slot details"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListenerSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get listener slots
	slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// remove times and dates with no availability
	response["slots"] = UTIL.FilterAvailableSlots(slots)
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerOrderCreate godoc
// @Tags Client Listener
// @Summary Create appointment order with client and listener
// @Router /client/listener/order [post]
// @Param body body model.ListenerOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListenerOrderCreate(w http.ResponseWriter, r *http.Request) {
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ListenerOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["client_id"]})
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

	// get listener details
	listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"listener_id": body["listener_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if listener is valid
	if len(listener) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if listener is active
	if !strings.EqualFold(listener[0]["status"], CONSTANT.ListenerActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(body["listener_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["listener_id"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["type"] = CONSTANT.ListenerType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()
	// no paid amount and billing

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderClientAppointmentTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerOrderPaymentComplete godoc
// @Tags Client Listener
// @Summary Call after payment is completed for listener order
// @Router /client/listener/paymentcomplete [post]
// @Param body body model.ListenerOrderPaymentCompleteRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListenerOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ListenerOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if order is with listener
	if !strings.EqualFold(order[0]["type"], CONSTANT.ListenerType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create appointment between listener and client
	appointment := map[string]string{}
	appointment["order_id"] = body["order_id"]
	appointment["client_id"] = order[0]["client_id"]
	appointment["counsellor_id"] = order[0]["counsellor_id"]
	appointment["type"] = order[0]["type"]
	appointment["date"] = order[0]["date"]
	appointment["time"] = order[0]["time"]
	appointment["status"] = CONSTANT.AppointmentToBeStarted
	appointment["created_at"] = UTIL.GetCurrentTime().String()
	appointmentID, status, ok := DB.InsertWithUniqueID(CONSTANT.AppointmentsTable, CONSTANT.AppointmentDigits, appointment, "appointment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor_name, status, ok := DB.SelectProcess("select first_name , last_name from "+CONSTANT.ListenersTable+" where listener_id = ?", order[0]["counsellor_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	counsellor_fullname := counsellor_name[0]["first_name"] + " " + counsellor_name[0]["last_name"]

	client_name, status, ok := DB.SelectProcess("select first_name , last_name from "+CONSTANT.ClientsTable+" where client_id = ?", order[0]["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	client_fullname := client_name[0]["first_name"] + " " + client_name[0]["last_name"]

	qualitycheck_details := map[string]string{}
	qualitycheck_details["appointment_id"] = appointmentID
	qualitycheck_details["client_id"] = order[0]["client_id"]
	qualitycheck_details["client_name"] = client_fullname
	qualitycheck_details["counsellor_id"] = order[0]["counsellor_id"]
	qualitycheck_details["counsellor_name"] = counsellor_fullname
	qualitycheck_details["type"] = order[0]["type"]
	qualitycheck_details["date"] = order[0]["date"]
	qualitycheck_details["time"] = order[0]["time"]
	qualitycheck_details["status"] = CONSTANT.AppointmentToBeStarted
	qualitycheck_details["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok = DB.InsertWithUniqueID(CONSTANT.QualityCheckDetailsTable, CONSTANT.AppointmentDigits, qualitycheck_details, "qualitycheck_details_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// change order status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok = DB.UpdateSQL(CONSTANT.OrderClientAppointmentTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	// update listener slots
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": order[0]["counsellor_id"],
			"date":          order[0]["date"],
		},
		map[string]string{
			order[0]["time"]: CONSTANT.SlotBooked,
		},
	)

	// sent notitifications
	listener, _, _ := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "timezone"}, map[string]string{"listener_id": order[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone", "timezone"}, map[string]string{"client_id": order[0]["client_id"]})

	// send appointment booking notification to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleClientHeading, CONSTANT.ClientAppointmentScheduleClientContent, order[0]["client_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().String(), CONSTANT.NotificationSent, appointmentID,
	)

	// send appointment reminder notification to client before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		appointmentID,
	)

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBookClientEmailBody,
			map[string]string{
				"###therpist_name###": listener[0]["first_name"],
				"###date_time###":     UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookClientTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// send appointment booking notification, message to listener
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentScheduleCounsellorContent,
			map[string]string{
				"###Date###": order[0]["date"],
				"###Time###": UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["counsellor_id"],
		CONSTANT.ListenerType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	// Send to appointment Reminder SMS to client
	// send at 15 min before of appointment
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###userName###":  listener[0]["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// Send to appointment Reminder SMS to counsellor
	// send at 15 min before of appointment
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": listener[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###userName###":  client[0]["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		listener[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// send appointment reminder notification to listener before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": listener[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["counsellor_id"],
		CONSTANT.ListenerType,
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	emaildata := Model.EmailBodyMessageModel{
		Name: listener[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBookCounsellorEmailBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date_time###":   UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookCounsellorTitle,
		emailBody,
		listener[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentScheduleCounsellorTextMessage,
	// 		map[string]string{
	// 			"###counsellor_name###": listener[0]["first_name"],
	// 			"###client_name###":     client[0]["first_name"],
	// 			"###date###":            UTIL.ConvertTimezone(UTIL.BuildDateTime(order[0]["date"], order[0]["time"]), listener[0]["timezone"]).Format(CONSTANT.ReadbleDateFormat),
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	listener[0]["phone"],
	// 	CONSTANT.LaterSendTextMessage,
	// )

	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
