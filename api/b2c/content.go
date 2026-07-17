package b2c

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
	"strconv"
	"strings"
)

func GetWebsiteContent(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// var encryptedResponse = make(map[string]any)

	var SQLQuery, contentSQLQuery string
	args := []any{}
	contentArgs := []any{}

	wheres := []string{}

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	if !(len(body["content_mode"]) != 0 && len(body["content_mode"]) < 2) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ContentMoodIsRequiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	wheres = append(wheres, " content_mode = ? ")
	contentArgs = append(contentArgs, body["content_mode"])

	if len(body["resource_id"]) > 0 {
		wheres = append(wheres, " resource_id = ? ")
		contentArgs = append(contentArgs, body["resource_id"])
	}
	if len(body["category_id"]) > 0 {
		wheres = append(wheres, " FIND_IN_SET("+body["category_id"]+", category_id) ")
	}
	if len(body["name"]) > 0 {
		wheres = append(wheres, " content like '%%"+body["name"]+"%%'")
	}
	if len(body["type"]) > 0 {
		wheres = append(wheres, " type = ? ")
		contentArgs = append(contentArgs, body["type"])
	}

	wheres = append(wheres, " status = "+CONSTANT.ContentActive+" ") // only active therapists
	contentSQLQuery += " where " + strings.Join(wheres, " and ")

	SQLQuery = "select * from " + CONSTANT.ContentsInWebTable + contentSQLQuery + " "
	args = append(args, contentArgs...)

	sortBy := " created_at " // default ordering by
	orderBy := " desc "

	SQLQuery += " order by " + sortBy + orderBy

	contents, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.ContentForWebPerPageUser)+" offset "+strconv.Itoa((UTIL.GetPageNumber(body["page"])-1)*CONSTANT.ContentForWebPerPageUser), args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	contentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from ("+SQLQuery+") as a", args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, content := range contents {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		content["photo"] = endPointURL

		urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		content["background_photo"] = endPointURLBackgroundPhoto

		if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType {
			urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
			content["share_content"] = endPointURLShareContent
		}

		if content["type"] != "3" {
			urlContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
			_, endPointURLContent := UTIL.GetBaseURLAndEndpointFromURL(urlContent)
			content["content"] = endPointURLContent
		}

		if content["type"] == CONSTANT.ArticleContentType {
			if len(content["article_page_id"]) != 0 {
				id, _ := strconv.Atoi(content["article_page_id"])
				articlePage, err := UTIL.GetWordPressPageByID(CONFIG.WordPressURL, id)
				if err != nil {
					content["content"] = ""
				} else {
					content["content"] = articlePage.Content.Rendered
				}
			} else {
				articlePage, err := UTIL.GetWordPressPageBySlug(CONFIG.WordPressURL, content["content"])
				if err != nil {
					content["content"] = ""
				} else {
					content["content"] = articlePage.Content.Rendered

					DB.UpdateSQL(CONSTANT.ContentsInWebTable, map[string]string{"content_id": content["content_id"]}, map[string]string{"article_page_id": strconv.Itoa(articlePage.ID)})
				}
			}
		}
	}

	response["contents"] = contents
	response["contents_count"] = contentsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(contentsCount[0]["ctn"], CONSTANT.ContentForWebPerPageUser))
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	// encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	// if encrypt == "" {
	// 	UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetResourceCategoryForWeb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// var encryptedResponse = make(map[string]any)

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	contentCategories, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentCategoriesInWebTable)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	contentResource, status, ok := DB.SelectProcess("select * from " + CONSTANT.ResourceCategoriesInWebTable)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["categories"] = contentCategories
	response["resources"] = contentResource

	// encrypt, _ := EncryptPayload(response, CONSTANT.ENCRYPTION_SECRET_KEY_FOR_WEB, CONSTANT.ENCRYPTION_SECRET_IV_FOR_WEB)
	// if encrypt == "" {
	// 	UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// encryptedResponse["data"] = encrypt

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
