package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// ListContentCategory godoc
// @Tags Miscellaneous
// @Summary Get all available content categories
// @Router /content-category [get]
// @Produce json
// @Success 200
func ListContentCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get content categories
	contentCategories, status, ok := DB.SelectSQL(CONSTANT.ContentCategoriesTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["content_categories"] = contentCategories
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
