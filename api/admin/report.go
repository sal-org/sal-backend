package admin

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

		heading = []string{"Client Name", "Gender", "Age", "Company Name", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date & Time", "Therapist Start", "Therapist End", "Client Start", "Client End", "Mod. At", "Status"}
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

			var startTime, endTime, modAt, clientStartTime, clientEndTime, partnerName, status, location string

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
				} else if appointment["ended_at"] == "" || appointment["client_ended_at"] == "" {
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

		var partnerName, location string

		heading = []string{"Client Name", "Company Name", "Age", "Gender", "Location", "CreatedAt"}

		clients, status, ok := DB.SelectProcess("select first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, created_at from " + CONSTANT.ClientsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
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
				partnerName = "None"
			}

			if client["location"] == "40.04" {
				location = ""
			} else {
				location = client["location"]
			}

			data = append(data, []string{client["first_name"] + " " + client["last_name"], partnerName, client["age"], client["gender"], location, UTIL.ConvertTimezone(UTIL.ConvertToTime(client["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)})
		}

	case "12": // client mood report

		var partnerName, location string

		heading = []string{"Client Name", "Mood", "Age", "Gender", "Company Name", "Location", "Department"}

		moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodResultsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(moods, "client_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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
			} else {
				partnerName = "None"
			}

			if clientsMap[mood["client_id"]]["location"] == "40.04" {
				location = ""
			} else {
				location = clientsMap[mood["client_id"]]["location"]
			}

			moodTitle := DB.QueryRowSQL("select title from "+CONSTANT.MoodsTable+" where id = ?", mood["mood_id"])

			data = append(data, []string{mood["name"], moodTitle, mood["age"], mood["gender"], partnerName, location, clientsMap[mood["client_id"]]["department"]})
		}
	case "13": // client asessment report

		var partnerName, location, department string

		heading = []string{"Client Name", "Assessment", "Age", "Gender", "Company Name", "Location", "Department", "CreatedAt"}

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
				} else {
					partnerName = "None"
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

			data = append(data, []string{assessment["name"], assessmentTitle, assessment["age"], assessment["gender"], partnerName, location, department, UTIL.ConvertTimezone(UTIL.ConvertToTime(assessment["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)})
		}
	case "14": // counsellor record report
		var noShow string

		heading = []string{"Counsellor Name", "Client Name", "Age", "Gender", "Location", "Department", "No Show", "Session Mode", "Session No", "Session Date", "In-Time", "Out-Time", "Mental Health Scale", "Therapeutic Goal"}

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
			data = append(data, []string{counsellorsMap[counsellorRecord["counsellor_id"]]["first_name"] + " " + counsellorsMap[counsellorRecord["counsellor_id"]]["last_name"], counsellorRecord["client_first_name"] + " " + counsellorRecord["client_last_name"], counsellorRecord["client_age"], counsellorRecord["client_gender"], counsellorRecord["client_location"], counsellorRecord["client_department"], noShow, counsellorRecord["session_mode"], counsellorRecord["session_no"], counsellorRecord["session_date"], counsellorRecord["in_time"], counsellorRecord["out_time"], counsellorRecord["mental_health"], counsellorRecord["therapeutic_goal"]})
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
