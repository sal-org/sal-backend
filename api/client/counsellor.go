package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"

	UTIL "salbackend/util"
	"strconv"
	"strings"
)

// CounsellorProfile godoc
// @Tags Client Counsellor
// @Summary Get counsellor details
// @Router /client/counsellor [get]
// @Param counsellor_id query string true "Counsellor ID to get details"
// @Produce json
// @Success 200
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
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorSlots godoc
// @Tags Client Counsellor
// @Summary Get counsellor slots
// @Router /client/counsellor/slots [get]
// @Param counsellor_id query string true "Counsellor ID to get slot details"
// @Produce json
// @Success 200
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

// CounsellorOrderCreate godoc
// @Tags Client Counsellor
// @Summary Create appointment order with client and counsellor
// @Router /client/counsellor/order [post]
// @Param body body model.CounsellorOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
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

	// get client details
	client, status, ok := DB.SelectSQL(CONSTANT.ClientsTable, []string{"first_name", "last_name", "email", "phone", "status"}, map[string]string{"client_id": body["client_id"]})
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
	counsellor, status, ok := DB.SelectSQL(CONSTANT.CounsellorsTable, []string{"first_name", "price", "price_3", "price_5", "status"}, map[string]string{"counsellor_id": body["counsellor_id"]})
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
	order["order_type"] = CONSTANT.OrderAppointmentType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()

	price := counsellor[0]["price"] // default 1 session price
	if strings.EqualFold(body["no_session"], "3") {
		price = counsellor[0]["price_3"]
	} else if strings.EqualFold(body["no_session"], "5") {
		price = counsellor[0]["price_5"]
	}

	// appointment actual price should not be free
	if strings.EqualFold(price, "0") {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CounsellorSessionsPriceNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	order["slots_bought"] = body["no_session"]
	order["actual_amount"] = price

	if len(body["coupon_code"]) > 0 {
		// get coupon details
		coupon, status, ok := DB.SelectProcess("select * from "+CONSTANT.CouponsTable+" where coupon_code = ? and status = 1 and start_by < '"+UTIL.GetCurrentTime().String()+"' and end_by > '"+UTIL.GetCurrentTime().String()+"' and order_type = "+CONSTANT.OrderAppointmentType+" and (client_id = ? or client_id is null) and (counsellor_id = ? or counsellor_id is null) order by created_at desc limit 1", body["coupon_code"], body["client_id"], body["counsellor_id"])
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CouponInCorrectMessage, CONSTANT.ShowDialog, response)
			return
		}
		if len(coupon[0]["valid_for_order"]) > 0 {
			// get number of client counsellor orders
			noOrders := DB.RowCount(CONSTANT.OrdersTable, " client_id = ? and type = "+CONSTANT.CounsellorType+" and order_type = "+CONSTANT.OrderAppointmentType+" and status > "+CONSTANT.OrderWaiting, body["client_id"])
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

	paidAmount, _ := strconv.ParseFloat(billing["paid_amount"], 64)
	response["payu_payment"] = UTIL.GetPayUPaymentObject(CONFIG.PAYUMerchatKey, CONFIG.PAYUSalt, orderID, "Appointment with "+counsellor[0]["first_name"], client[0]["first_name"], client[0]["last_name"], client[0]["email"], client[0]["phone"], "https://hwmpf9h476.execute-api.ap-south-1.amazonaws.com/prod/client/testpayu", "https://hwmpf9h476.execute-api.ap-south-1.amazonaws.com/prod/client/testpayu", paidAmount)
	response["billing"] = billing
	response["order_id"] = orderID
	response["prices"] = counsellor[0]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// CounsellorOrderPaymentComplete godoc
// @Tags Client Counsellor
// @Summary Call after payment is completed for counsellor order
// @Router /client/counsellor/paymentcomplete [post]
// @Param body body model.CounsellorOrderPaymentCompleteRequest true "Request Body"
// @Produce json
// @Success 200
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

	// get order details
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
	// check if order is appointment
	if !strings.EqualFold(order[0]["order_type"], CONSTANT.OrderAppointmentType) {
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
	invoice["type"] = order[0]["type"]
	invoice["order_type"] = order[0]["order_type"]
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

	// create appointment slots between counsellor and client
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

	// create appointment between counsellor and client
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

// TestPAYU .
func TestPAYU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("TestPAYU", err)
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	defer r.Body.Close()
	fmt.Println("TestPAYU", string(b))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
