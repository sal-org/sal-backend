package cron

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
)

// SlotsFill - fill next day slots for counsellor/listener
func SlotsFill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	slotDate := UTIL.GetCurrentTime().AddDate(0, 0, 60)
	schedules, status, ok := DB.SelectProcess("select counsellor_id, `0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`, `22`, `23` from " + CONSTANT.SchedulesTable + " where weekday = " + strconv.Itoa(int(slotDate.Weekday())) + "")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for _, schedule := range schedules {
		schedule["date"] = slotDate.Format("2006-01-02")
		DB.InsertSQL(CONSTANT.SlotsTable, schedule)
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
