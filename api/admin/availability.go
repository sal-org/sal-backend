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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []any{}
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

	counsellorIDs := UTIL.ExtractValuesFromArrayMap(availabilitys, "counsellor_id")

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
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
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	body := Model.AvailabilityUpdateRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	slotsDate := map[string]string{}

	slotsDate["counsellor_id"] = body.CounsellorID
	slotsDate["date"] = body.Date
	slotsDate["from_time"] = body.FromTime
	slotsDate["to_time"] = body.ToTime
	slotsDate["partner_name"] = body.CompanyName
	slotsDate["partner_location"] = body.CompanyLocation
	slotsDate["counselling_room"] = body.RoomNo
	slotsDate["counselling_address"] = body.Address
	slotsDate["schedules_status"] = CONSTANT.InPersonSlotsInProgress
	slotsDate["status"] = body.Status
	slotsDate["0"] = body.Zero
	slotsDate["1"] = body.One
	slotsDate["2"] = body.Two
	slotsDate["3"] = body.Three
	slotsDate["4"] = body.Four
	slotsDate["5"] = body.Five
	slotsDate["6"] = body.Six
	slotsDate["7"] = body.Seven
	slotsDate["8"] = body.Eight
	slotsDate["9"] = body.Nine
	slotsDate["10"] = body.Ten
	slotsDate["11"] = body.Eleven
	slotsDate["12"] = body.Twelve
	slotsDate["13"] = body.Thirteen
	slotsDate["14"] = body.Fourteen
	slotsDate["15"] = body.Fifteen
	slotsDate["16"] = body.Sixteen
	slotsDate["17"] = body.Seventeen
	slotsDate["18"] = body.Eighteen
	slotsDate["19"] = body.Nineteen
	slotsDate["20"] = body.Twenty
	slotsDate["21"] = body.TwentyOne
	slotsDate["22"] = body.TwentyTwo
	slotsDate["23"] = body.TwentyThree
	slotsDate["24"] = body.TwentyFour
	slotsDate["25"] = body.TwentyFive
	slotsDate["26"] = body.TwentySix
	slotsDate["27"] = body.TwentySeven
	slotsDate["28"] = body.TwentyEight
	slotsDate["29"] = body.TwentyNine
	slotsDate["30"] = body.Thirteen
	slotsDate["31"] = body.ThirtyOne
	slotsDate["32"] = body.ThirtyTwo
	slotsDate["33"] = body.ThirtyThree
	slotsDate["34"] = body.ThirtyFour
	slotsDate["35"] = body.ThirtyFive
	slotsDate["36"] = body.ThirtySix
	slotsDate["37"] = body.ThirtySeven
	slotsDate["38"] = body.ThirtyEight
	slotsDate["39"] = body.ThirtyNine
	slotsDate["40"] = body.Fourteen
	slotsDate["41"] = body.FortyOne
	slotsDate["42"] = body.FortyTwo
	slotsDate["43"] = body.FortyThree
	slotsDate["44"] = body.FortyFour
	slotsDate["45"] = body.FortyFive
	slotsDate["46"] = body.FortySix
	slotsDate["47"] = body.FortySeven

	slots := map[string]string{}

	slots["available"] = body.Status
	slots["0"] = body.Zero
	slots["1"] = body.One
	slots["2"] = body.Two
	slots["3"] = body.Three
	slots["4"] = body.Four
	slots["5"] = body.Five
	slots["6"] = body.Six
	slots["7"] = body.Seven
	slots["8"] = body.Eight
	slots["9"] = body.Nine
	slots["10"] = body.Ten
	slots["11"] = body.Eleven
	slots["12"] = body.Twelve
	slots["13"] = body.Thirteen
	slots["14"] = body.Fourteen
	slots["15"] = body.Fifteen
	slots["16"] = body.Sixteen
	slots["17"] = body.Seventeen
	slots["18"] = body.Eighteen
	slots["19"] = body.Nineteen
	slots["20"] = body.Twenty
	slots["21"] = body.TwentyOne
	slots["22"] = body.TwentyTwo
	slots["23"] = body.TwentyThree
	slots["24"] = body.TwentyFour
	slots["25"] = body.TwentyFive
	slots["26"] = body.TwentySix
	slots["27"] = body.TwentySeven
	slots["28"] = body.TwentyEight
	slots["29"] = body.TwentyNine
	slots["30"] = body.Thirty
	slots["31"] = body.ThirtyOne
	slots["32"] = body.ThirtyTwo
	slots["33"] = body.ThirtyThree
	slots["34"] = body.ThirtyFour
	slots["35"] = body.ThirtyFive
	slots["36"] = body.ThirtySix
	slots["37"] = body.ThirtySeven
	slots["38"] = body.ThirtyEight
	slots["39"] = body.ThirtyNine
	slots["40"] = body.Forty
	slots["41"] = body.FortyOne
	slots["42"] = body.FortyTwo
	slots["43"] = body.FortyThree
	slots["44"] = body.FortyFour
	slots["45"] = body.FortyFive
	slots["46"] = body.FortySix
	slots["47"] = body.FortySeven

	if len(body.ID) > 0 {

		inpersonSlotList, status, ok := DB.SelectSQL(CONSTANT.InPersonSLotsScheduleTable, []string{"*"}, map[string]string{"id": body.ID})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		DB.UpdateSQL(CONSTANT.InPersonSLotsScheduleTable, map[string]string{"id": body.ID}, slotsDate)

		if body.Status == "0" {
			appointment, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": body.CounsellorID, "date": body.Date})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellor details
			counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": body.CounsellorID})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			if len(counsellor) == 0 {
				counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body.CounsellorID})
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
						"",
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
						"###date###":     body.Date,
						"###location###": body.RoomNo + ", " + body.Address,
					},
				),
				appointment[0]["counsellor_id"],
				"1",
				UTIL.GetCurrentTime().String(),
				CONSTANT.NotificationSent,
				body.ID,
				"",
			)

			// therapist Cancel the Appointment to send text message to therpaist
			UTIL.SendMessage(
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistInPersonAppointmentCancellationTherapistTextMeassage,
					map[string]string{
						"###date###":     body.Date,
						"###location###": body.RoomNo + ", " + body.Address,
					},
				),
				CONSTANT.TransactionalRouteTextMessage,
				counsellor[0]["phone"],
				UTIL.BuildDateTime(appointment[0]["date"], appointment[0]["time"]).UTC().String(),
				body.ID,
				CONSTANT.InstantSendTextMessage,
			)

			// send email for therapist
			emaildata := Model.EmailBodyMessageModel{
				Name: counsellor[0]["first_name"],
				Message: UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistInPersonAppointmentCancellationTherapistBody,
					map[string]string{
						"###date###":     body.Date,
						"###location###": body.RoomNo + ", " + body.Address,
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

		if body.RoomNo != inpersonSlotList[0]["roomNo"] || body.Address != inpersonSlotList[0]["address"] {
			inpersonAppointments, status, ok := DB.SelectSQL(CONSTANT.InPersonAppointmentsTable, []string{"*"}, map[string]string{"counsellor_id": body.CounsellorID, "date": body.Date})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			for _, appoint := range inpersonAppointments {
				DB.UpdateSQL(CONSTANT.InPersonAppointmentsTable, map[string]string{"appointment_id": appoint["appointment_id"]}, map[string]string{"counselling_room": body.RoomNo, "counselling_address": body.Address})
			}

		}

		for key, val := range slots {

			DB.ExecuteSQL("update "+CONSTANT.InPersonSLotsTable+" set `"+key+"` = "+val+" where counsellor_id = ? and date = ? and company_name = ? and company_location = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", body.CounsellorID, body.Date, body.CompanyName, body.CompanyLocation) // dont update already booked slots

		}

	} else {
		// newly added schedule
		DB.InsertSQL(CONSTANT.InPersonSLotsScheduleTable, slotsDate)
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorConnectWithCorporateGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []any{}
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

	counsellorIDs := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "counsellor_id")

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
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
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["in_person_connect_count"] = inPersonConnectCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(inPersonConnectCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetCounsellorName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where status = '1' and corporate_therpist != 0 ) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where status = '1' and corporate_therpist != 0 )")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorConnectWithCorporateAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body := Model.InPersonCounsellorConnectWithCorporateAddRequest{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	corporateCounsellor := map[string]string{}
	corporateCounsellor["counsellor_id"] = body.CounsellorID
	corporateCounsellor["partner_name"] = body.PartnerName
	corporateCounsellor["partner_location"] = body.PartnerLocation
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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	connectID, ok := UTIL.Required(r.FormValue("connect_id"), "Connect ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, connectID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body := Model.InPersonCounsellorConnectWithCorporateUpdateRequest{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if connect_id exists
	if !DB.CheckIfExists(CONSTANT.InPersonCounsellorConnectWithCorporateTable, map[string]string{"connect_id": connectID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	corporateCounsellor := map[string]string{}
	corporateCounsellor["counsellor_id"] = body.CounsellorID
	corporateCounsellor["partner_name"] = body.PartnerName
	corporateCounsellor["partner_location"] = body.PartnerLocation
	corporateCounsellor["status"] = body.Status
	corporateCounsellor["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, map[string]string{"connect_id": r.FormValue("connect_id")}, corporateCounsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
