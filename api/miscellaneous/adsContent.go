package miscellaneous

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
	"time"
)

// ListMood godoc
// @Tags Miscellaneous
// @Summary Get all moods
// @Router /adsContent [get]
// @Security JWTAuth
// @Produce json
// @Success 200
func AdsContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	adsContent, status, ok := DB.SelectSQL(CONSTANT.AdsContentTable, []string{"title", "target", "image"}, map[string]string{"status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, adsCont := range adsContent {
		url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, adsCont["image"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
		adsCont["image"] = endPointURL
	}

	response["ads"] = adsContent
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetDocumentList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get moods
	getDocuments, status, ok := DB.SelectSQL(CONSTANT.CounsellorDocumentListTable, []string{"*"}, map[string]string{"status": "1"})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, document := range getDocuments {
		url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, document["document"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
		document["document"] = endPointURL
	}

	response["documents"] = getDocuments
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetCounsellorRecordFromMainCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	category, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordFormCategoryTable + " where status = 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	emotionalState, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordFormEmotionalState + " where status = 1")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["main_category"] = category
	response["emotional_state"] = emotionalState
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func GetCounsellorRecordFromSubCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	subCategory, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordFormSubCategoryTable+" where status = 1 and category_id = ? ", r.FormValue("category_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["sub_category"] = subCategory
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func GetCounsellorClientRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var lastClient []map[string]string

	var status string
	var ok bool

	// if len(r.FormValue("counsellor_id")) != 0 {
	// 	lastClient, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ? and counsellor_id = ?  order by session_date desc", r.FormValue("client_id"), r.FormValue("counsellor_id"))
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// } else {
	// 	lastClient, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc limit 5", r.FormValue("client_id"))
	// 	if !ok {
	// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 		return
	// 	}
	// }

	lastClient, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc limit 10", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["lastest_record"] = lastClient
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CheckCounsellorClientRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	if len(r.FormValue("appointment_id")) != 0 {
		// get last cleints
		lastClient, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(lastClient) == 0 {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
			return
		}

	} else {
		// get last cleints
		lastClient, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id"), "client_id": r.FormValue("client_id"), "session_date": r.FormValue("date")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(lastClient) == 0 {
			UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
			return
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
			return
		}
	}

}

func CheckGetCounsellorClientRecordForNewest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var clients []map[string]string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// newVersionClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ?  order by created_at desc limit 5", r.FormValue("appointment_id"))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	newVersionClientRecordForm, status, ok := DB.SelectSQL(CONSTANT.CounsellorRecordsFormLastestVersionTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(newVersionClientRecordForm) > 0 {
		// get client details
		clients, status, ok = DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + newVersionClientRecordForm[0]["client_id"] + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		countOldClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc", newVersionClientRecordForm[0]["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		countNewClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and status = '2' order by modified_at desc", newVersionClientRecordForm[0]["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalCount := 0

		if len(countOldClient) > 0 {
			total, _ := strconv.Atoi(countOldClient[0]["ctn"])
			totalCount += total
		}

		if len(countNewClient) > 0 {
			total, _ := strconv.Atoi(countNewClient[0]["ctn"])
			totalCount += total
		}

		clients[0]["total_sessions"] = strconv.Itoa(totalCount)

	} else {
		newVersionClientRecordForm, status, ok = DB.SelectSQL(CONSTANT.CounsellorRecordsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(newVersionClientRecordForm) > 0 {

			if len(newVersionClientRecordForm[0]["client_id"]) != 0 {
				clients, status, ok = DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + newVersionClientRecordForm[0]["client_id"] + "'")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				countOldClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc", newVersionClientRecordForm[0]["client_id"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				countNewClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and status = '2' order by modified_at desc", newVersionClientRecordForm[0]["client_id"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				totalCount := 0

				if len(countOldClient) > 0 {
					total, _ := strconv.Atoi(countOldClient[0]["ctn"])
					totalCount += total
				}

				if len(countNewClient) > 0 {
					total, _ := strconv.Atoi(countNewClient[0]["ctn"])
					totalCount += total
				}

				clients[0]["total_sessions"] = strconv.Itoa(totalCount)
			}
		}

	}

	if len(clients) > 0 {
		response["client_details"] = clients[0]
	} else {
		response["client_details"] = map[string]string{}
	}

	response["client_record"] = newVersionClientRecordForm
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetCounsellorClientRecordForNewestVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	var totalLastClientRecord []map[string]string

	var status string
	var ok bool

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + r.FormValue("client_id") + "'")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	countOldClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsTable+" where client_id = ?  order by session_date desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	countNewClient, status, ok := DB.SelectProcess("select count(record_id) as ctn from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and status = '2' order by modified_at desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	totalCount := 0

	if len(countOldClient) > 0 {
		total, _ := strconv.Atoi(countOldClient[0]["ctn"])
		totalCount += total
	}

	if len(countNewClient) > 0 {
		total, _ := strconv.Atoi(countNewClient[0]["ctn"])
		totalCount += total
	}

	clients[0]["total_sessions"] = strconv.Itoa(totalCount)

	SQLQueryLastest := "select * from " + CONSTANT.CounsellorRecordsFormLastestVersionTable + " where client_id = ? and status = 2"

	sortByLastest := " modified_at " // default ordering by modified_at
	orderByLastest := " desc "
	if strings.EqualFold(r.FormValue("sort_by"), "1") {
		sortByLastest = " modified_at "
	}
	if strings.EqualFold(r.FormValue("sort_by"), "2") {
		sortByLastest = " mental_health_scale "
	}
	if strings.EqualFold(r.FormValue("order_by"), "1") {
		orderByLastest = " asc "
	}
	SQLQueryLastest += " order by " + sortByLastest + orderByLastest

	SQLQuery := "select * from " + CONSTANT.CounsellorRecordsTable + " where client_id = ?"

	sortBy := " created_at " // default ordering by created_at
	orderBy := " desc "
	if strings.EqualFold(r.FormValue("sort_by"), "1") {
		sortBy = " created_at "
	}
	if strings.EqualFold(r.FormValue("sort_by"), "2") {
		sortBy = " mental_health "
	}
	if strings.EqualFold(r.FormValue("order_by"), "1") {
		orderBy = " asc "
	}
	SQLQuery += " order by " + sortBy + orderBy

	lastClient, status, ok := DB.SelectProcess(SQLQueryLastest+" limit "+strconv.Itoa(CONSTANT.CounsellorsRecordFormPerPage), r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(lastClient) == 0 {
		lastClient, status, ok = DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsRecordFormPerPage), r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalLastClientRecord = append(totalLastClientRecord, lastClient...)

	} else if len(lastClient) < 10 {

		totalLastClientRecord = append(totalLastClientRecord, lastClient...)

		oldLastClient, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(10-len(lastClient)), r.FormValue("client_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalLastClientRecord = append(totalLastClientRecord, oldLastClient...)

	} else {

		totalLastClientRecord = append(totalLastClientRecord, lastClient...)
	}

	appointmentIDs := UTIL.ExtractValuesFromArrayMap(totalLastClientRecord, "appointment_id")

	// get client details
	appointmentDetails, status, ok := DB.SelectProcess("select time, appointment_id from " + CONSTANT.AppointmentsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "') ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	inpersonAppointmentDetails, status, ok := DB.SelectProcess("select time, appointment_id from " + CONSTANT.InPersonAppointmentsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "') ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	appointmentDetails = append(appointmentDetails, inpersonAppointmentDetails...)

	appointmentDetailsMap := UTIL.ConvertMapToKeyMap(appointmentDetails, "appointment_id")

	for index, value := range totalLastClientRecord {
		if _, ok := appointmentDetailsMap[value["appointment_id"]]; ok {

			totalLastClientRecord[index]["appointment_time"] = appointmentDetailsMap[value["appointment_id"]]["time"]
		} else {
			totalLastClientRecord[index]["appointment_time"] = ""
		}
	}

	response["client_details"] = clients[0]
	// response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(countClient[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
	response["lastest_record"] = totalLastClientRecord
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetLastHistoryRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})
	var sessionDetails = make(map[string]string)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	totalSessionTakenByTheClientWithSameTherapist := 0

	newVersionClientRecordForm, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and counsellor_id = ? and no_show = 0 and incomplete_session = 0 and status = '2'  order by modified_at desc limit 3", r.FormValue("client_id"), r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	prevClientRecordFormTotal, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.CounsellorRecordsTable+" where client_id = ? and counsellor_id = ? and noshow = 0 order by session_date desc", r.FormValue("client_id"), r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	newClientRecordFormTotal, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and counsellor_id = ? and no_show = 0 and incomplete_session = 0 and status = '2' order by modified_at desc", r.FormValue("client_id"), r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	prevTotalSessions, _ := strconv.Atoi(prevClientRecordFormTotal[0]["ctn"])

	newTotalSessions, _ := strconv.Atoi(newClientRecordFormTotal[0]["ctn"])

	totalSessionTakenByTheClientWithSameTherapist = prevTotalSessions + newTotalSessions

	if len(newVersionClientRecordForm) == 0 {
		newVersionClientRecordForm, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where client_id = ? and counsellor_id = ? order by created_at desc limit 1", r.FormValue("client_id"), r.FormValue("counsellor_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		sessionDetails["previousTotalSession"] = strconv.Itoa(totalSessionTakenByTheClientWithSameTherapist)
		sessionDetails["totalSession"] = "0"
		sessionDetails["takenSession"] = "0"
	} else {
		newVersionClientRecordFormTotalSession, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and counsellor_id = ? and no_show = 0 and incomplete_session = 0 and status = '2' order by modified_at desc", r.FormValue("client_id"), r.FormValue("counsellor_id"))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalSessionByClient := 0
		totaltakenSessionByClient := 0

		if len(newVersionClientRecordFormTotalSession) > 0 {
			if newVersionClientRecordFormTotalSession[0]["goals_achieved"] != "Yes" {
				totalSessions, _ := strconv.Atoi(newVersionClientRecordFormTotalSession[0]["total_session_needed"])
				takenSessions, _ := strconv.Atoi(newVersionClientRecordFormTotalSession[0]["taken_sessions"])

				totalSessionByClient = totalSessionByClient + totalSessions
				totaltakenSessionByClient = totaltakenSessionByClient + takenSessions + 1
			}
		}

		sessionDetails["previousTotalSession"] = strconv.Itoa(totalSessionTakenByTheClientWithSameTherapist)
		sessionDetails["totalSession"] = strconv.Itoa(totalSessionByClient)
		sessionDetails["takenSession"] = strconv.Itoa(totaltakenSessionByClient)
	}

	if len(newVersionClientRecordForm) > 0 {
		if newVersionClientRecordForm[0]["session_mode"] == "In-Person" {
			appointmentDetails, status, ok := DB.SelectProcess("select time, appointment_id from " + CONSTANT.InPersonAppointmentsTable + " where appointment_id = '" + newVersionClientRecordForm[0]["appointment_id"] + "' ")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			newVersionClientRecordForm[0]["appointment_time"] = appointmentDetails[0]["time"]
		} else {
			appointmentDetails, status, ok := DB.SelectProcess("select time, appointment_id from " + CONSTANT.AppointmentsTable + " where appointment_id = '" + newVersionClientRecordForm[0]["appointment_id"] + "' ")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			newVersionClientRecordForm[0]["appointment_time"] = appointmentDetails[0]["time"]
		}
	}

	checkPointsList, status, ok := DB.SelectProcess("select * from " + CONSTANT.TherapistCheckMHScalePointsTable + " ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["check_points_list"] = checkPointsList
	response["last_client_record"] = newVersionClientRecordForm
	response["session_details"] = sessionDetails
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorClientRecordForNewestVersion(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var response = make(map[string]interface{})

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorRecordFormNewVersionAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := Model.NewVersionOfCounsellorRecordFormRequest{}

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

	var noshow, incompleteSession, mentalHealth string
	var newVersionClientRecordForm []map[string]string

	if len(body.MentalHealth) > 0 {
		mentalHealth = body.MentalHealth
	} else {
		mentalHealth = "0"
	}

	mentalStatus, _ := strconv.Atoi(mentalHealth)

	if len(body.NoShow) > 0 {
		noshow = body.NoShow
	} else {
		noshow = "0"
	}

	if len(body.IncompleteSession) > 0 {
		incompleteSession = body.IncompleteSession
	} else {
		incompleteSession = "0"
	}

	// if len(body["session_for"]) > 0 {
	// 	sessionFor = body["session_for"]
	// 	sessionType = body["session_type"]
	// } else {
	// 	sessionFor = ""
	// 	sessionType = ""
	// }

	if len(body.RecordID) > 0 {

		newVersionClientRecordForm, _, _ = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where record_id = ? and status = '2'", body.RecordID)

		if len(newVersionClientRecordForm) > 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Record already exists for the given record ID", CONSTANT.ShowDialog, response)
			return
		}

		counsellorRecord := map[string]string{}
		counsellorRecord["session_type"] = body.SessionType
		counsellorRecord["family_relation"] = body.FamilyRelation
		counsellorRecord["out_time"] = body.OutTime
		counsellorRecord["no_show"] = noshow
		counsellorRecord["incomplete_session"] = incompleteSession
		counsellorRecord["presenting_concers"] = body.PresentingConcerns
		counsellorRecord["mental_health"] = mentalHealth
		counsellorRecord["mental_health_check"] = body.MentalHealthCheck
		counsellorRecord["downgrading_high_risk_case"] = body.DowngradingHighRiskCase
		counsellorRecord["is_clinical_psychologist_required_reason"] = body.IsClinicalPsychologistRequiredReason
		counsellorRecord["psychiatric_intervention_required_reason"] = body.PsychiatricInterventionRequiredReason
		counsellorRecord["category"] = body.Category
		counsellorRecord["sub_category"] = body.SubCategory
		counsellorRecord["emotional_state"] = body.EmotionalState
		counsellorRecord["therapy_notes"] = body.TherapyNotes
		counsellorRecord["goals_achieved"] = body.GoalsAchieved
		counsellorRecord["goals_achieved_reason"] = body.GoalsAchievedReason
		counsellorRecord["total_session_needed"] = body.TotalSessionNeeded
		counsellorRecord["taken_sessions"] = body.TakenSessions
		counsellorRecord["next_session_plan"] = body.NextSessionPlan
		counsellorRecord["next_follow_up_date"] = body.NextFollowUpDate
		counsellorRecord["client_notes"] = body.ClientNotes
		counsellorRecord["self_work_material"] = body.SelfWorkMaterial
		counsellorRecord["assessment"] = body.Assessment
		counsellorRecord["status"] = CONSTANT.CounsellorRecordFormCompleted
		counsellorRecord["modified_at"] = UTIL.GetCurrentTime().String()

		// _, status, ok := DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.CounsellorRecordDigits, counsellorRecord, "record_id")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		DB.UpdateSQL(CONSTANT.CounsellorRecordsFormLastestVersionTable,
			map[string]string{
				"record_id": body.RecordID,
			},
			counsellorRecord,
		)

		newVersionClientRecordForm, _, _ = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where record_id = ? ", body.RecordID)

	} else {

		checkCounsellorRecordFormExistsOrNot, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? ", body.AppointmentID)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(checkCounsellorRecordFormExistsOrNot) > 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Record already exists for the given appointment ID", CONSTANT.ShowDialog, response)
			return
		}

		counsellorRecord := map[string]string{}
		counsellorRecord["counsellor_id"] = body.CounsellorID
		counsellorRecord["client_id"] = body.ClientID
		counsellorRecord["appointment_id"] = body.AppointmentID
		counsellorRecord["session_for"] = body.SessionFor
		counsellorRecord["family_relation"] = body.FamilyRelation
		counsellorRecord["session_type"] = body.SessionType
		counsellorRecord["session_mode"] = body.SessionMood
		counsellorRecord["session_date"] = body.SessionDate
		counsellorRecord["in_time"] = body.InTime
		counsellorRecord["out_time"] = body.OutTime
		counsellorRecord["no_show"] = noshow
		counsellorRecord["incomplete_session"] = body.IncompleteSession
		counsellorRecord["presenting_concers"] = body.PresentingConcerns
		counsellorRecord["mental_health"] = mentalHealth
		counsellorRecord["mental_health_check"] = body.MentalHealthCheck
		counsellorRecord["downgrading_high_risk_case"] = body.DowngradingHighRiskCase
		counsellorRecord["is_clinical_psychologist_required_reason"] = body.IsClinicalPsychologistRequiredReason
		counsellorRecord["psychiatric_intervention_required_reason"] = body.PsychiatricInterventionRequiredReason
		counsellorRecord["category"] = body.Category
		counsellorRecord["sub_category"] = body.SubCategory
		counsellorRecord["emotional_state"] = body.EmotionalState
		counsellorRecord["therapy_notes"] = body.TherapyNotes
		counsellorRecord["goals_achieved"] = body.GoalsAchieved
		counsellorRecord["goals_achieved_reason"] = body.GoalsAchievedReason
		counsellorRecord["total_session_needed"] = body.TotalSessionNeeded
		counsellorRecord["taken_sessions"] = body.TakenSessions
		counsellorRecord["next_session_plan"] = body.NextSessionPlan
		counsellorRecord["next_follow_up_date"] = body.NextFollowUpDate
		counsellorRecord["client_notes"] = body.ClientNotes
		counsellorRecord["self_work_material"] = body.SelfWorkMaterial
		counsellorRecord["assessment"] = body.Assessment
		counsellorRecord["status"] = CONSTANT.CounsellorRecordFormCompleted
		counsellorRecord["created_at"] = UTIL.GetCurrentTime().String()
		counsellorRecord["modified_at"] = UTIL.GetCurrentTime().String()

		recordID, status, ok := DB.InsertWithUniqueID(CONSTANT.CounsellorRecordsFormLastestVersionTable, CONSTANT.CounsellorRecordDigits, counsellorRecord, "record_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		newVersionClientRecordForm, status, ok = DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where record_id = ? ", recordID)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "email"}, map[string]string{"counsellor_id": newVersionClientRecordForm[0]["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "email"}, map[string]string{"therapist_id": newVersionClientRecordForm[0]["counsellor_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + newVersionClientRecordForm[0]["client_id"] + "'")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	sessionDate := UTIL.BuildOnlyDate(newVersionClientRecordForm[0]["session_date"])

	nextFollowDate := ""

	if len(newVersionClientRecordForm[0]["next_follow_up_date"]) != 0 {
		nextFollowDate = UTIL.BuildOnlyDate(newVersionClientRecordForm[0]["next_follow_up_date"])
	}

	if noshow == "1" {
		noshow = "Yes"
	} else {
		noshow = "No"
	}

	data := Model.EmailDataForCounsellorRecordForLastestVersion{
		TherapistName:                         counsellor[0]["first_name"] + " " + counsellor[0]["last_name"],
		Client_First_Name:                     clients[0]["first_name"],
		Client_Last_Name:                      clients[0]["last_name"],
		Client_Gender:                         clients[0]["gender"],
		Client_Age:                            clients[0]["age"],
		SessionFor:                            newVersionClientRecordForm[0]["session_for"],
		SessionType:                           newVersionClientRecordForm[0]["session_type"],
		FamilyRelation:                        newVersionClientRecordForm[0]["family_relation"],
		SessionMode:                           newVersionClientRecordForm[0]["session_mode"],
		SessionDate:                           sessionDate,
		InTime:                                newVersionClientRecordForm[0]["in_time"],
		OutTime:                               newVersionClientRecordForm[0]["out_time"],
		NoShow:                                noshow,
		PresentingConcerns:                    body.PresentingConcerns,
		MentalHealthScale:                     mentalHealth,
		MentalHealthCheck:                     body.MentalHealthCheck,
		DowngradingHighRiskCase:               body.DowngradingHighRiskCase,
		IsClinicalPsychologistRequiredReason:  body.IsClinicalPsychologistRequiredReason,
		PsychiatricInterventionRequiredReason: body.PsychiatricInterventionRequiredReason,
		Category:                              body.Category,
		SubCategory:                           body.SubCategory,
		EmotionalState:                        body.EmotionalState,
		GoalsAchieved:                         body.GoalsAchieved,
		GoalsAchievedReason:                   body.GoalsAchievedReason,
		TotalSessionNeeded:                    body.TotalSessionNeeded,
		TakenSessions:                         newVersionClientRecordForm[0]["taken_sessions"],
		TherapyNotes:                          body.TherapyNotes,
		NextSessionPlan:                       body.NextSessionPlan,
		NextFollowDate:                        nextFollowDate,
		ClientNotes:                           body.ClientNotes,
		Assessment:                            body.Assessment,
		SelfWorkMaterial:                      body.SelfWorkMaterial,
	}

	filepath := "htmlfile/counsellorRecordNewVersion.html"

	emailbody := UTIL.GetHTMLTemplateForCounsellorRecordNewestVersion(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorRecordForClientTitle,
		emailbody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	if mentalStatus > 7 || len(body.PsychiatricInterventionRequiredReason) != 0 || len(body.IsClinicalPsychologistRequiredReason) != 0 {

		ccAddressess := []string{CONFIG.QCEmailID1, CONFIG.QCEmailID3}

		UTIL.SendEmailWithCCBB(
			CONSTANT.CounselloRecordClientEmergencyCaseTitle,
			emailbody,
			CONFIG.QCEmailID2,
			ccAddressess,
			CONSTANT.InstantSendEmailMessage,
		)
	}

	message, message1, message2, subjectLine := "", "", "", ""

	if len(body.NextFollowUpDate) != 0 {
		// 15 min push notification before appointment start

		if newVersionClientRecordForm[0]["session_mode"] == "In-Person" {
			UTIL.SendNotification(
				CONSTANT.ClientInPersonAppointmentFollowUpSessionToSuggestedReminderClientHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.ClientAppointmentFollowUpToSuggestedRemiderClientContent,
					map[string]string{
						"###therapist_name###": counsellor[0]["first_name"],
						"###follow_up_date###": nextFollowDate,
					},
				),
				newVersionClientRecordForm[0]["client_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body.NextFollowUpDate, "26").Add(-24*time.Hour).UTC().String(),
				CONSTANT.NotificationInProgress,
				newVersionClientRecordForm[0]["client_id"],
				"",
			)

			subjectLine = CONSTANT.CounsellorInPersonAppointmentDocumentForClientTitle
		} else {
			UTIL.SendNotification(
				CONSTANT.ClientVirtualAppointmentFollowUpSessionToSuggestedReminderClientHeading,
				UTIL.ReplaceNotificationContentInString(
					CONSTANT.ClientAppointmentFollowUpToSuggestedRemiderClientContent,
					map[string]string{
						"###therapist_name###": counsellor[0]["first_name"],
						"###follow_up_date###": nextFollowDate,
					},
				),
				newVersionClientRecordForm[0]["client_id"],
				CONSTANT.ClientType,
				UTIL.BuildDateTime(body.NextFollowUpDate, "26").Add(-24*time.Hour).UTC().String(),
				CONSTANT.NotificationInProgress,
				newVersionClientRecordForm[0]["client_id"],
				"",
			)

			subjectLine = CONSTANT.CounsellorVirtualAppointmentDocumentForClientTitle
		}

		message1 = UTIL.ReplaceNotificationContentInString(
			CONSTANT.TherapistAttachDocumentsFollowDateFooterClientBody,
			map[string]string{
				"###followupdate###": nextFollowDate,
			},
		)
	}

	if len(newVersionClientRecordForm[0]["client_id"]) != 0 {

		if len(body.SelfWorkMaterial) != 0 || len(body.ClientNotes) != 0 || len(body.Assessment) != 0 {

			client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "last_name", "email"}, map[string]string{"client_id": newVersionClientRecordForm[0]["client_id"]})

			var listofDocuments []Model.DocumentList
			if len(body.SelfWorkMaterial) != 0 {
				clientDocuments := strings.Split(body.SelfWorkMaterial, ",")

				for _, value := range clientDocuments {
					document, _, _ := DB.SelectSQL(CONSTANT.CounsellorDocumentListTable, []string{"document"}, map[string]string{"document_name": value})
					docu := Model.DocumentList{
						DocumentName: value,
						DocumentLink: document[0]["document"],
					}
					listofDocuments = append(listofDocuments, docu)
				}
			}

			// re := regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

			// message := re.ReplaceAllString(body["client_notes"], `<a href="$0">$0</a>`)

			var emaildata Model.EmailBodyMessageModelWithDocu

			if len(body.SelfWorkMaterial) != 0 {
				message = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachDocumentsWithClientBody,
					map[string]string{
						"###TherapistName###": counsellor[0]["first_name"],
						"###Date###":          sessionDate,
					},
				)

			} else {
				message = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachDocumentsWithOutClientBody,
					map[string]string{
						"###TherapistName###": counsellor[0]["first_name"],
						"###Date###":          sessionDate,
					},
				)
			}

			if len(body.Assessment) != 0 {
				message2 = UTIL.ReplaceNotificationContentInString(
					CONSTANT.TherapistAttachAssessmentWithClientBody,
					map[string]string{
						"###assessment_name###": body.Assessment,
					},
				)
			}

			if len(body.ClientNotes) != 0 {
				emaildata = Model.EmailBodyMessageModelWithDocu{
					Name:     client[0]["first_name"],
					Message:  message,
					Message1: "Your therapist has suggested the following guidelines:",
					Message2: body.ClientNotes,
					Message3: message2,
					Message4: message1,
				}
			} else {
				emaildata = Model.EmailBodyMessageModelWithDocu{
					Name:     client[0]["first_name"],
					Message:  message,
					Message1: "Your therapist has suggested the following guidelines:",
					Message2: "Notes not provided.",
					Message3: body.Links,
					Message4: message1,
				}
			}

			filepath_text := "htmlfile/emailbodywithassessment.html"

			emailBy := UTIL.GetHTMLTemplateForWithDocument(emaildata, filepath_text)

			UTIL.SendEmailWithDocument(client[0]["email"], emailBy, UTIL.ReplaceNotificationContentInString(
				subjectLine,
				map[string]string{
					"###date###": sessionDate,
				},
			), listofDocuments)
		}

	}

	// UTIL.SendEmail(
	// 	CONSTANT.CounsellorRecordForClientTitle,
	// 	emailbody,
	// 	client[0]["email"],
	// 	CONSTANT.InstantSendEmailMessage,
	// )

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
