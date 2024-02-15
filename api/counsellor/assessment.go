package counsellor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

// AssessmentsList godoc
// @Tags Counsellor Assessment
// @Summary List available assessments
// @Router /counsellor/assessments [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get all available assessments
	assessments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AssessmentsTable + " where status = " + CONSTANT.AssessmentActive + " order by `order` asc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get assessment latest result
	assessmentResults, status, ok := DB.SelectProcess("select final_score, assessment_id from "+CONSTANT.AssessmentResultsTable+" where id in (select max(id) from "+CONSTANT.AssessmentResultsTable+" where user_id = ? and status = "+CONSTANT.AssessmentResultActive+" group by assessment_id)", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["assessment_results"] = UTIL.ConvertArrayMapToKeyMapArray(assessmentResults, "assessment_id")
	response["assessments"] = assessments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentDetail godoc
// @Tags Counsellor Assessment
// @Summary Get assessment detail
// @Router /counsellor/assessment [get]
// @Param assessment_id query string true "Assessment ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get assessment questions
	questions, status, ok := DB.SelectProcess("select assessment_question_id, question from "+CONSTANT.AssessmentQuestionsTable+" where assessment_id = ? and status = "+CONSTANT.AssessmentQuestionActive+" order by `order` asc", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["questions"] = questions

	// get assessment question options
	questionOptions, status, ok := DB.SelectProcess("select assessment_question_id, assessment_question_option_id, `option`, `score` from "+CONSTANT.AssessmentQuestionOptionsTable+" where assessment_id = ? and status = "+CONSTANT.AssessmentQuestionOptionActive+" order by `order` asc", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["questions"] = questions

	// get assessment scores
	scores, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? order by `min` asc", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["questions"] = questions
	response["question_options"] = UTIL.ConvertArrayMapToKeyMapArray(questionOptions, "assessment_question_id")
	response["scores"] = scores
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentAdd godoc
// @Tags Counsellor Assessment
// @Summary Add counsellor assessment
// @Router /counsellor/assessment [post]
// @Param body body model.AssessmentAddRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body := MODEL.AssessmentAddRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(b, &body)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	finalScore := 0
	for _, detail := range body.Details {
		score, _ := strconv.Atoi(detail.Score)
		finalScore += score
	}

	// add assessment result
	assessmentResultID, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentResultsTable, CONSTANT.AssessmentResultsDigits, map[string]string{
		"user_id":       body.UserID,
		"name":          body.Name,
		"age":           body.Age,
		"gender":        body.Gender,
		"phone":         body.Phone,
		"assessment_id": body.AssessmentID,
		"feedback":      body.Feedback,
		"final_score":   strconv.Itoa(finalScore),
		"status":        CONSTANT.AssessmentResultActive,
		"created_at":    UTIL.GetCurrentTime().UTC().String(),
	}, "assessment_result_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, detail := range body.Details {
		DB.InsertSQL(CONSTANT.AssessmentResultDetailsTable, map[string]string{
			"assessment_result_id":          assessmentResultID,
			"assessment_question_id":        detail.AssessmentQuestionID,
			"assessment_question_option_id": detail.AssessmentQuestionOptionID,
			"score":                         detail.Score,
			"status":                        CONSTANT.AssessmentResultActive,
			"created_at":                    UTIL.GetCurrentTime().UTC().String(),
		})
	}

	response["result"] = DB.QueryRowSQL("select result from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? and min <= "+strconv.Itoa(finalScore)+" and max >= "+strconv.Itoa(finalScore), body.AssessmentID)
	response["assessment_result_id"] = assessmentResultID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentHistory godoc
// @Tags Counsellor Assessment
// @Summary GET counsellor assessment
// @Router /counsellor/assessment/history [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var results []string

	// get assessment past results
	assessmentResults, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentResultsTable+" where user_id = ?  and status = "+CONSTANT.AssessmentResultActive+" order by created_at desc", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// add swagger tag when you get assessment deatils  @Param assessment_id query string true "Assessment ID to get details"
	// get assessment details
	/*assessment, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentsTable+" where assessment_id = ?", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get assessments questions
	assessmentQuestions, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentQuestionsTable+" where assessment_id = ? order by `order`", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get assessments options
	assessmentOptions, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentQuestionOptionsTable+" where assessment_id = ? order by `order`", r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}*/

	assessmentResultIDs := UTIL.ExtractValuesFromArrayMap(assessmentResults, "assessment_result_id")

	assessmentIDs := UTIL.ExtractValuesFromArrayMap(assessmentResults, "assessment_id")

	assessmentDetails, status, ok := DB.SelectProcess("select assessment_id , title from " + CONSTANT.AssessmentsTable + " where assessment_id in ('" + strings.Join(assessmentIDs, "','") + "') order by created_at desc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get assessments result details
	/*assessmentResultDetails, status, ok := DB.SelectProcess("select * from " + CONSTANT.AssessmentResultDetailsTable + " where assessment_result_id in ('" + strings.Join(assessmentResultIDs, "','") + "') order by created_at desc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}*/

	for i, _ := range assessmentResultIDs {
		result := DB.QueryRowSQL("select result from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? and min <= ? and max >=  ? ", assessmentResults[i]["assessment_id"], assessmentResults[i]["final_score"], assessmentResults[i]["final_score"])

		results = append(results, result)
	}

	response["assessment_results"] = assessmentResults
	//response["assessment_result_details"] = UTIL.ConvertArrayMapToKeyMapArray(assessmentResultDetails, "assessment_result_id")
	response["result"] = results
	response["assessment"] = UTIL.ConvertArrayMapToKeyMapArray(assessmentDetails, "assessment_id")
	//response["assessment"] = assessment
	//response["assessment_questions"] = assessmentQuestions
	//response["assessment_options"] = UTIL.ConvertArrayMapToKeyMapArray(assessmentOptions, "assessment_question_id")
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentDownload godoc
// @Tags Counsellor Assessment
// @Summary Get assessment download
// @Router /counsellor/assessment/download [get]
// @Param assessment_result_id query string true "Logged in Assessment Result ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	//const assessment_id = "ywlxbz8yrlp942"

	var fileName, emailbody string

	if DB.CheckIfExists(CONSTANT.AssessmentPdfTable, map[string]string{"assessment_result_id": r.FormValue("assessment_result_id")}) {

		receipt, _, _ := DB.SelectSQL(CONSTANT.AssessmentPdfTable, []string{"*"}, map[string]string{"assessment_result_id": r.FormValue("assessment_result_id")})
		fileName = receipt[0]["pdf"]
	} else {

		assessment_result, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentResultsTable+" where assessment_result_id = ? ", r.FormValue("assessment_result_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		assessment_result_details, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentResultDetailsTable+" where assessment_result_id = ? ", r.FormValue("assessment_result_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		finalScore, _ := strconv.Atoi(assessment_result[0]["final_score"])

		//assign := assessment_result[0]["assessment_id"]

		if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp942" {

			var filePath string

			if finalScore >= 0 && finalScore <= 39 {

				filePath = "htmlfile/Assessment_AIS_Low.html"

			} else if finalScore >= 40 && finalScore <= 69 {

				fmt.Println(assessment_result[0]["score"])
				filePath = "htmlfile/Assessment_AIS_Mid.html"

			} else {

				filePath = "htmlfile/Assessment_AIS_High.html"

			}

			assessment_data := MODEL.AssessmentDownloadAIS{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentAIS(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp943" {

			var filePath string

			var question_options []string

			//var question_option map[string]string

			assessmentResultIDs := UTIL.ExtractValuesFromArrayMap(assessment_result_details, "assessment_question_option_id")

			for i := range assessmentResultIDs {

				options := DB.QueryRowSQL("select `option` from "+CONSTANT.AssessmentQuestionOptionsTable+" where assessment_question_option_id = ? ", assessment_result_details[i]["assessment_question_option_id"])
				question_options = append(question_options, options)

			}

			if finalScore >= 0 && finalScore <= 10 {

				filePath = "htmlfile/Assessment_BDI_II_Normal.html"

			} else if finalScore >= 11 && finalScore <= 18 {

				filePath = "htmlfile/Assessment_BDI_II_Mild.html"

			} else if finalScore >= 19 && finalScore <= 25 {

				filePath = "htmlfile/Assessment_BDI_II_Mod.html"

			} else {

				filePath = "htmlfile/Assessment_BDI_II_Severe.html"

			}

			assessment_data := MODEL.AssessmentDownloadBDIModel{
				Name:       assessment_result[0]["name"],
				Date:       UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:        assessment_result[0]["age"],
				Gender:     assessment_result[0]["gender"],
				Score:      assessment_result[0]["final_score"],
				Answer1:    assessment_result_details[0]["score"],
				Answer2:    assessment_result_details[1]["score"],
				Answer3:    assessment_result_details[2]["score"],
				Answer4:    assessment_result_details[3]["score"],
				Answer5:    assessment_result_details[4]["score"],
				Answer6:    assessment_result_details[5]["score"],
				Answer7:    assessment_result_details[6]["score"],
				Answer8:    assessment_result_details[7]["score"],
				Answer9:    assessment_result_details[8]["score"],
				Answer10:   assessment_result_details[9]["score"],
				Answer11:   assessment_result_details[10]["score"],
				Answer12:   assessment_result_details[11]["score"],
				Answer13:   assessment_result_details[12]["score"],
				Answer14:   assessment_result_details[13]["score"],
				Answer15:   assessment_result_details[14]["score"],
				Answer16:   assessment_result_details[15]["score"],
				Answer17:   assessment_result_details[16]["score"],
				Answer18:   assessment_result_details[17]["score"],
				Answer19:   assessment_result_details[18]["score"],
				Answer20:   assessment_result_details[19]["score"],
				Answer21:   assessment_result_details[20]["score"],
				Response1:  question_options[0],
				Response2:  question_options[1],
				Response3:  question_options[2],
				Response4:  question_options[3],
				Response5:  question_options[4],
				Response6:  question_options[5],
				Response7:  question_options[6],
				Response8:  question_options[7],
				Response9:  question_options[8],
				Response10: question_options[9],
				Response11: question_options[10],
				Response12: question_options[11],
				Response13: question_options[12],
				Response14: question_options[13],
				Response15: question_options[14],
				Response16: question_options[15],
				Response17: question_options[16],
				Response18: question_options[17],
				Response19: question_options[18],
				Response20: question_options[19],
				Response21: question_options[20],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentBDI(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}
		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp944" {

			var filePath string

			if finalScore >= 0 && finalScore <= 13 {

				filePath = "htmlfile/Assessment_PSS_Low.html"

			} else if finalScore >= 14 && finalScore <= 26 {

				filePath = "htmlfile/Assessment_PSS_Mid.html"

			} else {

				filePath = "htmlfile/Assessment_PSS_High.html"

			}

			assessment_data := MODEL.AssessmentDownloadAIS{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentAIS(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp945" {

			var filePath string

			if finalScore >= 0 && finalScore <= 49 {

				filePath = "htmlfile/Assessment_SRS_Low.html"

			} else if finalScore >= 50 && finalScore <= 69 {

				filePath = "htmlfile/Assessment_SRS_Lower_Middle.html"

			} else if finalScore >= 70 && finalScore <= 89 {

				filePath = "htmlfile/Assessment_SRS_Upper_Middle.html"

			} else {

				filePath = "htmlfile/Assessment_SRS_High.html"

			}

			assessment_data := MODEL.AssessmentDownloadSRSModel{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
				Answer11: assessment_result_details[10]["score"],
				Answer12: assessment_result_details[11]["score"],
				Answer13: assessment_result_details[12]["score"],
				Answer14: assessment_result_details[13]["score"],
				Answer15: assessment_result_details[14]["score"],
				Answer16: assessment_result_details[15]["score"],
				Answer17: assessment_result_details[16]["score"],
				Answer18: assessment_result_details[17]["score"],
				Answer19: assessment_result_details[18]["score"],
				Answer20: assessment_result_details[19]["score"],
				Answer21: assessment_result_details[20]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentSRS(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}
		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp947" {
			var filePath string

			if finalScore >= 0 && finalScore <= 20 {

				filePath = "htmlfile/GWB_Low.html"

			} else if finalScore >= 21 && finalScore <= 23 {

				filePath = "htmlfile/GWB_Mid.html"
			} else {
				filePath = "htmlfile/GWB_High.html"
			}

			assessment_data := MODEL.AssessmentDownloadGWBModel{
				Name:    assessment_result[0]["name"],
				Date:    UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:     assessment_result[0]["age"],
				Gender:  assessment_result[0]["gender"],
				Score:   assessment_result[0]["final_score"],
				Answer1: assessment_result_details[0]["score"],
				Answer2: assessment_result_details[1]["score"],
				Answer3: assessment_result_details[2]["score"],
				Answer4: assessment_result_details[3]["score"],
				Answer5: assessment_result_details[4]["score"],
				Answer6: assessment_result_details[5]["score"],
			}

			fmt.Println(assessment_data)

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentGWB(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}

			fmt.Println(emailbody)

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp948" {
			var filePath string

			if finalScore >= 12 && finalScore <= 17 {

				filePath = "htmlfile/CaregiverBurnoutIndex12_17.html"

			} else if finalScore >= 18 && finalScore <= 29 {

				filePath = "htmlfile/CaregiverBurnoutIndex18_29.html"

			} else if finalScore >= 30 && finalScore <= 41 {

				filePath = "htmlfile/CaregiverBurnoutIndex30_41.html"

			} else {

				filePath = "htmlfile/CaregiverBurnoutIndex42_48.html"

			}

			assessment_data := MODEL.AssessmentDownloadBurnOutModel{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
				Answer11: assessment_result_details[10]["score"],
				Answer12: assessment_result_details[11]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentBurnOut(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}
		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp949" {
			var filePath string

			if finalScore >= 10 && finalScore <= 25 {

				filePath = "htmlfile/SelfEsteemLow.html"

			} else if finalScore >= 26 && finalScore <= 29 {

				filePath = "htmlfile/SelfEsteemMedium.html"

			} else {

				filePath = "htmlfile/SelfEsteemHigh.html"

			}

			assessment_data := MODEL.AssessmentDownloadSelfEsteemModel{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentSelfEsteem(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}
		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp950" {

			var filePath string

			if finalScore >= 0 && finalScore <= 24 {

				filePath = "htmlfile/RelationshipMangLow.html"

			} else if finalScore >= 25 && finalScore <= 34 {

				filePath = "htmlfile/RelationshipMangAvg.html"

			} else {

				filePath = "htmlfile/RelationshipMangHigh.html"

			}

			assessment_data := MODEL.AssessmentDownloadSelfEsteemModel{
				Name:     assessment_result[0]["name"],
				Date:     UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:      assessment_result[0]["age"],
				Gender:   assessment_result[0]["gender"],
				Score:    assessment_result[0]["final_score"],
				Answer1:  assessment_result_details[0]["score"],
				Answer2:  assessment_result_details[1]["score"],
				Answer3:  assessment_result_details[2]["score"],
				Answer4:  assessment_result_details[3]["score"],
				Answer5:  assessment_result_details[4]["score"],
				Answer6:  assessment_result_details[5]["score"],
				Answer7:  assessment_result_details[6]["score"],
				Answer8:  assessment_result_details[7]["score"],
				Answer9:  assessment_result_details[8]["score"],
				Answer10: assessment_result_details[9]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentSelfEsteem(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}

		} else {

			var filePath string

			if finalScore >= 0 && finalScore <= 4 {

				filePath = "htmlfile/Assessment_GAD7_Minimal.html"

			} else if finalScore >= 5 && finalScore <= 9 {

				filePath = "htmlfile/Assessment_GAD7_Mild.html"

			} else if finalScore >= 10 && finalScore <= 14 {

				filePath = "htmlfile/Assessment_GAD7_Mod.html"

			} else {

				filePath = "htmlfile/Assessment_GAD7_Severe.html"

			}

			assessment_data := MODEL.AssessmentDownloadGAD7Model{
				Name:    assessment_result[0]["name"],
				Date:    UTIL.BuildDate(assessment_result[0]["created_at"]),
				Age:     assessment_result[0]["age"],
				Gender:  assessment_result[0]["gender"],
				Score:   assessment_result[0]["final_score"],
				Answer1: assessment_result_details[0]["score"],
				Answer2: assessment_result_details[1]["score"],
				Answer3: assessment_result_details[2]["score"],
				Answer4: assessment_result_details[3]["score"],
				Answer5: assessment_result_details[4]["score"],
				Answer6: assessment_result_details[5]["score"],
				Answer7: assessment_result_details[6]["score"],
				Answer8: assessment_result[0]["feedback"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentGAD7(assessment_data, filePath)
			if !ok {
				fmt.Println("html body not create ")
			}

		}

		created, ok := UTIL.GeneratePdfHeaderAndFooterFixted(emailbody, "pdffile/assessment1.pdf") // name created,

		if !ok {
			fmt.Println("Pdf is not created")
		}

		s3Path := "assessment"
		filename := "example1.pdf"

		name, uploaded := UTIL.UploadToS3File(CONFIG.S3Bucket, s3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(filename), CONSTANT.S3PublicRead, created)
		if !uploaded {
			fmt.Println("UploadFile")
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name

		assessment_pdf := map[string]string{}

		assessment_pdf["user_id"] = assessment_result[0]["user_id"]
		assessment_pdf["assessment_result_id"] = assessment_result[0]["assessment_result_id"]
		assessment_pdf["pdf"] = fileName
		assessment_pdf["created_at"] = UTIL.GetCurrentTime().String()

		_, status, ok = DB.InsertWithUniqueID(CONSTANT.AssessmentPdfTable, CONSTANT.ReceiptDigits, assessment_pdf, "assessment_pdf_id")

		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}
	//receipt := map[string]string{}

	response["media_url"] = CONFIG.MediaURL
	response["pdf_name"] = fileName

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
