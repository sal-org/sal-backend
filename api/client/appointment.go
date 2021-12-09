package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

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

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+")", r.FormValue("client_id"))
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
	invoice, status, ok := DB.SelectProcess("select order_id as id, payment_id, user_type  from " + CONSTANT.InvoicesTable + " where order_id in ('" + strings.Join(orderIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["order_detils"] = UTIL.ConvertMapToKeyMap(invoice, "id")
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

	// get appointment slots details
	appointmentSlots, status, ok := DB.SelectSQL(CONSTANT.AppointmentSlotsTable, []string{"*"}, map[string]string{"order_id": appointment[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get appointment order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"*"}, map[string]string{"order_id": appointment[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["appointment"] = appointment[0]
	response["appointment_slots"] = appointmentSlots[0]
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
	_, status, ok = DB.InsertWithUniqueID(CONSTANT.AppointmentsTable, CONSTANT.AppointmentDigits, appointment, "appointment_id")
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
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "timezone"}, map[string]string{"counsellor_id": appointmentSlot[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "phone", "timezone"}, map[string]string{"listener_id": appointmentSlot[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "phone", "timezone"}, map[string]string{"therapist_id": appointmentSlot[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone"}, map[string]string{"client_id": appointmentSlot[0]["client_id"]})

	// send appointment booking notification to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentScheduleClientContent,
			map[string]string{
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(body["date"], body["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		appointmentSlot[0]["client_id"],
		CONSTANT.ClientType,
	)

	// send appointment booking notification, message to counsellor
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentScheduleCounsellorContent,
			map[string]string{
				"###date_time###":   UTIL.ConvertTimezone(UTIL.BuildDateTime(body["date"], body["time"]), counsellor[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###client_name###": client[0]["first_name"],
			},
		),
		appointmentSlot[0]["counsellor_id"],
		counsellorType,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentScheduleCounsellorTextMessage,
			map[string]string{
				"###counsellor_name###": counsellor[0]["first_name"],
				"###client_name###":     client[0]["first_name"],
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(body["date"], body["time"]), counsellor[0]["timezone"]).Format(CONSTANT.ReadbleDateFormat),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		CONSTANT.LaterSendTextMessage,
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
	if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours() < 4 {
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
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"timezone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// send appointment reschedule notification to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentRescheduleClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleClientContent,
			map[string]string{
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(body["date"], body["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###counsellor_name###": counsellor[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
	)

	// send appointment reschedule notification to counsellor
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentRescheduleCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRescheduleCounsellorContent,
			map[string]string{},
		),
		appointment[0]["counsellor_id"],
		counsellorType,
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
	// // appointment can cancelled only after min reschedules
	// reschedules, _ := strconv.Atoi(appointment[0]["times_rescheduled"])
	// if reschedules < CONSTANT.MaximumAppointmentReschedule {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentCantCancelMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

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
			"status":      CONSTANT.AppointmentCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// refund amount
	// check if appointment is alteast after 4 hours
	// if below 4 hours, charges are 100% that means dont refund anything
	if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours() >= 4 {
		charges := CONSTANT.ClientAppointmentCancellationCharges
		// get invoice details
		invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount", "payment_id", "invoice_id"}, map[string]string{"order_id": appointment[0]["order_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(invoice) > 0 {
			// invoice is available => amount is paid, order is not free
			// get order details
			order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
			discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
			amountAfterDiscount := paidAmount - discount
			if amountAfterDiscount > 0 { // refund only if amount paid
				paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
				refundedAmount, _ := strconv.ParseFloat(DB.QueryRowSQL("select sum(refunded_amount) from "+CONSTANT.RefundsTable+" where invoice_id = '"+invoice[0]["invoice_id"]+"'"), 64)
				slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)
				cancellationCharges := (amountAfterDiscount / slotsBought) * charges
				refundAmount := (paidAmount / slotsBought) - cancellationCharges // remove from paid amount only, not from amount after discount, as discussed
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
	}

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "timezone"}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "timezone"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "timezone"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
		break

	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// send appointment cancel notification, email to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelClientContent,
			map[string]string{
				"###date_time###":       UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###counsellor_name###": counsellor[0]["first_name"],
				"###client_name###":     client[0]["first_name"],
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
	)

	UTIL.SendEmail(
		CONSTANT.ClientAppointmentCancelClientTitle,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelClientBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
			},
		),
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// send appointment cancel notification to counsellor
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentCancelCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancelCounsellorContent,
			map[string]string{
				"###date_time###":   UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]), counsellor[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###client_name###": client[0]["first_name"],
			},
		),
		appointment[0]["counsellor_id"],
		counsellorType,
	)

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
	})

	// send notifications
	// get counsellor name
	var counsellor []map[string]string
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentSlotsTable+" where appointment_slot_id = ?)", r.FormValue("appointment_slot_id"))
	switch counsellorType {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name"}, map[string]string{"counsellor_id": appointmentSlots[0]["counsellor_id"]})
		break
	case CONSTANT.ListenerType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name"}, map[string]string{"listener_id": appointmentSlots[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name"}, map[string]string{"therapist_id": appointmentSlots[0]["counsellor_id"]})
		break
	}
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name"}, map[string]string{"client_id": appointmentSlots[0]["client_id"]})

	// send appointment cancel notification to client
	UTIL.SendNotification(
		CONSTANT.ClientBulkAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientBulkAppointmentCancelClientContent,
			map[string]string{
				"###counsellor_name###": counsellor[0]["first_name"],
				"###client_name###":     client[0]["first_name"],
			},
		),
		appointmentSlots[0]["client_id"],
		CONSTANT.ClientType,
	)

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

	// get counsellor type and update their ratings
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", body["appointment_id"])

	switch counsellorType {
	case CONSTANT.CounsellorType:
		DB.ExecuteSQL("update "+CONSTANT.CounsellorsTable+" set total_rating = total_rating + 1, average_rating = (average_rating + ?)/2 where counsellor_id = ?", body["rating"], body["counsellor_id"])
		break
	case CONSTANT.ListenerType:
		DB.ExecuteSQL("update "+CONSTANT.ListenersTable+" set total_rating = total_rating + 1, average_rating = (average_rating + ?)/2 where listener_id = ?", body["rating"], body["counsellor_id"])
		break
	case CONSTANT.TherapistType:
		DB.ExecuteSQL("update "+CONSTANT.TherapistsTable+" set total_rating = total_rating + 1, average_rating = (average_rating + ?)/2 where therapist_id = ?", body["rating"], body["counsellor_id"])
		break
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
