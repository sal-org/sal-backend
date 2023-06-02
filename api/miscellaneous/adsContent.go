package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
)

// ListMood godoc
// @Tags Miscellaneous
// @Summary Get all moods
// @Router /adsContent [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func AdsContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	adsContent, status, ok := DB.SelectSQL(CONSTANT.AdsContentTable, []string{"title", "target", "image"}, map[string]string{"status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["ads"] = adsContent
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
