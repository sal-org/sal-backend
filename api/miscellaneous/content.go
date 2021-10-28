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
// @Security JWTAuth
// @Produce json
// @Success 200
func Content(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var categoryFilter string
	if len(r.FormValue("category_id")) > 0 {
		id, _ := strconv.Atoi(r.FormValue("category_id"))
		if id > 0 {
			categoryFilter = " and category_id = " + r.FormValue("category_id")
		}
	}

	// get latest videos
	videos, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.VideoContentType + categoryFilter + " and training = 0 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest audios
	audios, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.AudioContentType + categoryFilter + " and training = 0 and status = 1 order by created_at desc limit 20")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest articles
	articles, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where type = " + CONSTANT.ArticleContentType + categoryFilter + " and training = 0 and status = 1 order by created_at desc limit 20")
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

	response["videos"] = videos
	response["audios"] = audios
	response["articles"] = articles
	response["liked_content_ids"] = UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")
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

	response["videos"] = videos
	response["audios"] = audios
	response["articles"] = articles
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
