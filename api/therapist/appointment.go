package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strconv"
	"strings"
	"time"

	UTIL "salbackend/util"
)

// AppointmentsUpcoming godoc
// @Tags Therapist Appointment
// @Summary Get therapist upcoming appointments
// @Router /therapist/appointment/upcoming [get]
// @Param therapist_id query string true "Logged in therapist ID"
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

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, photo, date_of_birth, gender from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentsPast godoc
// @Tags Therapist Appointment
// @Summary Get therapist past appointments
// @Router /therapist/appointment/past [get]
// @Param therapist_id query string true "Logged in therapist ID"
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
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	/*appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status = ? order by date desc", r.FormValue("therapist_id"), CONSTANT.AppointmentCompleted)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}*/
	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentCancel godoc
// @Tags Therapist Appointment
// @Summary Cancel an appointment
// @Router /therapist/appointment [delete]
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
	if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update therapist slots
	// remove previous slot
	date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedules for a weekday
	schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update therapist availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			appointment[0]["time"]: UTIL.CheckIfScheduleAvailable(schedules, appointment[0]["time"]), // update availability to the latest one
		},
	)

	// update appointment date and time
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentCounsellorCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// add a slot to appointments
	DB.ExecuteSQL("update "+CONSTANT.AppointmentSlotsTable+" set slots_remaining = slots_remaining + 1 where order_id = ?", appointment[0]["order_id"])

	// add penalty for therapist for cancelling
	// add to therapist payments
	// get invoice details
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
		actualAmount, _ := strconv.ParseFloat(invoice[0]["actual_amount"], 64)
		discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
		amountAfterDiscount := actualAmount - discount
		if amountAfterDiscount > 0 { // add only if amount paid
			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

			amountFor1Session := amountAfterDiscount / slotsBought // for 1 counselling session
			cancellationCharges := amountFor1Session * CONSTANT.CounsellorCancellationCharges

			DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
				"counsellor_id": appointment[0]["counsellor_id"],
				"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
				"description":   "Cancellation",
				"amount":        strconv.FormatFloat(-cancellationCharges, 'f', 2, 64),
				"status":        CONSTANT.PaymentActive,
				"created_at":    UTIL.GetCurrentTime().String(),
			}, "payment_id")
		}
	}

	// send appointment cancel notification, email to client
	therapist, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "email"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientContent,
			map[string]string{
				"###therapist_name###": therapist[0]["first_name"],
				"###date_time###":      UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelCounsellorContent,
			map[string]string{
				"###date_time###":   UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###client_name###": client[0]["first_name"],
			},
		),
		appointment[0]["counsellor_id"],
		appointment[0]["type"],
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	filepath_text := "htmlfile/emailmessagebody.html"

	// send email for client
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientBodyEmailBody,
			map[string]string{
				"###therapist_name###": therapist[0]["first_name"],
				"###date_time###":      UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.CounsellorAppointmentCancelClientTitle,
		emailBody,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	emaildata1 := Model.EmailBodyMessageModel{
		Name:    therapist[0]["first_name"],
		Message: CONSTANT.CounsellorAppointmentCancelCounsellorBodyEmailBody,
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.CounsellorAppointmentCancelCounsellorTitle,
		emailBody1,
		therapist[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancellationToClientTextMessage,
			map[string]string{
				"###client_name###": client[0]["first_name"],
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		CONSTANT.LaterSendTextMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentStart godoc
// @Tags Therapist Appointment
// @Summary Start an appointment
// @Router /therapist/appointment/start [put]
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

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
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

	if agora[0]["uid"] == r.FormValue("uid") {

		sid, err := UTIL.AgoraRecordingCallStart(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["token"], agora[0]["resource_id"])
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.AgoraCallMessage, CONSTANT.ShowDialog, response)
			return
		}

		DB.UpdateSQL(CONSTANT.AgoraTable,
			map[string]string{
				"agora_id": agora[0]["agora_id"],
			},
			map[string]string{
				"sid":         sid,
				"status":      CONSTANT.AgoraCallStart1,
				"modified_at": UTIL.GetCurrentTime().String(),
			},
		)

	}

	// update appointment as started
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":     CONSTANT.AppointmentStarted,
			"started_at": UTIL.GetCurrentTime().String(),
		},
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentEnd godoc
// @Tags Therapist Appointment
// @Summary End an appointment
// @Router /therapist/appointment/end [put]
// @Param appointment_id query string true "Appointment ID to be ended"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentEnd(w http.ResponseWriter, r *http.Request) {
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

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	agora, status, ok := DB.SelectSQL(CONSTANT.AgoraTable, []string{"*"}, map[string]string{"appointment_id": appointment[0]["appointment_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if agora[0]["uid"] == r.FormValue("uid") {
		fileNameInMP4, fileNameInM3U8, err := UTIL.AgoraRecordingCallStop(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["resource_id"], agora[0]["sid"])
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.AgoraCallMessage, CONSTANT.ShowDialog, response)
			return
		}
		DB.UpdateSQL(CONSTANT.AgoraTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"fileNameInMp4":  fileNameInMP4,
				"fileNameInM3U8": fileNameInM3U8,
				"status":         CONSTANT.AgoraCallStop1,
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

	// update appointment as completed
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":   CONSTANT.AppointmentCompleted,
			"ended_at": UTIL.GetCurrentTime().String(),
		},
	)

	// add to therapist payments
	// get invoice details
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
		amountAfterDiscount := paidAmount + discount
		if amountAfterDiscount > 0 { // add only if amount paid
			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

			// these come from database
			// payoutPercentage, _ := strconv.ParseFloat(DB.QueryRowSQL("select payout_percentage from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"]), 64)

			amountToBePaid := (amountAfterDiscount / slotsBought) * CONSTANT.CounsellorPayoutPercentage / 100 // for 1 counselling session

			DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
				"counsellor_id": appointment[0]["counsellor_id"],
				"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
				"description":   "Consultation",
				"amount":        strconv.FormatFloat(amountToBePaid, 'f', 2, 64),
				"status":        CONSTANT.PaymentActive,
				"created_at":    UTIL.GetCurrentTime().String(),
			}, "payment_id")
		}
	}

	// send appointment ended notification and rating to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFeedbackHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentFeedbackContent,
			map[string]string{
				"###counsellor_name###": DB.QueryRowSQL("select first_name from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"]),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
