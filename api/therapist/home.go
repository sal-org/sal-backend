package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"time"

	UTIL "salbackend/util"
)

// Home godoc
// @Tags Therapist Home
// @Summary Get home page content
// @Router /therapist/home [get]
// @Param therapist_id query string false "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	therapistID, ok := UTIL.Required(r.FormValue("therapist_id"), "Therapist ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, therapistID, CONSTANT.ShowDialog, response)
		return
	}

	if len(r.FormValue("therapist_id")) > 0 {
		active := DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"status": "1", "therapist_id": r.FormValue("therapist_id")})
		if !active {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// get latest content for recommended
	recommended, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 1 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest videos
	// videos, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.VideoContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // get latest audios
	// audios, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.AudioContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // get latest articles
	// articles, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.ArticleContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

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

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
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
			// Get client information for the selected appointment.
			clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", latestUpcoming["client_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			latestUpcoming["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
			latestUpcoming["client_photo"] = clients[0]["photo"]

			response["appointments"] = latestUpcoming
		}

	}

	// if len(appointments) > 0 {

	// 	if appointments[0]["date"] == UTIL.GetCurrentTime().Add(330*time.Minute).Format("2006-01-02") {
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
	// 				clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", appointments[1]["client_id"])
	// 				if !ok {
	// 					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 					return
	// 				}

	// 				// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 				// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 				// counsellors[0]["photo"] = endPointURL

	// 				appointments[1]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 				appointments[1]["client_photo"] = clients[0]["photo"]

	// 				// virtual appointments
	// 				response["appointments"] = appointments[1]
	// 			} else {
	// 				response["appointments"] = make(map[string]string)
	// 			}
	// 		} else {
	// 			clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", appointments[0]["client_id"])
	// 			if !ok {
	// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 				return
	// 			}

	// 			// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 			// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 			// counsellors[0]["photo"] = endPointURL

	// 			appointments[0]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 			appointments[0]["client_photo"] = clients[0]["photo"]
	// 			// virtual appointments
	// 			response["appointments"] = appointments[0]
	// 		}
	// 	} else {
	// 		clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", appointments[0]["client_id"])
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 		// counsellors[0]["photo"] = endPointURL

	// 		appointments[0]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 		appointments[0]["client_photo"] = clients[0]["photo"]
	// 		// appointments[0]["counsellor_education"] = counsellors[0]["education"]

	// 		// virtual appointments
	// 		response["appointments"] = appointments[0]
	// 	}
	// } else {
	// 	response["appointments"] = make(map[string]string)
	// }

	// get upcoming appointments both to be started and started
	inpersonAppointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time, company_name, company_location, counselling_address from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
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
			// Get client information for the selected appointment.
			clientsInPerson, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", upcomingInPerson["client_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			upcomingInPerson["client_name"] = clientsInPerson[0]["first_name"] + " " + clientsInPerson[0]["last_name"]
			upcomingInPerson["client_photo"] = clientsInPerson[0]["photo"]

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
	// 				clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", inpersonAppointments[1]["client_id"])
	// 				if !ok {
	// 					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 					return
	// 				}

	// 				// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 				// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 				// counsellors[0]["photo"] = endPointURL

	// 				inpersonAppointments[1]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 				inpersonAppointments[1]["client_photo"] = clients[0]["photo"]

	// 				// virtual appointments
	// 				response["inperson_appointments"] = inpersonAppointments[1]
	// 			} else {
	// 				response["inperson_appointments"] = make(map[string]string)
	// 			}
	// 		} else {
	// 			clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", inpersonAppointments[0]["client_id"])
	// 			if !ok {
	// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 				return
	// 			}

	// 			// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 			// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 			// counsellors[0]["photo"] = endPointURL

	// 			inpersonAppointments[0]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 			inpersonAppointments[0]["client_photo"] = clients[0]["photo"]

	// 			// virtual appointments
	// 			response["inperson_appointments"] = inpersonAppointments[0]
	// 		}
	// 	} else {
	// 		clients, status, ok := DB.SelectProcess("select client_id as id, first_name, last_name, photo from "+CONSTANT.ClientsTable+" where client_id = ? ", inpersonAppointments[0]["client_id"])
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		// url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 		// counsellors[0]["photo"] = endPointURL

	// 		inpersonAppointments[0]["client_name"] = clients[0]["first_name"] + " " + clients[0]["last_name"]
	// 		inpersonAppointments[0]["client_photo"] = clients[0]["photo"]
	// 		// inpersonAppointments[0]["counsellor_education"] = counsellors[0]["education"]

	// 		// virtual appointments
	// 		response["inperson_appointments"] = inpersonAppointments[0]
	// 	}
	// } else {
	// 	response["inperson_appointments"] = make(map[string]string)
	// }

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where counsellor_id = ? and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc, time asc", r.FormValue("therapist_id"))
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

	response["recommended"] = recommended
	// response["videos"] = videos
	// response["audios"] = audios
	// response["articles"] = articles
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["urls"] = CONSTANT.URLs
	response["android_version"] = appInfo[0]["therapist_android_version"]
	response["ios_version"] = appInfo[0]["therapist_ios_version"]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
