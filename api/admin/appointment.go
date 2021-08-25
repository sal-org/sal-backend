package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	_ "salbackend/model"
	UTIL "salbackend/util"
)

func AppointmentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointments
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "state":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "client_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " client_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of appointments
	appointmentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.AppointmentsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["appointments"] = appointments
	response["appointments_count"] = appointmentsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(appointmentsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
