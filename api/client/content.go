package client

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
// @Tags Client Content
// @Summary Get contents
// @Router /client/content [get]
// @Param client_id query string false "Logged in client ID"
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

	response["videos"] = videos
	response["audios"] = audios
	response["articles"] = articles
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ContentLikeGet godoc
// @Tags Client Content
// @Summary Get liked contents
// @Router /client/content/like [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get liked content
	likedContent, status, ok := DB.SelectProcess("select content_id from "+CONSTANT.ContentLikesTable+" where person_id = ? order by created_at desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	contentIDs := UTIL.ExtractValuesFromArrayMap(likedContent, "content_id")

	// get latest videos
	videos, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where person_id = ?) and type = "+CONSTANT.VideoContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest audios
	audios, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where person_id = ?) and type = "+CONSTANT.AudioContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get latest articles
	articles, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where content_id in (select content_id from "+CONSTANT.ContentLikesTable+" where person_id = ?) and type = "+CONSTANT.ArticleContentType+" and training = 0 and status = 1 order by field(content_id, '"+strings.Join(contentIDs, "','")+"')", r.FormValue("client_id"))
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
// @Tags Client Content
// @Summary Like content
// @Router /client/content/like [post]
// @Param client_id query string true "Logged in client ID"
// @Param content_id query string true "Content ID to be liked"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	status, ok := DB.InsertSQL(CONSTANT.ContentLikesTable, map[string]string{
		"content_id": r.FormValue("content_id"),
		"person_id":  r.FormValue("client_id"),
		"created_at": UTIL.GetCurrentTime().String(),
	})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ContentLikeDelete godoc
// @Tags Client Content
// @Summary Unlike content
// @Router /client/content/like [delete]
// @Param client_id query string true "Logged in client ID"
// @Param content_id query string true "Content ID to be unliked"
// @Security JWTAuth
// @Produce json
// @Success 200
func ContentLikeDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	status, ok := DB.DeleteSQL(CONSTANT.ContentLikesTable, map[string]string{
		"content_id": r.FormValue("content_id"),
		"person_id":  r.FormValue("client_id"),
	})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
