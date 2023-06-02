package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	_ "salbackend/model"
	UTIL "salbackend/util"
)

func TherapistGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapists
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " (first_name like '%%"+val[0]+"%%' or last_name like '%%"+val[0]+"%%') ")
			}
		case "phone":
			if len(val[0]) > 0 {
				wheres = append(wheres, " phone = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "email":
			if len(val[0]) > 0 {
				wheres = append(wheres, " email = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "therapist_id":
			wheres = append(wheres, " therapist_id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	therapists, status, ok := DB.SelectProcess("select * from "+CONSTANT.TherapistsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of therapists
	therapistsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.TherapistsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["therapists"] = therapists
	response["therapists_count"] = therapistsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(therapistsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func TherapistUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add therapist
	therapist := map[string]string{}
	therapist["first_name"] = body["first_name"]
	therapist["last_name"] = body["last_name"]
	therapist["phone"] = body["phone"]
	therapist["email"] = body["email"]
	therapist["gender"] = body["gender"]
	therapist["price"] = body["price"]
	therapist["price_3"] = body["price_3"]
	therapist["price_5"] = body["price_5"]
	therapist["corporate_price"] = body["corporate_price"]
	therapist["education"] = body["education"]
	therapist["experience"] = body["experience"]
	therapist["about"] = body["about"]
	therapist["payout_percentage"] = body["payout_percentage"]
	therapist["payee_name"] = body["payee_name"]
	therapist["bank_account_no"] = body["bank_account_no"]
	therapist["ifsc"] = body["ifsc"]
	therapist["branch_name"] = body["branch_name"]
	therapist["bank_name"] = body["bank_name"]
	therapist["bank_account_type"] = body["bank_account_type"]
	therapist["pan"] = body["pan"]
	therapist["corporate_therpist"] = body["corporate_therpist"]
	therapist["status"] = body["status"]
	therapist["modified_by"] = body["modified_by"]
	therapist["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.TherapistsTable, map[string]string{"therapist_id": r.FormValue("therapist_id")}, therapist)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
