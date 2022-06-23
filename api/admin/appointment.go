package admin

import (
	"net/http"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"

	_ "salbackend/model"
	UTIL "salbackend/util"
)

func AppointmentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get appointments
	wheres := []string{}
	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "state":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "client_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " client_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "counsellor_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " counsellor_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		case "appointment_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " appointment_id = ? ")
				queryArgs = append(queryArgs, val[0])
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get counsellor, client, order ids to get details
	clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
	counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")
	orderIDs := UTIL.ExtractValuesFromArrayMap(appointments, "order_id")

	// get client details
	clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get counsellor details
	counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get order details
	orders, status, ok := DB.SelectProcess("select order_id, paid_amount, slots_bought, invoice_id from " + CONSTANT.OrderClientAppointmentTable + " where order_id in ('" + strings.Join(orderIDs, "','") + "')")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// get invoice ids to get refund details
	invoiceIDs := UTIL.ExtractValuesFromArrayMap(orders, "invoice_id")

	// get refund details
	refunds, status, ok := DB.SelectProcess("select invoice_id, sum(refunded_amount) as total_refunded_amount from " + CONSTANT.RefundsTable + " where invoice_id in ('" + strings.Join(invoiceIDs, "','") + "') group by invoice_id")
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of appointments
	appointmentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.AppointmentsTable+where, queryArgs...)
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["appointments"] = appointments
	response["clients"] = UTIL.ConvertMapToKeyMap(clients, "client_id")
	response["counsellors"] = UTIL.ConvertMapToKeyMap(counsellors, "id")
	response["orders"] = UTIL.ConvertMapToKeyMap(orders, "order_id")
	response["refunds"] = UTIL.ConvertMapToKeyMap(refunds, "invoice_id")
	response["appointments_count"] = appointmentsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(appointmentsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))
	response["media_url"] = CONFIG.MediaURL

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func AppointmentRefund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get appointment details
	appointment, status, ok := DB.SelectSQL(CONSTANT.AppointmentsTable, []string{"*"}, map[string]string{"appointment_id": r.FormValue("appointment_id")})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	// check if appointment is valid
	if len(appointment) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.AppointmentNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// refund amount
	invoice, status, ok := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"actual_amount", "discount", "paid_amount", "payment_id", "invoice_id"}, map[string]string{"order_id": appointment[0]["order_id"]})
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(invoice) > 0 {
		// invoice is available => amount is paid, order is not free
		paidAmount, _ := strconv.ParseFloat(invoice[0]["paid_amount"], 64)
		refundedAmount, _ := strconv.ParseFloat(DB.QueryRowSQL("select sum(refunded_amount) from "+CONSTANT.RefundsTable+" where invoice_id = '"+invoice[0]["invoice_id"]+"'"), 64)
		refundAmount, _ := strconv.ParseFloat(r.FormValue("refund_amount"), 64)
		if refundAmount <= 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Cant refund less than 0 amount", CONSTANT.ShowDialog, response)
			return
		}
		if refundAmount+refundedAmount <= paidAmount {
			// refunded amount will be less than paid amount
			DB.InsertWithUniqueID(CONSTANT.RefundsTable, CONSTANT.RefundDigits, map[string]string{
				"invoice_id":      invoice[0]["invoice_id"],
				"payment_id":      invoice[0]["payment_id"],
				"refunded_amount": strconv.FormatFloat(refundAmount, 'f', 2, 64),
				"status":          CONSTANT.RefundInProgress,
				"created_at":      UTIL.GetCurrentTime().String(),
			}, "refund_id")
		} else {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, strconv.FormatFloat(refundedAmount, 'f', 2, 64)+" has already been refunded from "+invoice[0]["paid_amount"]+". So, "+r.FormValue("refund_amount")+" cant be refunded.", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "No invoice found", CONSTANT.ShowDialog, response)
		return
	}

	// update appointment status to refunded
	DB.UpdateSQL(CONSTANT.AppointmentsTable,
		map[string]string{
			"appointment_id": r.FormValue("appointment_id"),
		},
		map[string]string{
			"status":      CONSTANT.AppointmentAdminCancelled,
			"modified_at": UTIL.GetCurrentTime().String(),
		},
	)

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "Refund of "+r.FormValue("refund_amount")+" is initiated", CONSTANT.ShowDialog, response)
}
