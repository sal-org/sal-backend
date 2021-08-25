package admin

import (
	"encoding/csv"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"time"

	UTIL "salbackend/util"
)

func ReportGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=reports.csv")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired("Bearer " + r.FormValue("access_token")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	writer := csv.NewWriter(w)

	switch r.FormValue("id") {
	case "1":

		writer.Write([]string{"ID", "Quote", "Mood ID"})
		quotes, status, ok := DB.SelectProcess("select * from " + CONSTANT.QuotesTable + " order by id desc ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, quote := range quotes {
			writer.Write([]string{quote["id"], quote["quote"], quote["mood_id"]})
		}
	case "2": // onboarding report

		writer.Write([]string{"First Name", "Last Name", "Gender", "Email", "Phone", "Type", "Created At"})
		startBy, _ := time.Parse("2006-01-02", r.FormValue("start_by"))
		endBy, _ := time.Parse("2006-01-02", r.FormValue("end_by"))
		counsellors, status, ok := DB.SelectProcess("(select first_name, last_name, gender, email, phone, 'Counsellor' as `type`, created_at from " + CONSTANT.CounsellorsTable + " where status = " + CONSTANT.CounsellorActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Listener' as `type`, created_at from " + CONSTANT.ListenersTable + " where status = " + CONSTANT.ListenerActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Therapist' as `type`, created_at from " + CONSTANT.TherapistsTable + " where status = " + CONSTANT.TherapistActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, counsellor := range counsellors {
			writer.Write([]string{counsellor["first_name"], counsellor["last_name"], counsellor["gender"], counsellor["email"], counsellor["phone"], counsellor["type"], counsellor["created_at"]})
		}

	}

	writer.Flush()
}
