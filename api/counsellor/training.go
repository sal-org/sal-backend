package counsellor

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
)

// Training godoc
// @Tags Counsellor Training
// @Summary Get counsellor training content
// @Router /counsellor/training [get]
// @Param counsellor_id query string false "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Training(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get latest training content
	training, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 1 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["training"] = training
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
