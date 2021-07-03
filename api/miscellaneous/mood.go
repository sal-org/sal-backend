package miscellaneous

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

// ListMood godoc
// @Tags Miscellaneous
// @Summary Get all moods
// @Router /mood [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func ListMood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	moods, status, ok := DB.SelectSQL(CONSTANT.MoodsTable, []string{"*"}, map[string]string{})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["moods"] = moods
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
