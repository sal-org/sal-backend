package admin

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	MODEL "salbackend/model"
	UTIL "salbackend/util"
	"strconv"
	"strings"
)

func CouponGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get coupons
	wheres := []string{}
	queryArgs := []any{}
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

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// // read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // check for required fields
	// fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CouponAddRequiredFields)
	// if len(fieldCheck) > 0 {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.CouponAddRequestInAdminPanel{}

	if err := UTIL.DecodeAndValidate(w, r, http.MethodPost, &body); err != nil {

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

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))

	// add coupon
	coupon := map[string]string{}
	coupon["coupon_code"] = body.CouponCode
	coupon["description"] = body.Description
	coupon["client_id"] = body.ClientID
	coupon["counsellor_id"] = body.CounsellorID
	coupon["therapist_id"] = body.TherapistID
	coupon["discount"] = body.Discount
	coupon["minimum_order_value"] = body.MinimumOrderValue
	coupon["maximum_discount_value"] = body.MaximumDiscountValue
	coupon["valid_for_order"] = body.ValidForOrder
	coupon["type"] = body.Type
	coupon["order_type"] = body.OrderType
	coupon["start_by"] = body.StartBy
	coupon["end_by"] = body.EndBy
	coupon["status"] = CONSTANT.CouponActive
	coupon["created_by"] = id
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

	var response = make(map[string]any)

	couponID, ok := UTIL.Required(r.FormValue("id"), "ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, couponID, CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	// body, ok := UTIL.ReadRequestBody(r)
	// if !ok {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	body := MODEL.CouponUpdateRequestInAdminPanel{}

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
	if !DB.CheckIfExists(CONSTANT.CouponsTable, map[string]string{"id": couponID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	id, _, _ := UTIL.ParseJWTAccessToken(r.Header.Get("Authorization"))


	// update coupon
	coupon := map[string]string{}
	coupon["coupon_code"] = body.CouponCode
	coupon["description"] = body.Description
	coupon["client_id"] = body.ClientID
	coupon["counsellor_id"] = body.CounsellorID
	coupon["therapist_id"] = body.TherapistID
	coupon["discount"] = body.Discount
	coupon["minimum_order_value"] = body.MinimumOrderValue
	coupon["maximum_discount_value"] = body.MaximumDiscountValue
	coupon["valid_for_order"] = body.ValidForOrder
	coupon["type"] = body.Type
	coupon["order_type"] = body.OrderType
	coupon["start_by"] = body.StartBy
	coupon["end_by"] = body.EndBy
	coupon["status"] = body.Status
	coupon["modified_by"] = id
	coupon["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok := DB.UpdateSQL(CONSTANT.CouponsTable, map[string]string{"id": r.FormValue("id")}, coupon)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
