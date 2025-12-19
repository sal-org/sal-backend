package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
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
		"instruction": body.Instruction,
		"source":      body.Source,
		"reference":   body.Reference,
		"order":       body.Order,
		"status":      body.Status,
		"created_at":  UTIL.GetCurrentTime().UTC().String(),
	}, "assessment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, score := range body.Scores {
		// add assessment score
		status, ok := DB.InsertSQL(CONSTANT.AssessmentScoresTable, map[string]string{
			"assessment_id": assessmentID,
			"min":           score.MinScore,
			"max":           score.MaxScore,
			"result":        score.Result,
		})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	}

	for _, question := range body.Questions {
		// add assessment question
		assessmentQuestionID, status, ok := DB.InsertWithUniqueID(CONSTANT.AssessmentQuestionsTable, CONSTANT.AssessmentQuestionDigits, map[string]string{
			"assessment_id": assessmentID,
			"question":      question.Question,
			"order":         question.Order,
			"status":        question.Status,
			"created_at":    UTIL.GetCurrentTime().UTC().String(),
		}, "assessment_question_id")
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

func AssessmentUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body := MODEL.AssessmentUpdateRequestInAdminPanel{}
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
	status, ok := DB.UpdateSQL(CONSTANT.AssessmentsTable, map[string]string{"assessment_id": body.AssessmentID}, map[string]string{
		"title":       body.Title,
		"subtitle":    body.SubTitles,
		"photo":       body.Photo,
		"duration":    body.Duration,
		"type":        body.Type,
		"instruction": body.Instruction,
		"source":      body.Source,
		"reference":   body.Reference,
		"order":       body.Order,
		"status":      body.Status,
		"modified_at": UTIL.GetCurrentTime().UTC().String(),
	})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, score := range body.Scores {
		// add assessment score
		status, ok := DB.UpdateSQL(CONSTANT.AssessmentScoresTable, map[string]string{"assessment_id": body.AssessmentID, "id": score.ScoreID}, map[string]string{
			"min":    score.MinScore,
			"max":    score.MaxScore,
			"result": score.Result,
		})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	}

	for _, question := range body.Questions {
		// add assessment question
		status, ok := DB.UpdateSQL(CONSTANT.AssessmentQuestionsTable, map[string]string{"assessment_id": body.AssessmentID, "assessment_question_id": question.AssessmentQuestionID}, map[string]string{
			"question":    question.Question,
			"order":       question.Order,
			"status":      question.Status,
			"modified_at": UTIL.GetCurrentTime().UTC().String(),
		})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, options := range question.Options {

			// add assessment result
			status, ok := DB.UpdateSQL(CONSTANT.AssessmentQuestionOptionsTable, map[string]string{"assessment_question_option_id": options.AssessmentQuestionOptionID}, map[string]string{
				"option":      options.Option,
				"score":       options.Score,
				"order":       options.Order,
				"status":      options.Status,
				"modified_at": UTIL.GetCurrentTime().UTC().String(),
			})
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func AssessmentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get clients
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " title = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " assessment_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	assessments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentsTable+where+" order by 'order' desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// var

	assessmentswithDetails := []map[string]interface{}{}

	for _, m := range assessments {
		conv := make(map[string]interface{})
		for k, v := range m {
			conv[k] = v
		}
		assessmentswithDetails = append(assessmentswithDetails, conv)
	}

	for _, assessment := range assessmentswithDetails {
		// get assessment scores
		scores, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentScoresTable+" where assessment_id = ? order by min asc", assessment["assessment_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		assessment["scores"] = scores
	}

	for _, assessment := range assessmentswithDetails {
		// get assessment scores
		questions, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentQuestionsTable+" where assessment_id = ?", assessment["assessment_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		assessmentswithQuestion := []map[string]interface{}{}

		for _, m := range questions {
			conv := make(map[string]interface{})
			for k, v := range m {
				conv[k] = v
			}
			assessmentswithQuestion = append(assessmentswithQuestion, conv)
		}

		for _, question := range assessmentswithQuestion {
			// get assessment question options
			options, status, ok := DB.SelectProcess("select * from "+CONSTANT.AssessmentQuestionOptionsTable+" where assessment_question_id = ? and assessment_id = ? order by 'order' asc", question["assessment_question_id"], question["assessment_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			question["options"] = options
		}
		assessment["questions"] = assessmentswithQuestion
	}

	// get total number of clients
	assessmentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.AssessmentsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["assessments"] = assessmentswithDetails
	response["assessments_count"] = assessmentsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(assessmentsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToAssessmentUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.MiscellaneousS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
