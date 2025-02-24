package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

func ListSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var SQLQuery, therapistSQLQuery string
	args := []interface{}{}
	therapistArgs := []interface{}{}
	counsellorlist := []map[string]interface{}{}

	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// build therapist query
	therapistSQLQuery = "select therapist_id as id, first_name, last_name, pronoun, total_rating, average_rating, photo, price, multiple_sessions, education, experience, therapeutic_approach, about, " + CONSTANT.TherapistType + " as type, slot_type from " + CONSTANT.TherapistsTable

	wheres := []string{}
	if len(body["topic"]) > 0 { // get therapists with specified topic
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.CounsellorTopicsTable+" where topic_id = ?) ")
		therapistArgs = append(therapistArgs, body["topic"])
	}
	if len(body["language"]) > 0 { // get therapists with specified language
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.CounsellorLanguagesTable+" where language_id = ?) ")
		therapistArgs = append(therapistArgs, body["language"])
	}
	if len(body["date"]) > 0 { // get therapists available in specified date
		wheres = append(wheres, " therapist_id in (select counsellor_id from "+CONSTANT.SlotsTable+" where date = ? and available = 1) ")
		therapistArgs = append(therapistArgs, body["date"])
	}
	if len(body["name"]) > 0 {
		wheres = append(wheres, " (first_name like '%%"+body["name"]+"%%' or last_name like '%%"+body["name"]+"%%') ")
	}
	if len(body["price"]) > 0 { // get therapists available in specified price range
		prices := strings.Split(body["price"], ",") // min,max price range
		wheres = append(wheres, " price >= ? and price <= ? ")
		therapistArgs = append(therapistArgs, prices[0], prices[1])
	}
	if len(body["experience"]) > 0 { // get counsellors available in specified price range
		experiences := strings.Split(body["experience"], ",") // min,max price range
		wheres = append(wheres, " experience >= ? and experience <= ? ")
		min, _ := strconv.ParseFloat(experiences[0], 64)
		max, _ := strconv.ParseFloat(experiences[1], 64)
		therapistArgs = append(therapistArgs, min, max)
	}
	wheres = append(wheres, " status = "+CONSTANT.TherapistActive+" and corporate_therpist != '2' ") // only active therapists
	therapistSQLQuery += " where " + strings.Join(wheres, " and ")

	if len(body["type"]) > 0 { // get only certain types
		types := strings.Split(body["type"], ",")
		for _, t := range types {
			if strings.EqualFold(t, CONSTANT.TherapistType) {
				SQLQuery = therapistSQLQuery
				args = therapistArgs
			}
		}
	} else if len(body["price"]) > 0 {
		SQLQuery = " ( " + therapistSQLQuery + " ) "
		args = append(args, therapistArgs...)

	} else { // union if all needed
		SQLQuery = " ( " + therapistSQLQuery + " ) "
		args = append(args, therapistArgs...)
	}

	sortBy := " average_rating " // default ordering by rating
	orderBy := " desc "
	if strings.EqualFold(body["sort_by"], "1") {
		sortBy = " price "
	}
	if strings.EqualFold(body["sort_by"], "3") {
		sortBy = " age_group "
	}
	if strings.EqualFold(body["order_by"], "1") {
		orderBy = " asc "
	}
	SQLQuery += " order by " + sortBy + orderBy

	// get counsellors|listeners|therapists
	counsellors, status, ok := DB.SelectProcess(SQLQuery+" limit "+strconv.Itoa(CONSTANT.CounsellorsListPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(body["page"])-1)*CONSTANT.CounsellorsListPerPageClient), args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellors|listeners|therapists count
	counsellorsCount, status, ok := DB.SelectProcess("select count(*) as ctn from ("+SQLQuery+") as a", args...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, value := range counsellors {

		var counsellor = make(map[string]interface{})

		slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and available = 1  and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' and date < '"+UTIL.GetCurrentTime().AddDate(0, 0, 15).Format("2006-01-02")+"' order by date asc", value["id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(slots) > 0 {
			todaySlots := UTIL.FilterAvailableSlots(slots)

			if len(todaySlots) > 0 {

				counsellor["next_available"] = todaySlots[0]
				counsellor["available"] = "1"
			} else {
				counsellor["next_available"] = todaySlots[0]
				counsellor["available"] = "1"
			}
		} else {
			counsellor["next_available"] = []map[string]string{}
			counsellor["available"] = "0"
		}

		lang := []string{}

		// get therapist languages
		therapistLang, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", value["id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for i := 0; i < len(therapistLang); i++ {
			value := therapistLang[i]["language"]
			lang = append(lang, value)

		}

		expert := []string{}
		// get therapist topics
		topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", value["id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for i := 0; i < len(topics); i++ {
			value := topics[i]["topic"]
			expert = append(expert, value)
		}

		counsellor["id"] = value["id"]
		counsellor["name"] = value["first_name"] + " " + value["last_name"]
		counsellor["pronoun"] = value["pronoun"]
		counsellor["total_rate"] = value["total_rating"]
		counsellor["average_rate"] = value["average_rating"]
		counsellor["photo"] = value["photo"]
		counsellor["education"] = value["education"]
		counsellor["experience"] = value["experience"]
		counsellor["therapeutic_approach"] = value["therapeutic_approach"]
		counsellor["about"] = value["about"]
		counsellor["type"] = value["type"]
		counsellor["expertise"] = expert
		counsellor["languages"] = lang

		counsellorlist = append(counsellorlist, counsellor)

	}

	response["counsellors"] = counsellorlist
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.CounsellorsListPerPageClient))
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
