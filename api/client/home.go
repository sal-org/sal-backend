package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
)

// Home godoc
// @Tags Client Home
// @Summary Get home page content
// @Router /client/home [get]
// @Param client_id query string false "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var recommended []map[string]string
	var ok bool
	var status string

	if len(r.FormValue("client_id")) > 0 {

		client, status, ok := DB.SelectProcess("select topic_ids from "+CONSTANT.ClientsTable+" where client_id = ? ", r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get latest content for recommended
		recommended, status, ok = DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where category_id in (" + client[0]["topic_ids"] + ") and training = 0 and status = 1 order by created_at desc limit 20")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		recommended, status, ok = DB.SelectProcess("select * from " + CONSTANT.ContentsTable + " where training = 0 and status = 1 order by created_at desc limit 20")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

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

	// get random quote
	quote, status, ok := DB.SelectProcess("select * from " + CONSTANT.QuotesTable + " order by rand() limit 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["recommended"] = recommended
	// response["videos"] = videos
	// response["audios"] = audios
	// response["articles"] = articles
	response["quote"] = quote[0]["quote"]
	response["media_url"] = CONFIG.MediaURL
	response["urls"] = CONSTANT.URLs
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
