package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

func CouponGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	// if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
	// 	return
	// }

	// get coupons
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "active":
			if len(val[0]) > 0 {
				if strings.EqualFold(val[0], "1") {
					wheres = append(wheres, " status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and '"+UTIL.GetCurrentTime().String()+"' < end_by  ")
				} else {
					wheres = append(wheres, " (status = 0 or start_by > '"+UTIL.GetCurrentTime().String()+"' or '"+UTIL.GetCurrentTime().String()+"' > end_by)  ")
				}
			}
		case "id":
			wheres = append(wheres, " id = ? ")
			queryArgs = append(queryArgs, val[0])
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	coupons, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+where+" order by id desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of coupons
	couponsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.CouponsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["coupons"] = coupons
	response["coupons_count"] = couponsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(couponsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CouponAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CouponAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// add coupon
	coupon := map[string]string{}
	coupon["coupon_code"] = body["coupon_code"]
	coupon["description"] = body["description"]
	coupon["client_id"] = body["client_id"]
	coupon["counsellor_id"] = body["counsellor_id"]
	coupon["therapist_id"] = body["therapist_id"]
	coupon["discount"] = body["discount"]
	coupon["minimum_order_value"] = body["minimum_order_value"]
	coupon["maximum_discount_value"] = body["maximum_discount_value"]
	coupon["valid_for_order"] = body["valid_for_order"]
	coupon["type"] = body["type"]
	coupon["order_type"] = body["order_type"]
	coupon["start_by"] = body["start_by"]
	coupon["end_by"] = body["end_by"]
	coupon["status"] = CONSTANT.CouponActive
	coupon["created_by"] = body["created_by"]
	coupon["created_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.InsertSQL(CONSTANT.CouponsTable, coupon)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CouponUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// update coupon
	coupon := map[string]string{}
	coupon["coupon_code"] = body["coupon_code"]
	coupon["description"] = body["description"]
	coupon["client_id"] = body["client_id"]
	coupon["counsellor_id"] = body["counsellor_id"]
	coupon["therapist_id"] = body["therapist_id"]
	coupon["discount"] = body["discount"]
	coupon["minimum_order_value"] = body["minimum_order_value"]
	coupon["maximum_discount_value"] = body["maximum_discount_value"]
	coupon["valid_for_order"] = body["valid_for_order"]
	coupon["type"] = body["type"]
	coupon["order_type"] = body["order_type"]
	coupon["start_by"] = body["start_by"]
	coupon["end_by"] = body["end_by"]
	coupon["status"] = body["status"]
	coupon["modified_by"] = body["modified_by"]
	coupon["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CouponsTable, map[string]string{"id": r.FormValue("id")}, coupon)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
