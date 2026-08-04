package therapist

import (
	"encoding/csv"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strconv"
	"strings"
	"time"

	UTIL "salbackend/util"
)

// PaymentsGet godoc
// @Tags Therapist Payments
// @Summary Get payments for counsellor
// @Router /therapist/payment [get]
// @Param therapist_id query string true "Logged in therapist ID"
// @Param order_by query string false "Order by - 1(asc), 2(desc) - should be sent along with sort_by"
// @Param page query string false "Page number"
// @Security JWTAuth
// @Produce json
// @Success 200
func PaymentsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]any)

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	therapistID, ok := UTIL.Required(r.FormValue("therapist_id"), "Therapist ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, therapistID, CONSTANT.ShowDialog, response)
		return
	}

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	now := UTIL.GetCurrentTime()

	startOfDayUTC := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	dateRange := "created_at >='" + startOfDayUTC.Format("2006-01-02 15:04:05") + "'"

	orderBy := " desc "

	if strings.EqualFold(r.FormValue("order_by"), "1") {
		orderBy = " asc "
	}

	switch r.FormValue("filter_by") {
	case "1": // get payments for therapist
		dateRange = "created_at >='" + startOfDayUTC.Format("2006-01-02 15:04:05") + "'"

	case "7": // get payments for therapist in last 7 days
		dateRange = "created_at >='" + startOfDayUTC.AddDate(0, 0, -7).Format("2006-01-02 15:04:05") + "'"

	case "30": // get payments for therapist in last 30 days
		dateRange = "created_at >='" + startOfDayUTC.AddDate(0, 0, -30).Format("2006-01-02 15:04:05") + "'"

	case "90": // get payments for therapist in last 90 days
		dateRange = "created_at >='" + startOfDayUTC.AddDate(0, 0, -90).Format("2006-01-02 15:04:05") + "'"
	default:
		dateRange = "created_at >='" + startOfDayUTC.Format("2006-01-02 15:04:05") + "'"
	}

	// get payments for therapist
	payments, status, ok := DB.SelectProcess("select * from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" and "+dateRange+" order by created_at "+orderBy+" limit "+strconv.Itoa(CONSTANT.CounsellorsPaymentsPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsPaymentsPerPageClient), r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	for index, payment := range payments {

		if len(payment["appointment_id"]) > 0 {
			var statusMessage string
			appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where appointment_id = ? ", payment["appointment_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellors details
			counsellorRecords, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where appointment_id = ? ", payment["appointment_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			if len(counsellorRecords) == 0 {

				counsellorRecordsNewVersion, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? ", payment["appointment_id"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				if len(counsellorRecordsNewVersion) == 0 {
					payments[index]["is_counsellor_record_filled"] = "Pending"
				} else {

					if counsellorRecordsNewVersion[0]["status"] == "2" {
						payments[index]["is_counsellor_record_filled"] = "Completed"
					} else {
						payments[index]["is_counsellor_record_filled"] = "Pending"
					}
				}
			} else {
				payments[index]["is_counsellor_record_filled"] = "Completed"
			}

			switch appointments[0]["status"] {
			case "3":
				if appointments[0]["started_at"] == "" && appointments[0]["ended_at"] == "" {
					statusMessage = getAppointmentStatusInText("8")
					payments[index]["is_counsellor_record_filled"] = "Not Applicable"
				} else if appointments[0]["client_started_at"] == "" && appointments[0]["client_ended_at"] == "" {
					statusMessage = getAppointmentStatusInText("7")
				} else if len(appointments[0]["client_started_at"]) != 0 && len(appointments[0]["client_ended_at"]) == 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["client_started_at"]) == 0 && len(appointments[0]["client_ended_at"]) != 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["started_at"]) != 0 && len(appointments[0]["ended_at"]) == 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["started_at"]) == 0 && len(appointments[0]["ended_at"]) != 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else {
					if UTIL.BuildToDteTime(appointments[0]["ended_at"]).Sub(UTIL.BuildToDteTime(appointments[0]["started_at"])).Minutes() > 10 {
						if UTIL.BuildToDteTime(appointments[0]["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointments[0]["client_started_at"])).Minutes() < 10 {
							statusMessage = getAppointmentStatusInText("19")
						} else {
							statusMessage = getAppointmentStatusInText("3")
						}
					} else {
						statusMessage = getAppointmentStatusInText("14")
					}
				}

			case "4":
				statusMessage = getAppointmentStatusInText("4")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			case "5":
				statusMessage = getAppointmentStatusInText("5")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
				// case "6":
				// 	statusMessage = getAppointmentStatusInText("6")
				// 	payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			case "1":
				statusMessage = getAppointmentStatusInText("8")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			default:
				statusMessage = getAppointmentStatusInText(appointments[0]["status"])
			}

			payments[index]["status_text"] = statusMessage
			payments[index]["date"] = appointments[0]["date"]
			payments[index]["time"] = appointments[0]["time"]
		} else {
			payments[index]["is_counsellor_record_filled"] = "-"
			payments[index]["status_text"] = "-"
			payments[index]["date"] = "-"
			payments[index]["time"] = "-"
		}
	}

	// get payments for therapist
	totalPayments, status, ok := DB.SelectProcess("select sum(amount) as total from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" and created_at >= ? and created_at <= ? ", r.FormValue("therapist_id"), UTIL.FirstDayOfMonth(UTIL.GetCurrentTime()).Format("2006-01-02 15:04:05"), UTIL.LastDayOfMonth(UTIL.GetCurrentTime()).Format("2006-01-02 15:04:05"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number payments for therapist
	paymentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" and "+dateRange, r.FormValue("therapist_id"))
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// old code for getting payments for therapist

	// orderBy := " desc "

	// if strings.EqualFold(r.FormValue("order_by"), "1") {
	// 	orderBy = " asc "
	// }

	// // get payments for therapist
	// payments, status, ok := DB.SelectProcess("select * from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" order by created_at "+orderBy+" limit "+strconv.Itoa(CONSTANT.CounsellorsPaymentsPerPageClient)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.CounsellorsPaymentsPerPageClient), r.FormValue("therapist_id"))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// // get total number payments for therapist
	// paymentsCount, status, ok := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid, r.FormValue("therapist_id"))
	// if !ok {
	// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	response["payments"] = payments
	response["payments_count"] = paymentsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(paymentsCount[0]["ctn"], CONSTANT.CounsellorsPaymentsPerPageClient))
	response["total_amount"] = totalPayments[0]["total"]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func PaymentsDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startBy, _ := time.Parse("2006-01-02", r.FormValue("start_by"))
	endBy, _ := time.Parse("2006-01-02", r.FormValue("end_by"))

	fileName := "payment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	var response = make(map[string]any)
	data := [][]string{}

	// check if access token is valid, not expired
	if !UTIL.CheckIfAccessTokenExpired(r.Header.Get("Authorization")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, response)
		return
	}

	therapistID, ok := UTIL.Required(r.FormValue("therapist_id"), "Therapist ID")
	if !ok {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, therapistID, CONSTANT.ShowDialog, response)
		return
	}

	// check if therapist_id exists
	if !DB.CheckIfExists(CONSTANT.TherapistsTable, map[string]string{"therapist_id": therapistID}) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid id", CONSTANT.ShowDialog, response)
		return
	}

	heading := []string{"Client Name", "Date", "Time", "Counsellor Record", "Status", "Amount"}

	// get payments for therapist
	payments, status, ok := DB.SelectProcess("select * from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" and created_at >= ? and created_at <= ? order by created_at desc", r.FormValue("therapist_id"), startBy.UTC().String(), endBy.UTC().String())
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(payments) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.StatusCodeNoDataFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get payments for therapist
	totalPayments, status, ok := DB.SelectProcess("select sum(amount) as total from "+CONSTANT.PaymentsTable+" where counsellor_id = ? and status = "+CONSTANT.PaymentValid+" and created_at >= ? and created_at <= ? ", r.FormValue("therapist_id"), startBy.UTC().String(), endBy.UTC().String())
	if !ok {
		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		return
	}

	// for _, payment := range payments {
	// 	data = append(data, []string{payment["heading"], payment["description"], payment["amount"]})
	// }

	for index, payment := range payments {

		if len(payment["appointment_id"]) > 0 {
			var statusMessage string
			appointments, status, ok := DB.SelectProcess("select * from "+CONSTANT.AppointmentsTable+" where appointment_id = ? ", payment["appointment_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			// get counsellors details
			counsellorRecords, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where appointment_id = ? ", payment["appointment_id"])
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			if len(counsellorRecords) == 0 {

				counsellorRecordsNewVersion, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsFormLastestVersionTable+" where appointment_id = ? ", payment["appointment_id"])
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				if len(counsellorRecordsNewVersion) == 0 {
					payments[index]["is_counsellor_record_filled"] = "Pending"
				} else {

					if counsellorRecordsNewVersion[0]["status"] == "2" {
						payments[index]["is_counsellor_record_filled"] = "Completed"
					} else {
						payments[index]["is_counsellor_record_filled"] = "Pending"
					}
				}
			} else {
				payments[index]["is_counsellor_record_filled"] = "Completed"
			}

			switch appointments[0]["status"] {
			case "3":
				if appointments[0]["started_at"] == "" && appointments[0]["ended_at"] == "" {
					statusMessage = getAppointmentStatusInText("8")
					payments[index]["is_counsellor_record_filled"] = "Not Applicable"
				} else if appointments[0]["client_started_at"] == "" && appointments[0]["client_ended_at"] == "" {
					statusMessage = getAppointmentStatusInText("7")
				} else if len(appointments[0]["client_started_at"]) != 0 && len(appointments[0]["client_ended_at"]) == 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["client_started_at"]) == 0 && len(appointments[0]["client_ended_at"]) != 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["started_at"]) != 0 && len(appointments[0]["ended_at"]) == 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else if len(appointments[0]["started_at"]) == 0 && len(appointments[0]["ended_at"]) != 0 {
					statusMessage = getAppointmentStatusInText("14")
				} else {
					if UTIL.BuildToDteTime(appointments[0]["ended_at"]).Sub(UTIL.BuildToDteTime(appointments[0]["started_at"])).Minutes() > 10 {
						if UTIL.BuildToDteTime(appointments[0]["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointments[0]["client_started_at"])).Minutes() < 10 {
							statusMessage = getAppointmentStatusInText("19")
						} else {
							statusMessage = getAppointmentStatusInText("3")
						}
					} else {
						statusMessage = getAppointmentStatusInText("14")
					}
				}

			case "4":
				statusMessage = getAppointmentStatusInText("4")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			case "5":
				statusMessage = getAppointmentStatusInText("5")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
				// case "6":
				// 	statusMessage = getAppointmentStatusInText("6")
				// 	payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			case "1":
				statusMessage = getAppointmentStatusInText("8")
				payments[index]["is_counsellor_record_filled"] = "Not Applicable"
			default:
				statusMessage = getAppointmentStatusInText(appointments[0]["status"])
			}

			payments[index]["status_text"] = statusMessage
			payments[index]["date"] = appointments[0]["date"]
			payments[index]["time"] = UTIL.GetTimeFromTimeSlotIN12Hour(appointments[0]["time"])

			data = append(data, []string{payment["heading"], payment["date"], payment["time"], payment["is_counsellor_record_filled"], payment["status_text"], payment["amount"]})
		} else {
			payments[index]["is_counsellor_record_filled"] = ""
			payments[index]["status_text"] = ""
			payments[index]["date"] = ""
			payments[index]["time"] = ""

			data = append(data, []string{payment["heading"], "", "", "", "", payment["amount"]})
		}

	}

	w.Header().Set("Content-Type", "text/csv")

	data = append(data, []string{"Total Amount", "", "", "", "", totalPayments[0]["total"]})

	writer := csv.NewWriter(w)

	writer.Write(heading)

	for _, d := range data {
		writer.Write(d)
	}

	writer.Flush()
}
