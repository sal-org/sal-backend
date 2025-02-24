package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strings"
	"time"

	UTIL "salbackend/util"
)

// TherapistProfile godoc
// @Tags Client Therapist
// @Summary Get therapist details
// @Router /client/therapist [get]
// @Param therapist_id query string true "Therapist ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func TherapistProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	lang := []string{}
	expert := []string{}
	counsellor := map[string]interface{}{}

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

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "pronoun", "total_rating", "average_rating", "photo", "video", "education", "experience", "therapeutic_approach", "about", "slot_type"}, map[string]string{"therapist_id": body["therapist_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapist languages
	therapistLang, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", body["therapist_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get therapist topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", body["therapist_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for i := 0; i < len(therapistLang); i++ {
		value := therapistLang[i]["language"]
		lang = append(lang, value)

	}

	for i := 0; i < len(topics); i++ {
		value := topics[i]["topic"]
		expert = append(expert, value)
	}

	counsellor["name"] = therapist[0]["first_name"] + " " + therapist[0]["last_name"]
	counsellor["pronoun"] = therapist[0]["pronoun"]
	counsellor["total_rate"] = therapist[0]["total_rating"]
	counsellor["average_rate"] = therapist[0]["average_rating"]
	counsellor["photo"] = therapist[0]["photo"]
	counsellor["video"] = therapist[0]["video"]
	counsellor["education"] = therapist[0]["education"]
	counsellor["experience"] = therapist[0]["experience"]
	counsellor["therapeutic_approach"] = therapist[0]["therapeutic_approach"]
	counsellor["about"] = therapist[0]["about"]
	counsellor["slot_type"] = therapist[0]["slot_type"]
	counsellor["languages"] = lang
	counsellor["topics"] = expert

	response["therapist"] = counsellor
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// TherapistSlots godoc
// @Tags Client Therapist
// @Summary Get therapist slots
// @Router /client/therapist/slots [get]
// @Param therapist_id query string true "Therapist ID to get slot details"
// @Security JWTAuth
// @Produce json
// @Success 200
func TherapistSlots(w http.ResponseWriter, r *http.Request) {
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

	// get therapist slots
	slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and available = '1' and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", body["therapist_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// remove times and dates with no availability
	response["slots"] = UTIL.FilterAvailableSlots(slots)
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TherapistOrderCreateForCorporateRequiredFields)
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

	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["therapist_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) > 0 && !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
		return
	}
	counsellorType = CONSTANT.TherapistType

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(body["therapist_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["therapist_id"]
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

	// client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "phone", "email", "timezone"}, map[string]string{"client_id": order[0]["client_id"]})

	// send email to client
	filepath_text := "htmlfile/appointmentConfirmation.html"

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
