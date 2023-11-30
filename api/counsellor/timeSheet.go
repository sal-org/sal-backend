package counsellor

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	UTIL "salbackend/util"
)

func PartnerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	partnersName, status, ok := DB.SelectProcess("select partner_name from " + CONSTANT.CorporatePartnersTable + " where status = 1 order by id desc ")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["partners"] = partnersName

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func GetPartnerAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// // check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	partnerAddress, status, ok := DB.SelectProcess("select address from "+CONSTANT.CorporatePartnersAddressTable+" where partner_name = ? and status = 1 order by id desc ", r.FormValue("client_name"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["address"] = partnerAddress

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AddMyTimeSheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// add partner
	myTimeSheet := map[string]string{}
	myTimeSheet["counsellor_id"] = body["counsellor_id"]
	myTimeSheet["partner_name"] = body["partner_name"]
	myTimeSheet["location"] = body["location"]
	myTimeSheet["inTime"] = body["inTime"]
	// myTimeSheet["outTime"] = body["outTime"]
	myTimeSheet["status"] = CONSTANT.PartnerActive
	myTimeSheet["created_at"] = UTIL.GetCurrentTime().String()
	// status, ok := DB.InsertSQL(CONSTANT.CorporatePartnersTable, myTimeSheet)
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }
	_, status, ok := DB.InsertWithUniqueID(CONSTANT.CounsellorMyTimeSheetTable, CONSTANT.CounsellorRecordDigits, myTimeSheet, "sheet_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func MyTimeSheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	params := map[string]string{}

	params["counsellor_id"] = r.FormValue("counsellor_id")

	mysheet, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorMyTimeSheetTable+" where counsellor_id = ? and status = 1 order by created_at desc", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// mysheet, status, ok := DB.SelectSQL(CONSTANT.CounsellorMyTimeSheetTable, []string{"*"}, params)
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	if len(mysheet) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	response["mysheet"] = mysheet[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func UpdateMyTimeSheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update myTimeSheet details
	myTimeSheet := map[string]string{}
	if len(body["partner_name"]) > 0 {
		myTimeSheet["partner_name"] = body["partner_name"]
	}
	if len(body["location"]) > 0 {
		myTimeSheet["location"] = body["location"]
	}
	if len(body["inTime"]) > 0 {
		myTimeSheet["inTime"] = body["inTime"]
	}
	if len(body["outTime"]) > 0 {
		myTimeSheet["outTime"] = body["outTime"]
	}

	myTimeSheet["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CounsellorMyTimeSheetTable, map[string]string{"sheet_id": r.FormValue("sheet_id")}, myTimeSheet)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	mysheet, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorMyTimeSheetTable+" where sheet_id = ? and status = 1 order by created_at desc", r.FormValue("sheet_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"email"}, map[string]string{"counsellor_id": mysheet[0]["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(counsellor) == 0 {
		// get therapist details
		counsellor, status, ok = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"email"}, map[string]string{"therapist_id": mysheet[0]["counsellor_id"]})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	data := Model.EmailDataForCounsellorVisit{
		Client_Name:  mysheet[0]["partner_name"],
		Client_Location:   mysheet[0]["location"],
		InTime:      mysheet[0]["inTime"],
		OutTime:        mysheet[0]["outTime"],
	}

	filepath := "htmlfile/CounsellorVisit.html"

	emailbody := UTIL.GetHTMLTemplateForCounsellorVisit(data, filepath)

	UTIL.SendEmail(
		CONSTANT.CounsellorVisitForClientTitle,
		emailbody,
		counsellor[0]["email"],
		CONSTANT.InstantSendEmailMessage,
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
