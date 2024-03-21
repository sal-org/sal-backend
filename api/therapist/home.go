package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
)

// Home godoc
// @Tags Therapist Home
// @Summary Get home page content
// @Router /therapist/home [get]
// @Param therapist_id query string false "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get latest content for recommended
	recommended, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 1 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest videos
	// videos, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.VideoContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // get latest audios
	// audios, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.AudioContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // get latest articles
	// articles, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.ArticleContentType + " and training = 0 and status = 1 order by created_at desc limit 20")
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	response["recommended"] = recommended
	// response["videos"] = videos
	// response["audios"] = audios
	// response["articles"] = articles
	response["media_url"] = CONFIG.MediaURL
	response["urls"] = CONSTANT.URLs
	response["android_version"] = CONSTANT.TherapistAndroidVersion
	response["ios_version"] = CONSTANT.TherapistIOSVersion
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
