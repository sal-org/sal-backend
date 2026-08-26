package client

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

// Home godoc
// @Tags Client Home
// @Summary Get home page content
// @Router /client/home [get]
// @Param client_id query string false "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)
	var recommended []map[string]string
	var ok bool
	var status string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	accessCode := ""

	clientID, ok := UTIL.Required(r.FormValue("client_id"), "Client ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	platform, ok := UTIL.Required(r.FormValue("platform"), "Platform")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, platform, CONSTANT.ShowDialog, response)
		return
	}

	version, ok := UTIL.Required(r.FormValue("version"), "Version")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, version, CONSTANT.ShowDialog, response)
		return
	}

	timezone, ok := UTIL.Required(r.FormValue("timezone"), "Timezone")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, timezone, CONSTANT.ShowDialog, response)
		return
	}

	active := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"status": "1", "client_id": r.FormValue("client_id")})
	if !active {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.InValidIDError, CONSTANT.ShowDialog, response)
		return
	}

	client, status, ok := DB.SelectProcess("select topic_ids, email from "+CONSTANT.ClientsTable+" where client_id = ? and status = 1", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest content for recommended
	recommended, status, ok = DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where category_id in (" + client[0]["topic_ids"] + ") and training = 0 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(client[0]["email"], "@")

	switch domainName[1] {
	case "db.com":
		accessCode = "2332"
	case "ageasfederal.com":
		accessCode = "2523"
	default:
		accessCode = "1234"
	}

	updateClient := map[string]string{}

	// update last login time and platform
	if len(r.FormValue("platform")) > 0 {
		updateClient["platform"] = r.FormValue("platform")
	}
	if len(r.FormValue("version")) > 0 {
		updateClient["version"] = r.FormValue("version")
	}
	if len(r.FormValue("timezone")) > 0 {
		updateClient["timezone"] = r.FormValue("timezone")
	}
	updateClient["last_active_time"] = UTIL.GetCurrentTime().String()

	status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("client_id")}, updateClient)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time from "+CONSTANT.AppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(appointments) == 0 {

		response["appointments"] = make(map[string]string)

	} else {

		loc, _ := time.LoadLocation("Asia/Kolkata")

		now := time.Now().In(loc)

		var latestUpcoming map[string]string
		var latestUpcomingTime time.Time
		for _, appointment := range appointments {
			appointmentDate := appointment["date"]

			slot, err := strconv.Atoi(appointment["time"])
			if err != nil {
				continue
			}

			slot = slot + 1 // Increment the slot by 1 to get the end time of the appointment

			// Convert slot number to hour/minute.
			hour := slot / 2
			minute := (slot % 2) * 30

			appointmentTime := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				hour,
				minute,
				0,
				0,
				loc,
			)

			// Use appointment's actual date.
			date, err := time.ParseInLocation("2006-01-02", appointmentDate, loc)
			if err != nil {
				continue
			}

			appointmentTime = time.Date(
				date.Year(),
				date.Month(),
				date.Day(),
				hour,
				minute,
				0,
				0,
				loc,
			)

			// Remove/skip expired appointments.
			if !appointmentTime.After(now) {
				continue
			}

			// Pick the earliest upcoming appointment.
			if latestUpcoming == nil || appointmentTime.Before(latestUpcomingTime) {
				latestUpcoming = appointment
				latestUpcomingTime = appointmentTime
			}
		}

		if latestUpcoming == nil {
			response["appointments"] = make(map[string]string)

		} else {

			// Get counsellor information for the selected appointment.
			counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + latestUpcoming["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + latestUpcoming["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + latestUpcoming["counsellor_id"] + "'))")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
			// counsellors[0]["photo"] = endPointURL

			latestUpcoming["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
			latestUpcoming["counsellor_photo"] = counsellors[0]["photo"]
			latestUpcoming["counsellor_education"] = counsellors[0]["education"]

			response["appointments"] = latestUpcoming
		}

	}

	// if len(appointments) > 0 {

	// 	if appointments[0]["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
	// 		localTime := 0
	// 		loc, _ := time.LoadLocation("Asia/Kolkata")
	// 		now := time.Now().In(loc)
	// 		if now.Minute() >= 30 {
	// 			localTime = now.Hour()*2 + 1
	// 		} else {
	// 			localTime = now.Hour() * 2
	// 		}
	// 		appointmentTime, _ := strconv.Atoi(appointments[0]["time"])

	// 		if appointmentTime+1 < localTime {
	// 			if len(appointments) > 1 {
	// 				counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[1]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[1]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[1]["counsellor_id"] + "'))")
	// 				if !ok {
	// 					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 					return
	// 				}

	// 				// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 				// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 				// counsellors[0]["photo"] = endPointURL

	// 				appointments[1]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 				appointments[1]["counsellor_photo"] = counsellors[0]["photo"]
	// 				appointments[1]["counsellor_education"] = counsellors[0]["education"]

	// 				// virtual appointments
	// 				response["appointments"] = appointments[1]
	// 			} else {
	// 				response["appointments"] = make(map[string]string)
	// 			}
	// 		} else {
	// 			counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[0]["counsellor_id"] + "'))")
	// 			if !ok {
	// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 				return
	// 			}

	// 			url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 			counsellors[0]["photo"] = endPointURL

	// 			appointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 			appointments[0]["counsellor_photo"] = counsellors[0]["photo"]
	// 			appointments[0]["counsellor_education"] = counsellors[0]["education"]
	// 			// virtual appointments
	// 			response["appointments"] = appointments[0]
	// 		}
	// 	} else {
	// 		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[0]["counsellor_id"] + "'))")
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 		// counsellors[0]["photo"] = endPointURL

	// 		appointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 		appointments[0]["counsellor_photo"] = counsellors[0]["photo"]
	// 		appointments[0]["counsellor_education"] = counsellors[0]["education"]

	// 		// virtual appointments
	// 		response["appointments"] = appointments[0]
	// 	}
	// } else {
	// 	response["appointments"] = make(map[string]string)
	// }

	// get upcoming appointments both to be started and started
	inpersonAppointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time, company_name, company_location, counselling_address from "+CONSTANT.InPersonAppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(inpersonAppointments) == 0 {

		response["inperson_appointments"] = make(map[string]string)

	} else {

		loc, _ := time.LoadLocation("Asia/Kolkata")

		now := time.Now().In(loc)

		var upcomingInPerson map[string]string
		var upcomingInPersonTime time.Time

		for _, inpersonAppointment := range inpersonAppointments {
			appointmentDate := inpersonAppointment["date"]

			slot, err := strconv.Atoi(inpersonAppointment["time"])
			if err != nil {
				continue
			}

			slot = slot + 1 // Increment the slot by 1 to get the end time of the appointment

			// Convert slot number to hour/minute.
			hour := slot / 2
			minute := (slot % 2) * 30

			appointmentTime := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				hour,
				minute,
				0,
				0,
				loc,
			)

			// Use appointment's actual date.
			date, err := time.ParseInLocation("2006-01-02", appointmentDate, loc)
			if err != nil {
				continue
			}

			appointmentTime = time.Date(
				date.Year(),
				date.Month(),
				date.Day(),
				hour,
				minute,
				0,
				0,
				loc,
			)

			// Remove/skip expired appointments.
			if !appointmentTime.After(now) {
				continue
			}

			// Pick the earliest upcoming appointment.
			if upcomingInPerson == nil || appointmentTime.Before(upcomingInPersonTime) {
				upcomingInPerson = inpersonAppointment
				upcomingInPersonTime = appointmentTime
			}
		}

		if upcomingInPerson == nil {
			response["inperson_appointments"] = make(map[string]string)
		} else {
			// Get counsellor information for the selected appointment.
			counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + upcomingInPerson["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + upcomingInPerson["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + upcomingInPerson["counsellor_id"] + "'))")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
			// counsellors[0]["photo"] = endPointURL

			upcomingInPerson["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
			upcomingInPerson["counsellor_photo"] = counsellors[0]["photo"]
			upcomingInPerson["counsellor_education"] = counsellors[0]["education"]

			response["inperson_appointments"] = upcomingInPerson
		}

	}

	// if len(inpersonAppointments) > 0 {

	// 	if inpersonAppointments[0]["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
	// 		localTime := 0
	// 		timeNow := UTIL.GetCurrentTime()
	// 		timeNow = timeNow.Add(330 * time.Minute)
	// 		if timeNow.Minute() >= 30 {
	// 			localTime = timeNow.Hour()*2 + 1
	// 		} else {
	// 			localTime = timeNow.Hour() * 2
	// 		}
	// 		appointmentTime, _ := strconv.Atoi(inpersonAppointments[0]["time"])

	// 		if appointmentTime+1 < localTime {

	// 			if len(inpersonAppointments) > 1 {
	// 				counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[1]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[1]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[1]["counsellor_id"] + "'))")
	// 				if !ok {
	// 					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 					return
	// 				}

	// 				url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 				_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 				counsellors[0]["photo"] = endPointURL

	// 				inpersonAppointments[1]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 				inpersonAppointments[1]["counsellor_photo"] = counsellors[0]["photo"]
	// 				inpersonAppointments[1]["counsellor_education"] = counsellors[0]["education"]

	// 				// virtual appointments
	// 				response["inperson_appointments"] = inpersonAppointments[1]
	// 			} else {
	// 				response["inperson_appointments"] = make(map[string]string)
	// 			}
	// 		} else {
	// 			counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[0]["counsellor_id"] + "'))")
	// 			if !ok {
	// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 				return
	// 			}

	// 			// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 			// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 			// counsellors[0]["photo"] = endPointURL

	// 			inpersonAppointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 			inpersonAppointments[0]["counsellor_photo"] = counsellors[0]["photo"]
	// 			inpersonAppointments[0]["counsellor_education"] = counsellors[0]["education"]

	// 			// virtual appointments
	// 			response["inperson_appointments"] = inpersonAppointments[0]
	// 		}
	// 	} else {
	// 		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[0]["counsellor_id"] + "'))")
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 		// counsellors[0]["photo"] = endPointURL

	// 		inpersonAppointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
	// 		inpersonAppointments[0]["counsellor_photo"] = counsellors[0]["photo"]
	// 		inpersonAppointments[0]["counsellor_education"] = counsellors[0]["education"]

	// 		// virtual appointments
	// 		response["inperson_appointments"] = inpersonAppointments[0]
	// 	}
	// } else {
	// 	response["inperson_appointments"] = make(map[string]string)
	// }

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc, time asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(events) > 0 {
		// get counsellor details
		// get upcoming booked events
		// urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, events[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// _, endPointURLPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		// events[0]["photo"] = endPointURLPhoto

		// urlBackGroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, events[0]["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// _, endPointURLBackGroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackGroundPhoto)
		// events[0]["background_photo"] = endPointURLBackGroundPhoto
		response["upcoming_events"] = events[0]
	} else {
		response["upcoming_events"] = make(map[string]string)
	}

	webinarOrder, status, ok := DB.SelectProcess("select * from "+CONSTANT.WebinarsBookTable+" where status = 1 and client_id = ? ", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(webinarOrder) != 0 {

		// get webinar details
		webinar, status, ok := DB.SelectProcess("select * from "+CONSTANT.WebinarsTable+" where webinar_id = ? and status = 1", webinarOrder[0]["webinar_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(webinar) > 0 {
			webinarOrder[0]["title"] = webinar[0]["title"]
			webinarOrder[0]["time"] = webinar[0]["time"]
			webinarOrder[0]["date"] = webinar[0]["date"]
			webinarOrder[0]["mode"] = webinar[0]["mode"]

			response["webinar_order"] = webinarOrder[0]
		}
	} else {
		response["webinar_order"] = make(map[string]string)
	}

	// get webinar order details

	// get random quote
	quote, status, ok := DB.SelectProcess("select * from " + CONSTANT.QuotesTable + " order by rand() limit 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	appInfo, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppInfoTable + " where status = 1 ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// for _, content := range recommended {
	// 	urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
	// 	content["photo"] = endPointURL

	// 	urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
	// 	content["background_photo"] = endPointURLBackgroundPhoto

	// 	if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
	// 		urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
	// 		content["share_content"] = endPointURLShareContent
	// 	}

	// 	if len(content["counsellor_photo"]) > 0 {
	// 		urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
	// 		content["counsellor_photo"] = endPointURLCounsellorPhoto
	// 	}

	// 	if content["type"] != "3" {
	// 		urlContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		_, endPointURLContent := UTIL.GetBaseURLAndEndpointFromURL(urlContent)
	// 		content["content"] = endPointURLContent
	// 	}
	// }

	response["recommended"] = recommended
	response["access_code"] = accessCode
	response["quote"] = quote[0]["quote"]
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["android_version"] = appInfo[0]["client_android_version"]
	response["ios_version"] = appInfo[0]["client_ios_version"]
	response["urls"] = CONSTANT.URLs
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
