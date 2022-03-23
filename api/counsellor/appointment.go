package counsellor

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"
	"time"

	UTIL "salbackend/util"
)

// AppointmentsUpcoming godoc
// @Tags Counsellor Appointment
// @Summary Get counsellor upcoming appointments
// @Router /counsellor/appointment/upcoming [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsUpcoming(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("counsellor_id"))
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
// @Tags Counsellor Appointment
// @Summary Get counsellor past appointments
// @Router /counsellor/appointment/past [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsPast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	/*appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status = ? order by date desc", r.FormValue("counsellor_id"), CONSTANT.AppointmentCompleted)
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
// @Tags Counsellor Appointment
// @Summary Cancel an appointment
// @Router /counsellor/appointment [delete]
// @Param appointment_id query string true "Appointment ID to be cancelled"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.CounsellorType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update counsellor slots
	// remove previous slot
	date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedules for a weekday
	schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update counsellor availability
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

	// add penalty for counsellor for cancelling
	// add to counsellor payments
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
	counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientContent,
			map[string]string{
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	UTIL.SendEmail(
		CONSTANT.CounsellorAppointmentCancelClientTitle,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
			},
		),
		client[0]["email"],
		CONSTANT.LaterSendEmailMessage,
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
// @Tags Counsellor Appointment
// @Summary Start an appointment
// @Router /counsellor/appointment/start [put]
// @Param appointment_id query string true "Appointment ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.CounsellorType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
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
// @Tags Counsellor Appointment
// @Summary End an appointment
// @Router /counsellor/appointment/end [put]
// @Param appointment_id query string true "Appointment ID to be ended"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.CounsellorType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
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

	// add to counsellor payments
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
		paidAfterDiscount := paidAmount - discount
		if paidAfterDiscount > 0 { // add only if amount paid
			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

			payoutPercentage, _ := strconv.ParseFloat(DB.QueryRowSQL("select payout_percentage from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", appointment[0]["counsellor_id"]), 64)

			amountToBePaid := (paidAfterDiscount / slotsBought) * payoutPercentage / 100 // for 1 counselling session

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
				"###counsellor_name###": DB.QueryRowSQL("select first_name from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", appointment[0]["counsellor_id"]),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// Appointment Comment godoc
// @Tags Counsellor Appointment
// @Summary Comment an appointment
// @Router /counsellor/appointment/comment [put]
// @Param appointment_id query string true "Appointment ID to be ended"
// @Param body body model.CounsellorCommentRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func CounsellorComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorCommentFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok = DB.UpdateSQL(CONSTANT.AppointmentsTable, map[string]string{"appointment_id": r.FormValue("appointment_id")}, map[string]string{"commentforclient": body["commentforclient"], "attachments": body["attachments"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
