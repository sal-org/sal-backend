package client

import (
	"net/http"

	// CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"time"

	UTIL "salbackend/util"
	"strings"
	// "github.com/swaggo/swag"
)

func CorporateCounsellorOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	// get counsellor details
	var counsellorType string

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

	// this is for limit number of session client
	// domainName := strings.Split(client[0]["email"], "@")

	// if domainName[1] == "clovemind.com" {

	// 	appointmentUnLimit, status, ok := DB.SelectProcess("select * from "+CONSTANT.ClientCounsellingUnLimitTable+" where client_id = ? and status = 1 ", client[0]["client_id"])
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}

	// 	if len(appointmentUnLimit) == 0 {

	// 		appointmentLimit, status, ok := DB.SelectProcess("select * from "+CONSTANT.ClientCounsellingLimitTable+" where client_id = ? and status = 3 order by date desc", client[0]["client_id"])
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		if len(appointmentLimit) != 0 {

	// 			limit := len(appointmentLimit)

	// 			if limit%3 == 0 {

	// 				lastAppointmentDate, _ := time.Parse("2006-01-02", appointmentLimit[0]["date"])

	// 				startYear, startMonth, startDay := lastAppointmentDate.Date()
	// 				endYear, endMonth, endDay := time.Now().Date()

	// 				day := 0

	// 				if int(endDay-startDay) < 0 {
	// 					day = 1
	// 				}

	// 				// Calculate total months
	// 				totalMonths := (endYear-startYear)*12 + int(endMonth-startMonth) - day
	// 				if totalMonths <= 3 {
	// 					UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientAppointmentBookLimitOver, CONSTANT.ShowDialog, response)
	// 					return
	// 				}

	// 			}
	// 		}
	// 	}

	// 	// startSlot := (UTIL.GetCurrentTime().Add(330 * time.Minute).Hour()) * 2 // use next slot for removing expired time for today

	// 	// if (UTIL.GetCurrentTime().Add(330 * time.Minute).Minute()) >= 30 {
	// 	// 	startSlot = startSlot + 1
	// 	// }

	// 	upcomingAppointment, status, ok := DB.SelectProcess("select * from "+CONSTANT.ClientCounsellingLimitTable+" where client_id = ? and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and status = 1 order by date desc", client[0]["client_id"])
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}

	// 	if len(upcomingAppointment) != 0 {
	// 		timeInInt, _ := strconv.Atoi(upcomingAppointment[0]["time"])
	// 		timeInInt = timeInInt + 2

	// 		timeInString := strconv.Itoa(timeInInt)
	// 		appointmentTime := UTIL.BuildDateTime(upcomingAppointment[0]["date"], timeInString)

	// 		isUpcomingAppointment := appointmentTime.After(UTIL.GetCurrentTime().Add(330 * time.Minute))

	// 		if isUpcomingAppointment {
	// 			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientUpcomingAppointmentAlreadyExits, CONSTANT.ShowDialog, response)
	// 			return
	// 		}
	// 	}

	// }

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": body["listener_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if listener is valid
	if len(counsellor) == 0 {
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["listener_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 && !strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}
		counsellorType = CONSTANT.TherapistType
	} else {
		// check if listener is active
		if !strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotActiveMessage, CONSTANT.ShowDialog, response)
			return
		}
		counsellorType = CONSTANT.CounsellorType
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
	order["type"] = counsellorType
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

func CorporateCounsellorOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
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
	if strings.EqualFold(order[0]["type"], CONSTANT.ListenerType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": order[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(client[0]["email"], "@")

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

	// if domainName[1] == "clovemind.com" {

	// 	limitAppointment := map[string]string{}
	// 	limitAppointment["client_id"] = order[0]["client_id"]
	// 	limitAppointment["counsellor_id"] = order[0]["counsellor_id"]
	// 	limitAppointment["appointment_id"] = appointmentID
	// 	limitAppointment["date"] = order[0]["date"]
	// 	limitAppointment["time"] = order[0]["time"]
	// 	limitAppointment["status"] = CONSTANT.AppointmentToBeStarted
	// 	limitAppointment["created_at"] = UTIL.GetCurrentTime().String()

	// 	_, status, ok := DB.InsertWithUniqueID(CONSTANT.ClientCounsellingLimitTable, CONSTANT.LimitAppointmentDigits, limitAppointment, "order_id")
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// }

	var counsellor []map[string]string

	// sent notitifications
	switch order[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "phone", "email", "timezone"}, map[string]string{"counsellor_id": order[0]["counsellor_id"]})

	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "phone", "email", "timezone"}, map[string]string{"therapist_id": order[0]["counsellor_id"]})

	}

	counsellor_fullname := counsellor[0]["first_name"] + " " + counsellor[0]["last_name"]

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

	// change order status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()

	// client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"client_id": order[0]["client_id"]})

	// send email to client
	filepath_text := "htmlfile/appointmentConfirmation.html"

	_, status, ok = DB.InsertWithUniqueID(CONSTANT.QualityCheckDetailsTable, CONSTANT.AppointmentDigits, qualitycheck_details, "qualitycheck_details_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	DB.UpdateSQL(CONSTANT.OrderClientAppointmentTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": order[0]["counsellor_id"],
			"date":          order[0]["date"],
		},
		map[string]string{
			order[0]["time"]: CONSTANT.SlotBooked,
		},
	)

	// client notification

	// Booking confirmation
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleClientHeading, CONSTANT.ClientAppointmentScheduleClientContent, order[0]["client_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().Add(330*time.Minute).String(), CONSTANT.NotificationSent, appointmentID,
	)

	// 15 min push notification before appointment start
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

	// Counsellor Notification

	// send appointment booking notification, message to listener
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentScheduleCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentScheduleCounsellorContent,
			map[string]string{
				"###Date###": UTIL.BuildOnlyDate(order[0]["date"]),
				"###Time###": UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["counsellor_id"],
		order[0]["type"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	// send appointment reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentRemiderClientContent,
			map[string]string{
				"###user_name###": counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["counsellor_id"],
		order[0]["type"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		appointmentID,
	)

	// Client SMS

	// Client Booking Confirmation
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentConfirmationTextMessage,
			map[string]string{
				"###userName###":  client[0]["first_name"],
				"###user_Name###": counsellor[0]["first_name"],
				"###date###":      UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		appointmentID,
		CONSTANT.InstantSendTextMessage,
	)

	// 30 min reminder sms notification
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###userName###":  counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-5*time.Hour).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// 15 min reminder sms notification
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###userName###":  counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// Listener SMS

	// confirmation for therapists message
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentConfirmationTextMessage,
			map[string]string{
				"###userName###":  counsellor[0]["first_name"],
				"###user_Name###": client[0]["first_name"],
				"###date###":      UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		appointmentID,
		CONSTANT.InstantSendTextMessage,
	)

	// send at 15 min before of appointment
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": counsellor[0]["first_name"],
				"###userName###":  client[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// Client Email

	// domainName := strings.Split(client[0]["email"], "@")

	accessCode := ""

	if domainName[1] == "db.com" {
		accessCode = "2332"
	} else {
		accessCode = "1234"
	}

	// Payment receipt
	emaildata := Model.EmailBodyWithAccessCodeMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBookClientEmailBody,
			map[string]string{
				"###therpist_name###": counsellor[0]["first_name"],
				"###date###":          UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":          UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		AccessCode: accessCode,
	}

	emailBody := UTIL.GetHTMLTemplateForClientConfirmationWithAccessCodeText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookClientTitle,
		emailBody,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// Counsellor Email

	// Payment receipt
	emaildata1 := Model.EmailBodyMessageModel{
		Name: counsellor[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentBookCounsellorEmailBody,
			map[string]string{
				"###client_name###": client[0]["first_name"],
				"###date###":        UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":        UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, "htmlfile/emailmessagebody.html")
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookCounsellorTitle,
		emailBody1,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}

func InPersonCorporateCounsellorOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	// get counsellor details
	var counsellorType string

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

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": body["listener_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if listener is valid
	if len(counsellor) == 0 {
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["listener_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(counsellor) > 0 && !strings.EqualFold(counsellor[0]["status"], CONSTANT.TherapistActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
			return
		}
		counsellorType = CONSTANT.TherapistType
	} else {
		// check if listener is active
		if !strings.EqualFold(counsellor[0]["status"], CONSTANT.ListenerActive) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotActiveMessage, CONSTANT.ShowDialog, response)
			return
		}
		counsellorType = CONSTANT.CounsellorType
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentzInPersonSlotAvailable(body["listener_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check 2nd appoimtent with the same listener
	appointment2nd, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", body["client_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(appointment2nd) != 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.InPersonAppointmentAlreadyBooked, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["listener_id"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["type"] = counsellorType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()
	// no paid amount and billing

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.InPersonOrderClientAppointmentTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func InPersonCorporateCounsellorOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

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
	order, status, ok := DB.SelectSQL(CONSTANT.InPersonOrderClientAppointmentTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
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
	// if !strings.EqualFold(order[0]["type"], CONSTANT.ListenerType) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	address, status, ok := DB.SelectSQL(CONSTANT.InPersonSLotsScheduleTable, []string{"*"}, map[string]string{"date": order[0]["date"], "counsellor_id": order[0]["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
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
	appointment["company_name"] = address[0]["partner_name"]
	appointment["company_location"] = address[0]["partner_location"]
	appointment["counselling_room"] = address[0]["counselling_room"]
	appointment["counselling_address"] = address[0]["counselling_address"]
	appointment["status"] = CONSTANT.AppointmentToBeStarted
	appointment["created_at"] = UTIL.GetCurrentTime().String()
	appointmentID, status, ok := DB.InsertWithUniqueID(CONSTANT.InPersonAppointmentsTable, CONSTANT.AppointmentDigits, appointment, "appointment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var counsellor []map[string]string

	// sent notitifications
	switch order[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "phone", "email", "timezone"}, map[string]string{"counsellor_id": order[0]["counsellor_id"]})

	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "phone", "email", "timezone"}, map[string]string{"therapist_id": order[0]["counsellor_id"]})

	}

	// change order status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()

	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"client_id": order[0]["client_id"]})

	// send email to client
	filepath_text := "htmlfile/emailmessagebody.html"

	DB.UpdateSQL(CONSTANT.InPersonOrderClientAppointmentTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	DB.UpdateSQL(CONSTANT.InPersonSLotsTable,
		map[string]string{
			"counsellor_id":    order[0]["counsellor_id"],
			"date":             order[0]["date"],
			"company_name":     address[0]["partner_name"],
			"company_location": address[0]["partner_location"],
		},
		map[string]string{
			order[0]["time"]: CONSTANT.SlotBooked,
		},
	)

	// client notification

	// Booking confirmation
	UTIL.SendNotification(
		CONSTANT.ClientInPersonAppointmentConfirmationClientHeading, UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentConfirmationClientContent,
			map[string]string{
				"###firstName###": counsellor[0]["first_name"],
				"###date###":      UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###location###":  address[0]["counselling_room"] + ", " + address[0]["counselling_address"],
			},
		), order[0]["client_id"], CONSTANT.ClientType, UTIL.GetCurrentTime().Add(330*time.Minute).String(), CONSTANT.NotificationSent, appointmentID,
	)

	// 15 min push notification before appointment start
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentRemainderClientContent,
			map[string]string{
				"###thearpistName###": counsellor[0]["first_name"],
				"###location###":      address[0]["counselling_room"] + ", " + address[0]["counselling_address"],
			},
		),
		order[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		appointmentID,
	)

	// Counsellor Notification

	// send appointment booking notification, message to listener
	UTIL.SendNotification(
		CONSTANT.ClientInPersonAppointmentConfirmationTherapistHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentConfirmationTherapistContent,
			map[string]string{
				"###firstName###": client[0]["first_name"],
				"###location###":  address[0]["counselling_room"] + ", " + address[0]["counselling_address"],
				"###date###":      UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		order[0]["counsellor_id"],
		order[0]["type"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		CONSTANT.NotificationSent,
		appointmentID,
	)

	// send appointment reminder notification to counsellor before 15 min
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentReminderCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentRemainderTherapistContent,
			map[string]string{
				"###time###":       UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###clientName###": client[0]["first_name"],
			},
		),
		order[0]["counsellor_id"],
		order[0]["type"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
		CONSTANT.NotificationInProgress,
		appointmentID,
	)

	// Client SMS

	// Client Booking Confirmation
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentConfirmationClientTextMeassage,
			map[string]string{
				"###therapistName###": counsellor[0]["first_name"],
				"###date###":          UTIL.BuildOnlyDate(order[0]["date"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		appointmentID,
		CONSTANT.InstantSendTextMessage,
	)

	// 1 hr reminder sms notification
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			// need to change
			CONSTANT.ClientAppointmentReminderTextMessage,
			map[string]string{
				"###user_name###": client[0]["first_name"],
				"###userName###":  counsellor[0]["first_name"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-1*time.Hour).UTC().String(),
		appointmentID,
		CONSTANT.LaterSendTextMessage,
	)

	// // 15 min reminder sms notification
	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		// need to change
	// 		CONSTANT.ClientAppointmentReminderTextMessage,
	// 		map[string]string{
	// 			"###user_name###": client[0]["first_name"],
	// 			"###userName###":  counsellor[0]["first_name"],
	// 			"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	client[0]["phone"],
	// 	UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
	// 	appointmentID,
	// 	CONSTANT.LaterSendTextMessage,
	// )

	// Listener SMS

	// confirmation for therapists message
	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentConfirmationTextMessage,
			map[string]string{
				"###userName###":  counsellor[0]["first_name"],
				"###user_Name###": client[0]["first_name"],
				"###date###":      UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		counsellor[0]["phone"],
		UTIL.GetCurrentTime().Add(330*time.Minute).String(),
		appointmentID,
		CONSTANT.InstantSendTextMessage,
	)

	// send at 15 min before of appointment
	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		// need to change
	// 		CONSTANT.ClientAppointmentReminderTextMessage,
	// 		map[string]string{
	// 			"###user_name###": counsellor[0]["first_name"],
	// 			"###userName###":  client[0]["first_name"],
	// 			"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	counsellor[0]["phone"],
	// 	UTIL.BuildDateTime(order[0]["date"], order[0]["time"]).Add(-15*time.Minute).UTC().String(),
	// 	appointmentID,
	// 	CONSTANT.LaterSendTextMessage,
	// )

	// Client Email

	// Payment receipt
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentBookClientBody,
			map[string]string{
				"###therapistName###": counsellor[0]["first_name"],
				"###date###":          UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":          UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###location###":      address[0]["counselling_room"] + ", " + address[0]["counselling_address"],
			},
		),
	}

	emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientInPersonAppointmentBookClientTitle,
		emailBody,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	// Counsellor Email

	// Payment receipt
	emaildata1 := Model.EmailBodyMessageModel{
		Name: counsellor[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientInPersonAppointmentBookTherapistBody,
			map[string]string{
				"###clientName###": client[0]["first_name"],
				"###date###":       UTIL.BuildOnlyDate(order[0]["date"]),
				"###time###":       UTIL.GetTimeFromTimeSlotIN12Hour(order[0]["time"]),
				"###location###":   address[0]["counselling_room"] + ", " + address[0]["counselling_address"],
			},
		),
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientInPersonAppointmentBookTherapistTitle,
		emailBody1,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
