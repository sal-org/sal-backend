package qualitycheck

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	UTIL "salbackend/util"
	"strings"
	"time"
)

func GetAppointmentdetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	response := make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// startBy, _ := time.Parse("2006-01-02", r.FormValue("start_by"))
	// endBy, _ := time.Parse("2006-01-02", r.FormValue("end_by"))

	// get appointment details
	wheres := []string{}
	queryArgs := []interface{}{}

	if len(body["name"]) > 0 { // get counsellors with specified topic
		wheres = append(wheres, " counsellor_name like '%%"+body["name"]+"%%' ")
	}

	if len(body["counsellor_id"]) > 0 {
		wheres = append(wheres, " counsellor_id = ? ")
		queryArgs = append(queryArgs, body["counsellor_id"])
	}

	if len(body["type"]) > 0 { // get only certain types
		wheres = append(wheres, " type = ? ")
		queryArgs = append(queryArgs, body["type"])
	}

	if len(body["rating"]) > 0 {
		wheres = append(wheres, " rating in ('"+body["rating"]+"')")
	}

	if len(body["start_by"]) > 0 {
		startBy, _ := time.Parse("2006-01-02", body["start_by"])
		wheres = append(wheres, "created_at > ?")
		queryArgs = append(queryArgs, startBy.UTC().String())
	}

	if len(body["end_by"]) > 0 {
		endBy, _ := time.Parse("2006-01-02", body["end_by"])
		wheres = append(wheres, "created_at < ?")
		queryArgs = append(queryArgs, endBy.UTC().String())

	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}

	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.QualityCheckDetailsTable+where+" order by created_at desc ", queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// total number of appointment
	appointmentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.QualityCheckDetailsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["appointments"] = appointments
	response["appointments_count"] = appointmentsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func SendEmail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	response := make(map[string]interface{})

	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.QualityCheckEmailRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SendEmailForQuality(
		body["title"],
		body["body"],
		body["email_from"],
		body["email_to"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func VideoCallComment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	response := make(map[string]interface{})

	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.VideoCheckCommentRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable, map[string]string{"appointment_id": body["appointment_id"]}, map[string]string{"comments": body["comments"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
