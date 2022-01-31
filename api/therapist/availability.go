package therapist

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"
	"time"

	UTIL "salbackend/util"
)

// AvailabilityGet godoc
// @Tags Therapist Availability
// @Summary Get therapist availability hours
// @Router /therapist/availability [get]
// @Param therapist_id query string true "Therapist ID to get availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get therapist availability hours
	availability, status, ok := DB.SelectProcess("select * from "+CONSTANT.SchedulesTable+" where counsellor_id = ? order by weekday", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["availability"] = availability
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AvailabilityUpdate godoc
// @Tags Therapist Availability
// @Summary Update therapist availability hours
// @Router /therapist/availability [put]
// @Param therapist_id query string true "Therapist ID to update availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBodyInListMap(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": r.FormValue("therapists_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if counsellor[0]["price"] == "0" || counsellor[0]["multiple_sessions"] == "0" {
		UTIL.SetReponse(w, status, "Please update your price per sessions in 'My Profile' first.", CONSTANT.ShowDialog, response)
		return
	}

	for _, day := range body {
		if strings.EqualFold(day["status"], "0") {
			// delete schedule
			DB.DeleteSQL(CONSTANT.SchedulesTable, map[string]string{"id": day["id"]})
		} else {
			if len(day["id"]) > 0 {
				DB.UpdateSQL(CONSTANT.SchedulesTable, map[string]string{"id": day["id"]}, day)
			} else {
				// newly added schedule
				DB.InsertSQL(CONSTANT.SchedulesTable, day)
			}
		}
	}

	// get all dates for counsellor and group by weekday
	datesByWeekdays := map[int][]string{}
	availabileDates, status, ok := DB.SelectProcess("select date from "+CONSTANT.SlotsTable+" where counsellor_id = ? and `date` >= ?", r.FormValue("therapist_id"), UTIL.GetCurrentTime().AddDate(0, 0, -1).String())
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	dates := UTIL.ExtractValuesFromArrayMap(availabileDates, "date")
	// grouping by weekday
	for _, date := range dates {
		t, _ := time.Parse("2006-01-02", date)
		datesByWeekdays[int(t.Weekday())] = append(datesByWeekdays[int(t.Weekday())], date)
	}

	// group weekdays
	days := map[string]map[string]string{}
	for _, day := range body {
		if days[day["weekday"]] == nil {
			days[day["weekday"]] = map[string]string{
				"weekday": day["weekday"],
				"0":       "0",
				"1":       "0",
				"2":       "0",
				"3":       "0",
				"4":       "0",
				"5":       "0",
				"6":       "0",
				"7":       "0",
				"8":       "0",
				"9":       "0",
				"10":      "0",
				"11":      "0",
				"12":      "0",
				"13":      "0",
				"14":      "0",
				"15":      "0",
				"16":      "0",
				"17":      "0",
				"18":      "0",
				"19":      "0",
				"20":      "0",
				"21":      "0",
				"22":      "0",
				"23":      "0",
				"24":      "0",
				"25":      "0",
				"26":      "0",
				"27":      "0",
				"28":      "0",
				"29":      "0",
				"30":      "0",
				"31":      "0",
				"32":      "0",
				"33":      "0",
				"34":      "0",
				"35":      "0",
				"36":      "0",
				"37":      "0",
				"38":      "0",
				"39":      "0",
				"40":      "0",
				"41":      "0",
				"42":      "0",
				"43":      "0",
				"44":      "0",
				"45":      "0",
				"46":      "0",
				"47":      "0",
			}
		}
		if strings.EqualFold(day["status"], "1") && strings.EqualFold(day["availability_status"], "1") {
			for key, value := range day {
				if strings.EqualFold(value, CONSTANT.SlotAvailable) {
					days[day["weekday"]][key] = CONSTANT.SlotAvailable
				}
			}
		}
	}
	// calculate availability for a weekday
	for weekday, day := range days {
		availability := "0"
		if strings.EqualFold(day["status"], "1") && strings.EqualFold(day["availability_status"], "1") {
			for key, value := range day {
				_, err := strconv.Atoi(key)
				if err == nil && strings.EqualFold(value, CONSTANT.SlotAvailable) {
					availability = "1"
					break
				}
			}
		}
		days[weekday]["available"] = availability
	}

	// update weekday availability to respective dates
	// will run for 30 days * 24 * 2 hours = 1440 times - TODO needs to be optimised
	for _, day := range days { // 7 times
		weekday, _ := strconv.Atoi(day["weekday"])
		for _, date := range datesByWeekdays[weekday] { // respective weekday dates i.e., 4-5 times
			for key, value := range day { // 48 times
				if strings.EqualFold(day["status"], "0") || strings.EqualFold(day["availability_status"], "0") {
					value = CONSTANT.SlotUnavailable // not available
				}
				DB.ExecuteSQL("update "+CONSTANT.SlotsTable+" set `"+key+"` = "+value+" where counsellor_id = ? and date = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", r.FormValue("therapist_id"), date) // dont update already booked slots
			}
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
