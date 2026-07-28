package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	CONFIG "salbackend/config"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
)

func CounsellorGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellors
	wheres := []string{}
	queryArgs := []any{}
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

	// for _, counsellor := range counsellors {
	// 	url := UTIL.PreSignedS3URLToGetTheData(CONFIG.S3Bucket, counsellor["photo"], CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	// 	_, endPointURL := UTIL.GetBaseURLAndEndpointFromURL(url)
	// 	counsellor["photo"] = endPointURL
	// }

	response["counsellors"] = counsellors
	response["counsellors_count"] = counsellorsCount[0]["ctn"]
	response["media_url"] = CONFIG.MediaURLInCLOUDFRONT
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(counsellorsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CounsellorUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	counsellorID, ok := UTIL.Required(r.FormValue("counsellor_id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, counsellorID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.UpdateCounsellorProfileRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPut, &body); err != nil {

		switch err {
		case CONSTANT.ErrMethodNotAllowed:
			UTIL.SetReponse(w, CONSTANT.StatusMethodNotAllowed, err.Error(), CONSTANT.ShowDialog, response)
			return

		case CONSTANT.ErrInvalidContentType:
			UTIL.SetReponse(w, CONSTANT.StatusUnsupportedMediaType, err.Error(), CONSTANT.ShowDialog, response)
			return

		default:
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	// check if corporate_id exists
	if !DB.CheckIfExists(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": counsellorID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add counsellor
	counsellor := map[string]string{}
	counsellor["first_name"] = body.FirstName
	counsellor["last_name"] = body.LastName
	counsellor["phone"] = body.Phone
	counsellor["email"] = body.Email
	counsellor["gender"] = body.Gender
	counsellor["price"] = body.Price
	counsellor["price_3"] = body.Price3
	counsellor["price_5"] = body.Price5
	counsellor["corporate_price"] = body.CorporatePrice
	counsellor["education"] = body.Education
	counsellor["experience"] = body.Experience
	counsellor["about"] = body.About
	counsellor["payout_percentage"] = body.PayoutPercentage
	counsellor["payee_name"] = body.PayeeName
	counsellor["bank_account_no"] = body.BankAccountNo
	counsellor["ifsc"] = body.IFSC
	counsellor["branch_name"] = body.BranchName
	counsellor["bank_name"] = body.BankName
	counsellor["bank_account_type"] = body.BankAccountType
	counsellor["pan"] = body.PAN
	counsellor["corporate_therpist"] = body.CorporateTherapist
	counsellor["status"] = body.Status
	counsellor["modified_by"] = id
	counsellor["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CounsellorsTable, map[string]string{"counsellor_id": r.FormValue("counsellor_id")}, counsellor)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
