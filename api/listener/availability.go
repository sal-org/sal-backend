package listener

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"time"

	UTIL "salbackend/util"
)

// AvailabilityGet godoc
// @Tags Listener Availability
// @Summary Get listener availability hours
// @Router /listener/availability [get]
// @Param listener_id query string true "Listener ID to get availability details"
// @Produce json
// @Success 200
func AvailabilityGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get listener availability hours
	availability, status, ok := DB.SelectProcess("select * from "+CONSTANT.SchedulesTable+" where counsellor_id = ?", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["availability"] = availability
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// AvailabilityUpdate godoc
// @Tags Listener Availability
// @Summary Update listener availability hours
// @Router /listener/availability [put]
// @Param listener_id query string true "Listener ID to update availability details"
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

	for _, day := range body {
		DB.UpdateSQL(CONSTANT.SchedulesTable, map[string]string{"counsellor_id": r.FormValue("listener_id"), "weekday": day["weekday"]}, day)
	}

	// get all dates for listener and group by weekday
	datesByWeekdays := map[int][]string{}
	availabileDates, status, ok := DB.SelectProcess("select date from "+CONSTANT.SlotsTable+" where counsellor_id = ?", r.FormValue("listener_id"))
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

	// update weekday availability to respective dates
	// will run for 30 days * 24 * 2 hours = 1440 times - TODO needs to be optimised
	for _, day := range body { // 7 times
		weekday, _ := strconv.Atoi(day["weekday"])
		for _, date := range datesByWeekdays[weekday] { // respective weekday dates i.e., 4-5 times
			for key, value := range day { // 24 times
				DB.ExecuteSQL("update "+CONSTANT.SlotsTable+" set `"+key+"` = "+value+" where counsellor_id = ? and date = ? and `"+key+"` in ("+CONSTANT.SlotUnavailable+", "+CONSTANT.SlotAvailable+")", r.FormValue("listener_id"), date) // dont update already booked slots
			}
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
