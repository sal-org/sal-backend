package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

// Content godoc
// @Tags Content
// @Summary Get contents
// @Router /content [get]
// @Param user_id query string false "Logged in user ID (client_id/counsellor_id/listener_id/therapist_id)"
// @Param category_id query string false "Content category ID - false if required all"
// @Param mood_id query string false "Content mood ID - false if required all"
// @Security JWTAuth
// @Produce json
// @Success 200
func Content(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var categoryFilter string
	var moodFilter string
	if len(r.FormValue("category_id")) > 0 {
		id, _ := strconv.Atoi(r.FormValue("category_id"))
		if id > 0 {
			categoryFilter = " and category_id like '%" + r.FormValue("category_id") + "%'"
		}
	}

	if len(r.FormValue("mood_id")) > 0 {
		id, _ := strconv.Atoi(r.FormValue("mood_id"))
		if id > 0 {
			moodFilter = " and mood_id like '%" + r.FormValue("mood_id") + "%'"
		}
	}

	if len(r.FormValue("liked")) != 0 {

		likedContent, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		contentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

		contentType, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where content_id in ('" + strings.Join(contentIDs, "','") + "') and type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at desc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		videosCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from " + CONSTANT.ContentsTable + " where content_id in ('" + strings.Join(contentIDs, "','") + "') and type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1 ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, content := range contentType {
			urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
			content["photo"] = endPointURL

			urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
			content["background_photo"] = endPointURLBackgroundPhoto

			if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
				urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
				content["share_content"] = endPointURLShareContent
			}

			if len(content["counsellor_photo"]) > 0 {
				urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
				content["counsellor_photo"] = endPointURLCounsellorPhoto
			}
		}

		if r.FormValue("type") == "1" {
			contentVideo, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at asc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			for _, content := range contentVideo {
				urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
				content["photo"] = endPointURL

				urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
				content["background_photo"] = endPointURLBackgroundPhoto

				if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
					urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
					_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
					content["share_content"] = endPointURLShareContent
				}

				if len(content["counsellor_photo"]) > 0 {
					urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
					_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
					content["counsellor_photo"] = endPointURLCounsellorPhoto
				}
			}

			response["videos"] = contentVideo
			response["videos_count"] = videosCount[0]["ctn"]
			response["no_pages_videos"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		} else if r.FormValue("type") == "2" {
			response["audios"] = contentType
			response["audios_count"] = videosCount[0]["ctn"]
			response["no_pages_audios"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		} else {
			response["articles"] = contentType
			response["articles_count"] = videosCount[0]["ctn"]
			response["no_pages_articles"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		}

		response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

	} else if len(r.FormValue("type")) != 0 {
		contentType, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at desc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get liked content ids
		likedContent, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		videosCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from " + CONSTANT.ContentsTable + " where type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, content := range contentType {
			urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
			content["photo"] = endPointURL

			urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
			content["background_photo"] = endPointURLBackgroundPhoto

			if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
				urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
				content["share_content"] = endPointURLShareContent
			}

			if len(content["counsellor_photo"]) > 0 {
				urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
				content["counsellor_photo"] = endPointURLCounsellorPhoto
			}
		}

		if r.FormValue("type") == "1" {
			// contentVideo, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + r.FormValue("type") + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at asc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			likedContentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

			for i := range contentType {
				if UTIL.IsStringInSlice(contentType[i]["content_id"], likedContentIDs) {
					contentType[i]["liked"] = "1"
				} else {
					contentType[i]["liked"] = "0"
				}
			}

			response["videos"] = contentType
			response["videos_count"] = videosCount[0]["ctn"]
			response["no_pages_videos"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		} else if r.FormValue("type") == "2" {
			likedContentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

			for i := range contentType {
				if UTIL.IsStringInSlice(contentType[i]["content_id"], likedContentIDs) {
					contentType[i]["liked"] = "1"
				} else {
					contentType[i]["liked"] = "0"
				}
			}
			response["audios"] = contentType
			response["audios_count"] = videosCount[0]["ctn"]
			response["no_pages_audios"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		} else {

			likedContentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

			for i := range contentType {
				if UTIL.IsStringInSlice(contentType[i]["content_id"], likedContentIDs) {
					contentType[i]["liked"] = "1"
				} else {
					contentType[i]["liked"] = "0"
				}
			}

			response["articles"] = contentType
			response["articles_count"] = videosCount[0]["ctn"]
			response["no_pages_articles"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		}

		response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

	} else {
		// get latest videos
		videos, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.VideoContentType + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at asc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get latest audios
		audios, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.AudioContentType + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at desc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get latest articles
		articles, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.ArticleContentType + categoryFilter + moodFilter + " and training = 0 and status = 1 order by created_at desc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get liked content ids
		likedContent, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		likedContentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

		for i := range videos {
			if UTIL.IsStringInSlice(videos[i]["content_id"], likedContentIDs) {
				videos[i]["liked"] = "1"
			} else {
				videos[i]["liked"] = "0"
			}
		}

		for i := range audios {
			if UTIL.IsStringInSlice(audios[i]["content_id"], likedContentIDs) {
				audios[i]["liked"] = "1"
			} else {
				audios[i]["liked"] = "0"
			}
		}

		for i := range articles {
			if UTIL.IsStringInSlice(articles[i]["content_id"], likedContentIDs) {
				articles[i]["liked"] = "1"
			} else {
				articles[i]["liked"] = "0"
			}
		}

		// get total number of contents
		videosCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.VideoContentType + categoryFilter + moodFilter + " and training = 0 and status = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		audiosCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.AudioContentType + categoryFilter + moodFilter + " and training = 0 and status = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		articlesCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.ArticleContentType + categoryFilter + moodFilter + " and training = 0 and status = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, content := range videos {
			urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
			content["photo"] = endPointURL

			urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
			content["background_photo"] = endPointURLBackgroundPhoto

			if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
				urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
				content["share_content"] = endPointURLShareContent
			}

			if len(content["counsellor_photo"]) > 0 {
				urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
				content["counsellor_photo"] = endPointURLCounsellorPhoto
			}
		}

		for _, content := range audios {
			urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
			content["photo"] = endPointURL

			urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
			content["background_photo"] = endPointURLBackgroundPhoto

			if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
				urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
				content["share_content"] = endPointURLShareContent
			}

			if len(content["counsellor_photo"]) > 0 {
				urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
				_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
				content["counsellor_photo"] = endPointURLCounsellorPhoto
			}
		}

		response["videos"] = videos
		response["audios"] = audios
		response["articles"] = articles
		response["videos_count"] = videosCount[0]["ctn"]
		response["audios_count"] = audiosCount[0]["ctn"]
		response["articles_count"] = articlesCount[0]["ctn"]
		response["no_pages_videos"] = strconv.Itoa(UTIL.GetNumberOfPages(videosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		response["no_pages_audios"] = strconv.Itoa(UTIL.GetNumberOfPages(audiosCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		response["no_pages_articles"] = strconv.Itoa(UTIL.GetNumberOfPages(articlesCount[0]["ctn"], CONSTANT.ContentPerPageUser))
		response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")
	}

	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetContentUsedTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	contentType, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where title like '%" + r.FormValue("content_name") + "%'" + "and training = 0 and type = '" + r.FormValue("type") + "' and status = 1 order by created_at desc limit " + strconv.Itoa(CONSTANT.ContentPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ContentPerPageUser))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	contentCount, status, ok := DB.SelectProcess("select count(*) as ctn from " + CONSTANT.ContentsTable + " where title like '%" + r.FormValue("content_name") + "%'" + "and training = 0 and type = '" + r.FormValue("type") + "' and status = 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, content := range contentType {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		content["photo"] = endPointURL

		urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		content["background_photo"] = endPointURLBackgroundPhoto

		if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
			urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
			content["share_content"] = endPointURLShareContent
		}

		if len(content["counsellor_photo"]) > 0 {
			urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
			content["counsellor_photo"] = endPointURLCounsellorPhoto
		}
	}

	response["contents_name"] = contentType
	response["contents_count"] = contentCount[0]["ctn"]
	response["no_pages_contents"] = strconv.Itoa(UTIL.GetNumberOfPages(contentCount[0]["ctn"], CONSTANT.ContentPerPageUser))
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ContentLikeGet godoc
// @Tags Content
// @Summary Get liked contents
// @Router /content/like [get]
// @Param user_id query string true "Logged in user ID (client_id/counsellor_id/listener_id/therapist_id)"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get liked content ids
	likedContent, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	contentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

	// get latest liked videos
	videos, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ?) and type = "+CONSTANT.VideoContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest liked audios
	audios, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ?) and type = "+CONSTANT.AudioContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest liked articles
	articles, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ?) and type = "+CONSTANT.ArticleContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get liked content ids
	contentLiked, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where user_id = ? order by created_at desc", r.FormValue("user_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, content := range videos {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		content["photo"] = endPointURL

		urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		content["background_photo"] = endPointURLBackgroundPhoto

		if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
			urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
			content["share_content"] = endPointURLShareContent
		}

		if len(content["counsellor_photo"]) > 0 {
			urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
			content["counsellor_photo"] = endPointURLCounsellorPhoto
		}
	}

	for _, content := range audios {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		content["photo"] = endPointURL

		urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		content["background_photo"] = endPointURLBackgroundPhoto

		if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType || content["type"] == CONSTANT.ArticleContentType {
			urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
			content["share_content"] = endPointURLShareContent
		}

		if len(content["counsellor_photo"]) > 0 {
			urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
			content["counsellor_photo"] = endPointURLCounsellorPhoto
		}
	}

	response["videos"] = videos
	response["audios"] = audios
	response["articles"] = articles
	response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(contentLiked, "content_id")
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ContentLikeAdd godoc
// @Tags Content
// @Summary Like content
// @Router /content/like [post]
// @Param user_id query string true "Logged in user ID (client_id/counsellor_id/listener_id/therapist_id)"
// @Param content_id query string true "Content ID to be liked"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.InsertSQL(CONSTANT.ContentLikesTable, map[string]string{
		"content_id": r.FormValue("content_id"),
		"user_id":    r.FormValue("user_id"),
		"created_at": UTIL.GetCurrentTime().String(),
	})
	if !ok {
		UTIL.SetReponse(w, status, CONSTANT.ContentAlreadyLikedMessage, CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ContentLikeDelete godoc
// @Tags Content
// @Summary Unlike content
// @Router /content/like [delete]
// @Param user_id query string true "Logged in user ID (client_id/counsellor_id/listener_id/therapist_id)"
// @Param content_id query string true "Content ID to be unliked"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.DeleteSQL(CONSTANT.ContentLikesTable, map[string]string{
		"content_id": r.FormValue("content_id"),
		"user_id":    r.FormValue("user_id"),
	})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func IncreaseContentViewCount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	content, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where content_id = '" + r.FormValue("content_id") + "'")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(content) == 0 {
		UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		return
	}

	view, _ := strconv.Atoi(content[0]["views"])

	view = view + 1

	viewsInString := strconv.Itoa(view)

	userID := ""

	if len(r.FormValue("user_id")) != 0 {
		userID = r.FormValue("user_id")
	}

	DB.UpdateSQL(CONSTANT.ContentsTable,
		map[string]string{
			"content_id": content[0]["content_id"],
		},
		map[string]string{
			"views": viewsInString,
		},
	)

	DB.InsertSQL(CONSTANT.ContentUserViewsTable, map[string]string{
		"content_id": content[0]["content_id"],
		"user_id":    userID,
		"status":     "1",
		"created_at": UTIL.GetCurrentTime().String(),
	})

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
