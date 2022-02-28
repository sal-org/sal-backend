package client

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
// @Tags Client Assessment
// @Summary List available assessments
// @Router /client/assessments [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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

	var response = make(map[string]interface{})

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

	var response = make(map[string]interface{})

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
// @Tags Client Assessment
// @Summary Get assessment history
// @Router /client/assessment/history [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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
// @Tags Client Assessment
// @Summary Get assessment download
// @Router /client/assessment/download [get]
// @Param assessment_result_id query string true "Logged in Assessment Result ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var fileName string

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

	filePath := "htmlfile/assessment1.html"

	emailbody, ok := UTIL.GetHTMLTemplateForAssessmentAIS(assessment_data, filePath)
	if !ok {
		fmt.Println("html body not create ")
	}

	created, ok := UTIL.GeneratePdf(emailbody, "pdffile/assessment1.pdf") // name created,

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

	//receipt := map[string]string{}

	response["invoice_id"] = assessment_result[0]["assessment_result_id"]
	response["pdf"] = fileName

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
