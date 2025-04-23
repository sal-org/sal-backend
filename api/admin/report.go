package admin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strconv"
	"strings"
	"time"

	UTIL "salbackend/util"
)

func ReportGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=reports.csv")

	var response = make(map[string]interface{})

	startBy, _ := time.Parse("2006-01-02", r.FormValue("start_by"))
	endBy, _ := time.Parse("2006-01-02", r.FormValue("end_by"))

	heading := []string{}
	data := [][]string{}

	switch r.FormValue("id") {
	case "1": // appointment report

		heading = []string{"Client Name", "Gender", "Age", "Company Name", "IsFamilyMember", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date & Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Client Start", "Client End", "Mod. At", "Status"}
		appointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, appointment := range appointments {

			var startTime, endTime, modAt, clientStartTime, clientEndTime, partnerName, isFamilyMember, status, location string

			if appointment["started_at"] == "" {
				startTime = ""
			} else {
				startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["client_started_at"] == "" {
				clientStartTime = ""
			} else {
				clientStartTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_started_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["client_ended_at"] == "" {
				clientEndTime = ""
			} else {
				clientEndTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_ended_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["ended_at"] == "" {
				endTime = ""
			} else {
				endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["modified_at"] == "" {
				modAt = ""
			} else {
				modAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["status"] == "3" {
				if appointment["started_at"] == "" && appointment["ended_at"] == "" {
					status = getAppointmentStatusInText("8")
				} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
					status = getAppointmentStatusInText("7")
				} else {
					if UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes() > 11 {
						status = getAppointmentStatusInText("3")
					} else {
						status = getAppointmentStatusInText("14")
					}

				}
			} else if appointment["status"] == "4" {
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("12")
				} else {
					status = getAppointmentStatusInText("4")
				}
			} else if appointment["status"] == "5" {
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("13")
				} else {
					status = getAppointmentStatusInText("5")
				}
			} else {
				status = getAppointmentStatusInText(appointment["status"])
			}

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
				isFamilyMember = "No"
			} else {
				// get client details
				clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where client_id = '" + appointment["client_id"] + "'")
				if !ok {
					UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					return
				}

				partnerN := ""

				if clients[0]["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			if clientsMap[appointment["client_id"]]["location"] == "40.04" {
				location = ""
			} else {
				location = clientsMap[appointment["client_id"]]["location"]
			}

			data = append(data, []string{
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				clientsMap[appointment["client_id"]]["gender"],
				clientsMap[appointment["client_id"]]["age"],
				partnerName,
				isFamilyMember,
				location,
				clientsMap[appointment["client_id"]]["department"],
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				counsellorsMap[appointment["counsellor_id"]]["type"],
				UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment["date"], appointment["time"]), "0").Format(CONSTANT.ReadbleDateTimeFormat),
				appointment["times_rescheduled"],
				startTime,
				endTime,
				clientStartTime,
				clientEndTime,
				modAt,
				status,
			})
		}
	case "2": // sales report
		heading = []string{"User ID", "User Name", "User Type", "Total Individual Sessions amount", "Total SAL Cafe sessions amount", "Net Amount Received", "Refund Amount", "Cancellation Amount", "No Show Amount"}
		invoices, status, ok := DB.SelectProcess("select user_id, sum(CASE WHEN order_type = " + CONSTANT.OrderAppointmentType + " THEN paid_amount ELSE 0 END) as total_session_amount, sum(CASE WHEN order_type = " + CONSTANT.OrderEventBookType + " THEN paid_amount ELSE 0 END) as total_event_amount from " + CONSTANT.InvoicesTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' group by user_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get user ids to get details
		userIDs := UTIL.ExtractValuesFromArrayMap(invoices, "user_id")

		refunds, status, ok := DB.SelectProcess("select user_id, sum(refunded_amount) as total_refund_amount, sum(CASE WHEN type = " + CONSTANT.RefundCancellationType + " THEN refunded_amount ELSE 0 END) as total_cancel_refund_amount, sum(CASE WHEN type = " + CONSTANT.RefundNoShowType + " THEN refunded_amount ELSE 0 END) as total_no_show_refund_amount from (select r.refunded_amount, r.type, i.user_id from " + CONSTANT.RefundsTable + " r left join " + CONSTANT.InvoicesTable + " i on r.invoice_id = i.invoice_id where r.created_at > '" + startBy.UTC().String() + "' and r.created_at < '" + endBy.UTC().String() + "') as a group by user_id")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get user ids to get details
		userIDs = append(userIDs, UTIL.ExtractValuesFromArrayMap(invoices, "user_id")...)

		// get user details
		users, status, ok := DB.SelectProcess("(select client_id as id, first_name, last_name, 'Client' as type, email, phone from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "')) union(select counsellor_id as id, first_name, last_name, 'Counsellor' as type, email, phone from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(userIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type, email, phone from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(userIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type, email, phone from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(userIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		invoicesMap := UTIL.ConvertMapToKeyMap(invoices, "user_id")
		refundsMap := UTIL.ConvertMapToKeyMap(refunds, "user_id")

		for _, user := range users {
			data = append(data, []string{
				user["id"],
				user["first_name"],
				user["type"],
				invoicesMap[user["id"]]["total_session_amount"],
				invoicesMap[user["id"]]["total_event_amount"],
				strconv.FormatFloat(getInt(invoicesMap[user["id"]]["total_session_amount"])+getInt(invoicesMap[user["id"]]["total_event_amount"]), 'f', 2, 64),
				refundsMap[user["id"]]["total_refund_amount"],
				refundsMap[user["id"]]["total_cancel_refund_amount"],
				refundsMap[user["id"]]["total_no_show_refund_amount"],
			})
		}
	case "3": // booking report
		heading = []string{"Booking ID", "Client ID", "Client Name", "Counsellor ID", "Counsellor Name", "Counsellor Type", "Date & Time of Session Booking", "Total Sessions Bought", "Session Remaining", "Bulk Cancel (4 - cancel)"}
		bookings, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppointmentSlotsTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(bookings, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(bookings, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, booking := range bookings {
			data = append(data, []string{
				booking["order_id"],
				booking["client_id"],
				clientsMap[booking["client_id"]]["first_name"],
				booking["counsellor_id"],
				counsellorsMap[booking["counsellor_id"]]["first_name"],
				counsellorsMap[booking["counsellor_id"]]["type"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(booking["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
				booking["slots_bought"],
				booking["slots_remaining"],
				booking["status"],
			})
		}
	case "4": // sal cafe report
		heading = []string{"Booking ID", "Client ID", "Client Name", "Client Email", "Client Mobile", "Counsellor ID", "Counsellor Name", "Counsellor Type", "Topic", "Date & Time of Booking", "Paid Amount"}
		bookings, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderEventTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' and status = " + CONSTANT.OrderInProgress + " order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get user, event ids to get details
		userIDs := UTIL.ExtractValuesFromArrayMap(bookings, "user_id")
		eventOrderIDs := UTIL.ExtractValuesFromArrayMap(bookings, "event_order_id")

		// get event details
		events, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventTable + " where order_id in ('" + strings.Join(eventOrderIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor ids to get details
		userIDs = append(userIDs, UTIL.ExtractValuesFromArrayMap(events, "counsellor_id")...)

		// get user details
		users, status, ok := DB.SelectProcess("(select client_id as id, first_name, last_name, 'Client' as type, email, phone from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "')) union(select counsellor_id as id, first_name, last_name, 'Counsellor' as type, email, phone from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(userIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type, email, phone from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(userIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type, email, phone from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(userIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get topics
		topics, status, ok := DB.SelectProcess("select * from " + CONSTANT.TopicsTable)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		usersMap := UTIL.ConvertMapToKeyMap(users, "id")
		eventsMap := UTIL.ConvertMapToKeyMap(events, "order_id")
		topicsMap := UTIL.ConvertMapToKeyMap(topics, "id")

		for _, booking := range bookings {
			data = append(data, []string{
				booking["order_id"],
				booking["user_id"],
				usersMap[booking["user_id"]]["first_name"],
				usersMap[booking["user_id"]]["email"],
				usersMap[booking["user_id"]]["phone"],
				eventsMap[booking["event_order_id"]]["counsellor_id"],
				usersMap[eventsMap[booking["event_order_id"]]["counsellor_id"]]["first_name"],
				usersMap[eventsMap[booking["event_order_id"]]["counsellor_id"]]["type"],
				topicsMap[eventsMap[booking["event_order_id"]]["topic_id"]]["topic"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(booking["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
				booking["paid_amount"],
			})
		}
	case "5": // finance report
		heading = []string{"Invoice ID", "User ID", "User Type", "User Name", "User Email", "User Mobile", "Date & Time of Booking", "Paid Amount", "CGST Amount", "SGST Amount", "Total Tax Amount"}
		invoices, status, ok := DB.SelectProcess("select * from " + CONSTANT.InvoicesTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get user ids to get details
		userIDs := UTIL.ExtractValuesFromArrayMap(invoices, "user_id")

		// get user details
		users, status, ok := DB.SelectProcess("(select client_id as id, first_name, last_name, 'Client' as type, email, phone from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "')) union(select counsellor_id as id, first_name, last_name, 'Counsellor' as type, email, phone from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(userIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type, email, phone from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(userIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type, email, phone from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(userIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		usersMap := UTIL.ConvertMapToKeyMap(users, "id")

		for _, invoice := range invoices {
			data = append(data, []string{
				invoice["invoice_id"],
				invoice["user_id"],
				usersMap[invoice["user_id"]]["type"],
				usersMap[invoice["user_id"]]["first_name"],
				usersMap[invoice["user_id"]]["email"],
				usersMap[invoice["user_id"]]["phone"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(invoice["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
				invoice["paid_amount"],
				invoice["cgst"],
				invoice["sgst"],
				invoice["tax"],
			})
		}
	case "6": // payout report
		heading = []string{"Counsellor ID", "Counsellor Name", "Counsellor Type", "Heading", "Description", "Date & Time of Session", "Amount to be paid", "Beneficiary Name", "Bank Name", "Account Type", "IFSC Code", "Bank A/c Number"}
		payments, status, ok := DB.SelectProcess("select * from " + CONSTANT.PaymentsTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor ids to get details
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(payments, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type, payee_name, bank_account_no, ifsc, bank_name, bank_account_type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type, payee_name, bank_account_no, ifsc, bank_name, bank_account_type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, payment := range payments {
			data = append(data, []string{
				payment["counsellor_id"],
				counsellorsMap[payment["counsellor_id"]]["first_name"],
				counsellorsMap[payment["counsellor_id"]]["type"],
				payment["heading"],
				payment["description"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(payment["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
				payment["amount"],
				counsellorsMap[payment["counsellor_id"]]["payee_name"],
				counsellorsMap[payment["counsellor_id"]]["bank_name"],
				counsellorsMap[payment["counsellor_id"]]["bank_account_type"],
				counsellorsMap[payment["counsellor_id"]]["ifsc"],
				counsellorsMap[payment["counsellor_id"]]["bank_account_no"],
			})
		}
	case "7": // promo code report
		heading = []string{"Promo Code", "Description", "No. of times used", "Promo code used in amt.", "Amount Received"}
		invoices, status, ok := DB.SelectProcess("select coupon_code, count(*) as ctn, sum(discount) as used_amount, sum(paid_amount) as paid_amount from " + CONSTANT.InvoicesTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' and coupon_code != '' group by coupon_code")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get coupon codes to get details
		couponCodes := UTIL.ExtractValuesFromArrayMap(invoices, "coupon_code")

		coupons, status, ok := DB.SelectProcess("select * from " + CONSTANT.CouponsTable + " where coupon_code in ('" + strings.Join(couponCodes, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		couponsMap := UTIL.ConvertMapToKeyMap(coupons, "coupon_code")

		for _, invoice := range invoices {
			data = append(data, []string{
				invoice["coupon_code"],
				couponsMap[invoice["coupon_code"]]["description"],
				invoice["ctn"],
				invoice["used_amount"],
				invoice["paid_amount"],
			})
		}
	case "8": // push notification report
		heading = []string{"Notification Type", "Date", "Times"}
		notifications, status, ok := DB.SelectProcess("select date(created_at) as date, CASE WHEN notification_type = 1 THEN 'Promo' WHEN notification_type = 2 THEN 'Content' WHEN notification_type = 3 THEN 'Event' WHEN notification_type = 4 THEN 'Other' ELSE 'Nothing' END as notification_type, count(*) as ctn from " + CONSTANT.NotificationsBulkTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' group by date(created_at), notification_type")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, notification := range notifications {
			data = append(data, []string{
				notification["date"],
				notification["notification_type"],
				notification["ctn"],
			})
		}
	case "10": // onboarding report

		heading = []string{"First Name", "Last Name", "Gender", "Email", "Phone", "Type", "Created At"}
		counsellors, status, ok := DB.SelectProcess("(select first_name, last_name, gender, email, phone, 'Counsellor' as `type`, created_at from " + CONSTANT.CounsellorsTable + " where status = " + CONSTANT.CounsellorActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Listener' as `type`, created_at from " + CONSTANT.ListenersTable + " where status = " + CONSTANT.ListenerActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Therapist' as `type`, created_at from " + CONSTANT.TherapistsTable + " where status = " + CONSTANT.TherapistActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, counsellor := range counsellors {
			data = append(data, []string{counsellor["first_name"], counsellor["last_name"], counsellor["gender"], counsellor["email"], counsellor["phone"], counsellor["type"], UTIL.ConvertTimezone(UTIL.ConvertToTime(counsellor["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)})
		}
	case "11": // client onboarding report

		var partnerName, location, isFamilyMember string

		heading = []string{"Client Name", "Company Name", "IsFamilyMember", "Age", "Gender", "Location", "CreatedAt"}

		clients, status, ok := DB.SelectProcess("select asscoiate_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, created_at from " + CONSTANT.ClientsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, client := range clients {

			domainName := strings.Split(client["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
			} else {
				partnerN := ""

				if client["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + client["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			if client["location"] == "40.04" {
				location = ""
			} else {
				location = client["location"]
			}

			data = append(data, []string{client["first_name"] + " " + client["last_name"], partnerName, isFamilyMember, client["age"], client["gender"], location, UTIL.ConvertTimezone(UTIL.ConvertToTime(client["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)})
		}

	case "12": // client mood report

		var partnerName, isFamilyMember, location string

		heading = []string{"Client Name", "Mood", "Age", "Gender", "Company Name", "IsFamilyMember", "Location", "Department"}

		moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodResultsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(moods, "client_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")

		for _, mood := range moods {

			domainName := strings.Split(clientsMap[mood["client_id"]]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[mood["client_id"]]["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[mood["client_id"]]["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			if clientsMap[mood["client_id"]]["location"] == "40.04" {
				location = ""
			} else {
				location = clientsMap[mood["client_id"]]["location"]
			}

			moodTitle := DB.QueryRowSQL("select title from "+CONSTANT.MoodsTable+" where id = ?", mood["mood_id"])

			data = append(data, []string{mood["name"], moodTitle, mood["age"], mood["gender"], partnerName, isFamilyMember, location, clientsMap[mood["client_id"]]["department"]})
		}
	case "13": // client asessment report

		var partnerName, isFamilyMember, location, department string

		heading = []string{"Client Name", "Assessment", "Age", "Gender", "Company Name", "IsFamilyMember", "Location", "Department", "CreatedAt"}

		assessments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AssessmentResultsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(assessments, "user_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")

		for _, assessment := range assessments {

			if len(clientsMap[assessment["user_id"]]) != 0 {
				domainName := strings.Split(clientsMap[assessment["user_id"]]["email"], "@")

				title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

				if len(title) > 0 {
					partnerName = title[0]["partner_name"]
					isFamilyMember = "No"
				} else {
					partnerN := ""

					if clientsMap[assessment["user_id"]]["asscoiate_id"] != "" {
						// get client details
						clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[assessment["user_id"]]["asscoiate_id"] + "'")
						if !ok {
							UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
							return
						}

						domainName := strings.Split(clients[0]["email"], "@")

						title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

						partnerN = title[0]["partner_name"]
						isFamilyMember = "Yes"
					} else {
						partnerN = "None"
						isFamilyMember = "No"
					}

					partnerName = partnerN
				}

				if clientsMap[assessment["user_id"]]["location"] == "40.04" {
					location = ""
				} else {
					location = clientsMap[assessment["user_id"]]["location"]
				}

				department = clientsMap[assessment["user_id"]]["department"]

			} else {
				partnerName = "None"
				location = ""
				department = ""

			}

			assessmentTitle := DB.QueryRowSQL("select title from "+CONSTANT.AssessmentsTable+" where assessment_id = ?", assessment["assessment_id"])

			data = append(data, []string{assessment["name"], assessmentTitle, assessment["age"], assessment["gender"], partnerName, isFamilyMember, location, department, UTIL.ConvertTimezone(UTIL.ConvertToTime(assessment["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)})
		}
	case "14": // counsellor record report
		var noShow string

		heading = []string{"Session For", "Session Type", "Counsellor Name", "Client Name", "Age", "Gender", "Location", "Department", "No Show", "Session Mode", "Session No", "Session Date", "In-Time", "Out-Time", "Presenting Concerns", "Mental Health Scale", "Psychiatric Intervention", "Psychiatric Intervention Reason", "Therapy Notes", "Category", "Sub Category", "Emotional State", "Next Follow Date", "Client Notes", "Client Self-Work", "Client Assessment", "Therapeutic Plan"}

		counsellorRecords, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsTable + " where `session_date` >= '" + startBy.String() + "' and `session_date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecords, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, counsellorRecord := range counsellorRecords {
			if counsellorRecord["noshow"] == "0" {
				noShow = "No"
			} else {
				noShow = "Yes"
			}
			data = append(data, []string{counsellorRecord["session_for"], counsellorRecord["session_type"], counsellorsMap[counsellorRecord["counsellor_id"]]["first_name"] + " " + counsellorsMap[counsellorRecord["counsellor_id"]]["last_name"], counsellorRecord["client_first_name"] + " " + counsellorRecord["client_last_name"], counsellorRecord["client_age"], counsellorRecord["client_gender"], counsellorRecord["client_location"], counsellorRecord["client_department"], noShow, counsellorRecord["session_mode"], counsellorRecord["session_no"], counsellorRecord["session_date"], counsellorRecord["in_time"], counsellorRecord["out_time"], counsellorRecord["presenting_concerns"], counsellorRecord["mental_health"], counsellorRecord["psychiatric_intervention"], counsellorRecord["psychiatric_intervention_reason"], counsellorRecord["therapy_notes"], counsellorRecord["therapeutic_goal"], counsellorRecord["sub_category"], counsellorRecord["emotional_state"], counsellorRecord["next_follow_date"], counsellorRecord["client_notes"], counsellorRecord["client_documents"], counsellorRecord["assessment_tool"], counsellorRecord["therapy_plan"]})
		}

	case "15": // client virtual appointment rating
		heading = []string{"Counsellor Name", "Client Name", "Gender", "Date", "Rate", "Rate Type", "Comment"}

		appointments, status, ok := DB.SelectProcess("select counsellor_id, client_id, `date`, rating, rating_types, rating_comment  from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, appointment := range appointments {
			data = append(data, []string{
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				clientsMap[appointment["client_id"]]["first_name"],
				clientsMap[appointment["client_id"]]["gender"],
				appointment["date"],
				appointment["rating"],
				appointment["rating_types"],
				appointment["rating_comment"],
			})
		}
	case "16": // content
		heading = []string{"Title", "Type", "Category", "Mood", "Is Training Content", "Created By", "Uploaded At", "Status"}

		contents, status, ok := DB.SelectProcess("select *  from " + CONSTANT.ContentsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		moodIDs := UTIL.ExtractValuesFromArrayMap(contents, "mood_id")
		categoryIDs := UTIL.ExtractValuesFromArrayMap(contents, "category_id")

		// get mood details
		moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodsTable + " where id in ('" + strings.Join(moodIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get mood details
		categorys, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentCategoriesTable + " where id in ('" + strings.Join(categoryIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		moodsMap := UTIL.ConvertMapToKeyMap(moods, "id")
		categorysMap := UTIL.ConvertMapToKeyMap(categorys, "id")

		for _, content := range contents {

			var types, trainingType, status string
			if content["type"] == "1" {
				types = "Video"
			} else if content["type"] == "2" {
				types = "Audio"
			} else {
				types = "Article"
			}

			if content["training"] == "1" {
				trainingType = "Yes"
			} else {
				trainingType = "No"
			}

			if content["status"] == "1" {
				status = "Active"
			} else {
				status = "Inactive"
			}

			data = append(data, []string{
				content["title"],
				types,
				categorysMap[content["category_id"]]["category"],
				moodsMap[content["mood_id"]]["title"],
				trainingType,
				content["created_by"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(content["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
				status,
			})
		}

	case "17":  // counsellor time sheet
		heading = []string{"Counsellor Name", "Company Name", "Location", "In-Time", "Out-Time"}

		myTimeSheets, status, ok := DB.SelectProcess("select *  from " + CONSTANT.CounsellorMyTimeSheetTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(myTimeSheets, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, timeSheet := range myTimeSheets {

			data = append(data, []string{
				counsellorsMap[timeSheet["counsellor_id"]]["first_name"] + " " + counsellorsMap[timeSheet["counsellor_id"]]["last_name"],
				timeSheet["partner_name"],
				timeSheet["location"],
				timeSheet["inTime"],
				timeSheet["outTime"],
			})
		}
	case "18":  // total number of session
		heading = []string{"Client Name", "Company Name", "IsFamilyMember", "Total Number Of Session"}

		var partnerName, isFamilyMember string

		appointments, status, ok := DB.SelectProcess("select client_id, count(client_id) as total from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' and type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null) group by client_id ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")

		for _, appointment := range appointments {

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[appointment["client_id"]]["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			data = append(data, []string{
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				partnerName,
				isFamilyMember,
				appointment["total"],
			})
		}

	case "19": // total counsellor session
		heading = []string{"Counsellor Name", "Total Number Of Session"}

		appointments, status, ok := DB.SelectProcess("select counsellor_id, count(counsellor_id) as total from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' and type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null) group by counsellor_id ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, appointment := range appointments {

			data = append(data, []string{
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				appointment["total"],
			})
		}

	case "20": // total number session month and years wise
		heading = []string{"Client Name", "Company Name", "IsFamilyMember", "Total Number Of Session", "Month", "Year"}

		var partnerName, isFamilyMember string

		appointments, status, ok := DB.SelectProcess("select client_id, count(client_id) as total, monthname(date) as month, year(date) as year from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' and type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null) group by client_id having count(client_id) = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")

		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")

		for _, appointment := range appointments {

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[appointment["client_id"]]["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			data = append(data, []string{
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				partnerName,
				isFamilyMember,
				appointment["total"],
				appointment["month"],
				appointment["year"],
			})
		}

	case "21": // in-person appointment report

		heading = []string{"Client Name", "Gender", "Age", "Company Name", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date & Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Mod. At", "Status"}
		appointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonAppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, appointment := range appointments {

			var startTime, endTime, modAt, partnerName, status, location string

			if appointment["started_at"] == "" {
				startTime = ""
			} else {
				startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["ended_at"] == "" {
				endTime = ""
			} else {
				endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["modified_at"] == "" {
				modAt = ""
			} else {
				modAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
			}

			if appointment["status"] == "3" {
				if appointment["started_at"] == "" && appointment["ended_at"] == "" {
					status = getAppointmentStatusInText("8")
				} else if appointment["started_at"] != "" && appointment["ended_at"] == "" {
					status = getAppointmentStatusInText("14")
				} else {
					status = getAppointmentStatusInText("3")
				}
			} else if appointment["status"] == "4" {
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("12")
				} else {
					status = getAppointmentStatusInText("4")
				}
			} else if appointment["status"] == "5" {
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("13")
				} else {
					status = getAppointmentStatusInText("5")
				}
			} else {
				status = getAppointmentStatusInText(appointment["status"])
			}

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
			} else {
				partnerName = "None"
			}

			if clientsMap[appointment["client_id"]]["location"] == "40.04" {
				location = ""
			} else {
				location = clientsMap[appointment["client_id"]]["location"]
			}

			data = append(data, []string{
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				clientsMap[appointment["client_id"]]["gender"],
				clientsMap[appointment["client_id"]]["age"],
				partnerName,
				location,
				clientsMap[appointment["client_id"]]["department"],
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				counsellorsMap[appointment["counsellor_id"]]["type"],
				UTIL.ConvertTimezone(UTIL.BuildDateTime(appointment["date"], appointment["time"]), "0").Format(CONSTANT.ReadbleDateTimeFormat),
				appointment["times_rescheduled"],
				startTime,
				endTime,
				modAt,
				status,
			})
		}

	case "22": // inperson request appointment
		heading = []string{"Client Name", "Counsellor Name", "Company Name", "Location", "Date", "Status"}
		inPersonAppointmentsRequests, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonAppointmentRequestTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		status1 := ""

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(inPersonAppointmentsRequests, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(inPersonAppointmentsRequests, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, inpersonappointmentRequest := range inPersonAppointmentsRequests {

			if inpersonappointmentRequest["status"] == "1" {
				status1 = "InProgress"
			} else if inpersonappointmentRequest["status"] == "2" {
				status1 = "Completed"
			}

			date := UTIL.BuildDate(inpersonappointmentRequest["created_at"])

			data = append(data, []string{
				clientsMap[inpersonappointmentRequest["client_id"]]["first_name"] + " " + clientsMap[inpersonappointmentRequest["client_id"]]["last_name"],
				counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["first_name"] + " " + counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["last_name"],
				inpersonappointmentRequest["company_name"],
				inpersonappointmentRequest["company_location"],
				date,
				status1,
			})
		}

	case "23": // appointment request
		heading = []string{"Client Name", "Company Name", "IsFamilyMember", "Counsellor Name", "Date", "Status"}
		appointmentsRequests, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppointmentRequestTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		status1, partnerName, isFamilyMember := "", "", ""

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointmentsRequests, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointmentsRequests, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")

		for _, inpersonappointmentRequest := range appointmentsRequests {

			// get client details
			clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + inpersonappointmentRequest["client_id"] + "'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			domainName := strings.Split(clients[0]["email"], "@")

			title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(title) > 0 {
				partnerName = title[0]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clients[0]["asscoiate_id"] != "" {
					// get client details
					clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
					if !ok {
						UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						return
					}

					domainName := strings.Split(clients[0]["email"], "@")

					title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = title[0]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			if inpersonappointmentRequest["status"] == "1" {
				status1 = "InProgress"
			} else if inpersonappointmentRequest["status"] == "2" {
				status1 = "Completed"
			}

			date := UTIL.BuildDate(inpersonappointmentRequest["created_at"])

			data = append(data, []string{
				clientsMap[inpersonappointmentRequest["client_id"]]["first_name"] + " " + clientsMap[inpersonappointmentRequest["client_id"]]["last_name"],
				partnerName,
				isFamilyMember,
				counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["first_name"] + " " + counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["last_name"],
				date,
				status1,
			})
		}

	case "24": // inperson cafe appointment
		heading = []string{"Cafe Name", "Counsellor Name", "Client Name", "Attended", "Company Name", "Company Location", "Question1", "Question2", "Question3", "Question4", "Question5", "Cancellation Reason", "User Status", "Event Status", "Created At"}

		inPersonCafeAppointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderEventInPersonTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		userIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeAppointments, "user_id")
		eventOrderIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeAppointments, "event_order_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		orders, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventInPersonTable + " where order_id in ('" + strings.Join(eventOrderIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(orders, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		ordersMap := UTIL.ConvertMapToKeyMap(orders, "order_id")

		for _, inPersonCafeAppointment := range inPersonCafeAppointments {

			var attended, eventStatus, status string

			if inPersonCafeAppointment["attended"] == "1" {
				attended = "Yes"
			} else {
				attended = "No"
			}

			if inPersonCafeAppointment["status"] == "1" {
				status = "Active"
			} else if inPersonCafeAppointment["status"] == "4"  {
				status = "User Cancelled"
			}

			if ordersMap[inPersonCafeAppointment["event_order_id"]]["status"] == "1" {
				eventStatus = "Active"
			} else if ordersMap[inPersonCafeAppointment["event_order_id"]]["status"] == "2" {
				eventStatus = "Started"
			} else if ordersMap[inPersonCafeAppointment["event_order_id"]]["status"] == "3" {
				eventStatus = "Completed"
			} else {
				eventStatus = "Cancelled"
			}


			data = append(data, []string{
				ordersMap[inPersonCafeAppointment["event_order_id"]]["title"],
				counsellorsMap[ordersMap[inPersonCafeAppointment["event_order_id"]]["counsellor_id"]]["first_name"] + " " + counsellorsMap[ordersMap[inPersonCafeAppointment["event_order_id"]]["counsellor_id"]]["last_name"],
				clientsMap[inPersonCafeAppointment["user_id"]]["first_name"] + " " + clientsMap[inPersonCafeAppointment["user_id"]]["last_name"],
				attended,
				ordersMap[inPersonCafeAppointment["event_order_id"]]["company_name"],
				ordersMap[inPersonCafeAppointment["event_order_id"]]["company_location"],
				inPersonCafeAppointment["question1"],
				inPersonCafeAppointment["question2"],
				inPersonCafeAppointment["question3"],
				inPersonCafeAppointment["question4"],
				inPersonCafeAppointment["question5"],
				inPersonCafeAppointment["cancellation_reason"],
				status,
				eventStatus,
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(inPersonCafeAppointment["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat),
			})
		}

	case "25": // inperson cafe request
		heading = []string{"Client Name", "Cafe Name", "Total Seat", "Company Name", "Company Location", "Counsellor Name", "Status"}
		inPersonCafeRequests, status, ok := DB.SelectProcess("select * from " + CONSTANT.EventInPersonRequestTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeRequests, "client_id")
		orderIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeRequests, "order_id")

		// get client details
		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		orders, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderCounsellorEventInPersonTable + " where order_id in ('" + strings.Join(orderIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(orders, "counsellor_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		ordersMap := UTIL.ConvertMapToKeyMap(orders, "order_id")

		for _, inPersonCafeRequest := range inPersonCafeRequests {

			var status1 string

			if inPersonCafeRequest["status"] == "1" {
				status1 = "InProgress"
			} else if inPersonCafeRequest["status"] == "2" {
				status1 = "Completed"
			}

			data = append(data, []string{
				clientsMap[inPersonCafeRequest["client_id"]]["first_name"] + " " + clientsMap[inPersonCafeRequest["client_id"]]["last_name"],
				ordersMap[inPersonCafeRequest["order_id"]]["title"],
				ordersMap[inPersonCafeRequest["order_id"]]["total_seat"],
				ordersMap[inPersonCafeRequest["order_id"]]["company_name"],
				ordersMap[inPersonCafeRequest["order_id"]]["company_location"],
				counsellorsMap[ordersMap[inPersonCafeRequest["order_id"]]["counsellor_id"]]["first_name"] + " " + counsellorsMap[ordersMap[inPersonCafeRequest["order_id"]]["counsellor_id"]]["last_name"],
				status1,
			})
		}

	}

	if strings.EqualFold(r.FormValue("type"), "json") {
		response["headings"] = heading
		response["data"] = data
		UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
		return
	} else {
		writer := csv.NewWriter(w)

		writer.Write(heading)

		for _, d := range data {
			writer.Write(d)
		}

		writer.Flush()
		return
	}
}

// utils for reports
func getAppointmentStatusInText(status string) string {
	switch status {
	case CONSTANT.AppointmentToBeDuplicate:
		return "Duplicate Appointment"
	case CONSTANT.AppointmentToBeStarted:
		return "Both no-show"
	case CONSTANT.AppointmentStarted:
		return "Started"
	case CONSTANT.AppointmentCompleted:
		return "Completed"
	case CONSTANT.AppointmentUserCancelled:
		return "Client cancelled"
	case CONSTANT.AppointmentUserCancelledWithin4Hour:
		return "Client cancelled(to be paid)"
	case CONSTANT.AppointmentCounsellorCancelledWithin4Hour:
		return "Counsellor cancelled(to be charge)"
	case CONSTANT.AppointmentInTheReview:
		return "To be reviewed"
	case CONSTANT.AppointmentCounsellorCancelled:
		return "Counsellor cancelled"
	case CONSTANT.AppointmentAdminCancelled:
		return "Admin cancelled"
	case CONSTANT.AppointmentNoShowClient:
		return "Client no-show"
	case CONSTANT.AppointmentNoShowCounsellor:
		return "Counsellor no-show"
	case CONSTANT.AppointmentNoShowBoth:
		return "Both no-show"
	}
	return ""
}

func getInt(input string) float64 {
	out, _ := strconv.ParseFloat(input, 64)
	return out
}

func GetAppReport(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get contents
	wheresAppointment := []string{}
	wheresClient := []string{}
	wheresSession := []string{}

	appointmentTotal, clientTotal, appointmentsInPersonTotal, emeCaseVirtualTotal, emeCaseInPersonTotal, contentsTotal, moodsTotal, assessmentsTotal, totalRatingTotal, avgRatingTotal, appointmentsCancellationTotal, appointmentsNoShowTotal, familyMemberTotal, companyName, inpersonCafe, inpersonCafeAttend := "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""

	queryArgs := []interface{}{}
	for key, val := range r.URL.Query() {
		switch key {
		case "start_by":
			if len(val[0]) > 0 {
				wheresAppointment = append(wheresAppointment, " `date` >= ? ")
				wheresClient = append(wheresClient, " `created_at` >= ? ")
				wheresSession = append(wheresSession, " `session_date` >= ? ")
				startBy, _ := time.Parse("2006-01-02", val[0])
				queryArgs = append(queryArgs, startBy)
			}
		case "end_by":
			if len(val[0]) > 0 {
				wheresAppointment = append(wheresAppointment, " `date` <= ? ")
				wheresClient = append(wheresClient, " `created_at` <= ? ")
				wheresSession = append(wheresSession, " `session_date` <= ? ")
				endBy, _ := time.Parse("2006-01-02", val[0])
				queryArgs = append(queryArgs, endBy)
			}
		case "company":
			if len(val[0]) > 0 {
				companyName = val[0]
			}

		}
	}

	whereAppointment := ""
	if len(wheresAppointment) > 0 {
		whereAppointment = " where " + strings.Join(wheresAppointment, " and ")
	}

	whereClient := ""
	if len(wheresAppointment) > 0 {
		whereClient = " where " + strings.Join(wheresClient, " and ")
	}

	whereSession := ""
	if len(wheresAppointment) > 0 {
		whereSession = " where " + strings.Join(wheresSession, " and ")
	}

	if companyName != "" {

		clientIDsWithArray, status, ok := DB.SelectProcess("select client_id from " + CONSTANT.ClientsTable + " where email like '%" + companyName + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(clientIDsWithArray, "client_id")

		companyNames, status, ok := DB.SelectProcess("select partner_name from " + CONSTANT.CorporatePartnersTable + " where domain = '" + companyName + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if whereAppointment != "" && whereClient != "" && whereSession != "" {

			appointments, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '3' and client_id in ('"+strings.Join(clientIDs, "','")+"') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentTotal = appointments[0]["total"]

			clients, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ClientsTable+whereClient+" and status = 1 and email like '%"+companyName+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientTotal = clients[0]["total"]

			familyMember, status, ok := DB.SelectProcess("select asscoiate_id from "+CONSTANT.ClientsTable+whereClient+" and status = '1' and asscoiate_id != ''", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			asscoiateIDs := UTIL.ExtractValuesFromArrayMap(familyMember, "asscoiate_id")

			familyMembers, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ClientsTable+whereClient+" and status = 1 and client_id in ('"+strings.Join(asscoiateIDs, "','")+"') and email like '%"+companyName+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			familyMemberTotal = familyMembers[0]["total"]

			totalInpersonCafes, status, ok := DB.SelectProcess("select count(*) as total, order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+whereClient+" and status = '3' and company_name = '"+companyNames[0]["partner_name"]+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			inpersonCafe = totalInpersonCafes[0]["total"]

			inpersonCafes, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+whereClient+" and status = '3' and company_name = '"+companyNames[0]["partner_name"]+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			eventIDs := UTIL.ExtractValuesFromArrayMap(inpersonCafes, "order_id")

			inpersonCafeAttends, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.OrderEventInPersonTable+whereClient+" and event_order_id in ('"+strings.Join(eventIDs, "','")+"') and attended = '1'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			inpersonCafeAttend = inpersonCafeAttends[0]["total"]

			appointmentsInPerson, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.InPersonAppointmentsTable+whereAppointment+" and  type = '4' and status = '3' and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsInPersonTotal = appointmentsInPerson[0]["total"]

			appointmentsCancellation, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '4' and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsCancellationTotal = appointmentsCancellation[0]["total"]

			appointmentsNoShow, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and client_id in ('"+strings.Join(clientIDs, "','")+"') and (status = '1' or status = '3')  and (client_started_at is null and client_ended_at is null) ", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsNoShowTotal = appointmentsNoShow[0]["total"]

			emeCaseVirtual, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.CounsellorRecordsTable+whereSession+" and session_mode = 'Virtual' and mental_health >= '8' and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			emeCaseVirtualTotal = emeCaseVirtual[0]["total"]

			emeCaseInPerson, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.CounsellorRecordsTable+whereSession+" and session_mode = 'In Person' and mental_health >= '8' and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			emeCaseInPersonTotal = emeCaseInPerson[0]["total"]

			contents, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ContentsTable+whereClient+" and status = '1'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			contentsTotal = contents[0]["total"]

			moods, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.MoodResultsTable+whereAppointment+" and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			moodsTotal = moods[0]["total"]

			assessments, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AssessmentResultsTable+whereClient+" and user_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			assessmentsTotal = assessments[0]["total"]

			totalRating, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and rating is not null and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			totalRatingTotal = totalRating[0]["total"]

			totalAvg, status, ok := DB.SelectProcess("select avg(rating) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and rating is not null and client_id in ('"+strings.Join(clientIDs, "','")+"')", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			avgRating, _ := strconv.ParseFloat(totalAvg[0]["total"], 64)

			avgRatingTotal = strconv.FormatFloat(avgRating, 'f', 2, 64)

		} else {

			appointments, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '3' and client_id in ('" + strings.Join(clientIDs, "','") + "') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentTotal = appointments[0]["total"]

			clients, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = '1' and email like '%" + companyName + "'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientTotal = clients[0]["total"]

			familyMember, status, ok := DB.SelectProcess("select asscoiate_id from " + CONSTANT.ClientsTable + " where status = '1' and asscoiate_id != ''")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			asscoiateIDs := UTIL.ExtractValuesFromArrayMap(familyMember, "asscoiate_id")

			familyMembers, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = '1' and client_id in ('" + strings.Join(asscoiateIDs, "','") + "') and email like '%" + companyName + "'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			familyMemberTotal = familyMembers[0]["total"]

			totalInpersonCafes, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.OrderCounsellorEventInPersonTable + " where status = '3' and company_name = '" + companyNames[0]["partner_name"] + "'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			inpersonCafe = totalInpersonCafes[0]["total"]

			inpersonCafes, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+" where status = '3' and company_name = '"+companyNames[0]["partner_name"]+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			eventIDs := UTIL.ExtractValuesFromArrayMap(inpersonCafes, "order_id")

			inpersonCafeAttends, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.OrderEventInPersonTable + " where event_order_id in ('" + strings.Join(eventIDs, "','") + "') and attended = '1'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			inpersonCafeAttend = inpersonCafeAttends[0]["total"]

			appointmentsInPerson, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.InPersonAppointmentsTable + " where  type = '4' and status = '3' and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsInPersonTotal = appointmentsInPerson[0]["total"]

			appointmentsCancellation, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '4' and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsCancellationTotal = appointmentsCancellation[0]["total"]

			appointmentsNoShow, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and (status = '1' or status = '3')  and (client_started_at is null and client_ended_at is null) and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			appointmentsNoShowTotal = appointmentsNoShow[0]["total"]

			emeCaseVirtual, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'Virtual' and mental_health >= '8' and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			emeCaseVirtualTotal = emeCaseVirtual[0]["total"]

			emeCaseInPerson, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'In Person' and mental_health >= '8' and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			emeCaseInPersonTotal = emeCaseInPerson[0]["total"]

			contents, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ContentsTable + " where status = '1'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			contentsTotal = contents[0]["total"]

			moods, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.MoodResultsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			moodsTotal = moods[0]["total"]

			assessments, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AssessmentResultsTable + " where user_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			assessmentsTotal = assessments[0]["total"]

			totalRating, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			totalRatingTotal = totalRating[0]["total"]

			totalAvg, status, ok := DB.SelectProcess("select avg(rating) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null and client_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			avgRating, _ := strconv.ParseFloat(totalAvg[0]["total"], 64)

			avgRatingTotal = strconv.FormatFloat(avgRating, 'f', 2, 64)
		}
	} else if whereAppointment != "" && whereClient != "" && whereSession != "" {

		appointments, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentTotal = appointments[0]["total"]

		clients, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ClientsTable+whereClient+" and status = '1'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientTotal = clients[0]["total"]

		familyMember, status, ok := DB.SelectProcess("select asscoiate_id from "+CONSTANT.ClientsTable+whereClient+" and status = '1' and asscoiate_id != ''", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		asscoiateIDs := UTIL.ExtractValuesFromArrayMap(familyMember, "asscoiate_id")

		familyMembers, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ClientsTable+whereClient+" and status = '1' and client_id in ('"+strings.Join(asscoiateIDs, "','")+"') and email like '%"+companyName+"'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		familyMemberTotal = familyMembers[0]["total"]

		totalInpersonCafes, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.OrderCounsellorEventInPersonTable+whereClient+" and status = '3'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		inpersonCafe = totalInpersonCafes[0]["total"]

		inpersonCafes, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+whereClient+" and status = '3'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		eventIDs := UTIL.ExtractValuesFromArrayMap(inpersonCafes, "order_id")

		inpersonCafeAttends, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.OrderEventInPersonTable+whereClient+" and event_order_id in ('"+strings.Join(eventIDs, "','")+"') and attended = '1'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		inpersonCafeAttend = inpersonCafeAttends[0]["total"]

		appointmentsInPerson, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.InPersonAppointmentsTable+whereAppointment+" and  type = '4' and status = '3'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsInPersonTotal = appointmentsInPerson[0]["total"]

		appointmentsCancellation, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '4'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsCancellationTotal = appointmentsCancellation[0]["total"]

		appointmentsNoShow, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and (status = '1' or status = '3')  and (client_started_at is null and client_ended_at is null) ", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsNoShowTotal = appointmentsNoShow[0]["total"]

		emeCaseVirtual, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.CounsellorRecordsTable+whereSession+" and session_mode = 'Virtual' and mental_health >= '8'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		emeCaseVirtualTotal = emeCaseVirtual[0]["total"]

		emeCaseInPerson, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.CounsellorRecordsTable+whereSession+" and session_mode = 'In Person' and mental_health >= '8'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		emeCaseInPersonTotal = emeCaseInPerson[0]["total"]

		contents, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.ContentsTable+whereClient+" and status = '1'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		contentsTotal = contents[0]["total"]

		moods, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.MoodResultsTable+whereAppointment, queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		moodsTotal = moods[0]["total"]

		assessments, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AssessmentResultsTable+whereClient, queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		assessmentsTotal = assessments[0]["total"]

		totalRating, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and rating is not null", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalRatingTotal = totalRating[0]["total"]

		totalAvg, status, ok := DB.SelectProcess("select avg(rating) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and rating is not null", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		avgRating, _ := strconv.ParseFloat(totalAvg[0]["total"], 64)

		avgRatingTotal = strconv.FormatFloat(avgRating, 'f', 2, 64)

	} else {

		appointments, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentTotal = appointments[0]["total"]

		clients, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = 1")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientTotal = clients[0]["total"]

		familyMember, status, ok := DB.SelectProcess("select asscoiate_id from " + CONSTANT.ClientsTable + " where status = '1' and asscoiate_id != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		asscoiateIDs := UTIL.ExtractValuesFromArrayMap(familyMember, "asscoiate_id")

		familyMembers, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = '1' and client_id in ('" + strings.Join(asscoiateIDs, "','") + "') and email like '%" + companyName + "'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		familyMemberTotal = familyMembers[0]["total"]

		totalInpersonCafes, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.OrderCounsellorEventInPersonTable+" where status = '3' ", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		inpersonCafe = totalInpersonCafes[0]["total"]

		inpersonCafes, status, ok := DB.SelectProcess("select order_id from "+CONSTANT.OrderCounsellorEventInPersonTable+" where status = '3' ", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		eventIDs := UTIL.ExtractValuesFromArrayMap(inpersonCafes, "order_id")

		inpersonCafeAttends, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.OrderEventInPersonTable+" where event_order_id in ('"+strings.Join(eventIDs, "','")+"') and attended = '1'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		inpersonCafeAttend = inpersonCafeAttends[0]["total"]

		appointmentsInPerson, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.InPersonAppointmentsTable + " where  type = '4' and status = '3'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsInPersonTotal = appointmentsInPerson[0]["total"]

		appointmentsCancellation, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '4'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsCancellationTotal = appointmentsCancellation[0]["total"]

		appointmentsNoShow, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and (status = '1' or status = '3')  and (client_started_at is null and client_ended_at is null) ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		appointmentsNoShowTotal = appointmentsNoShow[0]["total"]

		emeCaseVirtual, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'Virtual' and mental_health >= '8'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		emeCaseVirtualTotal = emeCaseVirtual[0]["total"]

		emeCaseInPerson, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'In Person' and mental_health >= '8'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		emeCaseInPersonTotal = emeCaseInPerson[0]["total"]

		contents, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ContentsTable + " where status = '1'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		contentsTotal = contents[0]["total"]

		moods, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.MoodResultsTable)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		moodsTotal = moods[0]["total"]

		assessments, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AssessmentResultsTable)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		assessmentsTotal = assessments[0]["total"]

		totalRating, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		totalRatingTotal = totalRating[0]["total"]

		totalAvg, status, ok := DB.SelectProcess("select avg(rating) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		avgRating, _ := strconv.ParseFloat(totalAvg[0]["total"], 64)

		avgRatingTotal = strconv.FormatFloat(avgRating, 'f', 2, 64)
	}

	response["total_appointment"] = appointmentTotal
	response["total_client"] = clientTotal
	response["total_family_member"] = familyMemberTotal
	response["total_inperson_cafe"] = inpersonCafe
	response["total_inperson_cafe_attend"] = inpersonCafeAttend
	response["total_inperson_appointment"] = appointmentsInPersonTotal
	response["total_appointment_cancel"] = appointmentsCancellationTotal
	response["total_appointment_noshow"] = appointmentsNoShowTotal
	response["total_emecase_virtual"] = emeCaseVirtualTotal
	response["total_emecase_inperson"] = emeCaseInPersonTotal
	response["total_content"] = contentsTotal
	response["total_mood"] = moodsTotal
	response["total_assessment"] = assessmentsTotal
	response["total_rating"] = totalRatingTotal
	response["total_avg"] = avgRatingTotal

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)

}

func GetAppSummaryReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {

		appointmentTotal, clientTotal, appointmentsInPersonTotal, emeCaseVirtualTotal, emeCaseInPersonTotal, contentsTotal, moodsTotal, assessmentsTotal, totalRatingTotal, avgRatingTotal, appointmentsCancellationTotal, appointmentsNoShowTotal := "", "", "", "", "", "", "", "", "", "", "", ""

		appointments, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")

		appointmentTotal = appointments[0]["total"]

		clients, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = 1")

		clientTotal = clients[0]["total"]

		appointmentsInPerson, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.InPersonAppointmentsTable + " where  type = '4' and status = '3'")

		appointmentsInPersonTotal = appointmentsInPerson[0]["total"]

		appointmentsCancellation, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '4'")

		appointmentsCancellationTotal = appointmentsCancellation[0]["total"]

		appointmentsNoShow, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and (status = '1' or status = '3')  and (client_started_at is null and client_ended_at is null) ")

		appointmentsNoShowTotal = appointmentsNoShow[0]["total"]

		emeCaseVirtual, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'Virtual' and mental_health >= '8'")

		emeCaseVirtualTotal = emeCaseVirtual[0]["total"]

		emeCaseInPerson, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.CounsellorRecordsTable + " where session_mode = 'In Person' and mental_health >= '8'")

		emeCaseInPersonTotal = emeCaseInPerson[0]["total"]

		contents, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.ContentsTable + " where status = '1'")

		contentsTotal = contents[0]["total"]

		moods, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.MoodResultsTable)

		moodsTotal = moods[0]["total"]

		assessments, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.AssessmentResultsTable)

		assessmentsTotal = assessments[0]["total"]

		totalRating, _, _ := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null")

		totalRatingTotal = totalRating[0]["total"]

		totalAvg, _, _ := DB.SelectProcess("select avg(rating) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and rating is not null")

		avgRating, _ := strconv.ParseFloat(totalAvg[0]["total"], 64)

		avgRatingTotal = strconv.FormatFloat(avgRating, 'f', 2, 64)

		event := Model.AppSummaryReport{
			AppointmentTotal:              appointmentTotal,
			ClientTotal:                   clientTotal,
			AppointmentsInPersonTotal:     appointmentsInPersonTotal,
			EmeCaseVirtualTotal:           emeCaseVirtualTotal,
			EmeCaseInPersonTotal:          emeCaseInPersonTotal,
			ContentsTotal:                 contentsTotal,
			MoodsTotal:                    moodsTotal,
			AssessmentsTotal:              assessmentsTotal,
			TotalRatingTotal:              totalRatingTotal,
			AvgRatingTotal:                avgRatingTotal,
			AppointmentsCancellationTotal: appointmentsCancellationTotal,
			AppointmentsNoShowTotal:       appointmentsNoShowTotal,
		}
		data, _ := json.Marshal(event)
		fmt.Println(event)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}
}
