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

// TherapistProfile godoc
// @Tags Client Therapist
// @Summary Get therapist details
// @Router /client/therapist [get]
// @Param therapist_id query string true "Therapist ID to get details"
// @Produce json
// @Success 200
func TherapistProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo", "price", "education", "experience", "about", "slot_type"}, map[string]string{"therapist_id": r.FormValue("therapist_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get therapist languages
	languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get therapist topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get last 10 therapist apppointment reviews
	reviews, status, ok := DB.SelectProcess("select a.comment, a.rating, a.modified_at, c.first_name, c.last_name from "+CONSTANT.AppointmentsTable+" a, "+CONSTANT.ClientsTable+" c where a.client_id = c.client_id and a.counsellor_id = ? and a.status = "+CONSTANT.AppointmentCompleted+" and a.comment is not null and a.comment != '' order by a.modified_at desc limit 10 ", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor latest content
	contents, status, ok := DB.SelectProcess("select * from "+CONSTANT.ContentsTable+" where counsellor_id = ? and training = 0 and status = 1 order by created_at desc limit 20", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["therapist"] = therapist[0]
	response["languages"] = languages
	response["topics"] = topics
	response["reviews"] = reviews
	response["contents"] = contents
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// TherapistSlots godoc
// @Tags Client Therapist
// @Summary Get therapist slots
// @Router /client/therapist/slots [get]
// @Param therapist_id query string true "Therapist ID to get slot details"
// @Produce json
// @Success 200
func TherapistSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get therapist slots
	slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// remove times and dates with no availability
	response["slots"] = UTIL.FilterAvailableSlots(slots)
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// TherapistOrderCreate godoc
// @Tags Client Therapist
// @Summary Create appointment order with client and therapist
// @Router /client/therapist/order [post]
// @Param body body model.TherapistOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func TherapistOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TherapistOrderCreateRequiredFields)
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

	// get therapist details
	therapist, status, ok := DB.SelectSQL(CONSTANT.TherapistsTable, []string{"*"}, map[string]string{"therapist_id": body["therapist_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if cousellor is valid
	if len(therapist) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if therapist is active
	if !strings.EqualFold(therapist[0]["status"], CONSTANT.TherapistActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(body["therapist_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["therapist_id"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["type"] = CONSTANT.TherapistType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := therapist[0]["price"] // default 1 session price
	if strings.EqualFold(body["no_session"], "3") {
		price = therapist[0]["price_3"]
	} else if strings.EqualFold(body["no_session"], "5") {
		price = therapist[0]["price_5"]
	}

	// appointment actual price should not be free
	if len(price) == 0 || strings.EqualFold(price, "0") {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TherapistSessionsPriceNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	order["slots_bought"] = body["no_session"]
	order["actual_amount"] = price

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and end_by > '"+UTIL.GetCurrentTime().String()+"' and (order_type = "+CONSTANT.OrderAppointmentType+" or order_type is null) and (client_id = ? or client_id is null) and (therapist_id = ? or therapist_id is null) order by created_at desc limit 1", body["coupon_code"], body["client_id"], body["therapist_id"])
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
			noOrders := DB.RowCount(CONSTANT.InvoicesTable, " user_id = ?", body["client_id"])
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

	amount, _ := strconv.ParseFloat(order["paid_amount"], 64)
	order["paid_amount_razorpay"] = strconv.Itoa(int(math.Round(amount * 100)))
	response["paid_amount_razorpay"] = order["paid_amount_razorpay"]

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderClientAppointmentTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["billing"] = billing
	response["order_id"] = orderID
	response["prices"] = map[string]string{"price": therapist[0]["price"], "price_3": therapist[0]["price_3"], "price_5": therapist[0]["price_5"]}
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// TherapistOrderPaymentComplete godoc
// @Tags Client Therapist
// @Summary Call after payment is completed for therapist order
// @Router /client/therapist/paymentcomplete [post]
// @Param body body model.TherapistOrderPaymentCompleteRequest true "Request Body"
// @Produce json
// @Success 200
func TherapistOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TherapistOrderPaymentCompleteRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	order, status, ok := DB.SelectSQL(CONSTANT.OrderClientAppointmentTable, []string{"*"}, map[string]string{"order_id": body["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order is valid
	if len(order) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.OrderNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if order is with therapist
	if !strings.EqualFold(order[0]["type"], CONSTANT.TherapistType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
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
	invoice["order_type"] = CONSTANT.OrderAppointmentType
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

	// create appointment slots between therapist and client
	appointmentSlots := map[string]string{}
	appointmentSlots["order_id"] = body["order_id"]
	appointmentSlots["client_id"] = order[0]["client_id"]
	appointmentSlots["counsellor_id"] = order[0]["counsellor_id"]
	appointmentSlots["slots_bought"] = order[0]["slots_bought"]
	slotsBought, _ := strconv.Atoi(order[0]["slots_bought"])
	appointmentSlots["slots_remaining"] = strconv.Itoa(slotsBought - 1)
	appointmentSlots["status"] = CONSTANT.AppointmentSlotsActive
	appointmentSlots["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok = DB.InsertWithUniqueID(CONSTANT.AppointmentSlotsTable, CONSTANT.AppointmentSlotDigits, appointmentSlots, "appointment_slot_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// create appointment between therapist and client
	appointment := map[string]string{}
	appointment["order_id"] = body["order_id"]
	appointment["client_id"] = order[0]["client_id"]
	appointment["counsellor_id"] = order[0]["counsellor_id"]
	appointment["date"] = order[0]["date"]
	appointment["time"] = order[0]["time"]
	appointment["status"] = CONSTANT.AppointmentToBeStarted
	appointment["created_at"] = UTIL.GetCurrentTime().String()
	_, status, ok = DB.InsertWithUniqueID(CONSTANT.AppointmentsTable, CONSTANT.AppointmentDigits, appointment, "appointment_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// update order with invoice id and change status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	orderUpdate["invoice_id"] = invoiceID
	status, ok = DB.UpdateSQL(CONSTANT.OrderClientAppointmentTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	// update therapist slots
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
