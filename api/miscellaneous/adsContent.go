package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
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


	// add counsellorRecord details
	counsellorRecord := map[string]string{}
	counsellorRecord["counsellor_id"] = body["counsellor_id"]
	counsellorRecord["client_first_name"] = body["client_first_name"]
	counsellorRecord["client_last_name"] = body["client_last_name"]
	counsellorRecord["client_gender"] = body["client_gender"]
	counsellorRecord["client_age"] = body["client_age"]
	counsellorRecord["client_department"] = body["client_department"]
	counsellorRecord["client_location"] = body["client_location"]
	counsellorRecord["session_mode"] = body["session_mode"]
	counsellorRecord["session_no"] = body["session_no"]
	counsellorRecord["session_date"] = body["session_date"]
	counsellorRecord["in_time"] = body["in_time"]
	counsellorRecord["out_time"] = body["out_time"]
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
	
	
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"email"}, map[string]string{"counsellor_id": body["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"email"}, map[string]string{"therapist_id": body["counsellor_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	data := Model.EmailDataForCounsellorRecord{
		First_Name:  body["client_first_name"],
		Last_Name:   body["client_last_name"],
		Gender:      body["client_gender"],
		Age:        body["client_age"],
		Department:       body["client_department"],
		Location:       body["client_location"],
		SessionMode:       body["session_mode"],
		SessionNo:   body["session_no"],
		SessionDate:  body["session_date"],
		InTime:       body["in_time"],
		OutTime:      body["out_time"],
		TherapeuticGoal: body["therapeutic_goal"],
		TherapyPlan:      body["therapy_plan"],
		AssessmentTool:    body["assessment_tool"],
	}

	filepath := "htmlfile/CounsellorRecord.html"

	emailbody := UTIL.GetHTMLTemplateForCounsellorRecord(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorRecordForClientTitle,
		emailbody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
