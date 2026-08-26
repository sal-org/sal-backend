package client

import (
	"fmt"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	"slices"
	"strconv"

	UTIL "salbackend/util"
	"strings"
)

// ListSearch godoc
// @Tags Client Search
// @Summary Get counsellor/listener/therapist list with search filters
// @Router /client/search [get]
// @Param type query string false "Counsellor(1)/Listener(2)/Therapist(4) or dont send if all"
// @Param topic query string false "anxiety/anger/stress/depression/relationship/parenting/grief/motivation/life/others - send selected topic id"
// @Param language query string false "english/hindi/tamil/telugu/kannada/bengali/malayalam/marathi/gujarati/punjabi - send selected language id"
// @Param date query string false "Available on date (2020-02-27)"
// @Param price query string false "Price range - 100,200 (min,max)"
// @Param experience query string false "Experience range - 0,30 (min,max)"
// @Param sort_by query string false "Sort by - 1(price), 2(rating), 3(age_group)"
// @Param order_by query string false "Order by - 1(asc), 2(desc) - should be sent along with sort_by"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	clientID, ok := UTIL.Required(r.FormValue("client_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	// check if client_id exists
	if !DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": clientID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	var SQLQuery string // counsellorSQLQuery, listenerSQLQuery,
	args := []any{}
	// counsellorArgs := []interface{}{}
	// listenerArgs := []interface{}{}
	therapistArgs := []any{}

	// build counsellor query
	// counsellorSQLQuery = "select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable

	// wheres := []string{}
	// if len(r.FormValue("topic")) > 0 { // get counsellors with specified topic
	// 	wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
	// 	counsellorArgs = append(counsellorArgs, r.FormValue("topic"))
	// }
	// if len(r.FormValue("language")) > 0 { // get counsellors with specified language
	// 	wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
	// 	counsellorArgs = append(counsellorArgs, r.FormValue("language"))
	// }
	// if len(r.FormValue("date")) > 0 { // get counsellors available in specified date
	// 	wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
	// 	counsellorArgs = append(counsellorArgs, r.FormValue("date"))
	// }
	// if len(r.FormValue("price")) > 0 { // get counsellors available in specified price range
	// 	prices := strings.Split(r.FormValue("price"), ",") // min,max price range
	// 	wheres = append(wheres, " price >= ? and price <= ? ")
	// 	counsellorArgs = append(counsellorArgs, prices[0], prices[1])
	// }
	// if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
	// 	// Param experience query string false "Experience range - 0,30 (min,max)"
	// 	experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
	// 	wheres = append(wheres, " experience >= ? and experience <= ? ")
	// 	min, _ := strconv.ParseFloat(experiences[0], 64)
	// 	max, _ := strconv.ParseFloat(experiences[1], 64)
	// 	counsellorArgs = append(counsellorArgs, min, max)
	// }
	// wheres = append(wheres, " status = "+CONSTANT.CounsellorActive+" and corporate_therpist != '2' ") // only active counsellors
	// counsellorSQLQuery += " where " + strings.Join(wheres, " and ")

	// build listener query
	// listenerSQLQuery = "select listener_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, 0 as price, 0 as multiple_sessions, occupation, age_group, '' as therapeutic_approach, about, " + CONSTANT.ListenerType + " as type, slot_type from " + CONSTANT.ListenersTable

	// wheres = []string{}
	// if len(r.FormValue("topic")) > 0 { // get listeners with specified topic
	// 	wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
	// 	listenerArgs = append(listenerArgs, r.FormValue("topic"))
	// }
	// if len(r.FormValue("language")) > 0 { // get listeners with specified language
	// 	wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
	// 	listenerArgs = append(listenerArgs, r.FormValue("language"))
	// }
	// if len(r.FormValue("date")) > 0 { // get listeners available in specified date
	// 	wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
	// 	listenerArgs = append(listenerArgs, r.FormValue("date"))
	// }
	// /*if len(r.FormValue("age_group")) > 0 {
	// 	wheres = append(wheres, " listener_id in (select listener_id from "+CONSTANT.ListenersTable+" where age_group = ?) ")
	// 	listenerArgs = append(listenerArgs, r.FormValue("age_group"))
	// }*/
	// wheres = append(wheres, " status = "+CONSTANT.ListenerActive+" ") // only active listeners
	// listenerSQLQuery += " where " + strings.Join(wheres, " and ")

	// build therapist query
	// therapistSQLQuery = "select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable

	wheres := []string{}
	if len(r.FormValue("topic")) > 0 { // get therapists with specified topic
		wheres = append(wheres, " t.therapist_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		therapistArgs = append(therapistArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get therapists with specified language
		wheres = append(wheres, " t.therapist_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		therapistArgs = append(therapistArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get therapists available in specified date
		wheres = append(wheres, " t.therapist_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
		therapistArgs = append(therapistArgs, r.FormValue("date"))
	}
	if len(r.FormValue("name")) > 0 {
		wheres = append(wheres, " (t.first_name like '%%"+r.FormValue("name")+"%%' or t.last_name like '%%"+r.FormValue("name")+"%%') ")
	}
	if len(r.FormValue("price")) > 0 { // get therapists available in specified price range
		prices := strings.Split(r.FormValue("price"), ",") // min,max price range
		wheres = append(wheres, " t.price >= ? and t.price <= ? ")
		therapistArgs = append(therapistArgs, prices[0], prices[1])
	}
	if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
		experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
		wheres = append(wheres, " t.experience >= ? and t.experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		therapistArgs = append(therapistArgs, min, max)
	}
	wheres = append(wheres, " t.status = "+CONSTANT.TherapistActive+" and t.corporate_therpist != '2' and t.price >= 100") // only active therapists
	// therapistSQLQuery += " where " + strings.Join(wheres, " and ")

	slotColumns := make([]string, 48)

	for i := 0; i < 48; i++ {
		slotColumns[i] = fmt.Sprintf("a.`%d`", i)
	}

	where := strings.Join(wheres, " and ")

	query := fmt.Sprintf(`SELECT t.therapist_id, t.first_name, t.last_name, t.pronoun, t.total_rating, t.average_rating, t.photo, t.price, t.multiple_sessions, t.education, t.experience, t.therapeutic_approach, t.about, t.corporate_therpist, t.in_house_therapist, a.available, a.date, t.slot_type, t.video, %s FROM therapists t LEFT JOIN slots a ON a.counsellor_id = t.therapist_id AND a.date >= CURDATE() WHERE `+where+` ORDER BY t.therapist_id, a.date`, strings.Join(slotColumns, ", "))

	SQLQuery = query
	args = append(args, therapistArgs...)

	// // get counsellors|listeners|therapists
	// counsellors, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsListPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient), args...)
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get counsellors|therapists
	counsellor, status, ok := DB.SelectProcess(SQLQuery, args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var counsellorModel []MODEL.TherapistResponse

	counsellorModel = UTIL.BuildTherapistAvailability(counsellor, CONSTANT.TherapistType)

	counsellors := UTIL.RotateInHouseTherapist(r.FormValue("page"), CONSTANT.CounsellorsListPerPageClient, counsellorModel)

	response["counsellors"] = counsellors
	response["counsellors_count"] = strconv.Itoa(len(counsellorModel))
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(strconv.Itoa(len(counsellorModel)), CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ListSearchForCorporate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	var SQLQuery string
	args := []any{}

	therapistArgs := []any{}

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	clientID, ok := UTIL.Required(r.FormValue("client_id"), "Client ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	// check domain exists or not
	ok = DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
		return
	}

	// // build counsellor query
	// counsellorSQLQuery = "select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable
	// wheres := []string{}
	// if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
	// 	// Param experience query string false "Experience range - 0,30 (min,max)"
	// 	experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
	// 	wheres = append(wheres, " experience >= ? and experience <= ? ")
	// 	min, _ := strconv.ParseFloat(experiences[0], 64)
	// 	max, _ := strconv.ParseFloat(experiences[1], 64)
	// 	counsellorArgs = append(counsellorArgs, min, max)
	// }
	// wheres = append(wheres, " status = "+CONSTANT.CounsellorActive+" and corporate_therpist != 0 ") // only active counsellors
	// counsellorSQLQuery += " where " + strings.Join(wheres, " and ")

	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"email", "location"}, map[string]string{"client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	domainName := strings.Split(client[0]["email"], "@")

	checkClientRecordForm, status, ok := DB.SelectProcess("select taken_sessions, total_session_needed, mental_health, counsellor_id from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where client_id = ? and no_show = 0 and incomplete_session = 0 and status = '2' order by modified_at desc limit 5", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// isAddressExist := DB.CheckIfExists(CONSTANT.CorporatePartnersAddressTable, map[string]string{"address": client[0]["location"], "status": "1"})

	clientAddress, status, ok := DB.SelectProcess("select * from "+CONSTANT.CorporatePartnersAddressTable+" where address = ? and status = 1 order by created_at desc", client[0]["location"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if domainName[1] == "ageasfederal.com" {

		//	build therapist in person query
		// therapistSQLQuery = "select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable
		wheres := []string{}
		if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
			experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
			wheres = append(wheres, " t.experience >= ? and t.experience <= ? ")
			min, _ := strconv.ParseFloat(experiences[0], 64)
			max, _ := strconv.ParseFloat(experiences[1], 64)
			therapistArgs = append(therapistArgs, min, max)
		}

		slotColumns := make([]string, 48)

		for i := 0; i < 48; i++ {
			slotColumns[i] = fmt.Sprintf("a.`%d`", i)
		}

		wheres = append(wheres, " t.status = "+CONSTANT.TherapistActive+" and t.in_house_therapist = 1") // only active therapists
		where := strings.Join(wheres, " and ")

		query := fmt.Sprintf(`SELECT t.therapist_id, t.first_name, t.last_name, t.pronoun, t.total_rating, t.average_rating, t.photo, t.price, t.multiple_sessions, t.education, t.experience, t.therapeutic_approach, t.about, t.corporate_therpist, t.in_house_therapist, a.available, a.date, t.slot_type, t.video, %s FROM therapists t LEFT JOIN slots a ON a.counsellor_id = t.therapist_id AND a.date >= CURDATE() WHERE `+where+` ORDER BY t.therapist_id, a.date`, strings.Join(slotColumns, ", "))

		// } else { // union if all needed
		SQLQuery = query
		args = append(args, therapistArgs...)
		// }

		// sortBy := " average_rating " // default ordering by rating
		// orderBy := " desc "
		// if strings.EqualFold(r.FormValue("order_by"), "1") {
		// 	orderBy = " asc "
		// }
		// SQLQuery += " order by " + sortBy + orderBy

		// get counsellors|therapists
		counsellor, status, ok := DB.SelectProcess(SQLQuery, args...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		var counsellorModel []MODEL.TherapistResponse

		counsellorModel = UTIL.BuildTherapistAvailability(counsellor, CONSTANT.TherapistType)

		counsellors := UTIL.RotateInHouseTherapist(r.FormValue("page"), CONSTANT.CounsellorsListPerPageClient, counsellorModel)

		response["counsellors"] = counsellors
		response["counsellors_count"] = strconv.Itoa(len(counsellorModel))
		response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(strconv.Itoa(len(counsellorModel)), CONSTANT.CounsellorsListPerPageClient))
		response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	} else {

		counsellorID := []string{}

		if len(clientAddress) > 0 {

			counsellorConnectWithCompanyLocation, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonCounsellorConnectWithCorporateTable+" where partner_name = ? and partner_location = ? and status = '1' order by created_at desc", clientAddress[0]["partner_name"], clientAddress[0]["address"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			for _, value := range counsellorConnectWithCompanyLocation {

				if len(checkClientRecordForm) > 0 {

					if value["counsellor_id"] == checkClientRecordForm[0]["counsellor_id"] {
						mentalHealthScore, _ := strconv.Atoi(checkClientRecordForm[0]["mental_health"])

						if mentalHealthScore < 8 {
							counsellorID = append(counsellorID, value["counsellor_id"])
						}
					} else {
						counsellorID = append(counsellorID, value["counsellor_id"])
					}

				} else {
					counsellorID = append(counsellorID, value["counsellor_id"])
				}
			}

		}

		// build therapist query
		// therapistSQLQuery = "select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about, corporate_therpist, in_house_therapist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable
		wheres := []string{}
		if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
			experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
			wheres = append(wheres, " t.experience >= ? and t.experience <= ? ")
			min, _ := strconv.ParseFloat(experiences[0], 64)
			max, _ := strconv.ParseFloat(experiences[1], 64)
			therapistArgs = append(therapistArgs, min, max)
		}

		if len(counsellorID) > 0 {
			wheres = append(wheres, " t.therapist_id not in ('"+strings.Join(counsellorID, "', '")+"') ")
			// therapistArgs = append(therapistArgs, counsellorID)
		}

		slotColumns := make([]string, 48)

		for i := 0; i < 48; i++ {
			slotColumns[i] = fmt.Sprintf("a.`%d`", i)
		}

		wheres = append(wheres, " t.status = "+CONSTANT.TherapistActive+" and t.corporate_therpist != 0 ") // only active therapists
		// therapistSQLQuery += " where " + strings.Join(wheres, " and ")

		where := strings.Join(wheres, " and ")

		query := fmt.Sprintf(`SELECT t.therapist_id, t.first_name, t.last_name, t.pronoun, t.total_rating, t.average_rating, t.photo, t.price, t.multiple_sessions, t.education, t.experience, t.therapeutic_approach, t.about, t.corporate_therpist, t.in_house_therapist, a.available, a.date, t.slot_type, t.video, %s FROM therapists t LEFT JOIN slots a ON a.counsellor_id = t.therapist_id AND a.date >= CURDATE() WHERE `+where+` ORDER BY t.therapist_id, a.date`, strings.Join(slotColumns, ", "))

		// } else { // union if all needed
		SQLQuery = query
		args = append(args, therapistArgs...)
		// }

		// sortBy := " average_rating " // default ordering by rating
		// orderBy := " desc "
		// if strings.EqualFold(r.FormValue("order_by"), "1") {
		// 	orderBy = " asc "
		// }
		// SQLQuery += " order by " + sortBy + orderBy

		// // get counsellors|therapists
		// counsellors, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsListPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient), args...)
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		// get counsellors|therapists
		counsellor, status, ok := DB.SelectProcess(SQLQuery, args...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		var counsellorModel []MODEL.TherapistResponse

		counsellorModel = UTIL.BuildTherapistAvailability(counsellor, CONSTANT.TherapistType)

		counsellors := UTIL.RotateInHouseTherapist(r.FormValue("page"), CONSTANT.CounsellorsListPerPageClient, counsellorModel)

		response["counsellors"] = counsellors
		response["counsellors_count"] = strconv.Itoa(len(counsellorModel))
		response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(strconv.Itoa(len(counsellorModel)), CONSTANT.CounsellorsListPerPageClient))
		response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ListSearchForCorporateForAvailableSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	slotColumns := make([]string, 48)

	for i := 0; i < 48; i++ {
		slotColumns[i] = fmt.Sprintf("a.`%d`", i)
	}

	query := fmt.Sprintf(`SELECT t.therapist_id, t.first_name, t.last_name, t.pronoun, t.total_rating, t.average_rating, t.photo, t.price, t.multiple_sessions, t.education, t.experience, t.therapeutic_approach, t.about, t.corporate_therpist, t.in_house_therapist, t.slot_type, t.video, a.available, a.date, %s FROM therapists t LEFT JOIN slots a ON a.counsellor_id = t.therapist_id AND a.date >= CURDATE() WHERE t.status = 1 ORDER BY t.therapist_id, a.date`, strings.Join(slotColumns, ", "))

	// get counsellors|therapists
	counsellor, status, ok := DB.SelectProcess(query)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	var counsellorModel []MODEL.TherapistResponse

	counsellorModel = UTIL.BuildTherapistAvailability(counsellor, CONSTANT.TherapistType)

	counsellors := UTIL.RotateInHouseTherapist(r.FormValue("page"), CONSTANT.CounsellorsListPerPageClient, counsellorModel)

	response["counsellors"] = counsellors
	response["counsellors_count"] = strconv.Itoa(len(counsellorModel))
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(strconv.Itoa(len(counsellorModel)), CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ListSearchForCorporateInPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)
	var counsellorsCount []map[string]string
	var counsellors []map[string]string
	// filteredCounsellorSlots := map[string][]map[string]string{}
	// filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}
	var inPersonConnect []map[string]string
	var slots []map[string]string
	var nextSlots []map[string]string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	clientID, ok := UTIL.Required(r.FormValue("client_id"), "Client ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, clientID, CONSTANT.ShowDialog, response)
		return
	}

	// check domain exists or not
	ok = DB.CheckIfExists(CONSTANT.ClientsTable, map[string]string{"client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientCorEmailInvalid, CONSTANT.ShowDialog, response)
		return
	}

	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"email"}, map[string]string{"client_id": r.FormValue("client_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	emailIDAccess := []string{CONFIG.CorporateInpersonAccessEmailId1, CONFIG.CorporateInpersonAccessEmailId2, CONFIG.CorporateInpersonAccessEmailId3}

	if slices.Contains(emailIDAccess, client[0]["email"]) {

		partnerName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"*"}, map[string]string{"status": "1"})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		inPersonConn, status, ok := DB.SelectSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, []string{"*"}, map[string]string{"status": "1"})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIds := UTIL.ExtractValuesFromArrayMap(inPersonConn, "counsellor_id")

		mYSQL := "(select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0 ) union (select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0)"

		counsellors, status, ok = DB.SelectProcess(mYSQL + " limit " + strconv.Itoa(CONSTANT.CounsellorsListPerPageClient) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// extract counsellors|therapists ids
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

		for _, value := range inPersonConn {
			if slices.Contains(counsellorIDs, value["counsellor_id"]) {
				inPersonConnect = append(inPersonConnect, value)
			}
		}

		companyName := UTIL.ExtractValuesFromArrayMap(partnerName, "partner_name")

		companyLocation := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "partner_location")

		// get counsellors|therapists slots
		slots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name in ('" + strings.Join(companyName, "','") + "') and company_location in ('" + strings.Join(companyLocation, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		nextSlots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name in ('" + strings.Join(companyName, "','") + "') and company_location in ('" + strings.Join(companyLocation, "','") + "') and date > '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' and date < '" + UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02") + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// // get counsellors|therapists slots
		// slots, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and company_name in ('" + strings.Join(counsellorName, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		// // group counsellors|therapists slots
		// counsellorSlots := UTIL.ConvertArrayMapToKeyMapArray(slots, "counsellor_id")
		// // filteredCounsellorSlots := map[string][]map[string]string{}
		// // filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}

		// // var nextSlot []map[string]string
		// for counsellorID, counsellorSlot := range counsellorSlots {
		// 	filteredCounsellorSlots[counsellorID] = UTIL.FilterAvailableForInPersonSlots(counsellorSlot)
		// 	if len(filteredCounsellorSlots[counsellorID]) == 0 {
		// 		nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and available = 1 and company_name in ('"+strings.Join(counsellorName, "','")+"') and  date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", counsellorID)
		// 		if !ok {
		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 			return
		// 		}
		// 		filteredCounsellorSlotsNextAvaliable[counsellorID] = UTIL.FilterAvailableForInPersonSlots(nextSlots)
		// 	}
		// }

		// get counsellors|therapists count
		counsellorsCount, status, ok = DB.SelectProcess("select count(*) as ctn from (" + mYSQL + ") as a")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

	} else {

		domainName := strings.Split(client[0]["email"], "@")

		partnerName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"*"}, map[string]string{"domain": domainName[1]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(partnerName) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientInPersonAppointmentModelNotExistMessage, CONSTANT.ShowDialog, response)
			return
		}

		inPersonConn, status, ok := DB.SelectSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, []string{"*"}, map[string]string{"partner_name": partnerName[0]["partner_name"], "status": "1"})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIds := UTIL.ExtractValuesFromArrayMap(inPersonConn, "counsellor_id")

		mYSQL := "(select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0 ) union (select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0)"

		counsellors, status, ok = DB.SelectProcess(mYSQL + " limit " + strconv.Itoa(CONSTANT.CounsellorsListPerPageClient) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient))
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// extract counsellors|therapists ids
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

		for _, value := range inPersonConn {
			if slices.Contains(counsellorIDs, value["counsellor_id"]) {
				inPersonConnect = append(inPersonConnect, value)
			}
		}

		companyLocation := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "partner_location")

		// get counsellors|therapists slots
		slots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name = '" + partnerName[0]["partner_name"] + "' and company_location in ('" + strings.Join(companyLocation, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		nextSlots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name = '" + partnerName[0]["partner_name"] + "' and company_location in ('" + strings.Join(companyLocation, "','") + "') and date > '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' and date < '" + UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02") + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellors|therapists count
		counsellorsCount, status, ok = DB.SelectProcess("select count(*) as ctn from (" + mYSQL + ") as a")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	// for _, counsellor := range counsellors {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellor["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	counsellor["photo"] = endPointURL
	// }

	response["counsellors"] = counsellors
	response["slots"] = UTIL.FilterAvailableForInPersonSlots(slots)
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["next_available"] = UTIL.FilterAvailableForInPersonSlots(nextSlots)
	response["location"] = inPersonConnect
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// // for test after testing we can remove this

// func ListSearchForCorporateInPersonDuplication(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var response = make(map[string]any)
// 	var counsellorsCount []map[string]string
// 	var counsellors []map[string]string
// 	// filteredCounsellorSlots := map[string][]map[string]string{}
// 	// filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}
// 	var inPersonConnect []map[string]string
// 	var slots []map[string]string
// 	var nextSlots []map[string]string

// 	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"email"}, map[string]string{"client_id": r.FormValue("client_id")})
// 	if !ok {
// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 		return
// 	}

// 	if client[0]["email"] == "shivam.tiwari@clovemind.com" {

// 		partnerName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"*"}, map[string]string{"status": "1"})
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		inPersonConnect, status, ok = DB.SelectSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, []string{"*"}, map[string]string{"status": "3"})
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		counsellorIds := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "counsellor_id")

// 		mYSQL := "(select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0 ) union (select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0)"

// 		counsellors, status, ok = DB.SelectProcess(mYSQL + " limit " + strconv.Itoa(CONSTANT.CounsellorsListPerPageClient) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient))
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		// extract counsellors|therapists ids
// 		counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

// 		companyName := UTIL.ExtractValuesFromArrayMap(partnerName, "partner_name")

// 		companyLocation := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "partner_location")

// 		// get counsellors|therapists slots
// 		slots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and company_name in ('" + strings.Join(companyName, "','") + "') and available = 1 and company_location in ('" + strings.Join(companyLocation, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		nextSlots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and company_name in ('" + strings.Join(companyName, "','") + "') and available = 1 and company_location in ('" + strings.Join(companyLocation, "','") + "') and date > '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' and date < '" + UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02") + "'")
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		// // get counsellors|therapists slots
// 		// slots, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and company_name in ('" + strings.Join(counsellorName, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
// 		// if !ok {
// 		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 		// 	return
// 		// }

// 		// // group counsellors|therapists slots
// 		// counsellorSlots := UTIL.ConvertArrayMapToKeyMapArray(slots, "counsellor_id")
// 		// // filteredCounsellorSlots := map[string][]map[string]string{}
// 		// // filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}

// 		// // var nextSlot []map[string]string
// 		// for counsellorID, counsellorSlot := range counsellorSlots {
// 		// 	filteredCounsellorSlots[counsellorID] = UTIL.FilterAvailableForInPersonSlots(counsellorSlot)
// 		// 	if len(filteredCounsellorSlots[counsellorID]) == 0 {
// 		// 		nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and available = 1 and company_name in ('"+strings.Join(counsellorName, "','")+"') and  date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", counsellorID)
// 		// 		if !ok {
// 		// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 		// 			return
// 		// 		}
// 		// 		filteredCounsellorSlotsNextAvaliable[counsellorID] = UTIL.FilterAvailableForInPersonSlots(nextSlots)
// 		// 	}
// 		// }

// 		// get counsellors|therapists count
// 		counsellorsCount, status, ok = DB.SelectProcess("select count(*) as ctn from (" + mYSQL + ") as a")
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 	} else {

// 		domainName := strings.Split(client[0]["email"], "@")

// 		partnerName, status, ok := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"*"}, map[string]string{"domain": domainName[1]})
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		inPersonConnect, status, ok = DB.SelectSQL(CONSTANT.InPersonCounsellorConnectWithCorporateTable, []string{"*"}, map[string]string{"partner_name": partnerName[0]["partner_name"], "status": "3"})
// 		if !ok {
// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			return
// 		}

// 		if len(inPersonConnect) != 0 {
// 			counsellorIds := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "counsellor_id")

// 			mYSQL := "(select counsellor_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions , education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0 ) union (select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIds, "','") + "') and status = " + CONSTANT.CounsellorActive + " and corporate_therpist != 0)"

// 			counsellors, status, ok = DB.SelectProcess(mYSQL + " limit " + strconv.Itoa(CONSTANT.CounsellorsListPerPageClient) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient))
// 			if !ok {
// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 				return
// 			}

// 			// extract counsellors|therapists ids
// 			counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

// 			companyLocation := UTIL.ExtractValuesFromArrayMap(inPersonConnect, "partner_location")

// 			// get counsellors|therapists slots
// 			slots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name = '" + partnerName[0]["partner_name"] + "' and company_location in ('" + strings.Join(companyLocation, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
// 			if !ok {
// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 				return
// 			}

// 			nextSlots, status, ok = DB.SelectProcess("select * from " + CONSTANT.InPersonSLotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and available = 1 and company_name = '" + partnerName[0]["partner_name"] + "' and company_location in ('" + strings.Join(companyLocation, "','") + "') and date > '" + UTIL.GetCurrentTime().Format("2006-01-02") + "' and date < '" + UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02") + "'")
// 			if !ok {
// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 				return
// 			}

// 			// nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and available = 1 and company_name = '"+partnerName[0]["partner_name"]+"' and  date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", counsellorID)
// 			// if !ok {
// 			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			// 	return
// 			// }

// 			// group counsellors|therapists slots
// 			// counsellorSlots := UTIL.ConvertArrayMapToKeyMapArray(slots, "company_location")
// 			// // filteredCounsellorSlots := map[string][]map[string]string{}
// 			// // filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}

// 			// // var nextSlot []map[string]string
// 			// for counsellorID, counsellorSlot := range counsellorSlots {
// 			// 	fmt.Println(counsellorSlot)
// 			// 	filteredCounsellorSlots[counsellorID] = UTIL.FilterAvailableForInPersonSlots(counsellorSlot)
// 			// 	if len(filteredCounsellorSlots[counsellorID]) == 0 {
// 			// 		nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.InPersonSLotsTable+" where counsellor_id = ? and available = 1 and company_name = '"+partnerName[0]["partner_name"]+"' and  date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", counsellorID)
// 			// 		if !ok {
// 			// 			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 			// 			return
// 			// 		}
// 			// 		filteredCounsellorSlotsNextAvaliable[counsellorID] = UTIL.FilterAvailableForInPersonSlots(nextSlots)
// 			// 	}
// 			// }

// 			// get counsellors|therapists count
// 			counsellorsCount, status, ok = DB.SelectProcess("select count(*) as ctn from (" + mYSQL + ") as a")
// 			if !ok {
// 				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
// 				return
// 			}
// 		}

// 	}

// 	for _, counsellor := range counsellors {
// 		url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellor["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
// 		_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
// 		counsellor["photo"] = endPointURL
// 	}

// 	response["counsellors"] = counsellors
// 	response["slots"] = UTIL.FilterAvailableForInPersonSlots(slots)
// 	response["counsellors_count"] = counsellorsCount[0]["ctn"]
// 	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
// 	response["media_url"] = CONFIG.MediaURL
// 	response["next_available"] = UTIL.FilterAvailableForInPersonSlots(nextSlots)
// 	response["location"] = inPersonConnect
// 	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
// }
