package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

func ContentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "mood_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " mood_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "category_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " category_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "type":
			if len(val[0]) > 0 {
				wheres = append(wheres, " type = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "training":
			if len(val[0]) > 0 {
				wheres = append(wheres, " training = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "content_id":
			wheres = append(wheres, " content_id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	contents, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of contents
	contentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.ContentsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["contents"] = contents
	response["contents_count"] = contentsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(contentsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ContentAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ContentAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// add content
	content := map[string]string{}
	content["counsellor_id"] = body["counsellor_id"]
	content["title"] = body["title"]
	content["description"] = body["description"]
	content["photo"] = body["photo"]
	content["background_photo"] = body["background_photo"]
	content["share_content"] = body["share_content"]
	content["content"] = body["content"]
	content["type"] = body["type"]
	content["redirection"] = body["redirection"]
	content["category_id"] = body["category_id"]
	content["training"] = body["training"]
	content["mood_id"] = body["mood_id"]
	content["duration"] = body["duration"]
	content["status"] = CONSTANT.ContentActive
	content["created_by"] = body["created_by"]
	content["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.ContentsTable, CONSTANT.ContentDigits, content, "content_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ContentUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add content
	content := map[string]string{}
	content["counsellor_id"] = body["counsellor_id"]
	content["title"] = body["title"]
	content["description"] = body["description"]
	content["photo"] = body["photo"]
	content["background_photo"] = body["background_photo"]
	content["share_content"] = body["share_content"]
	content["content"] = body["content"]
	content["type"] = body["type"]
	content["redirection"] = body["redirection"]
	content["category_id"] = body["category_id"]
	content["training"] = body["training"]
	content["mood_id"] = body["mood_id"]
	content["created_by"] = body["created_by"]
	content["duration"] = body["duration"]
	content["status"] = body["status"]
	content["modified_by"] = body["modified_by"]
	content["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.ContentsTable, map[string]string{"content_id": r.FormValue("content_id")}, content)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UploadContentFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var fileName string
	// file upload
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("UploadContentFile", err)
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	if file != nil {
		defer file.Close()

		name, uploaded := UTIL.UploadToS3(CONFIG.S3Bucket, CONSTANT.ContentS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(handler.Filename), CONSTANT.S3PublicRead, file)
		if !uploaded {
			fmt.Println("UploadContentFile", err)
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
			return
		}
		fileName = name
	}

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.ContentS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
