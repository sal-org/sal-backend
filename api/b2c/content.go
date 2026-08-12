package b2c

import (
	"fmt"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
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
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.GetContentRequestInWebSite{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	wheres = append(wheres, " content_mode = ? ")
	contentArgs = append(contentArgs, body.ContentMood)

	if len(body.ResourceID) > 0 {
		wheres = append(wheres, " resource_id = ? ")
		contentArgs = append(contentArgs, body.ResourceID)
	}
	if len(body.CategoryID) > 0 {
		wheres = append(wheres, " FIND_IN_SET("+body.CategoryID+", category_id) ")
	}
	if len(body.Name) > 0 {
		wheres = append(wheres, " content like '%%"+body.Name+"%%'")
	}
	if len(body.Type) > 0 {
		wheres = append(wheres, " type = ? ")
		contentArgs = append(contentArgs, body.Type)
	}

	wheres = append(wheres, " status = "+CONSTANT.ContentActive+" ") // only active therapists
	contentSQLQuery += " where " + strings.Join(wheres, " and ")

	SQLQuery = "select content_id, title, subtitle, description, photo, background_photo, share_content, content, duration, type, redirection, category_id, resource_id, content_mode, article_page_id, created_by from " + CONSTANT.ContentsInWebTable + contentSQLQuery + " "
	args = append(args, contentArgs...)

	sortBy := " created_at " // default ordering by
	orderBy := " desc "

	SQLQuery += " order by " + sortBy + orderBy

	contents, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.ContentForWebPerPageUser)+" offset "+strconv.Itoa((UTIL.GetPageNumber(body.Page)-1)*CONSTANT.ContentForWebPerPageUser), args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	contentsCount, status, ok := DB.SelectProcess("select count(content_id) as ctn from ("+SQLQuery+") as a", args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, content := range contents {
		// urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// _, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		// content["photo"] = endPointURL

		// urlBackgroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// _, endPointURLBackgroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackgroundPhoto)
		// content["background_photo"] = endPointURLBackgroundPhoto

		// if content["type"] == CONSTANT.VideoContentType || content["type"] == CONSTANT.AudioContentType {
		// 	urlShareContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["share_content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// 	_, endPointURLShareContent := UTIL.GetBaseURLAndEndpointFromURL(urlShareContent)
		// 	content["share_content"] = endPointURLShareContent
		// }

		// if content["type"] != "3" {
		// 	urlContent := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, content["content"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		// 	_, endPointURLContent := UTIL.GetBaseURLAndEndpointFromURL(urlContent)
		// 	content["content"] = endPointURLContent
		// }

		if content["type"] == CONSTANT.ArticleContentType {
			if len(content["article_page_id"]) != 0 {
				id, _ := strconv.Atoi(content["article_page_id"])
				articlePage, err := UTIL.GetWordPressPageByID(CONFIG.WordPressURL, CONFIG.WordPressUsername, CONFIG.WordPressApplicationPassword, id)
				if err != nil {
					content["content"] = ""
					fmt.Println("Error fetching WordPress page by id:", err)
				} else {
					content["content"] = articlePage.Content.Rendered
				}
			} else {
				articlePage, err := UTIL.GetWordPressPageBySlug(CONFIG.WordPressURL, CONFIG.WordPressUsername, CONFIG.WordPressApplicationPassword, content["content"])
				if err != nil {
					content["content"] = ""
					fmt.Println("Error fetching WordPress page by slug:", err)
				} else {
					content["content"] = articlePage.Content.Rendered

					DB.UpdateSQL(CONSTANT.ContentsInWebTable, map[string]string{"content_id": content["content_id"]}, map[string]string{"article_page_id": strconv.Itoa(articlePage.ID)})
				}
			}

			// println("content[\"content\"]", content["content"])
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
