package counsellor

import (
	"math"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	UTIL "salbackend/util"
)

// EventsList godoc
// @Tags Counsellor Event
// @Summary Get upcoming and past counsellor events
// @Router /counsellor/events [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Security JWTAuth
// @Produce json
// @Success 200
func EventsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get upcoming booked events
	events, status, ok := DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where counsellor_id = ? and status in ("+CONSTANT.EventToBeStarted+", "+CONSTANT.EventStarted+") order by date asc, time asc", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["upcoming_events"] = events

	// get past booked events (get all booked event orders other than in progress, which is status > 1 (inprogress))
	events, status, ok = DB.SelectProcess("select * from "+CONSTANT.OrderCounsellorEventTable+" where counsellor_id = ? and status = "+CONSTANT.EventCompleted+" order by date desc, time desc", r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	response["past_events"] = events
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderCreate godoc
// @Tags Counsellor Event
// @Summary Book a slot in an event
// @Router /counsellor/event/order [post]
// @Param body body model.CounsellorEventOrderCreateRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorEventOrderCreateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor details
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"status"}, map[string]string{"counsellor_id": body["counsellor_id"]})
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

	// order object to be inserted
	order := map[string]string{}
	order["counsellor_id"] = body["counsellor_id"]
	order["title"] = body["title"]
	order["description"] = body["description"]
	order["topic_id"] = body["topic_id"]
	order["duration"] = body["duration"]
	order["price"] = body["price"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["status"] = CONSTANT.EventWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	// calculate bill
	billing := UTIL.GetBillingDetails(CONSTANT.EventPrice, "0")
	order["actual_amount"] = billing["actual_amount"]
	order["tax"] = billing["tax"]
	order["paid_amount"] = billing["paid_amount"]

	amount, _ := strconv.ParseFloat(order["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.Itoa(int(math.Round(amount * 100)))
	response["paid_amount_razorpay"] = order["paid_amount_razorpay"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderCounsellorEventTable, CONSTANT.EventDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// EventOrderPaymentComplete godoc
// @Tags Counsellor Event
// @Summary Call after payment is completed for event order
// @Router /counsellor/event/paymentcomplete [post]
// @Param body body model.CounsellorEventOrderPaymentCompleteRequest true "Request Body"
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CounsellorEventOrderPaymentCompleteRequiredFields)
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
	invoice["user_type"] = CONSTANT.CounsellorType
	invoice["order_type"] = CONSTANT.OrderEventType
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
