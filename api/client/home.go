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
func Home(w http.ResponseWriter, r *http.Request, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var encryptedResponse = make(map[string]interface{})
	var recommended []map[string]string
	var ok bool
	var status string

	accessCode := ""

	if len(body["client_id"]) > 0 {

		active := DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"status": "1", "client_id": body["client_id"]})
		if !active {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}

		client, status, ok := DB.SelectProcess("select topic_ids, email from "+CONSTANT.ClientsTable+" where client_id = ? and status = 1", body["client_id"])
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
		if len(body["platform"]) > 0 {
			updateClient["platform"] = body["platform"]
		}
		if len(body["version"]) > 0 {
			updateClient["version"] = body["version"]
		}
		if len(body["timezone"]) > 0 {
			updateClient["timezone"] = body["timezone"]
		}
		updateClient["last_active_time"] = UTIL.GetCurrentTime().String()

		status, ok = DB.UpdateSQL(CONSTANT.ClientsTable, map[string]string{"client_id": body["client_id"]}, updateClient)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get upcoming appointments both to be started and started
		appointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time from "+CONSTANT.AppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", body["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(appointments) > 0 {

			if appointments[0]["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
				localTime := 0
				loc, _ := time.LoadLocation("Asia/Kolkata")
				now := time.Now().In(loc)
				if now.Minute() >= 30 {
					localTime = now.Hour()*2 + 1
				} else {
					localTime = now.Hour() * 2
				}
				appointmentTime, _ := strconv.Atoi(appointments[0]["time"])

				if appointmentTime+1 < localTime {
					if len(appointments) > 1 {
						counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[1]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[1]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[1]["counsellor_id"] + "'))")
						if !ok {
							UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
							return
						}

						url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
						_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
						counsellors[0]["photo"] = endPointURL

						appointments[1]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
						appointments[1]["counsellor_photo"] = counsellors[0]["photo"]
						appointments[1]["counsellor_education"] = counsellors[0]["education"]

						// virtual appointments
						response["appointments"] = appointments[1]
					} else {
						response["appointments"] = make(map[string]string)
					}
				} else {
					counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[0]["counsellor_id"] + "'))")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
					_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
					counsellors[0]["photo"] = endPointURL

					appointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
					appointments[0]["counsellor_photo"] = counsellors[0]["photo"]
					appointments[0]["counsellor_education"] = counsellors[0]["education"]
					// virtual appointments
					response["appointments"] = appointments[0]
				}
			} else {
				counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + appointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + appointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + appointments[0]["counsellor_id"] + "'))")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
				counsellors[0]["photo"] = endPointURL

				appointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
				appointments[0]["counsellor_photo"] = counsellors[0]["photo"]
				appointments[0]["counsellor_education"] = counsellors[0]["education"]

				// virtual appointments
				response["appointments"] = appointments[0]
			}
		} else {
			response["appointments"] = make(map[string]string)
		}

		// get upcoming appointments both to be started and started
		inpersonAppointments, status, ok := DB.SelectProcess("select appointment_id, counsellor_id, client_id, date, time, company_name, company_location, counselling_address from "+CONSTANT.InPersonAppointmentsTable+" where client_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", body["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(inpersonAppointments) > 0 {

			if inpersonAppointments[0]["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
				localTime := 0
				timeNow := UTIL.GetCurrentTime()
				timeNow = timeNow.Add(330 * time.Minute)
				if timeNow.Minute() >= 30 {
					localTime = timeNow.Hour()*2 + 1
				} else {
					localTime = timeNow.Hour() * 2
				}
				appointmentTime, _ := strconv.Atoi(inpersonAppointments[0]["time"])

				if appointmentTime+1 < localTime {

					if len(inpersonAppointments) > 1 {
						counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[1]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[1]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[1]["counsellor_id"] + "'))")
						if !ok {
							UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
							return
						}

						url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
						_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
						counsellors[0]["photo"] = endPointURL

						inpersonAppointments[1]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
						inpersonAppointments[1]["counsellor_photo"] = counsellors[0]["photo"]
						inpersonAppointments[1]["counsellor_education"] = counsellors[0]["education"]

						// virtual appointments
						response["inperson_appointments"] = inpersonAppointments[1]
					} else {
						response["inperson_appointments"] = make(map[string]string)
					}
				} else {
					counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[0]["counsellor_id"] + "'))")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
					_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
					counsellors[0]["photo"] = endPointURL

					inpersonAppointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
					inpersonAppointments[0]["counsellor_photo"] = counsellors[0]["photo"]
					inpersonAppointments[0]["counsellor_education"] = counsellors[0]["education"]

					// virtual appointments
					response["inperson_appointments"] = inpersonAppointments[0]
				}
			} else {
				counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable + " where counsellor_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select listener_id as id, first_name, last_name, photo, '' as education, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable + " where listener_id = ('" + inpersonAppointments[0]["counsellor_id"] + "')) union (select therapist_id as id, first_name, last_name, photo, education, " + CONSTANT.TherapistType + " as type from " + CONSTANT.TherapistsTable + " where therapist_id = ('" + inpersonAppointments[0]["counsellor_id"] + "'))")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellors[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
				counsellors[0]["photo"] = endPointURL

				inpersonAppointments[0]["counsellor_name"] = counsellors[0]["first_name"] + " " + counsellors[0]["last_name"]
				inpersonAppointments[0]["counsellor_photo"] = counsellors[0]["photo"]
				inpersonAppointments[0]["counsellor_education"] = counsellors[0]["education"]

				// virtual appointments
				response["inperson_appointments"] = inpersonAppointments[0]
			}
		} else {
			response["inperson_appointments"] = make(map[string]string)
		}

		// get upcoming booked events
		events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventInPersonTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventInPersonTable+" where user_id = ? and status = "+CONSTANT.OrderInProgress+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc, time asc", body["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(events) > 0 {
			// get counsellor details
			// get upcoming booked events
			urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, events[0]["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
			events[0]["photo"] = endPointURLPhoto

			urlBackGroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, events[0]["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLBackGroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackGroundPhoto)
			events[0]["background_photo"] = endPointURLBackGroundPhoto
			response["upcoming_events"] = events[0]
		} else {
			response["upcoming_events"] = make(map[string]string)
		}

		webinarOrder, status, ok := DB.SelectProcess("select * from "+CONSTANT.WebinarsBookTable+" where status = 1 and client_id = ? ", body["client_id"])
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

	} else {
		recommended, status, ok = DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 0 and status = 1 order by created_at desc limit 20")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	}

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

	for _, content := range recommended {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		content["photo"] = endPointURL

		urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		content["background_photo"] = endPointURLBackgroundPhoto

		if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType {
			urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
			content["share_content"] = endPointURLShareContent
		}
	}

	response["recommended"] = recommended
	response["access_code"] = accessCode
	response["quote"] = quote[0]["quote"]
	response["media_url"] = CONFIG.MediaURL
	response["android_version"] = appInfo[0]["client_android_version"]
	response["ios_version"] = appInfo[0]["client_ios_version"]
	response["urls"] = CONSTANT.URLs

	encrypt, _ := UTIL.EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY, CONSTANT.ENCRYPTION_SECRET_IV)
	if encrypt == "" {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, encryptedResponse)
}
