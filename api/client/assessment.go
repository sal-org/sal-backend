package client

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"regexp"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	VALIDATOR "salbackend/validator"
	"strconv"
	"strings"

	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

// AssessmentsList godoc
// @Tags Client Assessment
// @Summary List available assessments
// @Router /client/assessments [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	clientID, ok := VALIDATOR.Required(r.FormValue("client_id"), "Client ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	// get all available assessments
	assessments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AssessmentsTable + " where status = " + CONSTANT.AssessmentActive + " order by `order` asc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get assessment latest result
	assessmentResults, status, ok := DB.SelectProcess("select final_score, assessment_id from "+CONSTANT.AssessmentResultsTable+" where id in (select max(id) from "+CONSTANT.AssessmentResultsTable+" where user_id = ? and status = "+CONSTANT.AssessmentResultActive+" group by assessment_id)", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, assessment := range assessments {
		url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, assessment["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
		assessment["photo"] = endPointURL

		// get assessment questions
		questions, status, ok := DB.SelectProcess("select count(assessment_question_id) as question_count from "+CONSTANT.AssessmentQuestionsTable+" where assessment_id = ? and status = "+CONSTANT.AssessmentQuestionActive+" order by `order` asc", assessment["assessment_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		assessment["question_count"] = questions[0]["question_count"]
	}

	response["assessment_results"] = UTIL.ConvertArrayMapToKeyMapArray(assessmentResults, "assessment_id")
	response["assessments"] = assessments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentDetail godoc
// @Tags Client Assessment
// @Summary Get assessment detail
// @Router /client/assessment [get]
// @Param assessment_id query string true "Assessment ID to get details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	assessmentID, ok := VALIDATOR.Required(r.FormValue("assessment_id"), "Assessment ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, assessmentID, CONSTANT.ShowDialog, response)
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
// @Tags Client Assessment
// @Summary Add client assessment
// @Router /client/assessment [post]
// @Param body body model.AssessmentAddRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// read request body
	body := MODEL.AssessmentAddRequest{}
	b, err := io.ReadAll(r.Body)
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

	assessmentName := DB.QueryRowSQL("select title from "+CONSTANT.AssessmentsTable+" where assessment_id = ?", body.AssessmentID)

	finalScore := 0
	finalScoreInFloat := 0.0
	for _, detail := range body.Details {
		score, _ := strconv.Atoi(detail.Score)
		finalScore += score
	}

	assessmentName = strings.ToUpper(assessmentName)

	bol, _ := regexp.MatchString(assessmentName, "GRIT SCALE")

	if bol {
		finalScoreInFloat = float64(finalScore) / 12
		finalScoreInFloat = math.Round(finalScoreInFloat*100) / 100
	} else {
		finalScoreInFloat = math.Round(float64(finalScore)*100) / 100
	}

	// add assessment result
	assessmentResultID, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentResultsTable, CONSTANT.AssessmentResultsDigits, map[string]string{
		"user_id":       body.UserID,
		"name":          body.Name,
		"age":           body.Age,
		"gender":        body.Gender,
		"phone":         body.Phone,
		"assessment_id": body.AssessmentID,
		"final_score":   strconv.FormatFloat(finalScoreInFloat, 'f', 2, 64),
		"feedback":      body.Feedback,
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

	response["result"] = DB.QueryRowSQL("select result from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? and min <= "+strconv.FormatFloat(finalScoreInFloat, 'f', 2, 64)+" and max >= "+strconv.FormatFloat(finalScoreInFloat, 'f', 2, 64), body.AssessmentID)
	response["assessment_result_id"] = assessmentResultID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentHistory godoc
// @Tags Client Assessment
// @Summary Get assessment history
// @Router /client/assessment/history [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	clientID, ok := VALIDATOR.Required(r.FormValue("client_id"), "Client ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	var results []string

	// get assessment past results
	assessmentResults, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentResultsTable+" where user_id = ?  and status = "+CONSTANT.AssessmentResultActive+" order by created_at desc", r.FormValue("client_id"))
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

	for i := range assessmentResultIDs {
		result := DB.QueryRowSQL("select result from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? and min <= ? and max >=  ? ", assessmentResults[i]["assessment_id"], assessmentResults[i]["final_score"], assessmentResults[i]["final_score"])

		results = append(results, result)
	}

	for _, assessment := range assessmentDetails {
		url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, assessment["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
		assessment["photo"] = endPointURL
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
// @Tags Client Assessment
// @Summary Get assessment download
// @Router /client/assessment/download [get]
// @Param assessment_result_id query string true "Logged in Assessment Result ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	//const assessment_id = "ywlxbz8yrlp942"

	var fileName, emailbody string

	assessmentResultID, ok := VALIDATOR.Required(r.FormValue("assessment_result_id"), "Assessment Result ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, assessmentResultID, CONSTANT.ShowDialog, response)
		return
	}

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

		title := DB.QueryRowSQL("select title from "+CONSTANT.AssessmentsTable+" where assessment_id = ? ", assessment_result[0]["assessment_id"])

		// finalScore, _ := strconv.Atoi(assessment_result[0]["final_score"])

		final, _ := strconv.ParseFloat(assessment_result[0]["final_score"], 64)

		finalScore := math.Round(final*100) / 100

		//assign := assessment_result[0]["assessment_id"]

		if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp942" {

			var filePath string

			if finalScore >= 0 && finalScore <= 39 {

				filePath = "htmlfile/Assessment_AIS_Low.html"

			} else if finalScore >= 40 && finalScore <= 69 {

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentGWB(assessment_data, filePath)
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp951" {

			var filePath string

			if finalScore >= 0 && finalScore <= 24 {

				filePath = "htmlfile/EmotionalAwareness0_24.html"

			} else if finalScore >= 25 && finalScore <= 34 {

				filePath = "htmlfile/EmotionalAwareness25_34.html"

			} else {

				filePath = "htmlfile/EmotionalAwareness35_40.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp952" {

			var filePath string

			if finalScore >= 0 && finalScore <= 24 {

				filePath = "htmlfile/Social_Emotional_Awareness_0_24.html"

			} else if finalScore >= 25 && finalScore <= 34 {

				filePath = "htmlfile/Social_Emotional_Awareness_25_34.html"

			} else {

				filePath = "htmlfile/Social_Emotional_Awareness_35_40.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp953" {

			var filePath string

			if finalScore >= 0 && finalScore <= 24 {

				filePath = "htmlfile/Emotional_Management_0_24.html"

			} else if finalScore >= 25 && finalScore <= 34 {

				filePath = "htmlfile/Emotional_Management_25_34.html"

			} else {

				filePath = "htmlfile/Emotional_Management_35_40.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp954" {

			var filePath string

			if finalScore >= 0 && finalScore <= 13 {

				filePath = "htmlfile/LocusOfControlBelow14.html"

			} else {

				filePath = "htmlfile/LocusOfControlAbove14.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		} else if assessment_result[0]["assessment_id"] == "ywlxbz8yrlp955" {
			var filePath string

			if finalScore >= 6 && finalScore <= 24 {

				filePath = "htmlfile/GQ6_below_24.html"

			} else if finalScore >= 25 && finalScore <= 34 {

				filePath = "htmlfile/GQ6_below_25_To_34.html"
			} else if finalScore >= 35 && finalScore <= 40 {

				filePath = "htmlfile/GQ6_below_35_To_40.html"
			} else {
				filePath = "htmlfile/GQ6_below_41_To_42.html"
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


			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentGWB(assessment_data, filePath)
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}


		} else if title == "Grit Scale" || assessment_result[0]["assessment_id"] == "yu6h98081msq" {
			var filePath string

			if finalScore >= 1.0 && finalScore <= 2.0 {

				filePath = "htmlfile/GRIT_1.0_2.0.html"

			} else if finalScore >= 2.1 && finalScore <= 3.0 {

				filePath = "htmlfile/GRIT_2.1_3.0.html"

			} else if finalScore >= 3.1 && finalScore <= 4.0 {

				filePath = "htmlfile/GRIT_3.1_4.0.html"

			} else {

				filePath = "htmlfile/GRIT_4.1_5.0.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}
		} else if title == "UCLA Loneliness Scale" {
			var filePath string

			if finalScore >= 8.0 && finalScore <= 15.0 {

				filePath = "htmlfile/Loneliness_Scale_8_15.html"

			} else if finalScore >= 16.0 && finalScore <= 23.0 {

				filePath = "htmlfile/Loneliness_Scale_16_23.html"

			} else if finalScore >= 24.0 && finalScore <= 32.0 {

				filePath = "htmlfile/Loneliness_Scale_24_32.html"

			} else {

				filePath = "htmlfile/Loneliness_Scale_24_32.html"
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
				Answer8: assessment_result_details[7]["score"],
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentGAD7(assessment_data, filePath)
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}
		} else if title == "Doomscrolling Scale" || assessment_result[0]["assessment_id"] == "0m18gbm7vv13" {
			var filePath string

			if finalScore >= 0.0 && finalScore <= 16.0 {

				filePath = "htmlfile/doomscrolling_questionnaire_0_16.html"

			} else if finalScore >= 17.0 && finalScore <= 36.0 {

				filePath = "htmlfile/doomscrolling_questionnaire_17_36.html"

			} else if finalScore >= 37.0 && finalScore <= 60.0 {

				filePath = "htmlfile/doomscrolling_questionnaire_37_60.html"

			} else if finalScore >= 61.0 && finalScore <= 80.0 {

				filePath = "htmlfile/doomscrolling_questionnaire_61_80.html"

			} else {

				filePath = "htmlfile/doomscrolling_questionnaire_81_96.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}
		} else if title == "Assertiveness Scale" {
			var filePath string

			if finalScore >= 10.0 && finalScore <= 20.0 {

				filePath = "htmlfile/Self-Evaluation_for_Assertiveness_10_20.html"

			} else if finalScore >= 21.0 && finalScore <= 35.0 {

				filePath = "htmlfile/Self-Evaluation_for_Assertiveness_21_35.html"

			} else {

				filePath = "htmlfile/Self-Evaluation_for_Assertiveness_36_50.html"

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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}
		} else if title == "Psychological well-being scale" {
			var filePath string

			if finalScore >= 18.0 && finalScore <= 69.0 {

				filePath = "htmlfile/psychological_well-being_scale_18_69.html"

			} else if finalScore >= 70.0 && finalScore <= 94.0 {

				filePath = "htmlfile/psychological_well-being_scale_70_94.html"

			} else {

				filePath = "htmlfile/psychological_well-being_scale_95_126.html"

			}

			assessment_data := MODEL.AssessmentDownloadPSYCHOLOGICALWELLBEINGModel{
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
			}

			emailbody, ok = UTIL.GetHTMLTemplateForAssessmentPSYCHOLOGICALWELLBEING(assessment_data, filePath)
			if !ok {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
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
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
				return
			}

		}

		// created, ok := UTIL.GeneratePdfHeaderAndFooterFixted(emailbody, "pdffile/assessment1.pdf") // name created,

		// UTIL.PDFGenerater(emailbody, "pdffile/assessment1.pdf")

		created := UTIL.HtmlToPDFAssessment(emailbody)

		if created == nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
			return
		}

		s3Path := "assessment"
		filename := "example1.pdf"

		name, uploaded := UTIL.UploadToS3File(CONFIG.S3Bucket, s3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(filename), CONSTANT.S3PublicRead, created)
		if !uploaded {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.HTMLTemplateNotCreateMessage, CONSTANT.ShowDialog, response)
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

	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, fileName, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	fileName = endPointURL

	response["media_url"] = CONFIG.MediaURL
	response["pdf_name"] = fileName

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
