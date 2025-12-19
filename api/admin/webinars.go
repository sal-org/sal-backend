package admin

import (
	"net/http"
	"path/filepath"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	_ "salbackend/model"
	UTIL "salbackend/util"
)

func WebinarsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get events
	wheres := []string{}
	queryArgs := []interface{}{}
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
	response["events"] = events
	response["events_count"] = eventsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(eventsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func WebinarsAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add event
	webinar := map[string]string{}
	webinar["counsellor_name"] = body["counsellor_name"]
	webinar["title"] = body["title"]
	webinar["counsellor_photo"] = body["counsellor_photo"]
	webinar["counsellor_qualification"] = body["counsellor_qualification"]
	webinar["counsellor_about"] = body["counsellor_about"]
	webinar["counsellor_experience"] = body["counsellor_experience"]
	webinar["counsellor_rating"] = body["counsellor_rating"]
	webinar["about"] = body["about"]
	webinar["partner_name"] = body["partner_name"]
	webinar["why_attend"] = body["why_attend"]
	webinar["address"] = body["address"]
	webinar["photo"] = body["photo"]
	webinar["background_photo"] = body["background_photo"]
	webinar["date"] = body["date"]
	webinar["time"] = body["time"]
	webinar["duration"] = body["duration"]
	webinar["mode"] = body["mode"]
	webinar["status"] = body["status"]
	webinar["created_by"] = body["created_by"]
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

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add event
	webinar := map[string]string{}
	webinar["counsellor_name"] = body["counsellor_name"]
	webinar["title"] = body["title"]
	webinar["counsellor_photo"] = body["counsellor_photo"]
	webinar["counsellor_qualification"] = body["counsellor_qualification"]
	webinar["counsellor_about"] = body["counsellor_about"]
	webinar["counsellor_experience"] = body["counsellor_experience"]
	webinar["counsellor_rating"] = body["counsellor_rating"]
	webinar["about"] = body["about"]
	webinar["why_attend"] = body["why_attend"]
	webinar["partner_name"] = body["partner_name"]
	webinar["address"] = body["address"]
	webinar["photo"] = body["photo"]
	webinar["background_photo"] = body["background_photo"]
	webinar["date"] = body["date"]
	webinar["time"] = body["time"]
	webinar["duration"] = body["duration"]
	webinar["mode"] = body["mode"]
	webinar["status"] = body["status"]
	status, ok := DB.UpdateSQL(CONSTANT.WebinarsTable, map[string]string{"webinar_id": r.FormValue("webinar_id")}, webinar)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if body["status"] == CONSTANT.WebinarCancelledByAdmin {
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
						"###webinarName###": body["title"],
					},
				),
				booking["user_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body["date"], body["time"]).String(),
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

	var response = make(map[string]interface{})

	url, fileName := UTIL.PreSignedS3URLToUploadPut(CONFIG.S3Bucket, CONSTANT.EventS3Path, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion, filepath.Ext(r.FormValue("fileName")))

	response["file_name"] = fileName
	response["url"] = url
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
