package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

func ContentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get contents
	wheres := []string{}
	queryArgs := []any{}
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
	}

	response["contents"] = contents
	response["contents_count"] = contentsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(contentsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ContentAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := Model.ContentAddRequestInAdminPanel{}

	// // check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ContentAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

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

	counsellorPhoto := ""

	var counsellors []map[string]string

	if len(body.CounsellorID) != 0 {
		// get client details
		counsellor, _, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name, email", "photo"}, map[string]string{"counsellor_id": body.CounsellorID})
		if !ok {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			counsellor, _, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name, email", "photo"}, map[string]string{"listener_id": body.CounsellorID})
			if !ok {
				UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
				return
			}
		}

		if len(counsellor) == 0 {
			counsellor, _, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name, email", "photo"}, map[string]string{"therapist_id": body.CounsellorID})
			if !ok {
				UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
				return
			}
		}

		if len(counsellor) != 0 {
			counsellorPhoto = counsellor[0]["photo"]
			counsellors = append(counsellors, counsellor[0])
		}

	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add content
	content := map[string]string{}
	content["counsellor_id"] = body.CounsellorID
	content["title"] = body.Title
	content["description"] = body.Description
	content["photo"] = body.Photo
	content["background_photo"] = body.BackgroundPhoto
	content["share_content"] = body.ShareContent
	content["content"] = body.Content
	content["type"] = body.Type
	content["redirection"] = body.Redirection
	content["category_id"] = body.CategoryID
	content["training"] = body.Training
	content["mood_id"] = body.MoodID
	content["duration"] = body.Duration
	content["counsellor_photo"] = counsellorPhoto
	content["status"] = CONSTANT.ContentActive
	content["created_by"] = id
	content["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.ContentsTable, CONSTANT.ContentDigits, content, "content_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellors) != 0 {

		filepath_text := "htmlfile/emailmessagebody.html"
		// send email for therapist
		emaildata := Model.EmailBodyMessageModel{
			Name: counsellors[0]["first_name"],
			Message: UTIL.ReplaceNotificationContentInString(
				CONSTANT.CounsellorApprovedContentBody,
				map[string]string{
					"###content_name###": body.Title,
				},
			),
		}

		emailBody := UTIL.GetHTMLTemplateForCounsellorProfileText(emaildata, filepath_text)
		// email for therapist
		UTIL.SendEmail(
			CONSTANT.CounsellorApprovedContentTitle,
			emailBody,
			counsellors[0]["email"],
			CONSTANT.InstantSendEmailMessage,
		)
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ContentUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	contentID, ok := UTIL.Required(r.FormValue("content_id"), "Content ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, contentID, CONSTANT.ShowDialog, response)
		return
	}

	body := Model.ContentUpdateRequestInAdminPanel{}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

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

	// check if connect_id exists
	if !DB.CheckIfExists(CONSTANT.ContentsTable, map[string]string{"content_id": contentID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	var counsellors []map[string]string
	counsellorPhoto := ""

	if len(body.CounsellorID) != 0 {
		// get client details
		counsellor, _, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name, email", "photo"}, map[string]string{"counsellor_id": body.CounsellorID})
		if !ok {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		}

		if len(counsellor) == 0 {
			counsellor, _, ok = DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name, email", "photo"}, map[string]string{"listener_id": body.CounsellorID})
			if !ok {
				UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
				return
			}
		}

		if len(counsellor) == 0 {
			counsellor, _, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name, email", "photo"}, map[string]string{"therapist_id": body.CounsellorID})
			if !ok {
				UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
				return
			}
		}

		if len(counsellor) != 0 {
			counsellorPhoto = counsellor[0]["photo"]
			counsellors = append(counsellors, counsellor[0])
		}

	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add content
	content := map[string]string{}
	content["counsellor_id"] = body.CounsellorID
	content["title"] = body.Title
	content["description"] = body.Description
	content["photo"] = body.Photo
	content["background_photo"] = body.BackgroundPhoto
	content["share_content"] = body.ShareContent
	content["content"] = body.Content
	content["type"] = body.Type
	content["redirection"] = body.Redirection
	content["category_id"] = body.CategoryID
	content["training"] = body.Training
	content["counsellor_photo"] = counsellorPhoto
	content["mood_id"] = body.MoodID
	content["duration"] = body.Duration
	content["status"] = body.Status
	content["modified_by"] = id
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

	var response = make(map[string]any)

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

	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, fileName, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	fileName = endPointURL

	response["file"] = fileName
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.ContentS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
