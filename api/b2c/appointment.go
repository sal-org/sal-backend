package b2c

import (
	"math/rand"
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

// AppointmentsUpcoming godoc
// @Tags Client Appointment
// @Summary Get client upcoming appointments
// @Router /client/appointment/upcoming [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsUpcoming(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

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

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", body["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for i := 0; i < len(appointments); i++ {

		// get counsellor name
		counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", appointments[i]["appointment_id"])
		switch counsellorType {
		case CONSTANT.CounsellorType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "photo"}, map[string]string{"counsellor_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]
		case CONSTANT.ListenerType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "last_name", "photo"}, map[string]string{"listener_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]
		case CONSTANT.TherapistType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "photo"}, map[string]string{"therapist_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]

		}
	}

	// // get counsellor ids to get details
	// counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

	// // get counsellor/listener/therapist details
	// counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, photo, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	if encrypt == "" {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}

// AppointmentsPast godoc
// @Tags Client Appointment
// @Summary Get client past appointments
// @Router /client/appointment/past [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsPast(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

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

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"client_id": body["client_id"], "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for i := 0; i < len(appointments); i++ {

		// get counsellor name
		counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", appointments[i]["appointment_id"])
		switch counsellorType {
		case CONSTANT.CounsellorType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "photo"}, map[string]string{"counsellor_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]
		case CONSTANT.ListenerType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "last_name", "photo"}, map[string]string{"listener_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]
		case CONSTANT.TherapistType:
			counsellor, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "photo"}, map[string]string{"therapist_id": appointments[i]["counsellor_id"]})
			appointments[i]["counsellor_name"] = counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]
			appointments[i]["counsellor_photo"] = counsellor[0]["photo"]

		}
	}

	// response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	if encrypt == "" {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}

// AppointmentDetail godoc
// @Tags Client Appointment
// @Summary Get client appointment details
// @Router /client/appointment [get]
// @Param appointment_id query string true "Appointment ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentDetail(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var order []map[string]string

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
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

	response["appointment"] = appointment[0]
	response["order"] = order[0]
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	if encrypt == "" {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
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
func AppointmentReschedule(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.AppointmentRescheduleRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
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

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(appointment[0]["counsellor_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.RescheduleSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

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
			"appointment_id": body["appointment_id"],
		},
		map[string]string{
			"date":        body["date"],
			"time":        body["time"],
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)
	// update rescheduled times
	DB.ExecuteSQL("update "+CONSTANT.AppointmentsTable+" set times_rescheduled = times_rescheduled + 1 where appointment_id = ?", body["appointment_id"])

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", body["appointment_id"])
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "email"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
		// break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "email"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
		// break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
		// break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous message for client
	UTIL.RemoveMessage(body["appointment_id"], client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(body["appointment_id"], counsellor[0]["phone"])

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
		body["appointment_id"],
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
		body["appointment_id"],
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
		body["appointment_id"],
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
		body["appointment_id"],
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
func AppointmentCancel(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
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

	if time.Until(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Add(-time.Minute*330)).Hours() <= 1 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentCancellationAllowedWithIn1Hour, CONSTANT.ShowDialog, response)
		return
	}

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
			"appointment_id": body["appointment_id"],
		},
		map[string]string{
			"cancellation_reason": body["cancellation_reason"],
			"status":              CONSTANT.AppointmentUserCancelled,
			"modified_at":         UTIL.GetCurrentTime().String(),
		},
	)

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", body["appointment_id"])
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})

	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous message for client
	UTIL.RemoveMessage(body["appointment_id"], client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(body["appointment_id"], counsellor[0]["phone"])

	// send email
	filepath_text := "htmlfile/emailmessagebody.html"

	// send client email body
	emaildata1 := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelClientBody,
			map[string]string{
				"###therapist_name###": counsellor[0]["first_name"],
				"###date###":           UTIL.BuildOnlyDate(appointment[0]["date"]),
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

	// send email for therapist
	emaildata := Model.EmailBodyMessageModel{
		Name: counsellor[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelCounsellorEmailBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date###":        UTIL.BuildOnlyDate(appointment[0]["date"]),
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
				"###date###":      UTIL.BuildOnlyDate(appointment[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
		body["appointment_id"],
		CONSTANT.InstantSendTextMessage,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancellationToCounsellorTextMessage,
			map[string]string{
				"###userName###":  client[0]["first_name"],
				"###user_Name###": counsellor[0]["first_name"],
				"###date###":      UTIL.BuildOnlyDate(appointment[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
		body["appointment_id"],
		CONSTANT.InstantSendTextMessage,
	)
	// fmt.Println("Status:" + CONSTANT.StatusCodeOk)
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
func AppointmentRatingAdd(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if access token is valid, not expired
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

	avg := UTIL.AvgRatingFromula(counsellorRating, counsellorCount[0]["cnt"], "rating")

	switch counsellorType {
	case CONSTANT.CounsellorType:
		DB.ExecuteSQL("update "+CONSTANT.CounsellorsTable+" set total_rating = total_rating + 1, average_rating = ? where counsellor_id = ?", avg, body["counsellor_id"])
		// break
	case CONSTANT.ListenerType:
		DB.ExecuteSQL("update "+CONSTANT.ListenersTable+" set total_rating = total_rating + 1, average_rating = ? where listener_id = ?", avg, body["counsellor_id"])
		// break
	case CONSTANT.TherapistType:
		DB.ExecuteSQL("update "+CONSTANT.TherapistsTable+" set total_rating = total_rating + 1, average_rating = ? where therapist_id = ?", avg, body["counsellor_id"])
		// break
	}

	rate, _ := strconv.Atoi(body["rating"])
	if rate <= 3 {

		// get past completed appointments
		appointts, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		var counsellorName, comment string

		switch counsellorType {
		case CONSTANT.CounsellorType:
			counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", body["counsellor_id"])
			// break
		case CONSTANT.ListenerType:
			counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.ListenersTable+" where listener_id = ?", body["counsellor_id"])
			// break
		case CONSTANT.TherapistType:
			counsellorName = DB.QueryRowSQL("select first_name from "+CONSTANT.TherapistsTable+" where therapist_id = ?", body["counsellor_id"])
			// break
		}
		// send email to client
		filepath_text := "htmlfile/emailmessagebody.html"

		appointDate := UTIL.BuildOnlyDate(appointts[0]["date"])
		appointTime := UTIL.GetTimeFromTimeSlotIN12Hour(appointts[0]["time"])

		if len(body["rating_comment"]) == 0 {
			comment = "Empty comment"
		} else {
			comment = body["rating_comment"]
		}

		emaildata1 := Model.EmailBodyMessageModel{
			Name: "",
			Message: UTIL.ReplaceNotificationContentInString(
				CONSTANT.RatingTitleForInternalReviewBody,
				map[string]string{
					"###therapistname###": counsellorName,
					"###rating###":        body["rating"],
					"###date###":          appointDate,
					"###time###":          appointTime,
					"###content###":       comment,
				},
			),
		}

		emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)

		// email for client
		UTIL.SendEmail(
			UTIL.ReplaceNotificationContentInString(
				CONSTANT.RatingTitleForInternalReviewTitle,
				map[string]string{
					"###therapistName###": counsellorName,
					"###sessionDate###":   UTIL.BuildOnlyDate(appointts[0]["date"]),
				},
			),
			emailBody1,
			CONFIG.OnboardingEmailID,
			CONSTANT.InstantSendEmailMessage,
		)
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
func GenerateAgoraToken(w http.ResponseWriter, r *http.Request, body map[string]string) {

	var response = make(map[string]interface{})

	var encryptedResponse = make(map[string]interface{})

	var roleStr, agora_token, uidStr, channelName string

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

	uidStr = generateRandomID()

	switch body["session"] {
	case "1":
		exists := DB.CheckIfExists(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": body["appointment_id"]})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		channelName = body["appointment_id"]
		switch body["type"] {
		case "1":
			roleStr = CONSTANT.RolePublisher
		case "2":
			roleStr = CONSTANT.RoleSubscriber
		default:
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
	case "2":
		exists := DB.CheckIfExists(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": body["appointment_id"]})
		if !exists {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		channelName = body["appointment_id"]
		if body["type"] == "1" {
			roleStr = CONSTANT.RolePublisher
		} else if body["type"] == "2" {
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

	response["token"] = agora_token
	response["UID"] = uidStr

	encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	if encrypt == "" {
		UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)

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
func AppointmentStart(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	if len(appointment[0]["client_started_at"]) == 0 {
		// update appointment as started
		DB.UpdateSQL(CONSTANT.AppointmentsTable,
			map[string]string{
				"appointment_id": body["appointment_id"],
			},
			map[string]string{
				"status":            CONSTANT.AppointmentStarted,
				"client_started_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

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
func AppointmentEnd(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//check if access token is valid, not expired
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

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": body["appointment_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update appointment as completed
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": body["appointment_id"],
		},
		map[string]string{
			"status":          CONSTANT.AppointmentCompleted,
			"client_ended_at": UTIL.GetCurrentTime().String(),
		},
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
