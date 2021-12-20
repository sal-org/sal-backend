package therapist

import (
	"math"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

// EventsList godoc
// @Tags Therapist Event
// @Summary List available events
// @Router /therapist/events [get]
// @Security JWTAuth
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
// @Tags Therapist Event
// @Summary Get event details
// @Router /therapist/event [get]
// @Param order_id query string true "Event order ID to get details"
// @Security JWTAuth
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
	var counsellor []map[string]string
	switch event[0]["type"] {
	case CONSTANT.CounsellorType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"counsellor_id": event[0]["counsellor_id"]})
		break
	case CONSTANT.TherapistType:
		counsellor, _, _ = DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about"}, map[string]string{"therapist_id": event[0]["counsellor_id"]})
		break
	}
	if len(counsellor) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["event"] = event[0]
	response["counsellor"] = counsellor[0]
	if len(topics) > 0 && len(topics[0]) > 0 {
		response["topic"] = topics[0]["topic"]
	}
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventsBooked godoc
// @Tags Therapist Event
// @Summary Get booked upcoming and past events
// @Router /therapist/event/booked [get]
// @Param therapist_id query string true "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsBooked(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventTable+" where user_id = ? and status > "+CONSTANT.OrderWaiting+") and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") order by date asc, time asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where order_id in (select event_order_id from "+CONSTANT.OrderEventTable+" where user_id = ? and status > "+CONSTANT.OrderWaiting+") and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderCreate godoc
// @Tags Therapist Event
// @Summary Book a slot in an event
// @Router /therapist/event/order [post]
// @Param body body model.EventOrderCreateRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["user_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if therapist is valid
	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if therapist is active
	if !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistAccountDeletedMessage, CONSTANT.ShowDialog, response)
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
	order["user_id"] = body["user_id"]
	order["event_order_id"] = body["event_order_id"]
	order["user_type"] = CONSTANT.TherapistType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := event[0]["price"]

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and '"+UTIL.GetCurrentTime().String()+"' < end_by and (order_type = "+CONSTANT.OrderEventBookType+" or order_type = 0) and (counsellor_id = ? or counsellor_id = '') order by created_at desc limit 1", body["coupon_code"], body["user_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponInCorrectMessage, CONSTANT.ShowDialog, response)
			return
		}
		if !strings.EqualFold(coupon[0]["valid_for_order"], "0") { // coupon is valid for particular order
			// get total number of therapist appointment/event orders
			noOrders := DB.RowCount(CONSTANT.InvoicesTable, " user_id = ?", body["user_id"])
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
	order["cgst"] = billing["cgst"]
	order["sgst"] = billing["sgst"]
	order["paid_amount"] = billing["paid_amount"]

	amount, _ := strconv.ParseFloat(order["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.Itoa(int(math.Round(amount * 100)))
	response["paid_amount_razorpay"] = order["paid_amount_razorpay"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderEventTable, CONSTANT.OrderEventDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderPaymentComplete godoc
// @Tags Therapist Event
// @Summary Call after payment is completed for event order
// @Router /therapist/event/paymentcomplete [post]
// @Param body body model.EventOrderPaymentCompleteRequest true "Request Body"
// @Security JWTAuth
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderEventTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
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
	invoice["user_id"] = order[0]["user_id"]
	invoice["user_type"] = CONSTANT.TherapistType
	invoice["order_type"] = CONSTANT.OrderEventBookType
	invoice["actual_amount"] = order[0]["actual_amount"]
	invoice["tax"] = order[0]["tax"]
	invoice["cgst"] = order[0]["cgst"]
	invoice["sgst"] = order[0]["sgst"]
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
	status, ok = DB.UpdateSQL(CONSTANT.OrderEventTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	orderdetails, _, _ := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"title", "time", "date"}, map[string]string{"counsellor_id": order[0]["event_order_id"]})
	//counsellor, _, _ := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "phone", "timezone"}, map[string]string{"counsellor_id": orderdetails[0]["counsellor_id"]})
	client, _, _ := DB.SelectSQL(CONSTANT.ClientsTable, []string{"timezone"}, map[string]string{"therapist_id": order[0]["user_id"]})

	// send appointment booking notification to client
	UTIL.SendNotification(
		CONSTANT.ClientEventPaymentSucessClientHeading,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.ClientEventPaymentSucessClientContent,
			map[string]string{
				"###cafe_name###":   orderdetails[0]["title"],
				"###paid_amount###": order[0]["paid_amount"],
				"###date_time###":   UTIL.ConvertTimezone(UTIL.BuildDateTime(order[0]["date"], order[0]["time"]), client[0]["timezone"]).Format(CONSTANT.ReadbleDateTimeFormat),
			},
		),
		order[0]["user_id"],
		CONSTANT.TherapistType,
	)

	response["invoice_id"] = invoiceID
	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}

// EventsBlocked godoc
// @Tags Therapist Event
// @Summary Get blocked upcoming and past events
// @Router /therapist/event/block [get]
// @Param therapist_id query string true "Logged in therapist ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsBlocked(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming blocked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where counsellor_id = ? and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") order by date asc, time asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// get past blocked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where counsellor_id = ? and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventUpdate godoc
// @Tags Therapist Event
// @Summary Start and stop event
// @Router /therapist/event [put]
// @Param order_id query string true "Event order ID to update"
// @Param therapist_id query string true "Logged in therapist ID"
// @Param status query string true "Start(2), Stop(3)"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if strings.EqualFold(r.FormValue("status"), CONSTANT.EventStarted) {
		status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("order_id"), "counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.EventToBeStarted}, map[string]string{"status": CONSTANT.EventStarted, "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		status, ok := DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable, map[string]string{"order_id": r.FormValue("order_id"), "counsellor_id": r.FormValue("therapist_id"), "status": CONSTANT.EventStarted}, map[string]string{"status": CONSTANT.EventCompleted, "modified_at": UTIL.GetCurrentTime().String()})
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventBlockOrderCreate godoc
// @Tags Therapist Event
// @Summary Block a slot for an event
// @Router /therapist/event/block/order [post]
// @Param body body model.EventBlockOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventBlockOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventBlockOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"status"}, map[string]string{"therapist_id": body["counsellor_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if therapist is valid
	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if therapist is active
	if !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["counsellor_id"] = body["counsellor_id"]
	order["title"] = body["title"]
	order["description"] = body["description"]
	order["topic_id"] = body["topic_id"]
	order["type"] = CONSTANT.TherapistType
	order["duration"] = CONSTANT.EventDuration
	order["price"] = body["price"]
	order["photo"] = body["photo"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["status"] = CONSTANT.EventWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	// calculate bill
	billing := UTIL.GetBillingDetails(CONSTANT.EventPrice, "0")
	order["actual_amount"] = billing["actual_amount"]
	order["tax"] = billing["tax"]
	order["cgst"] = billing["cgst"]
	order["sgst"] = billing["sgst"]
	order["paid_amount"] = billing["paid_amount"]

	amount, _ := strconv.ParseFloat(order["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.Itoa(int(math.Round(amount * 100)))
	response["paid_amount_razorpay"] = order["paid_amount_razorpay"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderCounsellorEventTable, CONSTANT.EventDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	//body := HTMLDATA.GetHTMLTemplate(orderdetails)
	orderdetails, _, _ := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"counsellor_id", "title", "description", "photo", "topic_id", "date", "time", "duration", "price"}, map[string]string{"order_id": orderID})
	counsellordetails, _, _ := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name"}, map[string]string{"therapist_id": orderdetails[0]["counsellor_id"]})
	topic_name, _, _ := DB.SelectSQL(CONSTANT.TopicsTable, []string{"topic"}, map[string]string{"id": orderdetails[0]["topic_id"]})
	/*emailbody := UTIL.GetHTMLTemplateForEvent(orderdetails)

	UTIL.SendEmail(
		CONSTANT.NewEventWaitingForApprovalTitle,
		emailbody,
		CONSTANT.SameerEmailID,
		CONSTANT.InstantSendEmailMessage,
	)*/

	UTIL.SendEmail(
		CONSTANT.NewEventWaitingForApprovalTitle,
		UTIL.ReplaceNotificationContentInString(
			CONSTANT.EventWaitingForApprovalBody,
			map[string]string{
				"###First_name###":  counsellordetails[0]["first_name"],
				"###Last_name###":   counsellordetails[0]["last_name"],
				"###type###":        "Therapists",
				"###title###":       orderdetails[0]["title"],
				"###description###": orderdetails[0]["description"],
				"###photo###":       orderdetails[0]["photo"],
				"###topic_id###":    topic_name[0]["topic"],
				"###date###":        orderdetails[0]["date"],
				"###time###":        orderdetails[0]["time"],
				"###duration###":    orderdetails[0]["duration"],
				"###price###":       orderdetails[0]["price"],
			},
		),
		CONSTANT.SameerEmailID,
		CONSTANT.InstantSendEmailMessage,
	)

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventBlockOrderPaymentComplete godoc
// @Tags Therapist Event
// @Summary Call after payment is completed for event block order
// @Router /therapist/event/block/paymentcomplete [post]
// @Param body body model.EventBlockOrderPaymentCompleteRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
/*func EventBlockOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.EventBlockOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderCounsellorEventTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if event is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if event payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.EventWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedEventMessage, CONSTANT.ShowDialog, response)
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
	invoice["user_id"] = order[0]["counsellor_id"]
	invoice["order_id"] = body["order_id"]
	invoice["payment_method"] = body["payment_method"]
	invoice["payment_id"] = body["payment_id"]
	invoice["user_type"] = CONSTANT.TherapistType
	invoice["order_type"] = CONSTANT.OrderEventBlockType
	invoice["actual_amount"] = order[0]["actual_amount"]
	invoice["tax"] = order[0]["tax"]
	invoice["cgst"] = order[0]["cgst"]
	invoice["sgst"] = order[0]["sgst"]
	invoice["paid_amount"] = order[0]["paid_amount"]
	invoice["status"] = CONSTANT.InvoiceInProgress
	invoice["created_at"] = UTIL.GetCurrentTime().String()

	invoiceID, status, ok := DB.InsertWithUniqueID(CONSTANT.InvoicesTable, CONSTANT.InvoiceDigits, invoice, "invoice_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update event with invoice id and change status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.EventToBeStarted
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	orderUpdate["invoice_id"] = invoiceID
	status, ok = DB.UpdateSQL(CONSTANT.OrderCounsellorEventTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	response["invoice_id"] = invoiceID
	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
*/
