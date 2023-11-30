package listener

import (
	"fmt"
	"math/rand"
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

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("listener_id"))
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

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get past completed appointments
	appointments, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("listener_id"), "status": CONSTANT.AppointmentCompleted})
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
			"status":      CONSTANT.AppointmentCounsellorCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// send appointment cancel notification, email to client
	listener, _, _ := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "email", "phone"}, map[string]string{"listener_id": appointment[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	// remove all previous message for client
	UTIL.RemoveMessage(r.FormValue("appointment_id"), client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(r.FormValue("appointment_id"), listener[0]["phone"])

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientContent,
			map[string]string{
				"###therapist_name###": listener[0]["first_name"],
				"###date_time###":      UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelCounsellorHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelCounsellorContent,
			map[string]string{
				"###date_time###":   UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
				"###client_name###": client[0]["first_name"],
			},
		),
		appointment[0]["counsellor_id"],
		appointment[0]["type"],
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
	)

	filepath_text := "htmlfile/emailmessagebody.html"

	// send email for client
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientBodyEmailBody,
			map[string]string{
				"###therapist_name###": listener[0]["first_name"],
				"###date_time###":      UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
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
		Name:    listener[0]["first_name"],
		Message: CONSTANT.CounsellorAppointmentCancelCounsellorBodyEmailBody,
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.CounsellorAppointmentCancelCounsellorTitle,
		emailBody1,
		listener[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

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
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

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
// @Tags Listener Appointment
// @Summary Start an appointment
// @Router /listener/appointment/start [put]
// @Param appointment_id query string true "Appointment ID to be started"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

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
	// if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentToBeStarted) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyStartedMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get listener type
	listenerType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(listenerType, CONSTANT.ListenerType) {
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

	// var allUsers []string

	// var user = agora[0]["uid"]
	// fmt.Println(user)

	// allUsers = append(allUsers, agora[0]["uid1"])
	// allUsers = append(allUsers, agora[0]["uid"])

	if len(agora[0]["sid"]) == 0 {

		sid, err := UTIL.AgoraRecordingCallStart(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["token"], agora[0]["resource_id"])
		if err != nil {
			fmt.Println("cloud recording not started")
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

	if appointment[0]["started_at"] == "" {
		// update appointment as started
		DB.UpdateSQL(CONSTANT.AppointmentsTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"status":     CONSTANT.AppointmentStarted,
				"started_at": UTIL.GetCurrentTime().Local().String(),
			},
		)
	}

	// send appointment join the call notification to Client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentHasBeenStartedHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentHasBeenStartedContent,
			map[string]string{
				"###clientname###":    DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
				"###therapistname###": DB.QueryRowSQL("select first_name from "+CONSTANT.ListenersTable+" where listener_id = ?", appointment[0]["counsellor_id"]),
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

// AppointmentEnd godoc
// @Tags Listener Appointment
// @Summary End an appointment
// @Router /listener/appointment/end [put]
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
	// check if appointment is to be started
	// if !strings.EqualFold(appointment[0]["status"], CONSTANT.AppointmentStarted) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentDidntStartedMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get listener type
	listenerType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(listenerType, CONSTANT.ListenerType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
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
			"status":   CONSTANT.AppointmentCompleted,
			"ended_at": UTIL.GetCurrentTime().Local().String(),
		},
	)

	if len(agora[0]["fileNameInMp4"]) == 0 && len(agora[0]["fileNameInM3U8"]) == 0 {
		fileNameInMP4, fileNameInM3U8, err := UTIL.AgoraRecordingCallStop(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["resource_id"], agora[0]["sid"])
		if err != nil {
			fmt.Println("file is not created")
		}

		if fileNameInMP4 != "" {
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

	}

	// // send appointment ended notification and rating to client
	// UTIL.SendNotification(
	// 	CONSTANT.ClientAppointmentFeedbackHeading,
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentFeedbackContent,
	// 		map[string]string{
	// 			"###counsellor_name###": DB.QueryRowSQL("select first_name from "+CONSTANT.ListenersTable+" where listener_id = ?", appointment[0]["counsellor_id"]),
	// 		},
	// 	),
	// 	appointment[0]["client_id"],
	// 	CONSTANT.ClientType,
	// 	UTIL.GetCurrentTime().String(),
	// 	CONSTANT.NotificationSent,
	// 	r.FormValue("appointment_id"),
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
