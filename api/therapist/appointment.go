package therapist

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

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+")", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, photo, age, gender from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
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
	// get schedule for a day
	schedule, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(schedule) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update therapist availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			appointment[0]["time"]: schedule[0][appointment[0]["time"]], // update availability to the latest one
		},
	)

	// update appointment date and time
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentCancelled,
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

	// send appointment cancel notification to client
	therapist, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name"}, map[string]string{"listener_id": appointment[0]["therapist_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"timezone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientContent,
			map[string]string{
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleTimeFormat),
				"###counsellor_name###": therapist[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentStart godoc
// @Tags Therapist Appointment
// @Summary Start an appointment
// @Router /therapist/appointment/start [put]
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

	// update appointment as started
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status": CONSTANT.AppointmentStarted,
		},
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentEnd godoc
// @Tags Therapist Appointment
// @Summary End an appointment
// @Router /therapist/appointment/end [put]
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
	// check if appointment is to be started
	if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentDidntStartedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update appointment as completed
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status": CONSTANT.AppointmentCompleted,
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
		actualAmount, _ := strconv.ParseFloat(invoice[0]["actual_amount"], 64)
		discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
		amountAfterDiscount := actualAmount - discount
		if amountAfterDiscount > 0 { // add only if amount paid
			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

			payoutPercentage, _ := strconv.ParseFloat(DB.QueryRowSQL("select payout_percentage from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"]), 64)

			amountToBePaid := (amountAfterDiscount / slotsBought) * payoutPercentage / 100 // for 1 counselling session

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
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
