package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

func AvailabilityGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "slot_id":
			wheres = append(wheres, " id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	availabilitys, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsScheduleTable+where+" order by `date` desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of contents
	availabilitysCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InPersonSLotsScheduleTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["availabilitys"] = availabilitys
	response["availabilitys_count"] = availabilitysCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(availabilitysCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AvailabilityUpdate godoc
// @Tags Counsellor Availability
// @Summary Update counsellor availability hours
// @Router /counsellor/availability [put]
// @Param counsellor_id query string true "Counsellor ID to update availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// availabileDates, status, ok := DB.SelectProcess("select date,id from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and `date` = ?", r.FormValue("counsellor_id"), body["date"])
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	slotsDate := map[string]string{}

	slotsDate["counsellor_id"] = body["counsellor_id"]
	slotsDate["date"] = body["date"]
	slotsDate["from_time"] = body["fromTime"]
	slotsDate["to_time"] = body["toTime"]
	slotsDate["partner_name"] = body["companyName"]
	slotsDate["partner_location"] = body["companyLocation"]
	slotsDate["counselling_room"] = body["roomNo"]
	slotsDate["counselling_address"] = body["address"]
	slotsDate["status"] = body["status"]
	slotsDate["0"] = body["0"]
	slotsDate["1"] = body["1"]
	slotsDate["2"] = body["2"]
	slotsDate["3"] = body["3"]
	slotsDate["4"] = body["4"]
	slotsDate["5"] = body["5"]
	slotsDate["6"] = body["6"]
	slotsDate["7"] = body["7"]
	slotsDate["8"] = body["8"]
	slotsDate["9"] = body["9"]
	slotsDate["10"] = body["10"]
	slotsDate["11"] = body["11"]
	slotsDate["12"] = body["12"]
	slotsDate["13"] = body["13"]
	slotsDate["14"] = body["14"]
	slotsDate["15"] = body["15"]
	slotsDate["16"] = body["16"]
	slotsDate["17"] = body["17"]
	slotsDate["18"] = body["18"]
	slotsDate["19"] = body["19"]
	slotsDate["20"] = body["20"]
	slotsDate["21"] = body["21"]
	slotsDate["22"] = body["22"]
	slotsDate["23"] = body["23"]
	slotsDate["24"] = body["24"]
	slotsDate["25"] = body["25"]
	slotsDate["26"] = body["26"]
	slotsDate["27"] = body["27"]
	slotsDate["28"] = body["28"]
	slotsDate["29"] = body["29"]
	slotsDate["30"] = body["30"]
	slotsDate["31"] = body["31"]
	slotsDate["32"] = body["32"]
	slotsDate["33"] = body["33"]
	slotsDate["34"] = body["34"]
	slotsDate["35"] = body["35"]
	slotsDate["36"] = body["36"]
	slotsDate["37"] = body["37"]
	slotsDate["38"] = body["38"]
	slotsDate["39"] = body["39"]
	slotsDate["40"] = body["40"]
	slotsDate["41"] = body["41"]
	slotsDate["42"] = body["42"]
	slotsDate["43"] = body["43"]
	slotsDate["44"] = body["44"]
	slotsDate["45"] = body["45"]
	slotsDate["46"] = body["46"]
	slotsDate["47"] = body["47"]

	slots := map[string]string{}

	slots["available"] = body["status"]
	slots["0"] = body["0"]
	slots["1"] = body["1"]
	slots["2"] = body["2"]
	slots["3"] = body["3"]
	slots["4"] = body["4"]
	slots["5"] = body["5"]
	slots["6"] = body["6"]
	slots["7"] = body["7"]
	slots["8"] = body["8"]
	slots["9"] = body["9"]
	slots["10"] = body["10"]
	slots["11"] = body["11"]
	slots["12"] = body["12"]
	slots["13"] = body["13"]
	slots["14"] = body["14"]
	slots["15"] = body["15"]
	slots["16"] = body["16"]
	slots["17"] = body["17"]
	slots["18"] = body["18"]
	slots["19"] = body["19"]
	slots["20"] = body["20"]
	slots["21"] = body["21"]
	slots["22"] = body["22"]
	slots["23"] = body["23"]
	slots["24"] = body["24"]
	slots["25"] = body["25"]
	slots["26"] = body["26"]
	slots["27"] = body["27"]
	slots["28"] = body["28"]
	slots["29"] = body["29"]
	slots["30"] = body["30"]
	slots["31"] = body["31"]
	slots["32"] = body["32"]
	slots["33"] = body["33"]
	slots["34"] = body["34"]
	slots["35"] = body["35"]
	slots["36"] = body["36"]
	slots["37"] = body["37"]
	slots["38"] = body["38"]
	slots["39"] = body["39"]
	slots["40"] = body["40"]
	slots["41"] = body["41"]
	slots["42"] = body["42"]
	slots["43"] = body["43"]
	slots["44"] = body["44"]
	slots["45"] = body["45"]
	slots["46"] = body["46"]
	slots["47"] = body["47"]

	if len(body["id"]) > 0 {
		DB.UpdateSQL(CONSTANT.InPersonSLotsScheduleTable, map[string]string{"id": body["id"]}, slotsDate)

		if body["status"] == "0" {
			appointment, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": body["counsellor_id"], "date": body["date"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellor details
			counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": body["counsellor_id"]})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			if len(counsellor) == 0 {
				counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["counsellor_id"]})
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
			}

			filepath_text := "htmlfile/emailmessagebody.html"

			if len(appointment) != 0 {

				for _, appoint := range appointment {

					// var counsellor []map[string]string

					// switch appoint["type"] {
					// case CONSTANT.CounsellorType:
					// 	counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"counsellor_id": appoint["counsellor_id"]})
					// case CONSTANT.ListenerType:
					// 	counsellor, _, _ = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"listener_id": appoint["counsellor_id"]})
					// case CONSTANT.TherapistType:
					// 	counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "timezone", "phone", "email"}, map[string]string{"therapist_id": appoint["counsellor_id"]})

					// }

					DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable,
						map[string]string{
							"appointment_id": appoint["appointment_id"],
						},
						map[string]string{
							"status":      CONSTANT.AppointmentCounsellorCancelled,
							"modified_at": UTIL.GetCurrentTime().String(),
						},
					)

					client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "timezone", "email", "phone"}, map[string]string{"client_id": appoint["client_id"]})

					// remove all previous notifications
					UTIL.RemoveNotification(appoint["appointment_id"], appoint["client_id"])

					// remove all previous message for client
					UTIL.RemoveMessage(appoint["appointment_id"], client[0]["phone"])

					// remove all previous message for therpist
					UTIL.RemoveMessage(appoint["appointment_id"], counsellor[0]["phone"])

					// remove all previous notifications
					UTIL.RemoveNotification(appoint["appointment_id"], appoint["counsellor_id"])

					UTIL.SendNotification(
						CONSTANT.ClientInPersonAppointmentCancellationHeading,
						UTIL.ReplaceNotificationContentInString(
							CONSTANT.TherapistInPersonAppointmentCancellationClientContent,
							map[string]string{
								"###therapistName###": counsellor[0]["first_name"],
								"###timedate###":      UTIL.BuildDateTime(appoint["date"], appoint["time"]).Format(CONSTANT.ReadbleDateTimeFormat),
								"###location###":      appoint["counselling_room"] + ", " + appoint["counselling_address"],
							},
						),
						appoint["client_id"],
						CONSTANT.ClientType,
						UTIL.GetCurrentTime().String(),
						CONSTANT.NotificationSent,
						appoint["appointment_id"],
					)

					// therapist Cancel the Appointment to send text message to client
					UTIL.SendMessage(
						UTIL.ReplaceNotificationContentInString(
							CONSTANT.TherapistInPersonAppointmentCancellationClientTextMeassage,
							map[string]string{
								"###date###":          appoint["date"],
								"###therapistName###": counsellor[0]["first_name"],
							},
						),
						CONSTANT.TransactionalRouteTextMessage,
						client[0]["phone"],
						UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
						appoint["appointment_id"],
						CONSTANT.InstantSendTextMessage,
					)

					// send client email body
					emaildata1 := Model.EmailBodyMessageModel{
						Name: client[0]["first_name"],
						Message: UTIL.ReplaceNotificationContentInString(
							CONSTANT.TherapistInPersonAppointmentCancellationClientBody,
							map[string]string{
								"###therapist###": counsellor[0]["first_name"],
								"###date###":      appoint["date"],
								"###time###":      UTIL.GetTimeFromTimeSlotIN12Hour(appoint["time"]),
								"###location###":  appoint["counselling_room"] + ", " + appoint["counselling_address"],
							},
						),
					}

					emailBody1 := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata1, filepath_text)
					// email for client
					UTIL.SendEmail(
						CONSTANT.TherapistInPersonAppointmentCancellationClientTitle,
						emailBody1,
						client[0]["email"],
						CONSTANT.InstantSendEmailMessage,
					)
				}
			}

			// send appointment cancel notification to counsellor
			UTIL.SendNotification(
				CONSTANT.ClientInPersonAppointmentCancellationHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistInPersonAppointmentCancellationTherapistContent,
					map[string]string{
						"###date###":     body["date"],
						"###location###": body["roomNo"] + ", " + body["address"],
					},
				),
				appointment[0]["counsellor_id"],
				"1",
				UTIL.GetCurrentTime().String(),
				CONSTANT.NotificationSent,
				body["id"],
			)

			// therapist Cancel the Appointment to send text message to therpaist
			UTIL.SendMessage(
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistInPersonAppointmentCancellationTherapistTextMeassage,
					map[string]string{
						"###date###":     body["date"],
						"###location###": body["roomNo"] + ", " + body["address"],
					},
				),
				CONSTANT.TransactionalRouteTextMessage,
				counsellor[0]["phone"],
				UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
				body["id"],
				CONSTANT.InstantSendTextMessage,
			)

			// send email for therapist
			emaildata := Model.EmailBodyMessageModel{
				Name: counsellor[0]["first_name"],
				Message: UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistInPersonAppointmentCancellationTherapistBody,
					map[string]string{
						"###date###":     body["date"],
						"###location###": body["roomNo"] + ", " + body["address"],
					},
				),
			}

			emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
			// email for therapist
			UTIL.SendEmail(
				CONSTANT.ClientInPersonAppointmentCancellationTherapistTitle,
				emailBody,
				counsellor[0]["email"],
				CONSTANT.InstantSendEmailMessage,
			)
		}

	} else {
		// newly added schedule
		DB.InsertSQL(CONSTANT.InPersonSLotsScheduleTable, slotsDate)
	}

	for key, val := range slots {

		DB.ExecuteSQL("update "+CONSTANT.InPersonSLotsTable+" set `"+key+"` = "+val+" where counsellor_id = ? and date = ? and  `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", body["counsellor_id"], body["date"]) // dont update already booked slots

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorConnectWithCorporateGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " partner_name = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "location":
			if len(val[0]) > 0 {
				wheres = append(wheres, " partner_location = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "connect_id":
			wheres = append(wheres, " connect_id = ? ")
			queryArgs = append(queryArgs, val[0])

		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	inPersonConnect, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonCounsellorConnectWithCorporateTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of contents
	inPersonConnectCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InPersonCounsellorConnectWithCorporateTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["in_person_connect"] = inPersonConnect
	response["in_person_connect_count"] = inPersonConnectCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(inPersonConnectCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorConnectWithCorporateAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CorporateCounsellorAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	corporateCounsellor := map[string]string{}
	corporateCounsellor["counsellor_id"] = body["counsellor_id"]
	corporateCounsellor["partner_name"] = body["partner_name"]
	corporateCounsellor["partner_location"] = body["partner_location"]
	corporateCounsellor["status"] = CONSTANT.ContentActive
	corporateCounsellor["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.InPersonCounsellorConnectWithCorporateTable, CONSTANT.ContentDigits, corporateCounsellor, "connect_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorConnectWithCorporateUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	corporateCounsellor := map[string]string{}
	corporateCounsellor["counsellor_id"] = body["counsellor_id"]
	corporateCounsellor["partner_name"] = body["partner_name"]
	corporateCounsellor["partner_location"] = body["partner_location"]
	corporateCounsellor["status"] = body["status"]
	corporateCounsellor["created_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, map[string]string{"connect_id": r.FormValue("connect_id")}, corporateCounsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
