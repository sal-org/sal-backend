package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// ListMeta godoc
// @Tags Miscellaneous
// @Summary Get all available topics, languages
// @Router /meta [get]
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func ListMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get topics
	topics, status, ok := DB.SelectSQL(CONSTANT.TopicsTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get languages
	languages, status, ok := DB.SelectSQL(CONSTANT.LanguagesTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["topics"] = topics
	response["languages"] = languages
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
