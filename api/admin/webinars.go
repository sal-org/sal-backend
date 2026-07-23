package admin

import (
	"net/http"
	"net/url"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

func WebinarsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get events
	wheres := []string{}
	queryArgs := []any{}
	for key, val := range r.URL.Query() {
		switch key {
		case "title":
			if len(val[0]) > 0 {
				wheres = append(wheres, " title like '%%"+val[0]+"%%' ")
			}
		case "counsellor_name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_name like '%%"+val[0]+"%%' ")
			}
		case "webinar_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " webinar_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.WebinarsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of events
	eventsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.WebinarsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, event := range events {
		event["booked_count"] = "0"
		// get event booked count
		webinarBookedCount, status, ok := DB.SelectProcess("select count(*) as ctn from " + CONSTANT.WebinarsBookTable + " where webinar_id  = '" + event["webinar_id"] + "' and status > " + CONSTANT.OrderWaiting + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		event["booked_count"] = webinarBookedCount[0]["ctn"]
	}

	for _, event := range events {
		urlPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, event["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlPhoto)
		event["photo"] = endPointURLPhoto

		urlBackGroundPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, event["background_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLBackGroundPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlBackGroundPhoto)
		event["background_photo"] = endPointURLBackGroundPhoto

		urlCounsellorPhoto := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, event["counsellor_photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURLCounsellorPhoto := UTIL.GetBaseURLAndEndpointFromURL(urlCounsellorPhoto)
		event["counsellor_photo"] = endPointURLCounsellorPhoto
	}

	response["events"] = events
	response["events_count"] = eventsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func WebinarsAdd(w http.ResponseWriter, r *http.Request) {
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

	body := MODEL.AddWebinarSessionInAdminPanel{}

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

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add event
	webinar := map[string]string{}
	webinar["counsellor_name"] = body.CounsellorName
	webinar["title"] = body.Title
	webinar["counsellor_photo"] = body.CounsellorPhoto
	webinar["counsellor_qualification"] = body.CounsellorQualification
	webinar["counsellor_about"] = body.CounsellorAbout
	webinar["counsellor_experience"] = body.CounsellorExperience
	webinar["counsellor_rating"] = body.CounsellorRating
	webinar["about"] = body.About
	webinar["partner_name"] = body.PartnerName
	webinar["why_attend"] = body.WhyAttend
	webinar["address"] = body.Address
	webinar["photo"] = body.Photo
	webinar["background_photo"] = body.BackgroundPhoto
	webinar["date"] = body.Date
	webinar["time"] = body.Time
	webinar["duration"] = body.Duration
	webinar["mode"] = body.Mode
	webinar["status"] = body.Status
	webinar["created_by"] = id
	webinar["created_at"] = UTIL.GetCurrentTime().String()

	_, status, ok := DB.InsertWithUniqueID(CONSTANT.WebinarsTable, CONSTANT.WebinarsDigits, webinar, "webinar_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func WebinarsUpdate(w http.ResponseWriter, r *http.Request) {
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

	webinarID, ok := UTIL.Required(r.FormValue("webinar_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, webinarID, CONSTANT.ShowDialog, response)
		return
	}


	body := MODEL.AddWebinarSessionInAdminPanel{}

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

	// check if corporate_id exists
	if !DB.CheckIfExists(CONSTANT.WebinarsTable, map[string]string{"webinar_id": webinarID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	// id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// counsellor photo
	parsedCounsellorPhoto, err := url.Parse(body.CounsellorPhoto)
	if err != nil {
		panic(err)
	}
	// Remove query string
	parsedCounsellorPhoto.RawQuery = ""



	// event photo
	parsedPhoto, err := url.Parse(body.Photo)
	if err != nil {
		panic(err)
	}
	// Remove query string
	parsedPhoto.RawQuery = ""


	
	// event background photo
	parsedBackgroundPhoto, err := url.Parse(body.BackgroundPhoto)
	if err != nil {
		panic(err)
	}
	// Remove query string
	parsedBackgroundPhoto.RawQuery = ""

	// add event
	webinar := map[string]string{}
	webinar["counsellor_name"] = body.CounsellorName
	webinar["title"] = body.Title
	webinar["counsellor_photo"] = parsedCounsellorPhoto.String()
	webinar["counsellor_qualification"] = body.CounsellorQualification
	webinar["counsellor_about"] = body.CounsellorAbout
	webinar["counsellor_experience"] = body.CounsellorExperience
	webinar["counsellor_rating"] = body.CounsellorRating
	webinar["about"] = body.About
	webinar["why_attend"] = body.WhyAttend
	webinar["partner_name"] = body.PartnerName
	webinar["address"] = body.Address
	webinar["photo"] = parsedPhoto.String()
	webinar["background_photo"] = parsedBackgroundPhoto.String()
	webinar["date"] = body.Date
	webinar["time"] = body.Time
	webinar["duration"] = body.Duration
	webinar["mode"] = body.Mode
	webinar["status"] = body.Status
	status, ok := DB.UpdateSQL(CONSTANT.WebinarsTable, map[string]string{"webinar_id": r.FormValue("webinar_id")}, webinar)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if body.Status == CONSTANT.WebinarCancelledByAdmin {
		// send notification to clients who have booked the webinar
		webinarBookings, status, ok := DB.SelectProcess("select * from " + CONSTANT.WebinarsBookTable + " where webinar_id = '" + r.FormValue("webinar_id") + "' and status = " + CONSTANT.WebinarBooked + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		for _, booking := range webinarBookings {
			// send notification to client
			UTIL.SendNotification(
				CONSTANT.AdminCancelledWebinarHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.AdminCancelledWebinarContent,
					map[string]string{
						"###webinarName###": body.Title,
					},
				),
				booking["user_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body.Date, body.Time).String(),
				CONSTANT.NotificationSent,
				booking["order_id"],
				"",
			)
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PreSignedS3URLToUploadWebinar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.EventS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
