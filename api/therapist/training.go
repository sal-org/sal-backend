package therapist

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
)

// Training godoc
// @Tags Therapist Training
// @Summary Get therapist training content
// @Router /therapist/training [get]
// @Param therapist_id query string false "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Training(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get latest training content
	training, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 1 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// for _, content := range training {
	// 	urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
	// 	content["photo"] = endPointURL

	// 	urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
	// 	content["background_photo"] = endPointURLBackgroundPhoto

	// 	if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType {
	// 		urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
	// 		content["share_content"] = endPointURLShareContent
	// 	}

	// 	if content["type"] != "3" {
	// 		urlContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 		_, endPointURLContent := UTIL.GetBaseURLAndEndpointFromURL(urlContent)
	// 		content["content"] = endPointURLContent
	// 	}
	// }

	response["training"] = training
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
