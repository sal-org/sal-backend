package client

import (
	"math"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strconv"
	"strings"
)

// EventsList godoc
// @Tags Client Event
// @Summary List available events
// @Router /client/events [get]
// @Produce json
// @Success 200
func EventsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming events
	events, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventTable + " where status = " + CONSTANT.EventToBeStarted + " order by date asc, time asc")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventDetail godoc
// @Tags Client Event
// @Summary Get event details
// @Router /client/event [get]
// @Param order_id query string true "Event order ID to get details"
// @Produce json
// @Success 200
func EventDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"*"}, map[string]string{"order_id": r.FormValue("order_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get event topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id = ?", event[0]["topic_id"])
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get event counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"counsellor_id": event[0]["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["event"] = event[0]
	response["counsellor"] = counsellor[0]
	if len(topics[0]) > 0 {
		response["topic"] = topics[0]["topic"]
	}
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventsBooked godoc
// @Tags Client Event
// @Summary Get booked upcoming and past events
// @Router /client/event/booked [get]
// @Param client_id query string true "Logged in client ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsBooked(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderClientEventTable+" where client_id = ? and status > "+CONSTANT.OrderWaiting+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") order by date asc, time asc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderClientEventTable+" where client_id = ? and status > "+CONSTANT.OrderWaiting+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("client_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderCreate godoc
// @Tags Client Event
// @Summary Book a slot in an event
// @Router /client/event/order [post]
// @Param body body model.ClientEventOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ClientEventOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"*"}, map[string]string{"client_id": body["client_id"]})
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

	// get event details
	event, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"*"}, map[string]string{"order_id": body["event_order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if event is valid
	if len(event) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if event is active
	if !strings.EqualFold(event[0]["status"], CONSTANT.EventToBeStarted) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EventAlreadyStartedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["event_order_id"] = body["event_order_id"]
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := event[0]["price"]

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and end_by > '"+UTIL.GetCurrentTime().String()+"' and (order_type = "+CONSTANT.OrderEventType+" or order_type is null) and (client_id = ? or client_id is null) order by created_at desc limit 1", body["coupon_code"], body["client_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponInCorrectMessage, CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon[0]["valid_for_order"]) > 0 && !strings.EqualFold(coupon[0]["valid_for_order"], "0") { // coupon is valid for particular order
			// get total number of client appointment/event orders
			noOrders := DB.RowCount(CONSTANT.InvoicesTable, " uesr_id = ?", body["client_id"])
			// check if coupon applicable by order count and valid for order
			if !strings.EqualFold(coupon[0]["valid_for_order"], strconv.Itoa(noOrders+1)) { // add 1 to equal to valid for order value
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponNotApplicableMessage, CONSTANT.ShowDialog, response)
				return
			}
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
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, strings.ReplaceAll(CONSTANT.CouponMinimumAmountRequiredMessage, "###amount###", coupon[0]["minimum_order_value"]), CONSTANT.ShowDialog, response)
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

	amount, _ := strconv.ParseFloat(body["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.FormatFloat(math.Round(amount*100), 'f', 2, 64)

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderClientEventTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderPaymentComplete godoc
// @Tags Client Event
// @Summary Call after payment is completed for event order
// @Router /client/event/paymentcomplete [post]
// @Param body body model.ClientEventOrderPaymentCompleteRequest true "Request Body"
// @Produce json
// @Success 200
func EventOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ClientEventOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderClientEventTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	razorPayTransaction := UTIL.GetRazorpayPayment(body["payment_id"])
	if !strings.EqualFold(razorPayTransaction.Description, body["order_id"]) { // check if razorpay payment id is associated with correct order id
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// capture razorpay payment
	amountRazorpay, _ := strconv.ParseFloat(order[0]["paid_amount_razorpay"], 64)
	UTIL.CaptureRazorpayPayment(body["payment_id"], amountRazorpay)

	// create invoice for the order
	invoice := map[string]string{}
	invoice["order_id"] = body["order_id"]
	invoice["payment_method"] = body["payment_method"]
	invoice["payment_id"] = body["payment_id"]
	invoice["user_id"] = order[0]["client_id"]
	invoice["user_type"] = CONSTANT.ClientType
	invoice["order_type"] = CONSTANT.OrderEventType
	invoice["actual_amount"] = order[0]["actual_amount"]
	invoice["tax"] = order[0]["tax"]
	invoice["discount"] = order[0]["discount"]
	invoice["coupon_code"] = order[0]["coupon_code"]
	invoice["coupon_id"] = order[0]["coupon_id"]
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
	status, ok = DB.UpdateSQL(CONSTANT.OrderClientEventTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	response["invoice_id"] = invoiceID
	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
