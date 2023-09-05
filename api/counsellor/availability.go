package counsellor

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
// @Tags Counsellor Availability
// @Summary Get counsellor availability hours
// @Router /counsellor/availability [get]
// @Param counsellor_id query string true "Counsellor ID to get availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var availabilitMarge []map[string]string

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor availability hours
	availability, status, ok := DB.SelectProcess("select * from "+CONSTANT.SchedulesTable+" where counsellor_id = ? order by weekday", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	availabilitMarge = append(availabilitMarge, availability...)

	// get counsellor availability hours
	availabilityDates, status, ok := DB.SelectProcess("select * from "+CONSTANT.SchedulesDatesTable+" where counsellor_id = ? and dates >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by dates", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	availabilitMarge = append(availabilitMarge, availabilityDates...)

	response["availability"] = availabilitMarge
	// response["availability_dates"] = availabilityDates
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AvailabilityUpdate godoc
// @Tags Counsellor Availability
// @Summary Update counsellor availability hours
// @Router /counsellor/availability [put]
// @Param counsellor_id query string true "Counsellor ID to update availability details"
// @Security JWTAuth
// @Produce json
// @Success 200
func AvailabilityUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var dayInBool = false

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body, ok := UTIL.ReadRequestBodyInListMap(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"*"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if counsellor[0]["corporate_therpist"] != "2" {
		if counsellor[0]["price"] == "0" || counsellor[0]["multiple_sessions"] == "0" {
			UTIL.SetReponse(w, status, "Please update your price per sessions in 'My Profile' first.", CONSTANT.ShowDialog, response)
			return
		}
	}

	for _, day := range body {
		if day["dates"] == "" {
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
		} else {
			dayInBool = true
		}

	}

	if dayInBool {

		for _, day := range body {

			if strings.EqualFold(day["status"], "0") {

				DB.DeleteSQL(CONSTANT.SchedulesDatesTable, map[string]string{"id": day["id"]})

				availabileDates, status, ok := DB.SelectProcess("select date,id from "+CONSTANT.SlotsTable+" where counsellor_id = ? and `date` = ?", r.FormValue("counsellor_id"), day["dates"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				slots := map[string]string{}

				slots["available"] = "0"
				slots["0"] = "0"
				slots["1"] = "0"
				slots["2"] = "0"
				slots["3"] = "0"
				slots["4"] = "0"
				slots["5"] = "0"
				slots["6"] = "0"
				slots["7"] = "0"
				slots["8"] = "0"
				slots["9"] = "0"
				slots["10"] = "0"
				slots["11"] = "0"
				slots["12"] = "0"
				slots["13"] = "0"
				slots["14"] = "0"
				slots["15"] = "0"
				slots["16"] = "0"
				slots["17"] = "0"
				slots["18"] = "0"
				slots["19"] = "0"
				slots["20"] = "0"
				slots["21"] = "0"
				slots["22"] = "0"
				slots["23"] = "0"
				slots["24"] = "0"
				slots["25"] = "0"
				slots["26"] = "0"
				slots["27"] = "0"
				slots["28"] = "0"
				slots["29"] = "0"
				slots["30"] = "0"
				slots["31"] = "0"
				slots["32"] = "0"
				slots["33"] = "0"
				slots["34"] = "0"
				slots["35"] = "0"
				slots["36"] = "0"
				slots["37"] = "0"
				slots["38"] = "0"
				slots["39"] = "0"
				slots["40"] = "0"
				slots["41"] = "0"
				slots["42"] = "0"
				slots["43"] = "0"
				slots["44"] = "0"
				slots["45"] = "0"
				slots["46"] = "0"
				slots["47"] = "0"

				status, ok = DB.UpdateSQL(CONSTANT.SlotsTable, map[string]string{"id": availabileDates[0]["id"]}, slots)
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}
			} else {

				if DB.CheckIfExists(CONSTANT.SlotsTable, map[string]string{"date": day["dates"]}) {
					availabileDates, status, ok := DB.SelectProcess("select date,id from "+CONSTANT.SlotsTable+" where counsellor_id = ? and `date` = ?", r.FormValue("counsellor_id"), day["dates"])
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					slotsDate := map[string]string{}

					slotsDate["availability_status"] = day["availability_status"]
					slotsDate["counsellor_id"] = day["counsellor_id"]
					slotsDate["dates"] = day["dates"]
					slotsDate["format"] = day["format"]
					slotsDate["break"] = day["break"]
					slotsDate["status"] = day["status"]
					slotsDate["0"] = day["0"]
					slotsDate["1"] = day["1"]
					slotsDate["2"] = day["2"]
					slotsDate["3"] = day["3"]
					slotsDate["4"] = day["4"]
					slotsDate["5"] = day["5"]
					slotsDate["6"] = day["6"]
					slotsDate["7"] = day["7"]
					slotsDate["8"] = day["8"]
					slotsDate["9"] = day["9"]
					slotsDate["10"] = day["10"]
					slotsDate["11"] = day["11"]
					slotsDate["12"] = day["12"]
					slotsDate["13"] = day["13"]
					slotsDate["14"] = day["14"]
					slotsDate["15"] = day["15"]
					slotsDate["16"] = day["16"]
					slotsDate["17"] = day["17"]
					slotsDate["18"] = day["18"]
					slotsDate["19"] = day["19"]
					slotsDate["20"] = day["20"]
					slotsDate["21"] = day["21"]
					slotsDate["22"] = day["22"]
					slotsDate["23"] = day["23"]
					slotsDate["24"] = day["24"]
					slotsDate["25"] = day["25"]
					slotsDate["26"] = day["26"]
					slotsDate["27"] = day["27"]
					slotsDate["28"] = day["28"]
					slotsDate["29"] = day["29"]
					slotsDate["30"] = day["30"]
					slotsDate["31"] = day["31"]
					slotsDate["32"] = day["32"]
					slotsDate["33"] = day["33"]
					slotsDate["34"] = day["34"]
					slotsDate["35"] = day["35"]
					slotsDate["36"] = day["36"]
					slotsDate["37"] = day["37"]
					slotsDate["38"] = day["38"]
					slotsDate["39"] = day["39"]
					slotsDate["40"] = day["40"]
					slotsDate["41"] = day["41"]
					slotsDate["42"] = day["42"]
					slotsDate["43"] = day["43"]
					slotsDate["44"] = day["44"]
					slotsDate["45"] = day["45"]
					slotsDate["46"] = day["46"]
					slotsDate["47"] = day["47"]

					if len(day["id"]) > 0 {
						DB.UpdateSQL(CONSTANT.SchedulesDatesTable, map[string]string{"id": day["id"]}, slotsDate)
					} else {
						// newly added schedule
						DB.InsertSQL(CONSTANT.SchedulesDatesTable, slotsDate)
					}

					slots := map[string]string{}

					slots["available"] = day["availability_status"]
					slots["0"] = day["0"]
					slots["1"] = day["1"]
					slots["2"] = day["2"]
					slots["3"] = day["3"]
					slots["4"] = day["4"]
					slots["5"] = day["5"]
					slots["6"] = day["6"]
					slots["7"] = day["7"]
					slots["8"] = day["8"]
					slots["9"] = day["9"]
					slots["10"] = day["10"]
					slots["11"] = day["11"]
					slots["12"] = day["12"]
					slots["13"] = day["13"]
					slots["14"] = day["14"]
					slots["15"] = day["15"]
					slots["16"] = day["16"]
					slots["17"] = day["17"]
					slots["18"] = day["18"]
					slots["19"] = day["19"]
					slots["20"] = day["20"]
					slots["21"] = day["21"]
					slots["22"] = day["22"]
					slots["23"] = day["23"]
					slots["24"] = day["24"]
					slots["25"] = day["25"]
					slots["26"] = day["26"]
					slots["27"] = day["27"]
					slots["28"] = day["28"]
					slots["29"] = day["29"]
					slots["30"] = day["30"]
					slots["31"] = day["31"]
					slots["32"] = day["32"]
					slots["33"] = day["33"]
					slots["34"] = day["34"]
					slots["35"] = day["35"]
					slots["36"] = day["36"]
					slots["37"] = day["37"]
					slots["38"] = day["38"]
					slots["39"] = day["39"]
					slots["40"] = day["40"]
					slots["41"] = day["41"]
					slots["42"] = day["42"]
					slots["43"] = day["43"]
					slots["44"] = day["44"]
					slots["45"] = day["45"]
					slots["46"] = day["46"]
					slots["47"] = day["47"]

					status, ok = DB.UpdateSQL(CONSTANT.SlotsTable, map[string]string{"id": availabileDates[0]["id"]}, slots)
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

				}

			}

		}

	} else {

		// get all dates for counsellor and group by weekday
		datesByWeekdays := map[int][]string{}
		availabileDates, status, ok := DB.SelectProcess("select date from "+CONSTANT.SlotsTable+" where counsellor_id = ? and `date` >= ?", r.FormValue("counsellor_id"), UTIL.GetCurrentTime().AddDate(0, 0, -1).String())
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
					DB.ExecuteSQL("update "+CONSTANT.SlotsTable+" set `"+key+"` = "+value+" where counsellor_id = ? and date = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", r.FormValue("counsellor_id"), date) // dont update already booked slots
				}
			}
		}

	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
