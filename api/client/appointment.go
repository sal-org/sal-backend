package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
	"time"
)

// AppointmentsUpcoming godoc
// @Tags Client Appointment
// @Summary Get client upcoming appointments
// @Router /client/appointment/upcoming [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsUpcoming(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}
	// var localTime int
	// timeNow := UTIL.GetCurrentTime().Local()
	// if timeNow.Minute() >= 30 {
	// 	localTime = timeNow.Hour()*2 - 1
	// } else {
	// 	localTime = timeNow.Hour() * 2
	// }
	// fmt.Println(localTime)
	// local := strconv.Itoa(localTime)

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor ids to get details
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

	// get counsellor/listener/therapist details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, photo, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentSlotsUnused godoc
// @Tags Client Appointment
// @Summary Get client appointment slots
// @Router /client/appointment/slots [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentSlotsUnused(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get unused appointment slots
	appointmentSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentSlotsTable+" where client_id = ? and slots_remaining > 0 and status = "+CONSTANT.AppointmentSlotsActive, r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor ids to get details
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointmentSlots, "counsellor_id")

	// get counsellor/listener/therapist details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, photo, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	orderIDs := UTIL.ExtractValuesFromArrayMap(appointmentSlots, "order_id")
	invoice, status, ok := DB.SelectProcess("select order_id as id, payment_id, user_type, paid_amount, created_at, payment_method  from " + CONSTANT.InvoicesTable + " where order_id in ('" + strings.Join(orderIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["order_details"] = UTIL.ConvertMapToKeyMap(invoice, "id")
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["appointment_slots"] = appointmentSlots
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentsPast godoc
// @Tags Client Appointment
// @Summary Get client past appointments
// @Router /client/appointment/past [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsPast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"client_id": r.FormValue("client_id"), "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor ids to get details
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

	// get counsellor/listener/therapist details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, photo, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentDetail godoc
// @Tags Client Appointment
// @Summary Get client appointment details
// @Router /client/appointment [get]
// @Param appointment_id query string true "Appointment ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var invoice, order []map[string]string

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment order details
	order, status, ok = DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"*"}, map[string]string{"order_id": appointment[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(client) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(client[0]["email"], "@")

	ok = DB.CheckIfExists(CONSTANT.CorporatePartnersTable, map[string]string{"domain": domainName[1]})

	if !(appointment[0]["type"] == "2" || ok) {

		// get appointment slots details
		appointmentSlots, status, ok := DB.SelectSQL(CONSTANT.AppointmentSlotsTable, []string{"*"}, map[string]string{"order_id": appointment[0]["order_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		orderIDs := UTIL.ExtractValuesFromArrayMap(appointmentSlots, "order_id")
		invoice, status, ok = DB.SelectProcess("select order_id as id, payment_id, user_type, paid_amount, created_at, payment_method  from " + CONSTANT.InvoicesTable + " where order_id in ('" + strings.Join(orderIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		response["order_details"] = UTIL.ConvertMapToKeyMap(invoice, "id")
		response["appointment_slots"] = appointmentSlots[0]
	}

	response["appointment"] = appointment[0]
	response["order"] = order[0]
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentBook godoc
// @Tags Client Appointment
// @Summary Book an appointment
// @Router /client/appointment [post]
// @Param body body model.AppointmentBookRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentBook(w http.ResponseWriter, r *http.Request) {
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AppointmentBookRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment slot details
	appointmentSlot, status, ok := DB.SelectSQL(CONSTANT.AppointmentSlotsTable, []string{"*"}, map[string]string{"appointment_slot_id": body["appointment_slot_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment slot is valid
	if len(appointmentSlot) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if any slots remaining
	slotsRemaining, _ := strconv.Atoi(appointmentSlot[0]["slots_remaining"])
	if slotsRemaining <= 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.SlotCompletelyUsedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(appointmentSlot[0]["counsellor_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.RescheduleSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create appointment between counsellor and client
	appointment := map[string]string{}
	appointment["order_id"] = appointmentSlot[0]["order_id"]
	appointment["client_id"] = appointmentSlot[0]["client_id"]
	appointment["counsellor_id"] = appointmentSlot[0]["counsellor_id"]
	appointment["date"] = body["date"]
	appointment["time"] = body["time"]
	appointment["status"] = CONSTANT.AppointmentToBeStarted
	appointment["created_at"] = UTIL.GetCurrentTime().String()
	appointmentID, status, ok := DB.InsertWithUniqueID(CONSTANT.AppointmentsTable, CONSTANT.AppointmentDigits, appointment, "appointment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// reduce appointment slots by 1
	DB.ExecuteSQL("update "+CONSTANT.AppointmentSlotsTable+" set slots_remaining = slots_remaining - 1 where appointment_slot_id = ?", body["appointment_slot_id"])

	// update counsellor slots
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointmentSlot[0]["counsellor_id"],
			"date":          body["date"],
		},
		map[string]string{
			body["time"]: CONSTANT.SlotBooked,
		},
	)

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentSlotsTable+" where appointment_slot_id = ?)", body["appointment_slot_id"])
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"counsellor_id": appointmentSlot[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"listener_id": appointmentSlot[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"therapist_id": appointmentSlot[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"client_id": appointmentSlot[0]["client_id"]})

	// send appointment booking notification to therapist
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFollowUpSessionCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentFollowUpSessionCounsellorContent,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date_time###":   UTIL.BuildDateTime(body["date"], body["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		appointmentSlot[0]["counsellor_id"],
		counsellorType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	// // send appointment reminder notification to therapist before 15 min
	// UTIL.SendNotification(
	// 	CONSTANT.ClientAppointmentFollowUpSessionReminderClientHeading,
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentRemiderClientContent,
	// 		map[string]string{
	// 			"###user_name###": counsellor[0]["first_name"],
	// 			"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
	// 		},
	// 	),
	// 	appointmentSlot[0]["counsellor_id"],
	// 	counsellorType,
	// 	UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
	// 	CONSTANT.NotificationInProgress,
	// 	appointmentID,
	// )

	// send appointment booking notification to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFollowUpSessionClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentFollowUpSessionClientContent,
			map[string]string{
				"###therpist_name###": counsellor[0]["first_name"],
			},
		),
		appointmentSlot[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	// send appointment reminder notification to client before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFollowUpSessionReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		appointmentSlot[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		appointmentID,
	)

	// send appointment reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFollowUpSessionReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		appointmentSlot[0]["counsellor_id"],
		counsellorType,
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).String(),
		CONSTANT.NotificationInProgress,
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
				"###userName###":  counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
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
				"###user_name###": counsellor[0]["first_name"],
				"###userName###":  client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// filepath_text := "htmlfile/emailmessagebody.html"

	// // send email for counsellor
	// emaildata := Model.EmailBodyMessageModel{
	// 	Name: counsellor[0]["first_name"],
	// 	Message: UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentFollowUpSessionCounsellorEmailBody,
	// 		map[string]string{
	// 			"###client_name###": client[0]["first_name"],
	// 			"###date_time###":   UTIL.BuildDateTime(body["date"], body["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
	// 		},
	// 	),
	// }

	// emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// // email for counsellor
	// UTIL.SendEmail(
	// 	CONSTANT.ClientAppointmentFollowUpSessionCounsellorTitle,
	// 	emailBody,
	// 	counsellor[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBookClientEmailBody,
			map[string]string{
				"###therpist_name###": counsellor[0]["first_name"],
				"###date###":          body["date"],
				"###time###":          UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookCounsellorTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// send messsage to therpists for Appointment Confirmation
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentConfirmationTextMessage,
			map[string]string{
				"###userName###":  counsellor[0]["first_name"],
				"###user_Name###": client[0]["first_name"],
				"###date###":      body["date"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).UTC().String(),
		appointmentID,
		CONSTANT.InstantSendEmailMessage,
	)

	// send messsage to therpists for Appointment Confirmation
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentConfirmationTextMessage,
			map[string]string{
				"###userName###":  client[0]["first_name"],
				"###user_Name###": counsellor[0]["first_name"],
				"###date###":      body["date"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).UTC().String(),
		appointmentID,
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentReschedule godoc
// @Tags Client Appointment
// @Summary Reschedule an appointment
// @Router /client/appointment [put]
// @Param appointment_id query string true "Appointment ID to be rescheduled"
// @Param body body model.AppointmentRescheduleRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentReschedule(w http.ResponseWriter, r *http.Request) {
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AppointmentRescheduleRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is to be started
	if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment rescheduled times exceeded
	reschedules, _ := strconv.Atoi(appointment[0]["times_rescheduled"])
	if reschedules >= CONSTANT.MaximumAppointmentReschedule {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentCantRescheduleMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is alteast after 4 hours
	if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now().Add(330*time.Minute).UTC()).Hours() < 4 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.Reschedule4HoursMinimumMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(appointment[0]["counsellor_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.RescheduleSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update counsellor slots
	// remove previous slot
	// date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// // get schedules for a weekday
	// schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update counsellor availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			// this is for cancel slot menthod  UTIL.CheckIfScheduleAvailable(schedules, appointment[0]["time"])
			appointment[0]["time"]: CONSTANT.SlotAvailable, // update availability to the latest one
		},
	)

	// update slot
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          body["date"],
		},
		map[string]string{
			body["time"]: CONSTANT.SlotBooked,
		},
	)

	// update appointment date and time
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"date":        body["date"],
			"time":        body["time"],
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)
	// update rescheduled times
	DB.ExecuteSQL("update "+CONSTANT.AppointmentsTable+" set times_rescheduled = times_rescheduled + 1 where appointment_id = ?", r.FormValue("appointment_id"))

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "email"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "email"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])

	// remove all previous message for client
	UTIL.RemoveMessage(r.FormValue("appointment_id"), client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(r.FormValue("appointment_id"), counsellor[0]["phone"])

	// send appointment reschedule notification to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentRescheduleClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleClientContent,
			map[string]string{
				"###date_time###":     UTIL.BuildDateTime(body["date"], body["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###therapistname###": counsellor[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	// send appointment reminder notification to client before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		r.FormValue("appointment_id"),
	)

	// Send to appointment Reminder SMS to client
	// send at 15 min before of appointment
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###userName###":  counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.LaterSendTextMessage,
	)

	// send email
	filepath_text := "htmlfile/emailmessagebody.html"

	// send client email body
	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleClientEmailBody,
			map[string]string{
				"###date###":           body["date"],
				"###time###":           UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
				"###therpists_name###": counsellor[0]["first_name"],
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentRescheduleClientTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	// send appointment reschedule notification to counsellor
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentRescheduleCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleCounsellorContent,
			map[string]string{},
		),
		appointment[0]["counsellor_id"],
		counsellorType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	// send appointment reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		appointment[0]["counsellor_id"],
		counsellorType,
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		r.FormValue("appointment_id"),
	)

	// Send to appointment Reminder SMS to counsellor
	// send at 15 min before of appointment
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": counsellor[0]["first_name"],
				"###userName###":  client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).Add(-15*time.Minute).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.LaterSendTextMessage,
	)

	// emaildata := Model.EmailBodyMessageModel{
	// 	Name:    counsellor[0]["first_name"],
	// 	Message: CONSTANT.ClientAppointmentRescheduleCounsellorEmailBody,
	// }

	// send counsellor email body
	emaildata := Model.EmailBodyMessageModel{
		Name: counsellor[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleCounsellorEmailBody,
			map[string]string{
				"###first_name###": client[0]["first_name"],
				"###date###":       body["date"],
				"###time###":       UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentRescheduleClientTitle,
		emailBody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	//Send to Client
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleClientTextMeassge,
			map[string]string{
				"###clientName###":    client[0]["first_name"],
				"###therapistName###": counsellor[0]["first_name"],
				"###date###":          body["date"],
				"###time###":          UTIL.GetTimeFromTimeSlotIN12Hour(body["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.InstantSendTextMessage,
	)

	//send to counsellor
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleClientToCounsellorTextMeassge,
			map[string]string{
				"###counsellorName###": counsellor[0]["first_name"],
				"###clientName###":     client[0]["first_name"],
				"###date###":           body["date"],
				"###time###":           body["time"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(body["date"], body["time"]).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.InstantSendTextMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentCancel godoc
// @Tags Client Appointment
// @Summary Cancel an appointment
// @Router /client/appointment [delete]
// @Param appointment_id query string true "Appointment ID to be cancelled"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if appointment is to be started
	if strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentUserCancelled) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentCancelByUserMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if appointment is to be started
	if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}
	// // appointment can cancelled only after min reschedules
	// reschedules, _ := strconv.Atoi(appointment[0]["times_rescheduled"])
	// if reschedules < CONSTANT.MaximumAppointmentReschedule {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentCantCancelMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// update counsellor slots
	// remove previous slot
	// date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedules for a weekday
	// schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update counsellor availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			// this is for cancel slot menthod  UTIL.CheckIfScheduleAvailable(schedules, appointment[0]["time"])
			appointment[0]["time"]: CONSTANT.SlotAvailable, // update availability to the latest one
		},
	)

	// update appointment date and time
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentUserCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// refund amount
	// check if appointment is alteast after 4 hours
	// if below 4 hours, charges are 100% that means dont refund anything
	// if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours() >= 4 {
	// 	charges := CONSTANT.ClientAppointmentCancellationCharges
	// 	// get invoice details
	// 	invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount", "payment_id", "invoice_id"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// 	if len(invoice) > 0 {
	// 		// invoice is available => amount is paid, order is not free
	// 		// get order details
	// 		order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}
	// 		paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
	// 		// discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
	// 		// amountAfterDiscount := paidAmount - discount
	// 		if paidAmount > 0 { // refund only if amount paid
	// 			paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
	// 			refundedAmount, _ := strconv.ParseFloat(DB.QueryRowSQL("select sum(refunded_amount) from "+CONSTANT.RefundsTable+" where invoice_id = '"+invoice[0]["invoice_id"]+"'"), 64)
	// 			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)
	// 			cancellationCharges := (paidAmount / slotsBought) * charges
	// 			refundAmount := (paidAmount / slotsBought) - cancellationCharges // remove from paid amount only, not from amount after discount, as discussed
	// 			if refundAmount+refundedAmount <= paidAmount {
	// 				// refunded amount will be less than paid amount
	// 				DB.InsertWithUniqueID(CONSTANT.RefundsTable, CONSTANT.RefundDigits, map[string]string{
	// 					"invoice_id":      invoice[0]["invoice_id"],
	// 					"payment_id":      invoice[0]["payment_id"],
	// 					"refunded_amount": strconv.FormatFloat(refundAmount, 'f', 2, 64),
	// 					"status":          CONSTANT.RefundInProgress,
	// 					"created_at":      UTIL.GetCurrentTime().String(),
	// 				}, "refund_id")
	// 			}
	// 		}
	// 	}
	// }

	// add penalty for counsellor for cancelling
	// add to counsellor payments
	// get invoice details

	if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(UTIL.ConvertTimezone(UTIL.GetCurrentTime(), "330")).Hours() <= 4 {
		invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount"}, map[string]string{"order_id": appointment[0]["order_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(invoice) > 0 {
			// get order details
			order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
			discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
			amountBeforeDiscount := paidAmount + discount
			if amountBeforeDiscount > 0 { // add only if amount paid
				slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

				amountFor1Session := amountBeforeDiscount / slotsBought // for 1 counselling session
				cancellationCharges := amountFor1Session * CONSTANT.CounsellorCancellationCharges

				DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
					"counsellor_id": appointment[0]["counsellor_id"],
					"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
					"description":   "Cancellation",
					"amount":        strconv.FormatFloat(cancellationCharges, 'f', 2, 64),
					"status":        CONSTANT.PaymentActive,
					"created_at":    UTIL.GetCurrentTime().String(),
				}, "payment_id")
			}
		} else {

			//UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours()

			// get counsellor details
			counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"corporate_price"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			amount, _ := strconv.ParseFloat(counsellor[0]["corporate_price"], 64)

			cancellationCharges := amount * CONSTANT.CounsellorCancellationCharges

			DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
				"counsellor_id": appointment[0]["counsellor_id"],
				"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
				"description":   "Client Cancellation",
				"amount":        strconv.FormatFloat(cancellationCharges, 'f', 2, 64),
				"status":        CONSTANT.PaymentActive,
				"created_at":    UTIL.GetCurrentTime().String(),
			}, "payment_id")

		}
	}

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})

	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])

	// remove all previous message for client
	UTIL.RemoveMessage(r.FormValue("appointment_id"), client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(r.FormValue("appointment_id"), counsellor[0]["phone"])

	// send appointment cancel notification, email to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelClientContent,
			map[string]string{
				"###datetime###":      UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###therapistname###": counsellor[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	// send email
	filepath_text := "htmlfile/emailmessagebody.html"

	// send client email body
	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelClientBody,
			map[string]string{
				"###therapist_name###": counsellor[0]["first_name"],
				"###date###":           appointment[0]["date"],
				"###time###":           UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentCancelClientTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	// send appointment cancel notification to counsellor
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentCancelCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelCounsellorContent,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date_time###":   UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		appointment[0]["counsellor_id"],
		counsellorType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	// send email for therapist
	emaildata := Model.EmailBodyMessageModel{
		Name: counsellor[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelCounsellorEmailBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date###":        appointment[0]["date"],
				"###time###":        UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for therapist
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentCancelCounsellorTitle,
		emailBody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// Client Cancel the Appointment to send text message to counsellor
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancellationToCounsellorTextMessage,
			map[string]string{
				"###userName###":  counsellor[0]["first_name"],
				"###user_Name###": client[0]["first_name"],
				"###date###":      appointment[0]["date"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.InstantSendTextMessage,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancellationToCounsellorTextMessage,
			map[string]string{
				"###userName###":  client[0]["first_name"],
				"###user_Name###": counsellor[0]["first_name"],
				"###date###":      appointment[0]["date"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.InstantSendTextMessage,
	)
	// fmt.Println("Status:" + CONSTANT.StatusCodeOk)
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentBulkCancel godoc
// @Tags Client Appointment
// @Summary Cancel bulk appointments
// @Router /client/appointment/bulk [delete]
// @Param appointment_slot_id query string true "Appointment slot ID to be cancelled"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentBulkCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment slots details
	appointmentSlots, status, ok := DB.SelectSQL(CONSTANT.AppointmentSlotsTable, []string{"*"}, map[string]string{"appointment_slot_id": r.FormValue("appointment_slot_id"), "status": CONSTANT.AppointmentSlotsActive})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment slots is valid
	if len(appointmentSlots) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentSlotNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if any slots available
	slotsRemaining, _ := strconv.Atoi(appointmentSlots[0]["slots_remaining"])
	if slotsRemaining <= 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// refund amount
	charges := CONSTANT.ClientAppointmentBulkCancellationCharges
	// get invoice details
	invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount", "payment_id", "invoice_id"}, map[string]string{"order_id": appointmentSlots[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(invoice) > 0 {
		// invoice is available => amount is paid, order is not free
		paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
		discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
		amountAfterDiscount := paidAmount - discount
		if amountAfterDiscount > 0 { // refund only if amount paid
			paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
			refundedAmount, _ := strconv.ParseFloat(DB.QueryRowSQL("select sum(refunded_amount) from "+CONSTANT.RefundsTable+" where invoice_id = '"+invoice[0]["invoice_id"]+"'"), 64)
			slotsBought, _ := strconv.Atoi(appointmentSlots[0]["slots_bought"])
			slotsRemaining, _ = strconv.Atoi(appointmentSlots[0]["slots_remaining"])
			cancellationCharges := (amountAfterDiscount / float64(slotsBought)) * charges
			refundAmount := ((paidAmount / float64(slotsBought)) - cancellationCharges) * float64(slotsRemaining) // remove from paid amount only, not from amount after discount, as discussed
			if refundAmount+refundedAmount <= paidAmount {
				// refunded amount will be less than paid amount
				DB.InsertWithUniqueID(CONSTANT.RefundsTable, CONSTANT.RefundDigits, map[string]string{
					"invoice_id":      invoice[0]["invoice_id"],
					"payment_id":      invoice[0]["payment_id"],
					"refunded_amount": strconv.FormatFloat(refundAmount, 'f', 2, 64),
					"status":          CONSTANT.RefundInProgress,
					"created_at":      UTIL.GetCurrentTime().String(),
				}, "refund_id")
			}
		}
	}

	// update slots remaining to 0
	DB.UpdateSQL(CONSTANT.AppointmentSlotsTable, map[string]string{
		"appointment_slot_id": r.FormValue("appointment_slot_id"),
	}, map[string]string{
		"slots_remaining": "0",
		"status":          CONSTANT.AppointmentUserCancelled,
		"modified_at":     UTIL.GetCurrentTime().String(),
	})

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentSlotsTable+" where appointment_slot_id = ?)", r.FormValue("appointment_slot_id"))
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone"}, map[string]string{"counsellor_id": appointmentSlots[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone"}, map[string]string{"listener_id": appointmentSlots[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone"}, map[string]string{"therapist_id": appointmentSlots[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone"}, map[string]string{"client_id": appointmentSlots[0]["client_id"]})

	// send appointment cancel notification to client
	UTIL.SendNotification(
		CONSTANT.ClientBulkAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientBulkAppointmentCancelClientContent,
			map[string]string{
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		appointmentSlots[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_slot_id"),
	)

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBulkCancelClientEmailBody,
			map[string]string{},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBulkCancelClientTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// Client Cancel the Appointment to send text message to counsellor
	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentCancellationToCounsellorTextMessage,
	// 		map[string]string{
	// 			"###client_name###": client[0]["first_name"],
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	counsellor[0]["phone"],
	// 	CONSTANT.LaterSendTextMessage,
	// )

	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentCancellationTextMessage,
	// 		map[string]string{
	// 			"###client_name###":     client[0]["first_name"],
	// 			"###slot_bought###":     appointmentSlots[0]["slots_remaining"],
	// 			"###counsellor_name###": counsellor[0]["first_name"],
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	client[0]["phone"],
	// 	CONSTANT.LaterSendTextMessage,
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentRatingAdd godoc
// @Tags Client Appointment
// @Summary Rate the appointment
// @Router /client/appointment/rate [post]
// @Param body body model.AppointmentRatingAdd true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentRatingAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if access token is valid, not expired
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AppointmentRatingAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if client exists
	if !DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": body["client_id"]}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update rating to appointment
	status, ok := DB.UpdateSQL(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": body["appointment_id"], "client_id": body["client_id"], "counsellor_id": body["counsellor_id"]}, map[string]string{"rating": body["rating"], "rating_types": body["rating_types"], "rating_comment": body["rating_comment"], "modified_at": UTIL.GetCurrentTime().String()})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	status, ok = DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable, map[string]string{"appointment_id": body["appointment_id"], "client_id": body["client_id"], "counsellor_id": body["counsellor_id"]}, map[string]string{"rating": body["rating"], "rating_types": body["rating_types"], "rating_comment": body["rating_comment"], "modified_at": UTIL.GetCurrentTime().String()})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", body["appointment_id"])

	// get counsellor type and update their ratings
	counsellorRating, _, _ := DB.SelectProcess("select rating from "+CONSTANT.AppointmentsTable+" where counsellor_id = ?", body["counsellor_id"])

	counsellorCount, _, _ := DB.SelectProcess("select count(rating) as cnt from "+CONSTANT.AppointmentsTable+" where counsellor_id = ?", body["counsellor_id"])

	avg := UTIL.AvgRatingFromula(counsellorRating, counsellorCount[0]["cnt"])

	switch counsellorType {
	case CONSTANT.CounsellorType:
		DB.ExecuteSQL("update "+CONSTANT.CounsellorsTable+" set total_rating = total_rating + 1, average_rating = ? where counsellor_id = ?", avg, body["counsellor_id"])
		break
	case CONSTANT.ListenerType:
		DB.ExecuteSQL("update "+CONSTANT.ListenersTable+" set total_rating = total_rating + 1, average_rating = ? where listener_id = ?", avg, body["counsellor_id"])
		break
	case CONSTANT.TherapistType:
		DB.ExecuteSQL("update "+CONSTANT.TherapistsTable+" set total_rating = total_rating + 1, average_rating = ? where therapist_id = ?", avg, body["counsellor_id"])
		break
	}

	UTIL.SendNotification(
		CONSTANT.ClientAppointmentGivenFeedbackHeading,
		CONSTANT.ClientAppointmentFeedbackGivenContent,
		body["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

// DownloadReceipt godoc
// @Tags Client Appointment
// @Summary Get invoice download Receipt
// @Router /client/appointment/download [get]
// @Param invoice_id query string true "Logged in invoice ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func DownloadReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var fileName string

	invoice, status, ok := DB.SelectProcess("select * from "+CONSTANT.InvoicesTable+" where invoice_id = ? ", r.FormValue("invoice_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if DB.CheckIfExists(CONSTANT.ReceiptTable, map[string]string{"invoice_id": invoice[0]["invoice_id"]}) {

		receipt, _, _ := DB.SelectSQL(CONSTANT.ReceiptTable, []string{"*"}, map[string]string{"invoice_id": invoice[0]["invoice_id"]})
		fileName = receipt[0]["pdf"]
	} else {

		receiptdata := UTIL.BuildDate(invoice[0]["created_at"])

		order, _, _ := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought", "counsellor_id", "type"}, map[string]string{"order_id": invoice[0]["order_id"]})

		var sprice, tprice, discount, paid_amount string

		paid_Amount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)

		discount_value, _ := strconv.ParseFloat(invoice[0]["discount"], 64)

		slots_bought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

		actual_Amount := paid_Amount + discount_value

		price := actual_Amount / slots_bought

		sprice = strconv.FormatFloat(price, 'f', 2, 64)

		tprice = strconv.FormatFloat(actual_Amount, 'f', 2, 64)

		discount = strconv.FormatFloat(discount_value, 'f', 2, 64)

		paid_amount = strconv.FormatFloat(paid_Amount, 'f', 2, 64)

		// This is issue because counsellor anytime update yours profile

		/*if order[0]["type"] == CONSTANT.CounsellorType {

			counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"price", "multiple_sessions", "price_3", "price_5"}, map[string]string{"counsellor_id": order[0]["counsellor_id"]})

			switch order[0]["slots_bought"] {
			case "1":
				sprice = counsellor[0]["price"]
				tprice = counsellor[0]["price"]
			case "3":
				sprice = counsellor[0]["multiple_sessions"]
				tprice = counsellor[0]["price_3"]
			case "5":
				sprice = counsellor[0]["multiple_sessions"]
				tprice = counsellor[0]["price_5"]
			}
		}

		if order[0]["type"] == CONSTANT.TherapistType {

			counsellor, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"price", "multiple_sessions", "price_3", "price_5"}, map[string]string{"therapist_id": order[0]["counsellor_id"]})

			switch order[0]["slots_bought"] {
			case "1":
				sprice = counsellor[0]["price"]
				tprice = counsellor[0]["price"]
			case "3":
				sprice = counsellor[0]["multiple_sessions"]
				tprice = counsellor[0]["price_3"]
			case "5":
				sprice = counsellor[0]["multiple_sessions"]
				tprice = counsellor[0]["price_5"]
			}
		}*/

		data := Model.EmailDataForPaymentReceipt{
			Date:        receiptdata,
			ReceiptNo:   invoice[0]["id"],
			ReferenceNo: invoice[0]["payment_id"],
			SPrice:      sprice,
			Qty:         order[0]["slots_bought"],
			Total:       tprice,
			//SessionsType: CONSTANT.AppointmentSessionsTypeForReceipt,
			TPrice:   tprice,
			CouponC:  invoice[0]["coupon_code"],
			Discount: discount,
			TotalP:   paid_amount,
		}

		filePath := "htmlfile/receiptClove.html"

		emailbody, ok := UTIL.GetHTMLTemplateForReceipt(data, filePath)
		if !ok {
			fmt.Println("html body not create ")
		}

		created, ok := UTIL.GeneratePdf(emailbody, "pdffile/example1.pdf") // name created,

		if !ok {
			fmt.Println("Pdf is not created")
		}

		s3Path := "receipt"
		filename := "example1.pdf"

		name, uploaded := UTIL.UploadToS3File(CONFIG.S3Bucket, s3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(filename), CONSTANT.S3PublicRead, created)
		if !uploaded {
			fmt.Println("UploadFile")
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name

		receipt := map[string]string{}

		receipt["client_id"] = order[0]["client_id"]
		receipt["invoice_id"] = invoice[0]["invoice_id"]
		receipt["pdf"] = fileName
		receipt["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok := DB.InsertWithUniqueID(CONSTANT.ReceiptTable, CONSTANT.ReceiptDigits, receipt, "receipt_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}
	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// Cancellation Reason godoc
// @Tags Client Appointment
// @Summary Cancellation Region client
// @Router /client/appointment/cancellationreason [put]
// @Param appointment_id query string true "Appointment ID, Appointment Slot ID to Cancellation Reason"
// @Param sessions query string true "Single Sessions(2), Multiple Sessions(3)"
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

	if strings.EqualFold(r.FormValue("sessions"), CONSTANT.EventStarted) {
		exists := DB.CheckIfExists(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		status, ok := DB.UpdateSQL(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")}, map[string]string{"cancellation_reason": body["cancellation_reason"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {

		exists := DB.CheckIfExists(CONSTANT.AppointmentSlotsTable, map[string]string{"appointment_slot_id": r.FormValue("appointment_id")})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		status, ok := DB.UpdateSQL(CONSTANT.AppointmentSlotsTable, map[string]string{"appointment_slot_id": r.FormValue("appointment_id")}, map[string]string{"cancellation_reason": body["cancellation_reason"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// Generate Agora Token godoc
// @Tags Client Appointment
// @Summary Get Agora Token
// @Router /client/appointment/agoratoken [get]
// @Param appointment_id query string true "Appointment ID or Order ID is equal to Channel Name"
// @Param session query string true "Individual(1), Cafe(2)"
// @Param type query string true "Publisher(1), Subscriber(2)"
// @Param user_type query string true "Counsellor(1) , Client(2)"
// @Security JWTAuth
// @Produce json
// @Success 200
func GenerateAgoraToken(w http.ResponseWriter, r *http.Request) {

	var response = make(map[string]interface{})

	var roleStr, agora_token, uidStr, channelName string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	uidStr = generateRandomID()

	if r.FormValue("session") == "1" {
		exists := DB.CheckIfExists(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		channelName = r.FormValue("appointment_id")
		if r.FormValue("type") == "1" {
			roleStr = CONSTANT.RolePublisher
		} else if r.FormValue("type") == "2" {
			roleStr = CONSTANT.RoleSubscriber
		} else {
			roleStr = "attended"
		}

		//uidStr = generateRandomID()
		// For demonstration purposes the expiry time is set to 7200 seconds = 2 hours. This shows you the automatic token renew actions of the client.
		expireTimeInSeconds := uint32(7200)
		// Get current timestamp.
		currentTimestamp := uint32(time.Now().UTC().Unix())
		// Timestamp when the token expires.
		expireTimestamp := currentTimestamp + expireTimeInSeconds

		token, err := UTIL.GenerateAgoraRTCToken(channelName, roleStr, uidStr, expireTimestamp)
		if err != nil {
			UTIL.SetReponse(w, "", "", CONSTANT.ShowDialog, response)
			return
		}
		agora_token = token
	} else if r.FormValue("session") == "2" {
		exists := DB.CheckIfExists(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("appointment_id")})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		channelName = r.FormValue("appointment_id")
		if r.FormValue("type") == "1" {
			roleStr = CONSTANT.RolePublisher
		} else if r.FormValue("type") == "2" {
			roleStr = CONSTANT.RoleSubscriber
		} else {
			roleStr = "attended"
		}

		// For demonstration purposes the expiry time is set to 7200 seconds = 2 hours. This shows you the automatic token renew actions of the client.
		expireTimeInSeconds := uint32(7200)
		// Get current timestamp.
		currentTimestamp := uint32(time.Now().UTC().Unix())
		// Timestamp when the token expires.
		expireTimestamp := currentTimestamp + expireTimeInSeconds

		token, err := UTIL.GenerateAgoraRTCToken(channelName, roleStr, uidStr, expireTimestamp)
		if err != nil {
			UTIL.SetReponse(w, "500", "Server Error", CONSTANT.ShowDialog, response)
			return
		}
		agora_token = token
	}

	agora := map[string]string{}

	exists := DB.CheckIfExists(CONSTANT.AgoraTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !exists {

		expireTimeInSeconds := uint32(7200)
		// Get current timestamp.
		currentTimestamp := uint32(time.Now().UTC().Unix())
		// Timestamp when the token expires.
		expireTimestmp := currentTimestamp + expireTimeInSeconds

		uidSt := generateRandomID()
		channelNa := r.FormValue("appointment_id")

		tokenForResource, err := UTIL.GenerateAgoraRTCToken(channelNa, roleStr, uidSt, expireTimestmp)
		if err != nil {
			fmt.Println("Ressource Token not generated")
		}

		resourceid, err := UTIL.BasicAuthorization(channelNa, uidSt)
		if err != nil {
			fmt.Println("resource id not generated for recording file")
		}

		agora["appointment_id"] = channelNa
		agora["uid"] = uidSt
		agora["token"] = tokenForResource
		agora["resource_id"] = resourceid
		agora["status"] = CONSTANT.AgoraResourceID
		agora["created_at"] = UTIL.GetCurrentTime().String()
		_, status, ok := DB.InsertWithUniqueID(CONSTANT.AgoraTable, CONSTANT.AgoraDigits, agora, "agora_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// if r.FormValue("user_type") == "1" {
	// 	exists := DB.CheckIfExists(CONSTANT.AgoraTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	// 	if exists {
	// 		DB.UpdateSQL(CONSTANT.AgoraTable,
	// 			map[string]string{
	// 				"appointment_id": r.FormValue("appointment_id"),
	// 			},
	// 			map[string]string{
	// 				"uid":         uidStr,
	// 				"token":       agora_token,
	// 				"resource_id": "",
	// 				"status":      CONSTANT.AgoraResourceID2,
	// 				"modified_at": UTIL.GetCurrentTime().String(),
	// 			},
	// 		)
	// 	} else {
	// 		agora["appointment_id"] = channelName
	// 		agora["uid"] = uidStr
	// 		agora["token"] = agora_token
	// 		agora["resource_id"] = ""
	// 		agora["status"] = CONSTANT.AgoraResourceID
	// 		agora["created_at"] = UTIL.GetCurrentTime().String()
	// 		_, status, ok := DB.InsertWithUniqueID(CONSTANT.AgoraTable, CONSTANT.AgoraDigits, agora, "agora_id")
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 	}
	// } else if r.FormValue("user_type") == "2" {
	// 	resourceid, err := UTIL.BasicAuthorization(channelName, uidStr)
	// 	if err != nil {
	// 		fmt.Println("resource id not generated for recording file")
	// 		return
	// 	}
	// 	exists := DB.CheckIfExists(CONSTANT.AgoraTable, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	// 	if exists {
	// 		DB.UpdateSQL(CONSTANT.AgoraTable,
	// 			map[string]string{
	// 				"appointment_id": r.FormValue("appointment_id"),
	// 			},
	// 			map[string]string{
	// 				"uid1":         uidStr,
	// 				"token1":       agora_token,
	// 				"resource_id1": resourceid,
	// 				"status":       CONSTANT.AgoraResourceID2,
	// 				"modified_at":  UTIL.GetCurrentTime().String(),
	// 			},
	// 		)
	// 	} else {
	// 		agora["appointment_id"] = channelName
	// 		agora["uid1"] = uidStr
	// 		agora["token1"] = agora_token
	// 		agora["resource_id1"] = resourceid
	// 		agora["status"] = CONSTANT.AgoraResourceID
	// 		agora["created_at"] = UTIL.GetCurrentTime().String()
	// 		_, status, ok := DB.InsertWithUniqueID(CONSTANT.AgoraTable, CONSTANT.AgoraDigits, agora, "agora_id")
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 	}
	// }

	response["token"] = agora_token
	response["UID"] = uidStr

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func generateRandomID() string {
	const randomIDdigits = "123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = randomIDdigits[rand.Intn(len(randomIDdigits))]
	}
	return string(b)
}

// AppointmentStart godoc
// @Tags Client Appointment
// @Summary Start an appointment
// @Router /client/appointment/start [put]
// @Param appointment_id query string true "Appointment ID to be started"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	agora, status, ok := DB.SelectSQL(CONSTANT.AgoraTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(agora) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	if appointment[0]["client_started_at"] == "" {
		// update appointment as started
		DB.UpdateSQL(CONSTANT.AppointmentsTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"status":            CONSTANT.AppointmentStarted,
				"client_started_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

	if len(agora[0]["sid"]) == 0 {

		sid, err := UTIL.AgoraRecordingCallStart(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["token"], agora[0]["resource_id"])
		if err != nil {
			fmt.Println("cloud recording not start")
		}

		DB.UpdateSQL(CONSTANT.AgoraTable,
			map[string]string{
				"agora_id": agora[0]["agora_id"],
			},
			map[string]string{
				"sid":         sid,
				"status":      CONSTANT.AgoraCallStart2,
				"modified_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

	var counsellorName string

	switch appointment[0]["type"] {
	case "1":
		counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", appointment[0]["counsellor_id"])

	case "2":
		counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.ListenersTable+" where listener_id = ?", appointment[0]["counsellor_id"])

	case "4":
		counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"])

	}

	// send appointment join the call notification to Therapists
	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentHasBeenStartedHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentHasBeenStartedContent,
			map[string]string{
				"###therapistname###": counsellorName,
				"###clientname###":    DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
			},
		),
		appointment[0]["counsellor_id"],
		CONSTANT.CounsellorType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentEnd godoc
// @Tags Client Appointment
// @Summary End an appointment
// @Router /client/appointment/end [put]
// @Param appointment_id query string true "Appointment ID to be ended"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	agora, status, ok := DB.SelectSQL(CONSTANT.AgoraTable, []string{"*"}, map[string]string{"appointment_id": appointment[0]["appointment_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update appointment as completed
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":          CONSTANT.AppointmentCompleted,
			"client_ended_at": UTIL.GetCurrentTime().String(),
		},
	)

	if appointment[0]["client_ended_at"] != "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyCompletedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// type 2 is listerner
	if appointment[0]["type"] != "2" {
		if appointment[0]["started_at"] == "" && appointment[0]["ended_at"] == "" {
			invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount"}, map[string]string{"order_id": appointment[0]["order_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			if len(invoice) > 0 {
				// get order details
				order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
				paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
				discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
				amountBeforeDiscount := paidAmount + discount
				if amountBeforeDiscount > 0 { // add only if amount paid
					slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

					amountFor1Session := amountBeforeDiscount / slotsBought // for 1 counselling session
					cancellationCharges := amountFor1Session

					DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
						"counsellor_id": appointment[0]["counsellor_id"],
						"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
						"description":   "Therapist No Show",
						"amount":        strconv.FormatFloat(-cancellationCharges, 'f', 2, 64),
						"status":        CONSTANT.PaymentActive,
						"created_at":    UTIL.GetCurrentTime().String(),
					}, "payment_id")
				}
			} else {

				//UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours()

				// get counsellor details
				counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"corporate_price"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
					"counsellor_id": appointment[0]["counsellor_id"],
					"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
					"description":   "Therapist No Show",
					"amount":        "-" + counsellor[0]["corporate_price"],
					"status":        CONSTANT.PaymentActive,
					"created_at":    UTIL.GetCurrentTime().String(),
				}, "payment_id")

			}
		}
	}

	// codeStatus, err := UTIL.CallStatus(agora[0]["resource_id1"], agora[0]["sid1"])
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "Internal Server Error", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// fmt.Println(codeStatus)

	if len(agora[0]["fileNameInMp4"]) == 0 && len(agora[0]["fileNameInM3U8"]) == 0 {
		fileNameInMP4, fileNameInM3U8, _ := UTIL.AgoraRecordingCallStop(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["resource_id"], agora[0]["sid"])
		// if err != nil {
		// 	UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.AgoraCallMessage, CONSTANT.ShowDialog, response)
		// 	return
		// }
		DB.UpdateSQL(CONSTANT.AgoraTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"fileNameInMp4":  fileNameInMP4,
				"fileNameInM3U8": fileNameInM3U8,
				"status":         CONSTANT.AgoraCallStop2,
				"modified_at":    UTIL.GetCurrentTime().String(),
			},
		)

		DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"counsellor_mp4": fileNameInMP4,
				"status":         CONSTANT.QualityCheckLinkInsert,
				"modified_at":    UTIL.GetCurrentTime().String(),
			},
		)
	}

	var counsellor []map[string]string
	switch appointment[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "email"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "email"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
		break
	}

	// send appointment ended notification and rating to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFeedbackHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentFeedbackContent,
			map[string]string{
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CouponGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get counsellor details
	couponsCode, status, ok := DB.SelectProcess("select *  from "+CONSTANT.CouponsTable+" where client_id = ? or client_id = '' and status = '1' and start_by < '"+UTIL.GetCurrentTime().String()+"' and '"+UTIL.GetCurrentTime().String()+"' < end_by ", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["coupons"] = couponsCode

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
