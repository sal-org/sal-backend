package counsellor

import (
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"

	UTIL "salbackend/util"
)

// PaymentsGet godoc
// @Tags Counsellor Payments
// @Summary Get payments for counsellor
// @Router /counsellor/payment [get]
// @Param counsellor_id query string true "Logged in counsellor ID"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func PaymentsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get payments for counsellor
	payments, status, ok := DB.SelectProcess("select * from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" order by created_at desc limit "+strconv.Itoa(CONSTANT.CounsellorsPaymentsPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsPaymentsPerPageClient), r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number payments for counsellor
	paymentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid, r.FormValue("counsellor_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	response["payments"] = payments
	response["payments_count"] = paymentsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(paymentsCount[0]["ctn"], CONSTANT.CounsellorsPaymentsPerPageClient))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
