package client

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
	"strings"
)

// ListenerProfile godoc
// @Tags Client Listener
// @Summary Get listener details
// @Router /client/listener [get]
// @Param listener_id query string true "Listener ID to get details"
// @Produce json
// @Success 200
func ListenerProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get listener details
	listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"first_name", "last_name", "total_rating", "average_rating", "photo"}, map[string]string{"listener_id": r.FormValue("listener_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(listener) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get listener languages
	languages, status, ok := DB.SelectProcess("select language from "+CONSTANT.LanguagesTable+" where id in (select language_id from "+CONSTANT.CounsellorLanguagesTable+" where counsellor_id = ?)", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get listener topics
	topics, status, ok := DB.SelectProcess("select topic from "+CONSTANT.TopicsTable+" where id in (select topic_id from "+CONSTANT.CounsellorTopicsTable+" where counsellor_id = ?)", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get last 10 listener apppointment reviews
	reviews, status, ok := DB.SelectProcess("select a.comment, a.rating, a.modified_at, c.first_name, c.last_name from "+CONSTANT.AppointmentsTable+" a, "+CONSTANT.ClientsTable+" c where a.client_id = c.client_id and a.counsellor_id = ? and a.status = "+CONSTANT.AppointmentCompleted+" and a.comment is not null and a.comment != '' order by a.modified_at desc limit 10 ", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["listener"] = listener[0]
	response["languages"] = languages
	response["topics"] = topics
	response["reviews"] = reviews
	response["media_url"] = CONFIG.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerSlots godoc
// @Tags Client Listener
// @Summary Get listener slots
// @Router /client/listener/slots [get]
// @Param listener_id query string true "Listener ID to get slot details"
// @Produce json
// @Success 200
func ListenerSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get listener slots
	slots, status, ok := DB.SelectProcess("select * from "+CONSTANT.SlotsTable+" where counsellor_id = ? and date >= '"+UTIL.GetCurrentTime().Format("2006-01-02")+"' order by date asc", r.FormValue("listener_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// remove dates with no availability
	filteredSlots := []map[string]string{}
	for _, slot := range slots {
		if strings.EqualFold(slot["0"], "1") ||
			strings.EqualFold(slot["1"], "1") ||
			strings.EqualFold(slot["2"], "1") ||
			strings.EqualFold(slot["3"], "1") ||
			strings.EqualFold(slot["4"], "1") ||
			strings.EqualFold(slot["5"], "1") ||
			strings.EqualFold(slot["6"], "1") ||
			strings.EqualFold(slot["7"], "1") ||
			strings.EqualFold(slot["8"], "1") ||
			strings.EqualFold(slot["9"], "1") ||
			strings.EqualFold(slot["10"], "1") ||
			strings.EqualFold(slot["11"], "1") ||
			strings.EqualFold(slot["12"], "1") ||
			strings.EqualFold(slot["13"], "1") ||
			strings.EqualFold(slot["14"], "1") ||
			strings.EqualFold(slot["15"], "1") ||
			strings.EqualFold(slot["16"], "1") ||
			strings.EqualFold(slot["17"], "1") ||
			strings.EqualFold(slot["18"], "1") ||
			strings.EqualFold(slot["19"], "1") ||
			strings.EqualFold(slot["20"], "1") ||
			strings.EqualFold(slot["21"], "1") ||
			strings.EqualFold(slot["22"], "1") ||
			strings.EqualFold(slot["23"], "1") ||
			strings.EqualFold(slot["24"], "1") ||
			strings.EqualFold(slot["25"], "1") ||
			strings.EqualFold(slot["26"], "1") ||
			strings.EqualFold(slot["27"], "1") ||
			strings.EqualFold(slot["28"], "1") ||
			strings.EqualFold(slot["29"], "1") ||
			strings.EqualFold(slot["30"], "1") ||
			strings.EqualFold(slot["31"], "1") ||
			strings.EqualFold(slot["32"], "1") ||
			strings.EqualFold(slot["33"], "1") ||
			strings.EqualFold(slot["34"], "1") ||
			strings.EqualFold(slot["35"], "1") ||
			strings.EqualFold(slot["36"], "1") ||
			strings.EqualFold(slot["37"], "1") ||
			strings.EqualFold(slot["38"], "1") ||
			strings.EqualFold(slot["39"], "1") ||
			strings.EqualFold(slot["40"], "1") ||
			strings.EqualFold(slot["41"], "1") ||
			strings.EqualFold(slot["42"], "1") ||
			strings.EqualFold(slot["43"], "1") ||
			strings.EqualFold(slot["44"], "1") ||
			strings.EqualFold(slot["45"], "1") ||
			strings.EqualFold(slot["46"], "1") ||
			strings.EqualFold(slot["47"], "1") {

			filteredSlot := map[string]string{
				"date": slot["date"],
			}
			// show only times with availability
			for i := 0; i < 24; i++ {
				if strings.EqualFold(slot[strconv.Itoa(i)], "1") {
					filteredSlot[strconv.Itoa(i)] = "1"
				}
			}

			// TODO - remove expired time for today
			filteredSlots = append(filteredSlots, filteredSlot)
		}
	}

	response["slots"] = filteredSlots
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerOrderCreate godoc
// @Tags Client Listener
// @Summary Create appointment order with client and listener
// @Router /client/listener/order [post]
// @Param body body model.ListenerOrderCreateRequest true "Request Body"
// @Security JWTAuth
// @Produce json
// @Success 200
func ListenerOrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ListenerOrderCreateRequiredFields)
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

	// get listener details
	listener, status, ok := DB.SelectSQL(CONSTANT.ListenersTable, []string{"*"}, map[string]string{"listener_id": body["listener_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if cousellor is valid
	if len(listener) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	// check if listener is active
	if !strings.EqualFold(listener[0]["status"], CONSTANT.ListenerActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if slots available
	if !UTIL.CheckIfAppointmentSlotAvailable(body["listener_id"], body["date"], body["time"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ListenerSlotNotAvailableMessage, CONSTANT.ShowDialog, response)
		return
	}

	// order object to be inserted
	order := map[string]string{}
	order["client_id"] = body["client_id"]
	order["counsellor_id"] = body["listener_id"]
	order["date"] = body["date"]
	order["time"] = body["time"]
	order["type"] = CONSTANT.ListenerType
	order["status"] = CONSTANT.OrderWaiting
	order["created_at"] = UTIL.GetCurrentTime().String()
	// no paid amount and billing

	orderID, status, ok := DB.InsertWithUniqueID(CONSTANT.OrderClientAppointmentTable, CONSTANT.OrderDigits, order, "order_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["order_id"] = orderID
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

// ListenerOrderPaymentComplete godoc
// @Tags Client Listener
// @Summary Call after payment is completed for listener order
// @Router /client/listener/paymentcomplete [post]
// @Param body body model.ListenerOrderPaymentCompleteRequest true "Request Body"
// @Produce json
// @Success 200
func ListenerOrderPaymentComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, ok := UTIL.ReadRequestBody(r)
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ListenerOrderPaymentCompleteRequiredFields)
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
	// check if order is with listener
	if !strings.EqualFold(order[0]["type"], CONSTANT.ListenerType) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if order payment is already captured
	if !strings.EqualFold(order[0]["status"], CONSTANT.OrderWaiting) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, CONSTANT.PaymentCapturedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create appointment between listener and client
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

	// change order status
	orderUpdate := map[string]string{}
	orderUpdate["status"] = CONSTANT.OrderInProgress
	orderUpdate["modified_at"] = UTIL.GetCurrentTime().String()
	status, ok = DB.UpdateSQL(CONSTANT.OrderClientAppointmentTable,
		map[string]string{
			"order_id": body["order_id"],
		},
		orderUpdate,
	)

	// update listener slots
	DB.UpdateSQL(CONSTANT.SlotsTable,
		map[string]string{
			"counsellor_id": order[0]["counsellor_id"],
			"date":          order[0]["date"],
		},
		map[string]string{
			order[0]["time"]: CONSTANT.SlotBooked,
		},
	)

	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
}
