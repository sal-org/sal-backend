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

func CounsellorGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellors
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
		case "counsellor_id":
			wheres = append(wheres, " counsellor_id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	counsellors, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of counsellors
	counsellorsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.CounsellorsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellors"] = counsellors
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURL
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add counsellor
	counsellor := map[string]string{}
	counsellor["first_name"] = body["first_name"]
	counsellor["last_name"] = body["last_name"]
	counsellor["phone"] = body["phone"]
	counsellor["email"] = body["email"]
	counsellor["gender"] = body["gender"]
	counsellor["price"] = body["price"]
	counsellor["price_3"] = body["price_3"]
	counsellor["price_5"] = body["price_5"]
	counsellor["corporate_price"] = body["corporate_price"]
	counsellor["education"] = body["education"]
	counsellor["experience"] = body["experience"]
	counsellor["about"] = body["about"]
	counsellor["payout_percentage"] = body["payout_percentage"]
	counsellor["payee_name"] = body["payee_name"]
	counsellor["bank_account_no"] = body["bank_account_no"]
	counsellor["ifsc"] = body["ifsc"]
	counsellor["branch_name"] = body["branch_name"]
	counsellor["bank_name"] = body["bank_name"]
	counsellor["bank_account_type"] = body["bank_account_type"]
	counsellor["pan"] = body["pan"]
	counsellor["corporate_therpist"] = body["corporate_therpist"]
	counsellor["status"] = body["status"]
	counsellor["modified_by"] = body["modified_by"]
	counsellor["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")}, counsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
