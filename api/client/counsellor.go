package client

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	_ "salbackend/model" // for error response model
	UTIL "salbackend/util"
	"strconv"
	"strings"
)

// CounsellorProfile godoc
// @Tags Client
// @Summary Get counsellor details
// @Router /client/counsellor [get]
// @Param counsellor_id query string true "Counsellor ID to get details"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func CounsellorProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor languages
	languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get last 10 counsellor apppointment reviews
	reviews, status, ok := DB.SelectProcess("select a.comment, a.rating, a.modified_at, c.first_name, c.last_name from "+CONSTANT.AppointmentsTable+" a, "+CONSTANT.ClientsTable+" c where a.client_id = c.client_id and a.counsellor_id = ? and a.status = "+CONSTANT.AppointmentCompleted+" and a.comment is not null and a.comment != '' order by a.modified_at desc limit 10 ", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["counsellor"] = counsellor[0]
	response["languages"] = languages
	response["topics"] = topics
	response["reviews"] = reviews
	response["media_url"] = CONSTANT.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorSlots godoc
// @Tags Client
// @Summary Get counsellor slots
// @Router /client/counsellor/slots [get]
// @Param counsellor_id query string true "Counsellor ID to get slot details"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func CounsellorSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get counsellor slots
	slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["slots"] = slots
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorPrices godoc
// @Tags Client
// @Summary Get counsellor prices
// @Router /client/counsellor/prices [get]
// @Param counsellor_id query string true "Counsellor ID to get prices"
// @Param no_session query string false "Number of sessions to book (1,3,5) - default is 1, if you don't send"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func CounsellorPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get counsellor prices
	prices, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"price", "price_3", "price_5"}, map[string]string{"counsellor_id": r.FormValue("counsellor_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(prices) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	price := prices[0]["price"] // default 1 session price
	if strings.EqualFold(r.FormValue("no_session"), "3") {
		price = prices[0]["price_3"]
	} else if strings.EqualFold(r.FormValue("no_session"), "5") {
		price = prices[0]["price_5"]
	}

	if strings.EqualFold(price, "0") { // check if selected sessions price is not free
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorSessionsPriceNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = UTIL.GetBillingDetails(price, "0")
	response["slots"] = prices
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorOrderCreate godoc
// @Tags Client
// @Summary Create appointment order with client and counsellor
// @Router /client/counsellor/order [post]
// @Param body body model.CounsellorOrderCreateRequest true "Request Body"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func CounsellorOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if client id is valid
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"status"}, map[string]string{"client_id": body["client_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if client is valid
	if len(client) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if client is active
	if !strings.EqualFold(client[0]["status"], CONSTANT.ClientActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ClientNotAllowedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"price", "price_3", "price_5", "status"}, map[string]string{"counsellor_id": body["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if cousellor is valid
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if counsellor is active
	if !strings.EqualFold(counsellor[0]["status"], CONSTANT.CounsellorActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(body["counsellor_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["counsellor_id"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["type"] = CONSTANT.CounsellorType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := counsellor[0]["price"] // default 1 session price
	if strings.EqualFold(body["no_session"], "3") {
		price = counsellor[0]["price_3"]
	} else if strings.EqualFold(body["no_session"], "5") {
		price = counsellor[0]["price_5"]
	}

	order["slots_bought"] = body["no_session"]
	order["actual_amount"] = price

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and (client_id = ? or client_id is null) and (counsellor_id = ? or counsellor_id is null) order by created_at desc limit 1", body["coupon_code"], body["client_id"], body["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponInCorrectMessage, CONSTANT.ShowDialog, response)
			return
		}

		actualAmount, _ := strconv.ParseFloat(price, 64)
		minAmount, _ := strconv.ParseFloat(coupon[0]["minimum_order_value"], 64)
		if actualAmount >= minAmount { // coupon applicable only for minimum order value
			if strings.EqualFold(coupon[0]["type"], CONSTANT.CouponFlatType) { // flat
				order["discount"] = coupon[0]["discount"]
			} else if strings.EqualFold(coupon[0]["type"], CONSTANT.CouponPercentageType) { // percent
				discount, _ := strconv.ParseFloat(coupon[0]["discount"], 64) // percentage
				maxDiscount, _ := strconv.ParseFloat(coupon[0]["maximum_discount_value"], 64)
				discounted := actualAmount * discount / 100
				if discounted > maxDiscount {
					discounted = maxDiscount // maximum discount applied
				}
				order["discount"] = strconv.FormatFloat(discounted, 'f', 2, 64)
			}
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, strings.ReplaceAll(CONSTANT.CouponMinimumAmountRequiredMessage, "###amount###", price), CONSTANT.ShowDialog, response)
			return
		}

		order["coupon_code"] = body["coupon_code"]
		order["coupon_id"] = coupon[0]["id"]
	}

	// calculate bill
	billing := UTIL.GetBillingDetails(price, order["discount"])
	order["actual_amount"] = billing["actual_amount"]
	order["discount"] = billing["discount"]
	order["tax"] = billing["tax"]
	order["paid_amount"] = billing["paid_amount"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrdersTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorOrderPaymentComplete godoc
// @Tags Client
// @Summary Call after payment is completed for counsellor order
// @Router /client/counsellor/paymentcomplete [post]
// @Param body body model.CounsellorOrderPaymentCompleteRequest true "Request Body"
// @Produce json
// @Failure 400,500 {object} model.ErrorResponse
func CounsellorOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order
	order, status, ok := DB.SelectSQL(CONSTANT.OrdersTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if order is with counsellor
	if !strings.EqualFold(order[0]["type"], CONSTANT.CounsellorType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create invoice for the order
	invoice := map[string]string{}
	invoice["order_id"] = body["order_id"]
	invoice["payment_method"] = body["payment_method"]
	invoice["payment_id"] = body["payment_id"]
	invoice["client_id"] = order[0]["client_id"]
	invoice["counsellor_id"] = order[0]["counsellor_id"]
	invoice["type"] = CONSTANT.CounsellorType
	invoice["order_type"] = CONSTANT.OrderAppointmentType
	invoice["actual_amount"] = order[0]["actual_amount"]
	invoice["tax"] = order[0]["tax"]
	invoice["paid_amount"] = order[0]["paid_amount"]
	invoice["status"] = CONSTANT.InvoiceInProgress
	invoice["created_at"] = UTIL.GetCurrentTime().String()

	invoiceID, status, ok := DB.InsertWithUniqueID(CONSTANT.InvoicesTable, CONSTANT.InvoiceDigits, invoice, "invoice_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update order with invoice id and change status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	orderUpdate["invoice_id"] = invoiceID

	status, ok = DB.UpdateSQL(CONSTANT.OrdersTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	// update counsellor slots
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": order[0]["counsellor_id"],
			"date":          order[0]["date"],
		},
		map[string]string{
			order[0]["time"]: CONSTANT.SlotBooked,
		},
	)

	response["invoice_id"] = invoiceID
	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
