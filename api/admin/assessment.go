package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

func AssessmentAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body := MODEL.AssessmentAddRequestInAdminPanel{}
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

	// add assessment
	assessmentID, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentsTable, CONSTANT.AssessmentDigits, map[string]string{
		"title":       body.Title,
		"subtitle":    body.SubTitles,
		"photo":       body.Photo,
		"duration":    body.Duration,
		"type":        body.Type,
		"instruction": body.Intruction,
		"order":       body.Order,
		"status":      body.Status,
		"created_at":  UTIL.GetCurrentTime().UTC().String(),
	}, "assessment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, question := range body.Questions {
		// add assessment question
		assessmentQuestionID, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentQuestionsTable, CONSTANT.AssessmentQuestionDigits, map[string]string{
			"assessment_id": assessmentID,
			"question":      question.Question,
			"order":         question.Order,
			"status":        question.Status,
			"created_at":    UTIL.GetCurrentTime().UTC().String(),
		}, "assessment_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, options := range question.Options {

			// add assessment result
			_, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentQuestionOptionsTable, CONSTANT.AssessmentQuestionOptionDigits, map[string]string{
				"assessment_id":          assessmentID,
				"assessment_question_id": assessmentQuestionID,
				"option":                 options.Option,
				"score":                  options.Score,
				"order":                  options.Order,
				"status":                 options.Status,
				"created_at":             UTIL.GetCurrentTime().UTC().String(),
			}, "assessment_question_option_id")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		}
	}

	response["assessment_id"] = assessmentID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
