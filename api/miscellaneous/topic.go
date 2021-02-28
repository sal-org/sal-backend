package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// ListTopic godoc
// @Tags Miscellaneous
// @Summary Get all available topics
// @Router /topic [get]
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func ListTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get topics
	topics, status, ok := DB.SelectSQL(CONSTANT.TopicsTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["topics"] = topics
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
