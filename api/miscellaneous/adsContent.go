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

func GetCounsellorClientRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get last cleints
	lastClient, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"session_date", "mental_health"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["last_client"] = lastClient
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
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

	var noshow, mentalHealth, sendStatus string

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
	// add counsellorRecord details
	counsellorRecord := map[string]string{}
	counsellorRecord["counsellor_id"] = body["counsellor_id"]
	counsellorRecord["client_id"] = body["client_id"]
	counsellorRecord["client_first_name"] = body["client_first_name"]
	counsellorRecord["client_last_name"] = body["client_last_name"]
	counsellorRecord["client_gender"] = body["client_gender"]
	counsellorRecord["client_age"] = body["client_age"]
	counsellorRecord["client_department"] = body["client_department"]
	counsellorRecord["client_location"] = body["client_location"]
	counsellorRecord["noshow"] = noshow
	counsellorRecord["session_mode"] = body["session_mode"]
	counsellorRecord["session_no"] = body["session_no"]
	counsellorRecord["session_date"] = body["session_date"]
	counsellorRecord["in_time"] = body["in_time"]
	counsellorRecord["out_time"] = body["out_time"]
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

	if noshow == "1" {
		noshow = "Yes"
	} else {
		noshow = "No"
	}

	if len(body["client_id"]) != 0 && len(body["client_notes"]) != 0 {
		sendStatus = "Send"
	} else {
		sendStatus = "Not Send"
	}

	data := Model.EmailDataForCounsellorRecord{
		TherapistName:   counsellor[0]["first_name"] + " " + counsellor[0]["last_name"],
		First_Name:      body["client_first_name"],
		Last_Name:       body["client_last_name"],
		Gender:          body["client_gender"],
		Age:             body["client_age"],
		Department:      body["client_department"],
		Location:        body["client_location"],
		NoShow:          noshow,
		SessionMode:     body["session_mode"],
		SessionNo:       body["session_no"],
		SessionDate:     body["session_date"],
		InTime:          body["in_time"],
		OutTime:         body["out_time"],
		MentalHealth:    mentalHealth,
		TherapeuticGoal: body["therapeutic_goal"],
		TherapyPlan:     body["therapy_plan"],
		AssessmentTool:  body["assessment_tool"],
		ClientNotes:     body["client_notes"],
		ClientAttach:    body["client_documents"],
		SendingStatus:   sendStatus,
	}

	filepath := "htmlfile/CounsellorRecord.html"

	emailbody := UTIL.GetHTMLTemplateForCounsellorRecord(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorRecordForClientTitle,
		emailbody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	if mentalStatus > 6 {
		UTIL.SendEmail(
			CONSTANT.CounselloRecordClientEmergencyCaseTitle,
			emailbody,
			"shambhavi.alve@clovemind.com",
			CONSTANT.InstantSendEmailMessage,
		)
	}

	if len(body["client_id"]) != 0 && len(body["client_notes"]) != 0 {
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

		emaildata := Model.EmailBodyMessageModelWithDocu{
			Name:          client[0]["first_name"],
			TherapistName: counsellor[0]["first_name"],
			Date:          UTIL.BuildOnlyDate(body["session_date"]),
			Message:       body["client_notes"],
			Message1:      body["links"],
		}

		filepath_text := "htmlfile/emailbodywithlink.html"

		emailBy := UTIL.GetHTMLTemplateForWithDocument(emaildata, filepath_text)

		UTIL.SendEmailWithDocument(client[0]["email"], emailBy, CONSTANT.CounsellorDocumentForClientTitle, listofDocuments)
	}

	// UTIL.SendEmail(
	// 	CONSTANT.CounsellorRecordForClientTitle,
	// 	emailbody,
	// 	client[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
