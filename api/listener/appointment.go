package listener

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
// @Tags Listener Appointment
// @Summary Get listener upcoming appointments
// @Router /listener/appointment/upcoming [get]
// @Param listener_id query string true "Logged in listener ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsUpcoming(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+")", r.FormValue("listener_id"))
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

// AppointmentsPast godoc
// @Tags Listener Appointment
// @Summary Get listener past appointments
// @Router /listener/appointment/past [get]
// @Param listener_id query string true "Logged in listener ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsPast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("listener_id"), "status": CONSTANT.AppointmentCompleted})
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
// @Tags Listener Appointment
// @Summary Cancel an appointment
// @Router /listener/appointment [delete]
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

	// get listener type
	listenerType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(listenerType, CONSTANT.ListenerType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update listener slots
	// remove previous slot
	date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedules for a weekday
	schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update listener availability
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

	// send appointment cancel notification, email to client
	listener, _, _ := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
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
				"###counsellor_name###": listener[0]["first_name"],
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

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentStart godoc
// @Tags Listener Appointment
// @Summary Start an appointment
// @Router /listener/appointment/start [put]
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

	// get listener type
	listenerType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(listenerType, CONSTANT.ListenerType) {
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
// @Tags Listener Appointment
// @Summary End an appointment
// @Router /listener/appointment/end [put]
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

	// get listener type
	listenerType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(listenerType, CONSTANT.ListenerType) {
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

	// send appointment ended notification and rating to client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentFeedbackHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentFeedbackContent,
			map[string]string{
				"###counsellor_name###": DB.QueryRowSQL("select first_name from "+CONSTANT.ListenersTable+" where listener_id = ?", appointment[0]["counsellor_id"]),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		r.FormValue("appointment_id"),
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
