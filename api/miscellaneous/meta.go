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
// @Security JWTAuth
// @Produce json
// @Success 200
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

	// get content categories
	contentCategories, status, ok := DB.SelectSQL(CONSTANT.ContentCategoriesTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get rating types
	ratingTypes, status, ok := DB.SelectSQL(CONSTANT.RatingTypesTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["topics"] = topics
	response["languages"] = languages
	response["content_categories"] = contentCategories
	response["rating_types"] = ratingTypes
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
