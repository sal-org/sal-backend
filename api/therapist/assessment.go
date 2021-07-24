package therapist

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

// AssessmentsList godoc
// @Tags Therapist Assessment
// @Summary List available assessments
// @Router /therapist/assessments [get]
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

	response["assessments"] = assessments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AssessmentDetail godoc
// @Tags Therapist Assessment
// @Summary Get assessment detail
// @Router /therapist/assessment [get]
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
// @Tags Therapist Assessment
// @Summary Add therapist assessment
// @Router /therapist/assessment [post]
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
// @Tags Therapist Assessment
// @Summary Get assessment history
// @Router /therapist/assessment/history [get]
// @Param assessment_id query string true "Assessment ID to get details"
// @Param therapist_id query string true "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func AssessmentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get assessment past results
	assessmentResults, status, ok := DB.SelectProcess("select final_score, created_at from "+CONSTANT.AssessmentResultsTable+" where user_id = ? and assessment_id = ? and status = "+CONSTANT.AssessmentResultActive+" order by created_at desc", r.FormValue("therapist_id"), r.FormValue("assessment_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["assessment_results"] = assessmentResults
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
