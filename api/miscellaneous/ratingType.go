package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// ListRatingType godoc
// @Tags Miscellaneous
// @Summary Get ratings types
// @Router /rating-type [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func ListRatingType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get rating types
	ratingTypes, status, ok := DB.SelectSQL(CONSTANT.RatingTypesTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["rating_types"] = ratingTypes
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
