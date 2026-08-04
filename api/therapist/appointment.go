package therapist

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
// @Tags Therapist Appointment
// @Summary Get therapist upcoming appointments
// @Router /therapist/appointment/upcoming [get]
// @Param therapist_id query string true "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentsUpcoming(w http.ResponseWriter, r *http.Request) {
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

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Add(330*time.Minute).Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, photo, date_of_birth, gender, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// for _, client := range clients {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, client["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	client["photo"] = endPointURL
	// }

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func InPersonAppointmentsUpcoming(w http.ResponseWriter, r *http.Request) {
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

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming appointments both to be started and started
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentToBeStarted+", "+CONSTANT.AppointmentStarted+") and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, photo, date_of_birth, gender, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// slots, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsScheduleTable + " where status = '1' and date >= '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' and date < '" + UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02") + "' order by date asc")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	for _, appoint := range appointments {
		lastAppointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and client_id = ? and status = 3 order by date desc", appoint["counsellor_id"], appoint["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(lastAppointments) == 0 {
			appoint["last_session_date"] = ""
			appoint["session_type"] = "New"
		} else {
			appoint["last_session_date"] = lastAppointments[0]["date"]
			appoint["session_type"] = "Repeat"
		}
	}

	// for _, client := range clients {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, client["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	client["photo"] = endPointURL
	// }

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	// response["slots"] = slots
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
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

	var response = make(map[string]any)
	var appointments []map[string]string
	var statusMessage string
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

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	// get past completed appointments
	appointmentsCompleted, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.AppointmentCompleted})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	appointmentsReviewed, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status in ('"+CONSTANT.AppointmentStarted+"', '"+CONSTANT.AppointmentToBeStarted+"') and date <= ?", r.FormValue("therapist_id"), UTIL.GetCurrentTime().Add(330*time.Minute).Format("2006-01-02"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var appointmentsReviewedList []map[string]string

	for _, appointment := range appointmentsReviewed {
		if appointment["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
			localTime := 0
			loc, _ := time.LoadLocation("Asia/Kolkata")
			now := time.Now().In(loc)
			if now.Minute() >= 30 {
				localTime = now.Hour()*2 + 1
			} else {
				localTime = now.Hour() * 2
			}
			appointmentTime, _ := strconv.Atoi(appointment["time"])

			if appointmentTime < localTime {
				appointmentsReviewedList = append(appointmentsReviewedList, appointment)
			}

		} else {
			appointmentsReviewedList = append(appointmentsReviewedList, appointment)
		}
	}

	appointments = append(appointments, appointmentsCompleted...)
	appointments = append(appointments, appointmentsReviewedList...)

	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
	appointmentIDs := UTIL.ExtractValuesFromArrayMap(appointments, "appointment_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, photo, date_of_birth, gender, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellors details
	counsellorRecords, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorRecordsNewVersion, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsFormLastestVersionTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorRecordsIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecords, "appointment_id")
	counsellorRecordsNewVersionIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecordsNewVersion, "appointment_id")
	counsellorRecordMap := UTIL.ConvertMapToKeyMap(counsellorRecordsNewVersion, "appointment_id")
	counsellorRecordsIDs = append(counsellorRecordsIDs, counsellorRecordsNewVersionIDs...)

	idSet := make(map[string]struct{})
	for _, id := range counsellorRecordsIDs {
		idSet[id] = struct{}{}
	}

	for index, appointment := range appointments {
		idStr, ok := appointment["appointment_id"]
		if !ok {
			continue
		}

		if _, exists := idSet[idStr]; exists {
			if counsellorRecordMap[idStr]["status"] == "2" {
				appointments[index]["is_counsellor_record_filled"] = "yes"
			} else {
				if counsellorRecordMap[idStr]["status"] == "" {
					appointments[index]["is_counsellor_record_filled"] = "yes"
				} else {
					appointments[index]["is_counsellor_record_filled"] = "no"
				}
			}
			// appointments[index]["is_counsellor_record_filled"] = "yes"
			// fmt.Println("Matched row:", row)
		} else {
			appointments[index]["is_counsellor_record_filled"] = "no"
		}

		switch appointment["status"] {
		case "3":
			if appointment["started_at"] == "" && appointment["ended_at"] == "" {
				statusMessage = getAppointmentStatusInText("8")
			} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
				statusMessage = getAppointmentStatusInText("7")
			} else if len(appointment["client_started_at"]) != 0 && len(appointment["client_ended_at"]) == 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else if len(appointment["client_started_at"]) == 0 && len(appointment["client_ended_at"]) != 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else if len(appointment["started_at"]) != 0 && len(appointment["ended_at"]) == 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else if len(appointment["started_at"]) == 0 && len(appointment["ended_at"]) != 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else {
				if UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes() > 10 {
					if UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes() < 10 {
						statusMessage = getAppointmentStatusInText("19")
					} else {
						statusMessage = getAppointmentStatusInText("3")
					}
				} else {
					statusMessage = getAppointmentStatusInText("14")
				}
			}

		case "4":
			if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
				statusMessage = getAppointmentStatusInText("12")
			} else {
				statusMessage = getAppointmentStatusInText("4")
			}
		case "5":
			if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
				statusMessage = getAppointmentStatusInText("13")
			} else {
				statusMessage = getAppointmentStatusInText("5")
			}
		default:
			statusMessage = getAppointmentStatusInText(appointment["status"])
		}

		appointments[index]["status_text"] = statusMessage
	}

	// for _, client := range clients {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, client["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	client["photo"] = endPointURL
	// }

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func InPersonAppointmentsPast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)
	var appointments []map[string]string
	var statusMessage string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get past completed appointments
	// appointments, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.AppointmentCompleted})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	therapistID, ok := UTIL.Required(r.FormValue("therapist_id"), "Therapist ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, therapistID, CONSTANT.ShowDialog, response)
		return
	}

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	appointmentsCompleted, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and status in ("+CONSTANT.AppointmentCompleted+", "+CONSTANT.AppointmentNoShowClient+") order by date desc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	/*appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status = ? order by date desc", r.FormValue("therapist_id"), CONSTANT.AppointmentCompleted)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}*/

	// appointmentsCompleted, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.AppointmentCompleted})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	appointmentsReviewed, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and status in ('"+CONSTANT.AppointmentStarted+"', '"+CONSTANT.AppointmentToBeStarted+"') and date <= ?", r.FormValue("therapist_id"), UTIL.GetCurrentTime().Add(330*time.Minute).Format("2006-01-02"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var appointmentsReviewedList []map[string]string

	for _, appointment := range appointmentsReviewed {
		if appointment["date"] == UTIL.GetCurrentTime().Format("2006-01-02") {
			localTime := 0
			loc, _ := time.LoadLocation("Asia/Kolkata")
			now := time.Now().In(loc)
			if now.Minute() >= 30 {
				localTime = now.Hour()*2 + 1
			} else {
				localTime = now.Hour() * 2
			}
			appointmentTime, _ := strconv.Atoi(appointment["time"])

			if appointmentTime < localTime {
				appointmentsReviewedList = append(appointmentsReviewedList, appointment)
			}

		} else {
			appointmentsReviewedList = append(appointmentsReviewedList, appointment)
		}
	}

	appointments = append(appointments, appointmentsCompleted...)
	appointments = append(appointments, appointmentsReviewedList...)

	// get client ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
	appointmentIDs := UTIL.ExtractValuesFromArrayMap(appointments, "appointment_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, photo, date_of_birth, gender, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellors details
	counsellorRecords, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorRecordsNewVersion, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsFormLastestVersionTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellorRecordsIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecords, "appointment_id")
	counsellorRecordsNewVersionIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecordsNewVersion, "appointment_id")
	counsellorRecordMap := UTIL.ConvertMapToKeyMap(counsellorRecordsNewVersion, "appointment_id")
	counsellorRecordsIDs = append(counsellorRecordsIDs, counsellorRecordsNewVersionIDs...)

	idSet := make(map[string]struct{})
	for _, id := range counsellorRecordsIDs {
		idSet[id] = struct{}{}
	}

	for index, appointment := range appointments {
		idStr, ok := appointment["appointment_id"]
		if !ok {
			continue
		}

		if _, exists := idSet[idStr]; exists {
			if counsellorRecordMap[idStr]["status"] == "2" {
				appointments[index]["is_counsellor_record_filled"] = "yes"
			} else {
				if counsellorRecordMap[idStr]["status"] == "" {
					appointments[index]["is_counsellor_record_filled"] = "yes"
				} else {
					appointments[index]["is_counsellor_record_filled"] = "no"
				}
			}
			// appointments[index]["is_counsellor_record_filled"] = "yes"
			// fmt.Println("Matched row:", row)
		} else {
			appointments[index]["is_counsellor_record_filled"] = "no"
		}

		switch appointment["status"] {
		case "3":
			if appointment["started_at"] == "" && appointment["ended_at"] == "" {
				statusMessage = getAppointmentStatusInText("8")
			} else if len(appointment["started_at"]) != 0 && len(appointment["ended_at"]) == 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else if len(appointment["started_at"]) == 0 && len(appointment["ended_at"]) != 0 {
				statusMessage = getAppointmentStatusInText("14")
			} else {
				if UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes() > 10 {
					statusMessage = getAppointmentStatusInText("3")
				} else {
					statusMessage = getAppointmentStatusInText("14")
				}
			}

		case "4":
			if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
				statusMessage = getAppointmentStatusInText("12")
			} else {
				statusMessage = getAppointmentStatusInText("4")
			}
		case "5":
			if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
				statusMessage = getAppointmentStatusInText("13")
			} else {
				statusMessage = getAppointmentStatusInText("5")
			}
		case "7":
			statusMessage = getAppointmentStatusInText("7")
		default:
			statusMessage = getAppointmentStatusInText(appointment["status"])
		}

		appointments[index]["status_text"] = statusMessage
	}

	// for _, client := range clients {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, client["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	client["photo"] = endPointURL
	// }

	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["appointments"] = appointments
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
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

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	clientName, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"email"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(clientName[0]["email"], "@")

	companyName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1], "status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update therapist slots
	// remove previous slot
	// date, _ := time.Parse("2006-01-02", appointment[0]["date"])
	// get schedules for a weekday
	// schedules, status, ok := DB.SelectProcess("select `"+appointment[0]["time"]+"` from "+CONSTANT.SchedulesTable+" where counsellor_id = ? and weekday = ?", appointment[0]["counsellor_id"], strconv.Itoa(int(date.Weekday())))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// sometimes there will be no schedules. situation will be automatically taken care of below

	// update therapist availability
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": appointment[0]["counsellor_id"],
			"date":          appointment[0]["date"],
		},
		map[string]string{
			appointment[0]["time"]: CONSTANT.SlotUnavailable, // update availability to the latest one
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

	DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentCounsellorCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// if domainName[1] == "clovemind.com" {

	// 	DB.UpdateSQL(CONSTANT.ClientCounsellingLimitTable,
	// 		map[string]string{
	// 			"appointment_id": r.FormValue("appointment_id"),
	// 		},
	// 		map[string]string{
	// 			"status":      CONSTANT.AppointmentUserCancelled,
	// 			"modified_at": UTIL.GetCurrentTime().String(),
	// 		},
	// 	)
	// }

	if len(companyName) > 0 {
		// update appointment date and time
		// add a slot to appointments
		DB.ExecuteSQL("update "+CONSTANT.AppointmentSlotsTable+" set slots_remaining = slots_remaining + 1 where order_id = ?", appointment[0]["order_id"])
	}

	// add penalty for therapist for cancelling
	// add to therapist payments
	// get invoice details

	// if UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(UTIL.ConvertTimezone(UTIL.GetCurrentTime(), "330")).Hours() <= 4 {

	// 	invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// 	if len(invoice) > 0 {
	// 		// get order details
	// 		order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}
	// 		paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
	// 		discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
	// 		amountBeforeDiscount := paidAmount + discount
	// 		if amountBeforeDiscount > 0 { // add only if amount paid
	// 			slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

	// 			amountFor1Session := amountBeforeDiscount / slotsBought // for 1 counselling session
	// 			cancellationCharges := amountFor1Session

	// 			DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
	// 				"counsellor_id": appointment[0]["counsellor_id"],
	// 				"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
	// 				"description":   "Therapist Cancellation",
	// 				"amount":        strconv.FormatFloat(-cancellationCharges, 'f', 2, 64),
	// 				"status":        CONSTANT.PaymentActive,
	// 				"created_at":    UTIL.GetCurrentTime().String(),
	// 			}, "payment_id")
	// 		}
	// 	} else {

	// 		//UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Sub(time.Now()).Hours()

	// 		// get counsellor details
	// 		counsellor, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"corporate_price"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
	// 		if !ok {
	// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 			return
	// 		}

	// 		DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
	// 			"counsellor_id": appointment[0]["counsellor_id"],
	// 			"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
	// 			"description":   "Therapist Cancellation",
	// 			"amount":        "-" + counsellor[0]["corporate_price"],
	// 			"status":        CONSTANT.PaymentActive,
	// 			"created_at":    UTIL.GetCurrentTime().String(),
	// 		}, "payment_id")

	// 	}
	// }
	// invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	// if len(invoice) > 0 {
	// 	// get order details
	// 	order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"slots_bought"}, map[string]string{"order_id": appointment[0]["order_id"]})
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// 	paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
	// 	discount, _ := strconv.ParseFloat(invoice[0]["discount"], 64)
	// 	amountBeforeDiscount := paidAmount + discount
	// 	if amountBeforeDiscount > 0 { // add only if amount paid
	// 		slotsBought, _ := strconv.ParseFloat(order[0]["slots_bought"], 64)

	// 		amountFor1Session := amountBeforeDiscount / slotsBought // for 1 counselling session
	// 		cancellationCharges := amountFor1Session * CONSTANT.CounsellorCancellationCharges

	// 		DB.InsertWithUniqueID(CONSTANT.PaymentsTable, CONSTANT.PaymentsDigits, map[string]string{
	// 			"counsellor_id": appointment[0]["counsellor_id"],
	// 			"heading":       DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
	// 			"description":   "Cancellation",
	// 			"amount":        strconv.FormatFloat(-cancellationCharges, 'f', 2, 64),
	// 			"status":        CONSTANT.PaymentActive,
	// 			"created_at":    UTIL.GetCurrentTime().String(),
	// 		}, "payment_id")
	// 	}
	// }

	// send appointment cancel notification, email to client
	therapist, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "email", "phone"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[0]["client_id"]})

	// remove all previous notifications
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["client_id"])
	UTIL.RemoveNotification(r.FormValue("appointment_id"), appointment[0]["counsellor_id"])

	// remove all previous message for client
	UTIL.RemoveMessage(r.FormValue("appointment_id"), client[0]["phone"])

	// remove all previous message for therpist
	UTIL.RemoveMessage(r.FormValue("appointment_id"), therapist[0]["phone"])

	UTIL.SendNotification(
		CONSTANT.CounsellorAppointmentCancelClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientContent,
			map[string]string{
				"###therapist_name###": therapist[0]["first_name"],
				"###date_time###":      UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
		"",
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
		"",
	)

	filepath_text := "htmlfile/emailmessagebody.html"

	// send email for client
	emaildata := Model.EmailBodyMessageModel{
		Name: client[0]["first_name"],
		Message: UTIL.ReplaceNotificationContentInString(
			CONSTANT.CounsellorAppointmentCancelClientBodyEmailBody,
			map[string]string{
				"###therapist_name###": therapist[0]["first_name"],
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
		Name:    therapist[0]["first_name"],
		Message: CONSTANT.CounsellorAppointmentCancelCounsellorBodyEmailBody,
	}

	emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
	// email for counsellor
	UTIL.SendEmail(
		CONSTANT.CounsellorAppointmentCancelCounsellorTitle,
		emailBody1,
		therapist[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SendMessage(
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentCancellationToCounsellorTextMessage,
			map[string]string{
				"###userName###":  client[0]["first_name"],
				"###user_Name###": therapist[0]["first_name"],
				"###date###":      appointment[0]["date"],
				"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appointment[0]["time"]),
			},
		),
		CONSTANT.TransactionalRouteTextMessage,
		client[0]["phone"],
		UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
		r.FormValue("appointment_id"),
		CONSTANT.InstantSendTextMessage,
	)

	// appointment cancel 7 gays gap for email trigger
	appointmentCancel, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where counsellor_id = ? and status = '"+CONSTANT.AppointmentCounsellorCancelled+"' order by date desc", appointment[0]["counsellor_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(appointmentCancel) > 1 {
		format := "2006-01-02 15:04:05"
		previousDate, _ := time.Parse(format, appointmentCancel[1]["date"]+" 00:00:00")
		lastestDate, _ := time.Parse(format, appointmentCancel[0]["date"]+" 00:00:00")

		diff := lastestDate.Sub(previousDate)

		days := int(diff.Hours() / 24)

		if days < 8 {

			therapistC, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "email", "phone"}, map[string]string{"therapist_id": appointment[0]["counsellor_id"]})

			previousClient, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "last_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[1]["client_id"]})

			latestClient, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "last_name", "timezone", "email", "phone"}, map[string]string{"client_id": appointment[0]["client_id"]})

			htmlPath := "htmlfile/appointment7daysgap.html"

			data := Model.EmailDataForCounsellorCancellation{
				First_Name:            therapistC[0]["first_name"],
				Last_Name:             therapistC[0]["last_name"],
				Previous_Date:         appointmentCancel[1]["date"],
				Previous_Client_Name:  previousClient[1]["first_name"] + " " + previousClient[1]["last_name"],
				Previous_Client_Email: latestClient[0]["email"],
				Latest_Date:           appointmentCancel[0]["date"],
				Lastest_Client_Name:   latestClient[0]["first_name"] + " " + latestClient[0]["last_name"],
				Lastest_Client_Email:  latestClient[0]["email"],
			}

			emailbody := UTIL.GetHTMLTemplateForCounsellorCancellation(data, htmlPath)

			UTIL.SendEmail(
				CONSTANT.CounsellorCancelAppointmentTitle,
				emailbody,
				"corp.wellness@clovemind.com",
				CONSTANT.InstantSendEmailMessage,
			)
		}
	}

	// UTIL.SendMessage(
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.CounsellorAppointmentCancellationToClientTextMessage,
	// 		map[string]string{
	// 			"###client_name###": client[0]["first_name"],
	// 		},
	// 	),
	// 	CONSTANT.TransactionalRouteTextMessage,
	// 	client[0]["phone"],
	// 	CONSTANT.LaterSendTextMessage,
	// )

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

	var response = make(map[string]any)

	var roleStr, agora_token, uidStr, channelName string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	uidStr = generateRandomID(6)

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

		uidSt := generateRandomID(9)
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

func generateRandomID(maxlength int) string {
	const randomIDdigits = "123456789"
	b := make([]byte, maxlength)
	for i := range b {
		b[i] = randomIDdigits[rand.Intn(len(randomIDdigits))]
	}
	return string(b)
}

// AppointmentStart godoc
// @Tags Therapist Appointment
// @Summary Start an appointment
// @Router /therapist/appointment/start [put]
// @Param appointment_id query string true "Appointment ID to be started"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
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

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// agora, status, ok := DB.SelectSQL(CONSTANT.AgoraTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// if len(agora) == 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"asscoiate_id"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var sessionFor string

	if len(client[0]["asscoiate_id"]) == 0 {
		sessionFor = "Self"
	} else {
		sessionFor = "Family"
	}

	checkClientRecordForm, status, ok := DB.SelectProcess("select taken_sessions, total_session_needed from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? order by modified_at desc limit 5", appointment[0]["appointment_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(checkClientRecordForm) == 0 {

		counsellorRecordForm := map[string]string{}
		counsellorRecordForm["appointment_id"] = appointment[0]["appointment_id"]
		counsellorRecordForm["client_id"] = appointment[0]["client_id"]
		counsellorRecordForm["counsellor_id"] = appointment[0]["counsellor_id"]
		counsellorRecordForm["session_for"] = sessionFor
		counsellorRecordForm["session_mode"] = "Virtual"
		counsellorRecordForm["session_date"] = appointment[0]["date"]
		counsellorRecordForm["in_time"] = UTIL.GetIndiaCurrentTime()
		counsellorRecordForm["status"] = CONSTANT.AppointmentToBeStarted
		counsellorRecordForm["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.AppointmentDigits, counsellorRecordForm, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// domainName := strings.Split(client[0]["email"], "@")

	if len(appointment[0]["started_at"]) == 0 {
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
	}

	// if domainName[1] == "clovemind.com" {

	// 	DB.UpdateSQL(CONSTANT.ClientCounsellingLimitTable,
	// 		map[string]string{
	// 			"appointment_id": r.FormValue("appointment_id"),
	// 		},
	// 		map[string]string{
	// 			"status": CONSTANT.AppointmentStarted,
	// 		},
	// 	)
	// }

	// var allUsers []string

	// allUsers = append(allUsers, agora[0]["uid1"])
	// allUsers = append(allUsers, agora[0]["uid"])

	// if len(agora[0]["sid"]) == 0 {

	// 	sid, err := UTIL.AgoraRecordingCallStart(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["token"], agora[0]["resource_id"])
	// 	if err != nil {
	// 		fmt.Println("cloud recording not started")
	// 	}

	// 	DB.UpdateSQL(CONSTANT.AgoraTable,
	// 		map[string]string{
	// 			"agora_id": agora[0]["agora_id"],
	// 		},
	// 		map[string]string{
	// 			"sid":         sid,
	// 			"status":      CONSTANT.AgoraCallStart1,
	// 			"modified_at": UTIL.GetCurrentTime().String(),
	// 		},
	// 	)

	// }

	// send appointment join the call notification to Client
	UTIL.SendNotification(
		CONSTANT.ClientAppointmentHasBeenStartedHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientAppointmentHasBeenStartedContent,
			map[string]string{
				"###clientname###":    DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
				"###therapistname###": DB.QueryRowSQL("select first_name from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"]),
			},
		),
		appointment[0]["client_id"],
		CONSTANT.ClientType,
		UTIL.GetCurrentTime().String(),
		CONSTANT.NotificationSent,
		r.FormValue("appointment_id"),
		"",
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AppointmentInPersonStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"asscoiate_id"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get upcoming appointments both to be started and started
	isAppointmentsStarted, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonAppointmentsTable+" where counsellor_id = ? and status = "+CONSTANT.AppointmentStarted+" and date = '"+UTIL.GetCurrentTime().Format("2006-01-02")+"'", appointment[0]["counsellor_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(isAppointmentsStarted) > 0 {
		for _, v := range isAppointmentsStarted {
			// update appointment as ended
			DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable,
				map[string]string{
					"appointment_id": v["appointment_id"],
				},
				map[string]string{
					"status":   CONSTANT.AppointmentCompleted,
					"ended_at": UTIL.BuildToDteTime(v["started_at"]).Add(time.Minute * 60).String(), // one hour is added in the started_at time
				},
			)
		}
	}

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.InPersonOrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.InPersonAppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update appointment as started
	DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":     CONSTANT.AppointmentStarted,
			"started_at": UTIL.GetCurrentTime().String(),
		},
	)

	var sessionFor string

	if len(client[0]["asscoiate_id"]) == 0 {
		sessionFor = "Self"
	} else {
		sessionFor = "Family"
	}

	checkClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? order by modified_at desc limit 5", appointment[0]["appointment_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(checkClientRecordForm) == 0 {

		counsellorRecordForm := map[string]string{}
		counsellorRecordForm["appointment_id"] = appointment[0]["appointment_id"]
		counsellorRecordForm["client_id"] = appointment[0]["client_id"]
		counsellorRecordForm["counsellor_id"] = appointment[0]["counsellor_id"]
		counsellorRecordForm["session_for"] = sessionFor
		counsellorRecordForm["session_mode"] = "In-Person"
		counsellorRecordForm["session_date"] = appointment[0]["date"]
		counsellorRecordForm["in_time"] = UTIL.GetIndiaCurrentTime()
		counsellorRecordForm["status"] = CONSTANT.AppointmentToBeStarted
		counsellorRecordForm["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.AppointmentDigits, counsellorRecordForm, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// send appointment join the call notification to Client
	// UTIL.SendNotification(
	// 	CONSTANT.ClientAppointmentHasBeenStartedHeading,
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentHasBeenStartedContent,
	// 		map[string]string{
	// 			"###clientname###":    DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
	// 			"###therapistname###": DB.QueryRowSQL("select first_name from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", appointment[0]["counsellor_id"]),
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

func AppointmentInPersonNoShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.InPersonOrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.InPersonAppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"asscoiate_id", "first_name", "email"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	newVersionClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? order by modified_at desc limit 5", appointment[0]["appointment_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(newVersionClientRecordForm) == 0 {

		var sessionFor string

		if len(client[0]["asscoiate_id"]) == 0 {
			sessionFor = "Self"
		} else {
			sessionFor = "Family"
		}

		counsellorRecordForm := map[string]string{}
		counsellorRecordForm["appointment_id"] = appointment[0]["appointment_id"]
		counsellorRecordForm["client_id"] = appointment[0]["client_id"]
		counsellorRecordForm["counsellor_id"] = appointment[0]["counsellor_id"]
		counsellorRecordForm["session_for"] = sessionFor
		counsellorRecordForm["session_mode"] = "In-Person"
		counsellorRecordForm["session_date"] = appointment[0]["date"]
		counsellorRecordForm["in_time"] = UTIL.GetIndiaCurrentTime()
		counsellorRecordForm["status"] = CONSTANT.AppointmentToBeStarted
		counsellorRecordForm["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.AppointmentDigits, counsellorRecordForm, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		DB.UpdateSQL(CONSTANT.CounsellorRecordsFormLastestVersionTable,
			map[string]string{
				"record_id": newVersionClientRecordForm[0]["record_id"],
			},
			map[string]string{
				"out_time":    UTIL.GetIndiaCurrentTime(),
				"modified_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

	// update appointment as started
	DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentNoShowClient,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	// update counsellor availability
	DB.UpdateSQL(CONSTANT.InPersonSLotsTable,
		map[string]string{
			"counsellor_id":    appointment[0]["counsellor_id"],
			"date":             appointment[0]["date"],
			"company_name":     appointment[0]["company_name"],
			"company_location": appointment[0]["company_location"],
		},
		map[string]string{
			// this is for cancel slot menthod  UTIL.CheckIfScheduleAvailable(schedules, appointment[0]["time"])
			appointment[0]["time"]: CONSTANT.SlotAvailable, // update availability to the latest one
		},
	)

	// send appointment join the call notification to Client
	// UTIL.SendNotification(
	// 	CONSTANT.ClientAppointmentHasBeenStartedHeading,
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentHasBeenStartedContent,
	// 		map[string]string{
	// 			"###clientname###":    DB.QueryRowSQL("select first_name from "+CONSTANT.ClientsTable+" where client_id = ?", appointment[0]["client_id"]),
	// 			"###therapistname###": DB.QueryRowSQL("select first_name from "+CONSTANT.CounsellorsTable+" where counsellor_id = ?", appointment[0]["counsellor_id"]),
	// 		},
	// 	),
	// 	appointment[0]["client_id"],
	// 	CONSTANT.ClientType,
	// 	UTIL.GetCurrentTime().String(),
	// 	CONSTANT.NotificationSent,
	// 	r.FormValue("appointment_id"),
	// )

	// Payment receipt
	emaildata1 := Model.EmailBodyMessageWithNameModel{
		Name: client[0]["first_name"],
	}

	emailBody1 := UTIL.GetHTMLTemplateForClientNoShowText(emaildata1, "htmlfile/noShowByTheClient.html")
	// email for client
	UTIL.SendEmail(
		CONSTANT.ClientAppointmentBookCounsellorTitle,
		emailBody1,
		client[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AppointmentEnd godoc
// @Tags Therapist Appointment
// @Summary End an appointment
// @Router /therapist/appointment/end [put]
// @Param appointment_id query string true "Appointment ID to be ended"
// @Param uid query string true "User ID to be started"
// @Security JWTAuth
// @Produce json
// @Success 200
func AppointmentEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
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

	// get therapist type
	therapistType := DB.QueryRowSQL("select type from "+CONSTANT.OrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.AppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(therapistType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// agora, status, ok := DB.SelectSQL(CONSTANT.AgoraTable, []string{"*"}, map[string]string{"appointment_id": appointment[0]["appointment_id"]})
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"asscoiate_id"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	newVersionClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? order by modified_at desc limit 5", appointment[0]["appointment_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(newVersionClientRecordForm) == 0 {

		var sessionFor string

		if len(client[0]["asscoiate_id"]) == 0 {
			sessionFor = "Self"
		} else {
			sessionFor = "Family"
		}

		counsellorRecordForm := map[string]string{}
		counsellorRecordForm["appointment_id"] = appointment[0]["appointment_id"]
		counsellorRecordForm["client_id"] = appointment[0]["client_id"]
		counsellorRecordForm["counsellor_id"] = appointment[0]["counsellor_id"]
		counsellorRecordForm["session_for"] = sessionFor
		counsellorRecordForm["session_mode"] = "Virtual"
		counsellorRecordForm["session_date"] = appointment[0]["date"]
		counsellorRecordForm["in_time"] = UTIL.GetIndiaCurrentTime()
		counsellorRecordForm["status"] = CONSTANT.AppointmentToBeStarted
		counsellorRecordForm["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.AppointmentDigits, counsellorRecordForm, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		DB.UpdateSQL(CONSTANT.CounsellorRecordsFormLastestVersionTable,
			map[string]string{
				"record_id": newVersionClientRecordForm[0]["record_id"],
			},
			map[string]string{
				"out_time":    UTIL.GetIndiaCurrentTime(),
				"modified_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

	if len(r.FormValue("waiting")) > 0 {

		// update appointment as completed
		DB.UpdateSQL(CONSTANT.AppointmentsTable,
			map[string]string{
				"appointment_id": r.FormValue("appointment_id"),
			},
			map[string]string{
				"ended_at": UTIL.GetCurrentTime().String(),
			},
		)

	} else {

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

		// if appointment[0]["ended_at"] != "" {
		// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentAlreadyCompletedMessage, CONSTANT.ShowDialog, response)
		// 	return
		// }

		// if len(agora[0]["fileNameInMp4"]) == 0 && len(agora[0]["fileNameInM3U8"]) == 0 {
		// 	UTIL.AgoraRecordingCallStop(agora[0]["uid"], agora[0]["appointment_id"], agora[0]["resource_id"], agora[0]["sid"])

		// 	DB.UpdateSQL(CONSTANT.AgoraTable,
		// 		map[string]string{
		// 			"appointment_id": r.FormValue("appointment_id"),
		// 		},
		// 		map[string]string{
		// 			"fileNameInMp4":  "recordingfile/" + agora[0]["sid"] + "_" + agora[0]["appointment_id"] + "_0.mp4",
		// 			"fileNameInM3U8": "recordingfile/" + agora[0]["sid"] + "_" + agora[0]["appointment_id"] + ".m3u8",
		// 			"status":         CONSTANT.AgoraCallStop1,
		// 			"modified_at":    UTIL.GetCurrentTime().String(),
		// 		},
		// 	)

		// 	DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable,
		// 		map[string]string{
		// 			"appointment_id": r.FormValue("appointment_id"),
		// 		},
		// 		map[string]string{
		// 			"counsellor_mp4": "recordingfile/" + agora[0]["sid"] + "_" + agora[0]["appointment_id"] + "_0.mp4",
		// 			"status":         CONSTANT.QualityCheckLinkInsert,
		// 			"modified_at":    UTIL.GetCurrentTime().String(),
		// 		},
		// 	)

		// }

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AppointmentInPersonEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	appointmentID, ok := UTIL.Required(r.FormValue("appointment_id"), "Appointment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, appointmentID, CONSTANT.ShowDialog, response)
		return
	}

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor type
	counsellorType := DB.QueryRowSQL("select type from "+CONSTANT.InPersonOrderClientAppointmentTable+" where order_id in (select order_id from "+CONSTANT.InPersonAppointmentsTable+" where appointment_id = ?)", r.FormValue("appointment_id"))
	if !strings.EqualFold(counsellorType, CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update appointment as completed
	DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":   CONSTANT.AppointmentCompleted,
			"ended_at": UTIL.GetCurrentTime().String(),
		},
	)

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"asscoiate_id"}, map[string]string{"client_id": appointment[0]["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	newVersionClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? order by modified_at desc limit 5", appointment[0]["appointment_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(newVersionClientRecordForm) == 0 {

		var sessionFor string

		if len(client[0]["asscoiate_id"]) == 0 {
			sessionFor = "Self"
		} else {
			sessionFor = "Family"
		}

		counsellorRecordForm := map[string]string{}
		counsellorRecordForm["appointment_id"] = appointment[0]["appointment_id"]
		counsellorRecordForm["client_id"] = appointment[0]["client_id"]
		counsellorRecordForm["counsellor_id"] = appointment[0]["counsellor_id"]
		counsellorRecordForm["session_for"] = sessionFor
		counsellorRecordForm["session_mode"] = "In-Person"
		counsellorRecordForm["session_date"] = appointment[0]["date"]
		counsellorRecordForm["in_time"] = UTIL.GetIndiaCurrentTime()
		counsellorRecordForm["status"] = CONSTANT.AppointmentToBeStarted
		counsellorRecordForm["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.AppointmentDigits, counsellorRecordForm, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		DB.UpdateSQL(CONSTANT.CounsellorRecordsFormLastestVersionTable,
			map[string]string{
				"record_id": newVersionClientRecordForm[0]["record_id"],
			},
			map[string]string{
				"out_time":    UTIL.GetIndiaCurrentTime(),
				"modified_at": UTIL.GetCurrentTime().String(),
			},
		)
	}

	// send client for rating

	// UTIL.SendNotification(
	// 	CONSTANT.ClientAppointmentFeedbackHeading,
	// 	UTIL.ReplaceNotificationContentInString(
	// 		CONSTANT.ClientAppointmentFeedbackContent,
	// 		map[string]string{
	// 			"###counsellor_name###": DB.QueryRowSQL("select first_name from "+CONSTANT.TherapistsTable+" where therapist_id = ?", appointment[0]["counsellor_id"]),
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

// utils for reports
func getAppointmentStatusInText(status string) string {
	switch status {
	case CONSTANT.AppointmentToBeStarted:
		return "Both no-show"
	case CONSTANT.AppointmentStarted:
		return "To be Reviewed"
	case CONSTANT.AppointmentCompleted:
		return "Completed"
	case CONSTANT.AppointmentUserCancelled:
		return "Client cancelled"
	case CONSTANT.AppointmentUserCancelledWithin4Hour:
		return "Client cancelled(to be paid)"
	case CONSTANT.AppointmentCounsellorCancelledWithin4Hour:
		return "Therapist Cancelled - To be charged"
	case CONSTANT.AppointmentInTheReview:
		return "To be Reviewed"
	case CONSTANT.AppointmentCounsellorCancelled:
		return "Therapist Cancelled"
	case CONSTANT.AppointmentNoShowClient:
		return "Client No-Show"
	case CONSTANT.AppointmentNoShowCounsellor:
		return "Therapist No-Show"
	case CONSTANT.AppointmentIncompleteSession:
		return "Incomplete Session"
	case CONSTANT.AppointmentIncompleteSessionDueToTechIssue:
		return "Incomplete Session(tech issue)"
	default:
		return ""
	}
}
