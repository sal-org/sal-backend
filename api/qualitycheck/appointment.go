package qualitycheck

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
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

	// // total number of appointment
	// appointmentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.QualityCheckDetailsTable+where, queryArgs...)
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	response["appointments"] = appointments
	// response["appointments_count"] = appointmentsCount[0]["ctn"]
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

	// is_valid_email := UTIL.IsValidEmail(body["email_from"])
	// if is_valid_email == "" {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct email id", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// domainName := strings.Split(body["email_from"], "@")
	// if domainName[1] != "clovemind.com" {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct domain name like xyz@clovemind.com", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// is_valid_email = UTIL.IsValidEmail(body["email_to"])
	// if is_valid_email == "" {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Pls enter correct email id", CONSTANT.ShowDialog, response)
	// 	return
	// }

	to := strings.Split(body["email_to"], ",")

	receipts := Model.EmailRecipientModel{
		ToEmails:  to,
		CcEmails:  []string{},
		BccEmails: []string{"qualitycontrol@clovemind.com"},
	}

	// set toAddress section
	var toRecipients []*string
	for _, r := range receipts.ToEmails {
		toAddress := r
		toRecipients = append(toRecipients, &toAddress)
	}

	var ccRecipients []*string
	if len(receipts.CcEmails) > 0 {
		for _, c := range receipts.CcEmails {
			ccAddress := c
			ccRecipients = append(ccRecipients, &ccAddress)
		}
	}

	var bccRecipients []*string
	if len(receipts.BccEmails) > 0 {
		for _, b := range receipts.BccEmails {
			bccAddress := b
			bccRecipients = append(bccRecipients, &bccAddress)
		}
	}

	UTIL.SendEmailForQuality(
		body["title"],
		body["body"],
		body["email_from"],
		toRecipients,
		ccRecipients,
		bccRecipients,
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

	oK := DB.CheckIfExists(CONSTANT.QualityCheckDetailsTable, map[string]string{"appointment_id": body["appointment_id"]})
	if !oK {
		UTIL.SetReponse(w, "400", "", CONSTANT.ShowDialog, response)
		return
	}

	status, ok := DB.UpdateSQL(CONSTANT.QualityCheckDetailsTable, map[string]string{"appointment_id": body["appointment_id"]}, map[string]string{"comments": body["comments"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}
