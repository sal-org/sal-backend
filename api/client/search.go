package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
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

	var response = make(map[string]interface{})

	var SQLQuery, counsellorSQLQuery, listenerSQLQuery, therapistSQLQuery string
	args := []interface{}{}
	counsellorArgs := []interface{}{}
	listenerArgs := []interface{}{}
	therapistArgs := []interface{}{}

	// build counsellor query
	counsellorSQLQuery = "select counsellor_id as id, first_name, last_name, total_rating, average_rating, photo, price, multiple_sessions , education, experience, about, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable

	wheres := []string{}
	if len(r.FormValue("topic")) > 0 { // get counsellors with specified topic
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		counsellorArgs = append(counsellorArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get counsellors with specified language
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		counsellorArgs = append(counsellorArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get counsellors available in specified date
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
		counsellorArgs = append(counsellorArgs, r.FormValue("date"))
	}
	if len(r.FormValue("price")) > 0 { // get counsellors available in specified price range
		prices := strings.Split(r.FormValue("price"), ",") // min,max price range
		wheres = append(wheres, " price >= ? and price <= ? ")
		counsellorArgs = append(counsellorArgs, prices[0], prices[1])
	}
	if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
		// Param experience query string false "Experience range - 0,30 (min,max)"
		experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
		wheres = append(wheres, " experience >= ? and experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		counsellorArgs = append(counsellorArgs, min, max)
	}
	wheres = append(wheres, " status = "+CONSTANT.CounsellorActive+" and corporate_therpist != '2' ") // only active counsellors
	counsellorSQLQuery += " where " + strings.Join(wheres, " and ")

	// build listener query
	listenerSQLQuery = "select listener_id as id, first_name, last_name, total_rating, average_rating, photo, 0 as price, 0 as multiple_sessions, occupation, age_group, about, " + CONSTANT.ListenerType + " as type, slot_type from " + CONSTANT.ListenersTable

	wheres = []string{}
	if len(r.FormValue("topic")) > 0 { // get listeners with specified topic
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		listenerArgs = append(listenerArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get listeners with specified language
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		listenerArgs = append(listenerArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get listeners available in specified date
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
		listenerArgs = append(listenerArgs, r.FormValue("date"))
	}
	/*if len(r.FormValue("age_group")) > 0 {
		wheres = append(wheres, " listener_id in (select listener_id from "+CONSTANT.ListenersTable+" where age_group = ?) ")
		listenerArgs = append(listenerArgs, r.FormValue("age_group"))
	}*/
	wheres = append(wheres, " status = "+CONSTANT.ListenerActive+" ") // only active listeners
	listenerSQLQuery += " where " + strings.Join(wheres, " and ")

	// build therapist query
	therapistSQLQuery = "select therapist_id as id, first_name, last_name, total_rating, average_rating, photo, price, multiple_sessions, education, experience, about, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable

	wheres = []string{}
	if len(r.FormValue("topic")) > 0 { // get therapists with specified topic
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		therapistArgs = append(therapistArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get therapists with specified language
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		therapistArgs = append(therapistArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get therapists available in specified date
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
		therapistArgs = append(therapistArgs, r.FormValue("date"))
	}
	if len(r.FormValue("price")) > 0 { // get therapists available in specified price range
		prices := strings.Split(r.FormValue("price"), ",") // min,max price range
		wheres = append(wheres, " price >= ? and price <= ? ")
		therapistArgs = append(therapistArgs, prices[0], prices[1])
	}
	if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
		experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
		wheres = append(wheres, " experience >= ? and experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		therapistArgs = append(therapistArgs, min, max)
	}
	wheres = append(wheres, " status = "+CONSTANT.TherapistActive+" and corporate_therpist != '2' ") // only active therapists
	therapistSQLQuery += " where " + strings.Join(wheres, " and ")

	if len(r.FormValue("type")) > 0 { // get only certain types
		types := strings.Split(r.FormValue("type"), ",")
		for _, t := range types {
			if strings.EqualFold(t, CONSTANT.CounsellorType) {
				SQLQuery = counsellorSQLQuery
				args = counsellorArgs
			} else if strings.EqualFold(t, CONSTANT.ListenerType) {
				SQLQuery = listenerSQLQuery
				args = listenerArgs
			} else if strings.EqualFold(t, CONSTANT.TherapistType) {
				SQLQuery = therapistSQLQuery
				args = therapistArgs
			}
		}
	} else if len(r.FormValue("price")) > 0 {
		SQLQuery = " ( " + counsellorSQLQuery + " ) union ( " + therapistSQLQuery + " ) "
		args = append(args, counsellorArgs...)
		args = append(args, therapistArgs...)

	} else { // union if all needed
		SQLQuery = " ( " + counsellorSQLQuery + " ) union ( " + listenerSQLQuery + " ) union ( " + therapistSQLQuery + " ) "
		args = append(args, counsellorArgs...)
		args = append(args, listenerArgs...)
		args = append(args, therapistArgs...)
	}

	sortBy := " average_rating " // default ordering by rating
	orderBy := " desc "
	if strings.EqualFold(r.FormValue("sort_by"), "1") {
		sortBy = " price "
	}
	if strings.EqualFold(r.FormValue("sort_by"), "3") {
		sortBy = " age_group "
	}
	if strings.EqualFold(r.FormValue("order_by"), "1") {
		orderBy = " asc "
	}
	SQLQuery += " order by " + sortBy + orderBy

	// get counsellors|listeners|therapists
	counsellors, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsListPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient), args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// extract counsellors|listeners|therapists ids
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

	// get counsellors|listeners|therapists slots
	slots, status, ok := DB.SelectProcess("select * from " + CONSTANT.SlotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// group counsellors|listeners|therapists slots
	counsellorSlots := UTIL.ConvertArrayMapToKeyMapArray(slots, "counsellor_id")
	filteredCounsellorSlots := map[string][]map[string]string{}
	filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}
	var counsellorFirstList []map[string]string
	var counsellorSecondList []map[string]string
	// var nextSlot []map[string]string
	for counsellorID, counsellorSlot := range counsellorSlots {
		filteredCounsellorSlots[counsellorID] = UTIL.FilterAvailableSlots(counsellorSlot)
		if len(filteredCounsellorSlots[counsellorID]) == 0 {
			nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and available = 1  and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", counsellorID)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			filteredCounsellorSlotsNextAvaliable[counsellorID] = UTIL.FilterAvailableSlots(nextSlots)
		}
	}

	// get counsellors|listeners|therapists count
	counsellorsCount, status, ok := DB.SelectProcess("select count(*) as ctn from ("+SQLQuery+") as a", args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for key, value := range counsellors {

		slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and date = '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and available = 1 order by date asc", value["id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		filteredCounsellorSlot := map[string][]map[string]string{}

		filteredCounsellorSlot[value["id"]] = UTIL.FilterAvailableSlots(slots)
		if len(filteredCounsellorSlot[value["id"]]) != 0 {
			counsellorFirstList = append(counsellorFirstList, counsellors[key])
		} else {
			counsellorSecondList = append(counsellorSecondList, counsellors[key])
		}

	}

	counsellorFirstList = append(counsellorFirstList, counsellorSecondList...)

	// fmt.Println(counsellors)

	response["counsellors"] = counsellorFirstList
	response["slots"] = filteredCounsellorSlots
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURL
	response["next_available"] = filteredCounsellorSlotsNextAvaliable
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ListSearchForCorporate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var SQLQuery, counsellorSQLQuery, therapistSQLQuery string
	args := []interface{}{}

	counsellorArgs := []interface{}{}
	therapistArgs := []interface{}{}

	// build counsellor query
	counsellorSQLQuery = "select counsellor_id as id, first_name, last_name, total_rating, average_rating, photo, price, multiple_sessions , education, experience, about,corporate_therpist, " + CONSTANT.CounsellorType + " as type, slot_type from " + CONSTANT.CounsellorsTable
	wheres := []string{}
	if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
		// Param experience query string false "Experience range - 0,30 (min,max)"
		experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
		wheres = append(wheres, " experience >= ? and experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		counsellorArgs = append(counsellorArgs, min, max)
	}
	wheres = append(wheres, " status = "+CONSTANT.CounsellorActive+" and corporate_therpist != 0 ") // only active counsellors
	counsellorSQLQuery += " where " + strings.Join(wheres, " and ")

	// build therapist query
	therapistSQLQuery = "select therapist_id as id, first_name, last_name, total_rating, average_rating, photo, price, multiple_sessions, education, experience, about,corporate_therpist, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable
	wheres = []string{}
	if len(r.FormValue("experience")) > 0 { // get counsellors available in specified price range
		experiences := strings.Split(r.FormValue("experience"), ",") // min,max price range
		wheres = append(wheres, " experience >= ? and experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		therapistArgs = append(therapistArgs, min, max)
	}
	wheres = append(wheres, " status = "+CONSTANT.TherapistActive+" and corporate_therpist != 0 ") // only active therapists
	therapistSQLQuery += " where " + strings.Join(wheres, " and ")

	// } else { // union if all needed
	SQLQuery = " ( " + counsellorSQLQuery + " ) union ( " + therapistSQLQuery + " ) "
	args = append(args, counsellorArgs...)
	args = append(args, therapistArgs...)
	// }

	// sortBy := " average_rating " // default ordering by rating
	// orderBy := " desc "
	// if strings.EqualFold(r.FormValue("order_by"), "1") {
	// 	orderBy = " asc "
	// }
	// SQLQuery += " order by " + sortBy + orderBy

	// get counsellors|therapists
	counsellors, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsListPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsListPerPageClient), args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// extract counsellors|therapists ids
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellors, "id")

	// get counsellors|therapists slots
	slots, status, ok := DB.SelectProcess("select * from " + CONSTANT.SlotsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "') and date = '" + UTIL.GetCurrentTime().Format("2006-01-02") + "'")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// group counsellors|therapists slots
	counsellorSlots := UTIL.ConvertArrayMapToKeyMapArray(slots, "counsellor_id")
	filteredCounsellorSlots := map[string][]map[string]string{}
	filteredCounsellorSlotsNextAvaliable := map[string][]map[string]string{}
	// var nextSlot []map[string]string
	for counsellorID, counsellorSlot := range counsellorSlots {
		filteredCounsellorSlots[counsellorID] = UTIL.FilterAvailableSlots(counsellorSlot)
		if len(filteredCounsellorSlots[counsellorID]) == 0 {
			nextSlots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and available = 1 and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", counsellorID)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}
			filteredCounsellorSlotsNextAvaliable[counsellorID] = UTIL.FilterAvailableSlots(nextSlots)
		}
	}

	// get counsellors|therapists count
	counsellorsCount, status, ok := DB.SelectProcess("select count(*) as ctn from ("+SQLQuery+") as a", args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// fmt.Println(counsellors)

	response["counsellors"] = counsellors
	response["slots"] = filteredCounsellorSlots
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURL
	response["next_available"] = filteredCounsellorSlotsNextAvaliable
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
