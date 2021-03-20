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

	// get counsellor/listener details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
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

	// get counsellor/listener details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

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

	// get counsellor/listener details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, photo, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
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

	// get appointment order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrdersTable, []string{"*"}, map[string]string{"order_id": appointment[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentReschedule godoc
// @Tags Client Appointment
// @Summary Reschedule an appointment
// @Router /client/appointment [put]
// @Param appointment_id query string true "Appointment ID to be rescheduled"
// @Param body body model.AppointmentRescheduleRequest true "Request Body"
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

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(appointment[0]["counsellor_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.RescheduleSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update counsellor slots
	// remove previous slot
	date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedule for a day
	schedule, status, ok := DB.SelectSQL(CONSTANT.SchedulesTable, []string{appointment[0]["time"]}, map[string]string{"counsellor_id": appointment[0]["counsellor_id"], "weekday": strconv.Itoa(int(date.Weekday()))})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(schedule) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update counsellor availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			appointment[0]["time"]: schedule[0][appointment[0]["time"]], // update availability to the latest one
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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
