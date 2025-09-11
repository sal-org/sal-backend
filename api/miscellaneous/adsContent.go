package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
	"time"
)

// ListMood godoc
// @Tags Miscellaneous
// @Summary Get all moods
// @Router /adsContent [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func AdsContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	adsContent, status, ok := DB.SelectSQL(CONSTANT.AdsContentTable, []string{"title", "target", "image"}, map[string]string{"status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["ads"] = adsContent
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetDocumentList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	getDocuments, status, ok := DB.SelectSQL(CONSTANT.CounsellorDocumentListTable, []string{"*"}, map[string]string{"status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["documents"] = getDocuments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetCounsellorRecordFromMainCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	category, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordFormCategoryTable + " where status = 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	emotionalState, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordFormEmotionalState + " where status = 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["main_category"] = category
	response["emotional_state"] = emotionalState
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func GetCounsellorRecordFromSubCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	subCategory, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordFormSubCategoryTable+" where status = 1 and category_id = ? ", r.FormValue("category_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["sub_category"] = subCategory
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func GetCounsellorClientRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var lastClient []map[string]string

	var status string
	var ok bool

	if len(r.FormValue("counsellor_id")) != 0 {
		lastClient, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ? and counsellor_id = ?  order by session_date desc", r.FormValue("client_id"), r.FormValue("counsellor_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		lastClient, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc limit 5", r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	response["lastest_record"] = lastClient
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CheckCounsellorClientRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	if len(r.FormValue("appointment_id")) != 0 {
		// get last cleints
		lastClient, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(lastClient) == 0 {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
			return
		}

	} else {
		// get last cleints
		lastClient, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "client_id": r.FormValue("client_id"), "session_date": r.FormValue("date")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(lastClient) == 0 {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
			return
		}
	}

}

func CounsellorClientRecord(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorRecordAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	var noshow, psychiatricIntervention, mentalHealth, sendStatus, sessionFor, sessionType string

	if len(body["mental_health"]) > 0 {
		mentalHealth = body["mental_health"]
	} else {
		mentalHealth = "0"
	}

	mentalStatus, _ := strconv.Atoi(mentalHealth)

	if len(body["noshow"]) > 0 {
		noshow = body["noshow"]
	} else {
		noshow = "0"
	}

	if len(body["session_for"]) > 0 {
		sessionFor = body["session_for"]
		sessionType = body["session_type"]
	} else {
		sessionFor = ""
		sessionType = ""
	}

	// add counsellorRecord details
	counsellorRecord := map[string]string{}
	counsellorRecord["counsellor_id"] = body["counsellor_id"]
	counsellorRecord["client_id"] = body["client_id"]
	counsellorRecord["appointment_id"] = body["appointment_id"]
	counsellorRecord["client_first_name"] = body["client_first_name"]
	counsellorRecord["client_last_name"] = body["client_last_name"]
	counsellorRecord["client_gender"] = body["client_gender"]
	counsellorRecord["client_age"] = body["client_age"]
	counsellorRecord["client_department"] = body["client_department"]
	counsellorRecord["client_location"] = body["client_location"]
	counsellorRecord["session_for"] = sessionFor
	counsellorRecord["session_type"] = sessionType
	counsellorRecord["noshow"] = noshow
	counsellorRecord["session_mode"] = body["session_mode"]
	counsellorRecord["session_date"] = body["session_date"]
	counsellorRecord["in_time"] = body["in_time"]
	counsellorRecord["out_time"] = body["out_time"]
	counsellorRecord["presenting_concerns"] = body["presenting_concerns"]
	counsellorRecord["psychiatric_intervention"] = body["psychiatric_intervention"]
	counsellorRecord["psychiatric_intervention_reason"] = body["psychiatric_intervention_reason"]
	counsellorRecord["therapy_notes"] = body["therapy_notes"]
	counsellorRecord["sub_category"] = body["sub_category"]
	counsellorRecord["emotional_state"] = body["emotional_state"]
	counsellorRecord["next_follow_date"] = body["next_follow_date"]
	counsellorRecord["mental_health"] = mentalHealth
	counsellorRecord["client_notes"] = body["client_notes"]
	counsellorRecord["client_documents"] = body["client_documents"]
	counsellorRecord["links"] = body["links"]
	counsellorRecord["therapeutic_goal"] = body["therapeutic_goal"]
	counsellorRecord["therapy_plan"] = body["therapy_plan"]
	counsellorRecord["assessment_tool"] = body["assessment_tool"]
	counsellorRecord["created_at"] = UTIL.GetCurrentTime().String()
	counsellorRecord["status"] = CONSTANT.CounsellorActive

	_, status, ok := DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsTable, CONSTANT.CounsellorRecordDigits, counsellorRecord, "record_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "email"}, map[string]string{"counsellor_id": body["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "email"}, map[string]string{"therapist_id": body["counsellor_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	sessionDate := UTIL.BuildOnlyDate(body["session_date"])

	nextFollowDate := ""

	if len(body["next_follow_date"]) != 0 {
		nextFollowDate = UTIL.BuildOnlyDate(body["next_follow_date"])
	}

	if noshow == "1" {
		noshow = "Yes"
	} else {
		noshow = "No"
	}

	if body["psychiatric_intervention"] == "1" {
		psychiatricIntervention = "Yes"
	} else {
		psychiatricIntervention = "No"
	}

	if len(body["client_id"]) != 0 && len(body["client_notes"]) != 0 {
		sendStatus = "Send"
	} else {
		sendStatus = "Not Send"
	}

	data := Model.EmailDataForCounsellorRecord{
		TherapistName:                 counsellor[0]["first_name"] + " " + counsellor[0]["last_name"],
		SessionFor:                    body["session_for"],
		First_Name:                    body["client_first_name"],
		Last_Name:                     body["client_last_name"],
		Gender:                        body["client_gender"],
		Age:                           body["client_age"],
		NoShow:                        noshow,
		PresentingConcerns:            body["presenting_concerns"],
		PsychiatricIntervention:       psychiatricIntervention,
		PsychiatricInterventionReason: body["psychiatric_intervention_reason"],
		TherapyNotes:                  body["therapy_notes"],
		SubCategory:                   body["sub_category"],
		EmotionalState:                body["emotional_state"],
		NextFollowDate:                nextFollowDate,
		SessionMode:                   body["session_mode"],
		SessionDate:                   sessionDate,
		InTime:                        body["in_time"],
		OutTime:                       body["out_time"],
		MentalHealth:                  mentalHealth,
		TherapeuticGoal:               body["therapeutic_goal"],
		TherapyPlan:                   body["therapy_plan"],
		AssessmentTool:                body["assessment_tool"],
		ClientNotes:                   body["client_notes"],
		ClientAttach:                  body["client_documents"],
		SendingStatus:                 sendStatus,
	}

	filepath := "htmlfile/CounsellorRecord.html"

	emailbody := UTIL.GetHTMLTemplateForCounsellorRecord(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorRecordForClientTitle,
		emailbody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	if mentalStatus > 7 || body["psychiatric_intervention"] == "1" {

		UTIL.SendEmailWithCCBB(
			CONSTANT.CounselloRecordClientEmergencyCaseTitle,
			emailbody,
			CONFIG.QCEmailID2,
			CONFIG.QCEmailID1,
			CONSTANT.InstantSendEmailMessage,
		)
	}

	message, message1, message2, subjectLine := "", "", "", ""

	if len(body["next_follow_date"]) != 0 {
		// 15 min push notification before appointment start

		if body["session_mode"] == "In-Person" {
			UTIL.SendNotification(
				CONSTANT.ClientInPersonAppointmentFollowUpSessionToSuggestedReminderClientHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.ClientAppointmentFollowUpToSuggestedRemiderClientContent,
					map[string]string{
						"###therapist_name###": counsellor[0]["first_name"],
						"###follow_up_date###": nextFollowDate,
					},
				),
				body["client_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body["next_follow_date"], "26").Add(-24*time.Hour).UTC().String(),
				CONSTANT.NotificationInProgress,
				body["client_id"],
				"",
			)

			subjectLine = CONSTANT.CounsellorInPersonAppointmentDocumentForClientTitle
		} else {
			UTIL.SendNotification(
				CONSTANT.ClientVirtualAppointmentFollowUpSessionToSuggestedReminderClientHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.ClientAppointmentFollowUpToSuggestedRemiderClientContent,
					map[string]string{
						"###therapist_name###": counsellor[0]["first_name"],
						"###follow_up_date###": nextFollowDate,
					},
				),
				body["client_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body["next_follow_date"], "26").Add(-24*time.Hour).UTC().String(),
				CONSTANT.NotificationInProgress,
				body["client_id"],
				"",
			)

			subjectLine = CONSTANT.CounsellorVirtualAppointmentDocumentForClientTitle
		}

		message1 = UTIL.ReplaceNotificationContentInString(
			CONSTANT.TherapistAttachDocumentsFollowDateFooterClientBody,
			map[string]string{
				"###followupdate###": nextFollowDate,
			},
		)
	}

	if len(body["client_id"]) != 0 {

		if len(body["client_documents"]) != 0 || len(body["client_notes"]) != 0 || len(body["assessment_tool"]) != 0 {

			client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "last_name", "email"}, map[string]string{"client_id": body["client_id"]})

			var listofDocuments []Model.DocumentList
			if len(body["client_documents"]) != 0 {
				clientDocuments := strings.Split(body["client_documents"], ",")

				for _, value := range clientDocuments {
					document, _, _ := DB.SelectSQL(CONSTANT.CounsellorDocumentListTable, []string{"document"}, map[string]string{"document_name": value})
					docu := Model.DocumentList{
						DocumentName: value,
						DocumentLink: document[0]["document"],
					}
					listofDocuments = append(listofDocuments, docu)
				}
			}

			// re := regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

			// message := re.ReplaceAllString(body["client_notes"], `<a href="$0">$0</a>`)

			var emaildata Model.EmailBodyMessageModelWithDocu

			if len(body["client_documents"]) != 0 {
				message = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachDocumentsWithClientBody,
					map[string]string{
						"###TherapistName###": counsellor[0]["first_name"],
						"###Date###":          sessionDate,
					},
				)

			} else {
				message = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachDocumentsWithOutClientBody,
					map[string]string{
						"###TherapistName###": counsellor[0]["first_name"],
						"###Date###":          sessionDate,
					},
				)
			}

			if len(body["assessment_tool"]) != 0 {
				message2 = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachAssessmentWithClientBody,
					map[string]string{
						"###assessment_name###": body["assessment_tool"],
					},
				)
			}

			if len(body["client_notes"]) != 0 {
				emaildata = Model.EmailBodyMessageModelWithDocu{
					Name:     client[0]["first_name"],
					Message:  message,
					Message1: "Your therapist has suggested the following guidelines:",
					Message2: body["client_notes"],
					Message3: message2,
					Message4: message1,
				}
			} else {
				emaildata = Model.EmailBodyMessageModelWithDocu{
					Name:     client[0]["first_name"],
					Message:  message,
					Message1: "Your therapist has suggested the following guidelines:",
					Message2: body["client_notes"],
					Message3: body["links"],
					Message4: message1,
				}
			}

			filepath_text := "htmlfile/emailbodywithassessment.html"

			emailBy := UTIL.GetHTMLTemplateForWithDocument(emaildata, filepath_text)

			UTIL.SendEmailWithDocument(client[0]["email"], emailBy, UTIL.ReplaceNotificationContentInString(
				subjectLine,
				map[string]string{
					"###date###": sessionDate,
				},
			), listofDocuments)
		}

	}

	// UTIL.SendEmail(
	// 	CONSTANT.CounsellorRecordForClientTitle,
	// 	emailbody,
	// 	client[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
