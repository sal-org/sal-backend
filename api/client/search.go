package client

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model" // for error response model
	UTIL "salbackend/util"
	"strings"
)

// ListSearch godoc
// @Tags Client
// @Summary Get counsellor/listener list with search filters
// @Router /client/search [get]
// @Param type query string false "Counsellor(1)/Listener(2) or dont send if both"
// @Param topic query string false "anxiety/anger/stress/depression/relationship/parenting/grief/motivation/life/others - send selected topic id"
// @Param language query string false "english/hindi/tamil/telugu/kannada/bengali/malayalam/marathi/gujarati/punjabi - send selected language id"
// @Param date query string false "Available on date (2020-02-27)"
// @Param time query string false "Available on time (0-23 slots), in UTC, for the selected date"
// @Param price query string false "Price range - 100,200 (min,max)"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func ListSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var SQLQuery, counsellorSQLQuery, listenerSQLQuery string
	args := []interface{}{}
	counsellorArgs := []interface{}{}
	listenerArgs := []interface{}{}

	// build counsellor query
	counsellorSQLQuery = "select counsellor_id as id, first_name, last_name, total_rating, average_rating, photo, price, education, experience, about, " + CONSTANT.CounsellorType + " as type from " + CONSTANT.CounsellorsTable

	wheres := []string{}
	if len(r.FormValue("topic")) > 0 { // get counsellors with specified topic
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		counsellorArgs = append(counsellorArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get counsellors with specified language
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		counsellorArgs = append(counsellorArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get counsellors available in specified date time
		wheres = append(wheres, " counsellor_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and `"+r.FormValue("time")+"` = "+CONSTANT.SlotAvailable+") ")
		counsellorArgs = append(counsellorArgs, r.FormValue("date"))
	}
	if len(r.FormValue("price")) > 0 { // get counsellors available in specified date time
		prices := strings.Split(r.FormValue("price"), ",") // min,max price range
		wheres = append(wheres, " price >= ? and price <= ? ")
		counsellorArgs = append(counsellorArgs, prices[0], prices[1])
	}
	wheres = append(wheres, " status = "+CONSTANT.CounsellorActive+" ") // only active counsellors
	counsellorSQLQuery += " where " + strings.Join(wheres, " and ")

	// build listener query
	listenerSQLQuery = "select listener_id as id, first_name, last_name, total_rating, average_rating, photo, 0 as price, '' as education, '' as experience, '' as about, " + CONSTANT.ListenerType + " as type from " + CONSTANT.ListenersTable

	wheres = []string{}
	if len(r.FormValue("topic")) > 0 { // get listeners with specified topic
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		listenerArgs = append(listenerArgs, r.FormValue("topic"))
	}
	if len(r.FormValue("language")) > 0 { // get listeners with specified language
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		listenerArgs = append(listenerArgs, r.FormValue("language"))
	}
	if len(r.FormValue("date")) > 0 { // get listeners available in specified date time
		wheres = append(wheres, " listener_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and `"+r.FormValue("time")+"` = "+CONSTANT.SlotAvailable+") ")
		listenerArgs = append(listenerArgs, r.FormValue("date"))
	}
	wheres = append(wheres, " status = "+CONSTANT.ListenerActive+" ") // only active listeners
	listenerSQLQuery += " where " + strings.Join(wheres, " and ")

	if len(r.FormValue("type")) > 0 { // get only one type
		if strings.EqualFold(r.FormValue("type"), CONSTANT.CounsellorType) {
			SQLQuery = counsellorSQLQuery
			args = counsellorArgs
		} else if strings.EqualFold(r.FormValue("type"), CONSTANT.ListenerType) {
			SQLQuery = listenerSQLQuery
			args = listenerArgs
		}
	} else { // union if both needed
		SQLQuery = " ( " + counsellorSQLQuery + " ) union ( " + listenerSQLQuery + " ) "
		args = append(args, counsellorArgs...)
		args = append(args, listenerArgs...)
	}

	SQLQuery += " order by average_rating desc " // default ordering by rating

	// get counsellors|listeners
	counsellors, status, ok := DB.SelectProcess(SQLQuery, args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellors"] = counsellors
	response["media_url"] = CONSTANT.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
