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

	var response = make(map[string]interface{})

	startBy, _ := time.Parse("2006-01-02", r.FormValue("start_by"))
	endBy, _ := time.Parse("2006-01-02", r.FormValue("end_by"))

	heading := []string{}
	fileName := "report.csv"
	data := [][]string{}

	switch r.FormValue("id") {
	case "1": // appointment report

		fileName = "appointment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"

		isdaysLessThan31 := false
		if endBy.Sub(startBy).Hours() <= 24*366 {
			isdaysLessThan31 = true
		}

		if isdaysLessThan31 {
			heading = []string{"Client Name", "Gender", "Age", "Email", "Mobile No.", "Company Name", "IsFamilyMember", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date", "Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Client Start", "Client End", "Duration", "Mod. At", "Created At", "Status", "Days Ago", "Is Record Filled", "Mental Health Scale"}
		} else {
			heading = []string{"Client Name", "Gender", "Age", "Email", "Mobile No.", "Company Name", "IsFamilyMember", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date", "Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Client Start", "Client End", "Duration", "Mod. At", "Created At", "Status", "Is Record Filled", "Mental Health Scale"}
		}

		appointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")
		appointmentIDs := UTIL.ExtractValuesFromArrayMap(appointments, "appointment_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		counsellorRecord, status, ok := DB.SelectProcess("select mental_health, appointment_id, noshow from " + CONSTANT.CounsellorRecordsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "') ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// // calculate the days ago
		// appointmntsPast, status, ok := DB.SelectProcess("select date, appointment_id from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.AddDate(0, -1, 0).String() + "' and `date` < '" + endBy.String() + "' and status in ('2', '3') order by date desc")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorRecordMap := UTIL.ConvertMapToKeyMap(counsellorRecord, "appointment_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointment := range appointments {

			var days, duration, startTime, endTime, modAt, createdAt, clientStartTime, clientEndTime, partnerName, isFamilyMember, status, location, isRecordFilled, mentalHealth string

			// get client details
			// counsellorRecord, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where appointment_id = ? ", appointment["appointment_id"])
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			// if len(counsellorRecord) > 0 {
			// 	isRecordFilled = "Yes"
			// 	mentalHealth = counsellorRecord[0]["mental_health"]
			// } else {
			// 	isRecordFilled = "No"
			// 	mentalHealth = ""
			// }

			if appointment["started_at"] == "" {
				startTime = ""
			} else {
				startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["client_started_at"] == "" {
				clientStartTime = ""
			} else {
				clientStartTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["client_ended_at"] == "" {
				clientEndTime = ""
			} else {
				clientEndTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["ended_at"] == "" {
				endTime = ""
			} else {
				endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["modified_at"] == "" {
				modAt = ""
			} else {
				modAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["created_at"] == "" {
				createdAt = ""
			} else {
				createdAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			switch appointment["status"] {
			case "3":
				if appointment["started_at"] == "" && appointment["ended_at"] == "" {
					status = getAppointmentStatusInText("8")
					if appointment["client_started_at"] != "" && appointment["client_ended_at"] != "" {
						duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes())) + " minutes"
					} else {
						duration = "0 minutes"
					}
				} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
					status = getAppointmentStatusInText("7")
					duration = "0 minutes"
				} else {
					// calculate the duration

					if len(appointment["client_started_at"]) != 0 && len(appointment["client_ended_at"]) == 0 {
						duration = "5 minutes"
					} else if len(appointment["client_started_at"]) == 0 && len(appointment["client_ended_at"]) != 0 {
						duration = "5 minutes"
					} else {
						duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes())) + " minutes"
					}

					if isdaysLessThan31 {
						startDate := time.Date(startBy.Year(), startBy.Month()-1, 1, 0, 0, 0, 0, time.UTC)

						// calculate the days ago
						appointmntsPast := DB.QueryRowSQL("select date from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startDate.Format("2006-01-02") + "' and `date` < '" + appointment["date"] + "' and client_id = '" + appointment["client_id"] + "' and counsellor_id = '" + appointment["counsellor_id"] + "' and status in ('2', '3') order by date desc")
						if !ok {
							UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
							return
						}

						if appointmntsPast != "" {

							startDate, _ := time.Parse("2006-01-02", appointmntsPast)
							endDate, _ := time.Parse("2006-01-02", appointment["date"])

							start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
							end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

							// Calculate the difference in days
							day := int(end.Sub(start).Hours() / 24)

							days = strconv.Itoa(day) + " days ago"

						}
					}

					if appointment["client_started_at"] != "" && appointment["client_ended_at"] != "" && appointment["started_at"] != "" && appointment["ended_at"] != "" {
						if UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes() > 10 {
							if UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes() < 10 {
								status = getAppointmentStatusInText("19")
							} else {
								status = getAppointmentStatusInText("3")
							}
						} else {
							status = getAppointmentStatusInText("19")
						}
					} else {
						status = getAppointmentStatusInText("14")
					}

				}
			case "4":
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("12")
				} else {
					status = getAppointmentStatusInText("4")
				}
			case "5":
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("13")
				} else {
					status = getAppointmentStatusInText("5")
				}
			default:
				status = getAppointmentStatusInText(appointment["status"])
				if appointment["status"] == "2" {
					// duration
					if appointment["client_started_at"] != "" || appointment["client_started_at"] == "" {
						duration = "5 minutes"
					} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
						duration = "5 minutes"
					}

				}
			}

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				// get client details
				// clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where client_id = '" + appointment["client_id"] + "'")
				// if !ok {
				// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				// 	return
				// }

				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointment["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNamesMap[domainName[1]]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			if len(counsellorRecordMap[appointment["appointment_id"]]) > 0 {
				isRecordFilled = "Yes"
				mentalHealth = counsellorRecordMap[appointment["appointment_id"]]["mental_health"]
				// if counsellorRecordMap[appointment["appointment_id"]]["noshow"] == "1" {
				// 	// status = getAppointmentStatusInText("7")
				// 	duration = "0 minutes"
				// }
			} else {
				isRecordFilled = "No"
				mentalHealth = ""
			}

			if clientsMap[appointment["client_id"]]["location"] == "40.04" {
				location = ""
			} else {
				location = clientsMap[appointment["client_id"]]["location"]
			}

			if isdaysLessThan31 {
				data = append(data, []string{
					clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
					clientsMap[appointment["client_id"]]["gender"],
					clientsMap[appointment["client_id"]]["age"],
					clientsMap[appointment["client_id"]]["email"],
					clientsMap[appointment["client_id"]]["phone"],
					partnerName,
					isFamilyMember,
					location,
					clientsMap[appointment["client_id"]]["department"],
					counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
					counsellorsMap[appointment["counsellor_id"]]["type"],
					UTIL.BuildOnlyDateInDDMMYYYY(appointment["date"]),
					UTIL.GetTimeFromTimeSlot(appointment["time"]),
					appointment["times_rescheduled"],
					startTime,
					endTime,
					clientStartTime,
					clientEndTime,
					duration,
					modAt,
					createdAt,
					status,
					days,
					isRecordFilled,
					mentalHealth,
				})
			} else {
				data = append(data, []string{
					clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
					clientsMap[appointment["client_id"]]["gender"],
					clientsMap[appointment["client_id"]]["age"],
					clientsMap[appointment["client_id"]]["email"],
					clientsMap[appointment["client_id"]]["phone"],
					partnerName,
					isFamilyMember,
					location,
					clientsMap[appointment["client_id"]]["department"],
					counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
					counsellorsMap[appointment["counsellor_id"]]["type"],
					UTIL.BuildOnlyDateInDDMMYYYY(appointment["date"]),
					UTIL.GetTimeFromTimeSlot(appointment["time"]),
					appointment["times_rescheduled"],
					startTime,
					endTime,
					clientStartTime,
					clientEndTime,
					duration,
					modAt,
					createdAt,
					status,
					isRecordFilled,
					mentalHealth,
				})
			}

		}
	case "2": // sales report
		fileName = "sales_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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

		fileName = "booking_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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
				UTIL.ConvertTimezone(UTIL.ConvertToTime(booking["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				booking["slots_bought"],
				booking["slots_remaining"],
				booking["status"],
			})
		}
	case "4": // sal cafe report

		fileName = "sal_cafe_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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
				UTIL.ConvertTimezone(UTIL.ConvertToTime(booking["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				booking["paid_amount"],
			})
		}
	case "5": // finance report

		fileName = "finance_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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
				UTIL.ConvertTimezone(UTIL.ConvertToTime(invoice["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				invoice["paid_amount"],
				invoice["cgst"],
				invoice["sgst"],
				invoice["tax"],
			})
		}
	case "6": // payout report

		fileName = "payout_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Counsellor ID", "Counsellor Name", "Counsellor Type", "Client Name", "Description", "Date & Time of Session", "Amount to be paid"}
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
				counsellorsMap[payment["counsellor_id"]]["first_name"] + " " + counsellorsMap[payment["counsellor_id"]]["last_name"],
				counsellorsMap[payment["counsellor_id"]]["type"],
				payment["heading"],
				payment["description"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(payment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				payment["amount"],
			})
		}
	case "7": // promo code report

		fileName = "promo_code_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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

		fileName = "push_notification_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Notification Type", "Date", "Times"}
		notifications, status, ok := DB.SelectProcess("select date(created_at) as date, CASE WHEN notification_type = 1 THEN 'Promo' WHEN notification_type = 2 THEN 'Content' WHEN notification_type = 3 THEN 'Event' WHEN notification_type = 4 THEN 'Other' ELSE 'Nothing' END as notification_type, count(*) as ctn from " + CONSTANT.NotificationsBulkTable + " where created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "' group by date(created_at), notification_type")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, notification := range notifications {
			data = append(data, []string{
				notification["notification_type"],
				notification["date"],
				notification["ctn"],
			})
		}
	case "10": // onboarding report

		fileName = "onboarding_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"First Name", "Last Name", "Gender", "Email", "Phone", "Type", "Created At"}
		counsellors, status, ok := DB.SelectProcess("(select first_name, last_name, gender, email, phone, 'Counsellor' as `type`, created_at from " + CONSTANT.CounsellorsTable + " where status = " + CONSTANT.CounsellorActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Listener' as `type`, created_at from " + CONSTANT.ListenersTable + " where status = " + CONSTANT.ListenerActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "') union (select first_name, last_name, gender, email, phone, 'Therapist' as `type`, created_at from " + CONSTANT.TherapistsTable + " where status = " + CONSTANT.TherapistActive + " and created_at > '" + startBy.UTC().String() + "' and created_at < '" + endBy.UTC().String() + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		for _, counsellor := range counsellors {
			data = append(data, []string{counsellor["first_name"], counsellor["last_name"], counsellor["gender"], counsellor["email"], counsellor["phone"], counsellor["type"], UTIL.ConvertTimezone(UTIL.ConvertToTime(counsellor["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)})
		}
	case "11": // client onboarding report

		var partnerName, location string

		fileName = "client_onboarding_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Company Name", "Mobile No.", "Email", "Age", "Gender", "Location", "Last Active Time", "Last Login Time", "CreatedAt"}

		clients, status, ok := DB.SelectProcess("select asscoiate_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, last_active_time, last_login_time, created_at from " + CONSTANT.ClientsTable + " where asscoiate_id IS NULL and `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")

		for _, client := range clients {

			domainName := strings.Split(client["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
			} else {
				partnerName = "None"
			}

			if client["location"] == "40.04" {
				location = ""
			} else {
				location = client["location"]
			}

			data = append(data, []string{client["first_name"] + " " + client["last_name"], partnerName, client["phone"], client["email"], client["age"], client["gender"], location, UTIL.ConvertTimezone(UTIL.ConvertToTime(client["last_active_time"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat), UTIL.ConvertTimezone(UTIL.ConvertToTime(client["last_login_time"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat), UTIL.ConvertTimezone(UTIL.ConvertToTime(client["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)})
		}

	case "12": // client mood report

		var partnerName, isFamilyMember, location string

		fileName = "client_mood_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Mood", "Age", "Email ID", "Mobile No.", "Gender", "Company Name", "IsFamilyMember", "Location", "Department", "Notes"}

		moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodResultsTable + " where `date` >= '" + startBy.Format("2006-01-02") + "' and `date` <= '" + endBy.Format("2006-01-02") + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(moods, "client_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, gender, phone, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		moodsTitle, status, ok := DB.SelectProcess("select id, title from " + CONSTANT.MoodsTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")
		moodsTitleMap := UTIL.ConvertMapToKeyMap(moodsTitle, "id")

		for _, mood := range moods {

			domainName := strings.Split(clientsMap[mood["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[mood["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[mood["client_id"]]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[mood["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNamesMap[domainName[1]]["partner_name"]
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

			// appointments, status, ok := DB.SelectProcess("select count(*) as ctn from " + CONSTANT.AppointmentsTable + " where client_id = '" + mood["client_id"] + "' and status = " + CONSTANT.AppointmentCompleted + " and `date` >= '" + startBy.Format("2006-01-02") + "' and `date` <= '" + endBy.Format("2006-01-02") + "' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			//appointments[0]["ctn"]

			// moodTitle := DB.QueryRowSQL("select title from "+CONSTANT.MoodsTable+" where id = ? ", mood["mood_id"])

			moodTitle := moodsTitleMap[mood["mood_id"]]["title"]

			data = append(data, []string{clientsMap[mood["client_id"]]["first_name"] + " " + clientsMap[mood["client_id"]]["last_name"], moodTitle, mood["age"], clientsMap[mood["client_id"]]["email"], clientsMap[mood["client_id"]]["phone"], mood["gender"], partnerName, isFamilyMember, location, clientsMap[mood["client_id"]]["department"], mood["notes"]})
		}
	case "13": // client assessment report

		var partnerName, isFamilyMember, location, department string

		fileName = "client_assessment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Assessment", "Age", "Email ID", "Mobile No.", "Gender", "Company Name", "IsFamilyMember", "Location", "Department", "CreatedAt"}

		assessments, status, ok := DB.SelectProcess("select * from " + CONSTANT.AssessmentResultsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDs := UTIL.ExtractValuesFromArrayMap(assessments, "user_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		assessmentsTitle, status, ok := DB.SelectProcess("select title, assessment_id from " + CONSTANT.AssessmentsTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")
		assessmentsTitleMap := UTIL.ConvertMapToKeyMap(assessmentsTitle, "assessment_id")

		for _, assessment := range assessments {

			// var totalAppointment string

			if len(clientsMap[assessment["user_id"]]) != 0 {
				domainName := strings.Split(clientsMap[assessment["user_id"]]["email"], "@")

				// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

				if len(companyNamesMap[domainName[1]]) > 0 {
					partnerName = companyNamesMap[domainName[1]]["partner_name"]
					isFamilyMember = "No"
				} else {
					partnerN := ""

					if clientsMap[assessment["user_id"]]["asscoiate_id"] != "" {
						// get client details
						// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[assessment["user_id"]]["asscoiate_id"] + "'")
						// if !ok {
						// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						// 	return
						// }

						domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[assessment["user_id"]]["asscoiate_id"]]["email"], "@")

						// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

						partnerN = companyNamesMap[domainName[1]]["partner_name"]
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

				// startDate := UTIL.BuildOnlyDateInYYYYMMDD(assessment["created_at"])

				// appointments, status, ok := DB.SelectProcess("select count(*) as ctn from " + CONSTANT.AppointmentsTable + " where client_id = '" + assessment["user_id"] + "' and status = " + CONSTANT.AppointmentCompleted + " and `date` >= '" + startDate + "' and `date` <= '" + endBy.Format("2006-01-02") + "' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
				// if !ok {
				// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				// 	return
				// }

				// totalAppointment = appointments[0]["ctn"]

			} else {
				partnerName = "None"
				location = ""
				department = ""
				// totalAppointment = "0"

			}

			// assessmentTitle := DB.QueryRowSQL("select title from "+CONSTANT.AssessmentsTable+" where assessment_id = ?", assessment["assessment_id"])

			assessmentTitle := assessmentsTitleMap[assessment["assessment_id"]]["title"]

			data = append(data, []string{assessment["name"], assessmentTitle, assessment["age"], clientsMap[assessment["user_id"]]["email"], clientsMap[assessment["user_id"]]["phone"], assessment["gender"], partnerName, isFamilyMember, location, department, UTIL.ConvertTimezone(UTIL.ConvertToTime(assessment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)})
		}
	case "14": // counsellor record report
		var noShow string

		fileName = "counsellor_record_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Session For", "Session Type", "Company Name", "Counsellor Name", "Client Name", "Mental Health Scale", "Age", "Gender", "Location", "Department", "No Show", "Session Mode", "Session Date", "In-Time", "Out-Time", "Presenting Concerns", "Psychiatric Intervention", "Psychiatric Intervention Reason", "Therapy Notes", "Category", "Sub Category", "Emotional State", "Next Follow Date", "Client Notes", "Client Self-Work", "Client Assessment", "Therapeutic Plan"}

		counsellorRecords, status, ok := DB.SelectProcess("select * from " + CONSTANT.CounsellorRecordsTable + " where `session_date` >= '" + startBy.String() + "' and `session_date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecords, "counsellor_id")
		clientIDs := UTIL.ExtractValuesFromArrayMap(counsellorRecords, "client_id")

		// get counsellor details
		counsellors, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorIDs, "','") + "'))")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name,email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, counsellorRecord := range counsellorRecords {
			if counsellorRecord["noshow"] == "0" {
				noShow = "No"
			} else {
				noShow = "Yes"
			}

			partnerName := "None"

			// get client details
			if counsellorRecord["client_id"] != "" {
				// clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name,email from " + CONSTANT.ClientsTable + " where client_id = '" + counsellorRecord["client_id"] + "'")
				// if !ok {
				// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				// 	return
				// }
				// if len(clients) > 0 {

				domainName := strings.Split(clientsMap[counsellorRecord["client_id"]]["email"], "@")
				// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})
				if len(companyNamesMap[domainName[1]]) > 0 {
					partnerName = companyNamesMap[domainName[1]]["partner_name"]
				} else {
					if clientsMap[counsellorRecord["client_id"]]["asscoiate_id"] != "" {
						// get client details
						// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
						// if !ok {
						// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
						// 	return
						// }

						domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[counsellorRecord["client_id"]]["asscoiate_id"]]["email"], "@")

						// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

						if len(companyNamesMap[domainName[1]]) > 0 {
							partnerName = companyNamesMap[domainName[1]]["partner_name"]
						}
					}
				}

			}
			data = append(data, []string{counsellorRecord["session_for"], counsellorRecord["session_type"], partnerName, counsellorsMap[counsellorRecord["counsellor_id"]]["first_name"] + " " + counsellorsMap[counsellorRecord["counsellor_id"]]["last_name"], counsellorRecord["client_first_name"] + " " + counsellorRecord["client_last_name"], counsellorRecord["mental_health"], counsellorRecord["client_age"], counsellorRecord["client_gender"], counsellorRecord["client_location"], counsellorRecord["client_department"], noShow, counsellorRecord["session_mode"], counsellorRecord["session_date"], counsellorRecord["in_time"], counsellorRecord["out_time"], counsellorRecord["presenting_concerns"], counsellorRecord["psychiatric_intervention"], counsellorRecord["psychiatric_intervention_reason"], counsellorRecord["therapy_notes"], counsellorRecord["therapeutic_goal"], counsellorRecord["sub_category"], counsellorRecord["emotional_state"], counsellorRecord["next_follow_date"], counsellorRecord["client_notes"], counsellorRecord["client_documents"], counsellorRecord["assessment_tool"], counsellorRecord["therapy_plan"]})
		}

	case "15": // client virtual appointment rating

		fileName = "client_virtual_appointment_rating_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Counsellor Name", "Client Name", "Company Name", "IsFamilyMember", "Gender", "Date", "Rate", "Rate Type", "Comment"}

		appointments, status, ok := DB.SelectProcess("select counsellor_id, client_id, `date`, rating, rating_types, rating_comment  from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, gender from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointment := range appointments {

			partnerName := "None"
			isFamilyMember := "No"

			// clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name,email from " + CONSTANT.ClientsTable + " where client_id = '" + appointment["client_id"] + "'")
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointment["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerName = companyNamesMap[domainName[1]]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerName = "None"
					isFamilyMember = "No"
				}
			}

			data = append(data, []string{
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				partnerName,
				isFamilyMember,
				clientsMap[appointment["client_id"]]["gender"],
				appointment["date"],
				appointment["rating"],
				appointment["rating_types"],
				appointment["rating_comment"],
			})
		}
	case "16": // content

		fileName = "content_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Title", "Type", "Views", "Category", "Mood", "Is Training Content", "Created By", "Uploaded At", "Status"}

		contents, status, ok := DB.SelectProcess("select *  from " + CONSTANT.ContentsTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get mood details
		moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodsTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get category details
		categorys, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentCategoriesTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// // get counsellor, client ids to get details
		// moodIDs := UTIL.ExtractValuesFromArrayMap(contents, "mood_id")
		// categoryIDs := UTIL.ExtractValuesFromArrayMap(contents, "category_id")

		// // get mood details
		// moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodsTable + " where id in ('" + strings.Join(moodIDs, "','") + "')")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		moodsMap := UTIL.ConvertMapToKeyMap(moods, "id")
		categorysMap := UTIL.ConvertMapToKeyMap(categorys, "id")

		for _, content := range contents {

			var types, trainingType, status, moodTitle, category string

			splitMood := strings.Split(content["mood_id"], ",")

			splitMoodTitle := []string{}

			for _, v := range splitMood {
				if _, ok := moodsMap[v]; ok {
					splitMoodTitle = append(splitMoodTitle, moodsMap[v]["title"])
				}
			}

			// // get mood details
			// moods, status, ok := DB.SelectProcess("select * from " + CONSTANT.MoodsTable + " where id in ('" + strings.Join(splitMood, "','") + "')")
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }
			if len(splitMoodTitle) > 0 {
				moodTitle = strings.Join(splitMoodTitle, ", ")
			} else {
				moodTitle = "None"
			}

			splitCategory := strings.Split(content["category_id"], ",")

			splitCategoryTitle := []string{}

			for _, v := range splitCategory {
				if _, ok := categorysMap[v]; ok {
					splitCategoryTitle = append(splitCategoryTitle, categorysMap[v]["category"])
				}
			}

			// get mood details
			// categorys, status, ok := DB.SelectProcess("select * from " + CONSTANT.ContentCategoriesTable + " where id in ('" + strings.Join(splitCategory, "','") + "')")
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			if len(splitCategoryTitle) > 0 {
				category = strings.Join(splitCategoryTitle, ", ")
			} else {
				category = "None"
			}

			// fmt.Println("category", category)

			switch content["type"] {
			case "1":
				types = "Video"
			case "2":
				types = "Audio"
			default:
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
				content["views"],
				category,
				moodTitle,
				trainingType,
				content["created_by"],
				UTIL.ConvertTimezone(UTIL.ConvertToTime(content["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				status,
			})
		}

	case "17": // counsellor time sheet

		fileName = "counsellor_time_sheet_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Counsellor Name", "Company Name", "Location", "In-Date", "In-Time", "Out-Date", "Out-Time", "Duration"}

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

			var inTime, inDate, outTime, outDate, duration string

			inTimeDateAndTime := strings.Split(timeSheet["inTime"], " ")
			outTimeDateAndTime := strings.Split(timeSheet["outTime"], " ")

			if len(inTimeDateAndTime) > 1 {
				inTime = inTimeDateAndTime[1]
				inDate = inTimeDateAndTime[0]
			}

			if len(outTimeDateAndTime) > 1 {
				outTime = outTimeDateAndTime[1]
				outDate = outTimeDateAndTime[0]
			}

			duration = strconv.Itoa(int(UTIL.BuildToDteTime(timeSheet["outTime"]).Sub(UTIL.BuildToDteTime(timeSheet["inTime"])).Hours())) + " hours " + strconv.Itoa(int(UTIL.BuildToDteTime(timeSheet["outTime"]).Sub(UTIL.BuildToDteTime(timeSheet["inTime"])).Minutes())%60) + " minutes"

			data = append(data, []string{
				counsellorsMap[timeSheet["counsellor_id"]]["first_name"] + " " + counsellorsMap[timeSheet["counsellor_id"]]["last_name"],
				timeSheet["partner_name"],
				timeSheet["location"],
				inDate,
				inTime,
				outDate,
				outTime,
				duration,
			})
		}
	case "18": // total number of session

		fileName = "total_number_of_session_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointment := range appointments {

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[appointment["client_id"]]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointment["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNamesMap[domainName[1]]["partner_name"]
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

		fileName = "total_counsellor_session_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
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

		fileName = "total_number_of_session_month_year_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Company Name", "IsFamilyMember", "Total Number Of Session", "Month", "Year"}

		var partnerName, isFamilyMember string

		appointments, status, ok := DB.SelectProcess("select client_id, count(client_id) as total, monthname(date) as month, year(date) as year from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' and type = '4' and status = '3' and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null) group by client_id having count(client_id) = 1")
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointment := range appointments {

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[appointment["client_id"]]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointment["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNamesMap[domainName[1]]["partner_name"]
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

		fileName = "in_person_appointment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Gender", "Age", "Mobile No.", "Email", "Company Name", "Location", "Department", "Counsellor Name", "Counsellor Type", "Date", "Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Mod. At", "Status", "Is Record Filled", "Mental Health Scale", "Duration"}
		appointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonAppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}
		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointments, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointments, "counsellor_id")
		appointmentIDs := UTIL.ExtractValuesFromArrayMap(appointments, "appointment_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		counsellorRecord, status, ok := DB.SelectProcess("select mental_health, appointment_id, noshow from " + CONSTANT.CounsellorRecordsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "') ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		counsellorRecordMap := UTIL.ConvertMapToKeyMap(counsellorRecord, "appointment_id")
		count := 0

		for _, appointment := range appointments {

			var startTime, endTime, modAt, partnerName, status, location, isRecordFilled, mentalHealth, duration string

			duration = "0 minutes"

			if appointment["started_at"] == "" {
				startTime = ""
			} else {
				startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["ended_at"] == "" {
				endTime = ""
			} else {
				endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			if appointment["modified_at"] == "" {
				modAt = ""
			} else {
				modAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
			}

			switch appointment["status"] {
			case "1":
				status = getAppointmentStatusInText("1")
			case "2":
				status = getAppointmentStatusInText("2")
			case "3":
				count++
				status = getAppointmentStatusInText("3")
				duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes())) + " minutes"

			case "4":
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("12")
				} else {
					status = getAppointmentStatusInText("4")
				}
			case "5":
				if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
					status = getAppointmentStatusInText("13")
				} else {
					status = getAppointmentStatusInText("5")
				}
			default:
				status = getAppointmentStatusInText(appointment["status"])
			}

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
			} else {
				partnerName = "None"
			}

			if len(counsellorRecordMap[appointment["appointment_id"]]) > 0 {
				isRecordFilled = "Yes"
				mentalHealth = counsellorRecordMap[appointment["appointment_id"]]["mental_health"]
				if counsellorRecordMap[appointment["appointment_id"]]["noshow"] == "1" {
					status = getAppointmentStatusInText("7")
					duration = "0 minutes"
				}
			} else {
				isRecordFilled = "No"
				mentalHealth = ""
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
				clientsMap[appointment["client_id"]]["phone"],
				clientsMap[appointment["client_id"]]["email"],
				partnerName,
				location,
				clientsMap[appointment["client_id"]]["department"],
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				counsellorsMap[appointment["counsellor_id"]]["type"],
				UTIL.BuildOnlyDateInDDMMYYYY(appointment["date"]),
				UTIL.GetTimeFromTimeSlot(appointment["time"]),
				appointment["times_rescheduled"],
				startTime,
				endTime,
				modAt,
				status,
				isRecordFilled,
				mentalHealth,
				duration,
			})

		}
		fmt.Println(data)

	case "22": // inperson request appointment

		fileName = "in_person_appointment_request_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Email ID", "Mobile No.", "Counsellor Name", "Company Name", "Location", "Request At.", "Updated At.", "Status"}
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
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

			switch inpersonappointmentRequest["status"] {
			case "1":
				status1 = "InProgress"
			case "2":
				status1 = "Completed"
			}

			requestDate := UTIL.BuildDate(inpersonappointmentRequest["created_at"])

			var updatedDate string

			if len(inpersonappointmentRequest["modified_at"]) > 0 {
				updatedDate = UTIL.BuildDate(inpersonappointmentRequest["modified_at"])
			} else {
				updatedDate = ""
			}

			data = append(data, []string{
				clientsMap[inpersonappointmentRequest["client_id"]]["first_name"] + " " + clientsMap[inpersonappointmentRequest["client_id"]]["last_name"],
				clientsMap[inpersonappointmentRequest["client_id"]]["email"],
				clientsMap[inpersonappointmentRequest["client_id"]]["phone"],
				counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["first_name"] + " " + counsellorsMap[inpersonappointmentRequest["counsellor_id"]]["last_name"],
				inpersonappointmentRequest["company_name"],
				inpersonappointmentRequest["company_location"],
				requestDate,
				updatedDate,
				status1,
			})
		}

	case "23": // appointment request

		fileName = "appointment_request_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Email", "Mobile No.", "Company Name", "IsFamilyMember", "Counsellor Name", "Request At.", "Updated At.", "Status"}
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
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

		// get company name details
		companyNames, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyNames, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointmentRequest := range appointmentsRequests {

			// get client details
			// clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id = '" + inpersonappointmentRequest["client_id"] + "'")
			// if !ok {
			// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			// 	return
			// }

			domainName := strings.Split(clientsMap[appointmentRequest["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				partnerN := ""

				if clientsMap[appointmentRequest["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clients[0]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointmentRequest["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNamesMap[domainName[1]]["partner_name"]
					isFamilyMember = "Yes"
				} else {
					partnerN = "None"
					isFamilyMember = "No"
				}

				partnerName = partnerN
			}

			switch appointmentRequest["status"] {
			case "1":
				status1 = "InProgress"
			case "2":
				status1 = "Completed"
			}

			requestDate := UTIL.BuildDate(appointmentRequest["created_at"])
			var updatedDate string

			if len(appointmentRequest["modified_at"]) > 0 {
				updatedDate = UTIL.BuildDate(appointmentRequest["modified_at"])
			} else {
				updatedDate = ""
			}

			data = append(data, []string{
				clientsMap[appointmentRequest["client_id"]]["first_name"] + " " + clientsMap[appointmentRequest["client_id"]]["last_name"],
				clientsMap[appointmentRequest["client_id"]]["email"],
				clientsMap[appointmentRequest["client_id"]]["phone"],
				partnerName,
				isFamilyMember,
				counsellorsMap[appointmentRequest["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointmentRequest["counsellor_id"]]["last_name"],
				requestDate,
				updatedDate,
				status1,
			})
		}

	case "24": // inperson cafe appointment

		fileName = "in_person_cafe_appointment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Cafe Name", "Counsellor Name", "Client Name", "Client Age", "Client Gender", "Client Mobile No.", "Client Email", "Attended", "Company Name", "Company Location", "Question1", "Question2", "Question3", "Question4", "Question5", "Cancellation Reason", "User Status", "Event Status", "Created At"}

		inPersonCafeAppointments, status, ok := DB.SelectProcess("select * from " + CONSTANT.OrderEventInPersonTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		userIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeAppointments, "user_id")
		eventOrderIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeAppointments, "event_order_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(userIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		orders, status, ok := DB.SelectProcess("select counsellor_id, order_id, title, company_name, company_location from " + CONSTANT.OrderCounsellorEventInPersonTable + " where order_id in ('" + strings.Join(eventOrderIDs, "','") + "')")
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

			switch inPersonCafeAppointment["status"] {
			case "1":
				status = "Active"
			case "4":
				status = "User Cancelled"
			}

			switch ordersMap[inPersonCafeAppointment["event_order_id"]]["status"] {
			case "1":
				eventStatus = "Active"
			case "2":
				eventStatus = "Started"
			case "3":
				eventStatus = "Completed"
			default:
				eventStatus = "Cancelled"
			}

			data = append(data, []string{
				ordersMap[inPersonCafeAppointment["event_order_id"]]["title"],
				counsellorsMap[ordersMap[inPersonCafeAppointment["event_order_id"]]["counsellor_id"]]["first_name"] + " " + counsellorsMap[ordersMap[inPersonCafeAppointment["event_order_id"]]["counsellor_id"]]["last_name"],
				clientsMap[inPersonCafeAppointment["user_id"]]["first_name"] + " " + clientsMap[inPersonCafeAppointment["user_id"]]["last_name"],
				clientsMap[inPersonCafeAppointment["user_id"]]["age"],
				clientsMap[inPersonCafeAppointment["user_id"]]["gender"],
				clientsMap[inPersonCafeAppointment["user_id"]]["phone"],
				clientsMap[inPersonCafeAppointment["user_id"]]["email"],
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
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(inPersonCafeAppointment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
			})
		}

	case "25": // inperson cafe request

		fileName = "in_person_cafe_request_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Client Name", "Email ID", "Mobile No.", "Cafe Name", "Total Seat", "Company Name", "Company Location", "Counsellor Name", "Status", "Requested At", "Updated At"}
		inPersonCafeRequests, status, ok := DB.SelectProcess("select * from " + CONSTANT.EventInPersonRequestTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeRequests, "client_id")
		orderIDs := UTIL.ExtractValuesFromArrayMap(inPersonCafeRequests, "order_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

			switch inPersonCafeRequest["status"] {
			case "1":
				status1 = "InProgress"
			case "2":
				status1 = "Completed"
			}

			requestDate := UTIL.BuildDate(inPersonCafeRequest["created_at"])
			var updatedDate string
			if len(inPersonCafeRequest["modified_at"]) > 0 {
				updatedDate = UTIL.BuildDate(inPersonCafeRequest["modified_at"])
			} else {
				updatedDate = ""
			}

			data = append(data, []string{
				clientsMap[inPersonCafeRequest["client_id"]]["first_name"] + " " + clientsMap[inPersonCafeRequest["client_id"]]["last_name"],
				clientsMap[inPersonCafeRequest["client_id"]]["email"],
				clientsMap[inPersonCafeRequest["client_id"]]["phone"],
				ordersMap[inPersonCafeRequest["order_id"]]["title"],
				ordersMap[inPersonCafeRequest["order_id"]]["total_seat"],
				ordersMap[inPersonCafeRequest["order_id"]]["company_name"],
				ordersMap[inPersonCafeRequest["order_id"]]["company_location"],
				counsellorsMap[ordersMap[inPersonCafeRequest["order_id"]]["counsellor_id"]]["first_name"] + " " + counsellorsMap[ordersMap[inPersonCafeRequest["order_id"]]["counsellor_id"]]["last_name"],
				status1,
				requestDate,
				updatedDate,
			})
		}
	case "26": // family onboarding report

		fileName = "family_onboarding_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Company Name", "Client Name", "Family Member Relationship", "Family Member Name", "Family Member Email", "Family Member Phone", "Family Member Age", "Family Member Active Time", "Family Member Login At", "Family Member Status", "Family Member Created At"}
		familyOnboarding, status, ok := DB.SelectProcess("select client_id, asscoiate_id, relation, first_name, last_name, email, phone, year(curdate())-year(date_of_birth) as age, last_active_time, last_login_time, status, created_at from " + CONSTANT.ClientsTable + " where asscoiate_id != '' and `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(familyOnboarding, "asscoiate_id")
		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyName, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		companyNamesMap := UTIL.ConvertMapToKeyMap(companyName, "domain")

		for _, family := range familyOnboarding {
			var partnerName string

			domainName := strings.Split(clientsMap[family["asscoiate_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNamesMap[domainName[1]]) > 0 {
				partnerName = companyNamesMap[domainName[1]]["partner_name"]
			} else {
				partnerName = "None"
			}
			data = append(data, []string{
				partnerName,
				clientsMap[family["asscoiate_id"]]["first_name"] + " " + clientsMap[family["asscoiate_id"]]["last_name"],
				family["relation"],
				family["first_name"] + " " + family["last_name"],
				family["email"],
				family["phone"],
				strconv.Itoa(int(getInt(family["age"]))),
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(family["last_active_time"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(family["last_login_time"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
				family["status"],
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(family["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
			})
		}

	case "27": // virtual or in-person report

		fileName = "virtual_in_person_appointment_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"
		heading = []string{"Mode", "Client Name", "Gender", "Age", "Email", "Mobile No.", "Company Name", "IsFamilyMember", "Location", "Counsellor Name", "Date", "Time", "No. of Reschedule", "Therapist Start", "Therapist End", "Client Start", "Client End", "Duration", "Created At", "Status", "Is Record Filled", "Mental Health Scale"}

		var appointmentsCombinationReport []map[string]string

		appointments, status, ok := DB.SelectProcess("select client_id, counsellor_id, appointment_id, started_at, client_started_at, client_ended_at, created_at, status, ended_at, date, time, times_rescheduled, 'virtual' AS mode from " + CONSTANT.AppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(appointments) > 0 {
			appointmentsCombinationReport = append(appointmentsCombinationReport, appointments...)
		}

		appointmentsInPerson, status, ok := DB.SelectProcess("select client_id, counsellor_id, appointment_id, started_at, client_started_at, client_ended_at, created_at, status, ended_at, date, time, times_rescheduled, 'in-person' AS mode from " + CONSTANT.InPersonAppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		if len(appointmentsInPerson) > 0 {
			appointmentsCombinationReport = append(appointmentsCombinationReport, appointmentsInPerson...)
		}

		// get counsellor, client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(appointmentsCombinationReport, "client_id")
		counsellorIDs := UTIL.ExtractValuesFromArrayMap(appointmentsCombinationReport, "counsellor_id")
		appointmentIDs := UTIL.ExtractValuesFromArrayMap(appointmentsCombinationReport, "appointment_id")

		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
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

		// get client details
		counsellorRecord, status, ok := DB.SelectProcess("select mental_health, appointment_id, noshow from " + CONSTANT.CounsellorRecordsTable + " where appointment_id in ('" + strings.Join(appointmentIDs, "','") + "') ")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyName, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client details
		clientFamilyMembers, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where asscoiate_id  != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		extractClientIDsFromFamilyMembers := UTIL.ExtractValuesFromArrayMap(clientFamilyMembers, "asscoiate_id")

		// get client details
		extractClientFromFamilyMembers, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(extractClientIDsFromFamilyMembers, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		counsellorsMap := UTIL.ConvertMapToKeyMap(counsellors, "id")
		counsellorRecordMap := UTIL.ConvertMapToKeyMap(counsellorRecord, "appointment_id")
		companyNameMap := UTIL.ConvertMapToKeyMap(companyName, "domain")
		extractClientFromFamilyMembersMap := UTIL.ConvertMapToKeyMap(extractClientFromFamilyMembers, "client_id")

		for _, appointment := range appointmentsCombinationReport {

			var duration, mode, startTime, endTime, createdAt, clientStartTime, clientEndTime, partnerName, isFamilyMember, status, location, isRecordFilled, mentalHealth string

			if appointment["mode"] == "virtual" {

				mode = "Virtual"

				if appointment["started_at"] == "" {
					startTime = ""
				} else {
					startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["client_started_at"] == "" {
					clientStartTime = ""
				} else {
					clientStartTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["client_ended_at"] == "" {
					clientEndTime = ""
				} else {
					clientEndTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["client_ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["ended_at"] == "" {
					endTime = ""
				} else {
					endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["created_at"] == "" {
					createdAt = ""
				} else {
					createdAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				switch appointment["status"] {
				case "3":
					if appointment["started_at"] == "" && appointment["ended_at"] == "" {
						status = getAppointmentStatusInText("8")
						if appointment["client_started_at"] != "" && appointment["client_ended_at"] != "" {
							duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes())) + " minutes"
						}
					} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
						status = getAppointmentStatusInText("7")
						duration = "0 minutes"
					} else {
						// calculate the duration

						if len(appointment["client_started_at"]) != 0 && len(appointment["client_ended_at"]) == 0 {
							duration = "5 minutes"
						} else if len(appointment["client_started_at"]) == 0 && len(appointment["client_ended_at"]) != 0 {
							duration = "5 minutes"
						} else {
							duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["client_ended_at"]).Sub(UTIL.BuildToDteTime(appointment["client_started_at"])).Minutes())) + " minutes"
						}

						if UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes() > 11 {
							status = getAppointmentStatusInText("3")
						} else {
							status = getAppointmentStatusInText("14")
						}

					}
				case "4":
					if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
						status = getAppointmentStatusInText("12")
					} else {
						status = getAppointmentStatusInText("4")
					}
				case "5":
					if UTIL.BuildDateTime(appointment["date"], appointment["time"]).Sub(UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["modified_at"]), "330")).Hours() <= 4 {
						status = getAppointmentStatusInText("13")
					} else {
						status = getAppointmentStatusInText("5")
					}
				default:
					status = getAppointmentStatusInText(appointment["status"])
					if appointment["status"] == "2" {
						// duration
						if appointment["client_started_at"] != "" || appointment["client_started_at"] == "" {
							duration = "5 minutes"
						} else if appointment["client_started_at"] == "" && appointment["client_ended_at"] == "" {
							duration = "5 minutes"
						}

					}
				}
			} else {

				mode = "In-Person"

				clientStartTime = ""
				clientEndTime = ""

				if appointment["started_at"] == "" {
					startTime = ""
				} else {
					startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["ended_at"] == "" {
					endTime = ""
				} else {
					endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				if appointment["started_at"] != "" && appointment["ended_at"] != "" {
					status = getAppointmentStatusInText("3")
					duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes())) + " minutes"
				} else if appointment["started_at"] != "" && appointment["ended_at"] == "" {
					status = getAppointmentStatusInText("14")
					duration = "5 minutes"
				} else if appointment["started_at"] == "" && appointment["ended_at"] != "" {
					status = getAppointmentStatusInText("14")
					duration = "5 minutes"
				} else {
					status = getAppointmentStatusInText("14")
					duration = "0 minutes"
				}

				if appointment["created_at"] == "" {
					createdAt = ""
				} else {
					createdAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat)
				}

				status = getAppointmentStatusInText("8")
			}

			if len(counsellorRecordMap[appointment["appointment_id"]]) > 0 {
				isRecordFilled = "Yes"
				mentalHealth = counsellorRecordMap[appointment["appointment_id"]]["mental_health"]
				if counsellorRecordMap[appointment["appointment_id"]]["noshow"] == "1" {
					status = getAppointmentStatusInText("7")
					duration = "0 minutes"
				}
			} else {
				isRecordFilled = "No"
				mentalHealth = ""
			}

			domainName := strings.Split(clientsMap[appointment["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNameMap[domainName[1]]) > 0 {
				partnerName = companyNameMap[domainName[1]]["partner_name"]
				isFamilyMember = "No"
			} else {
				// // get client details
				// clients, status, ok := DB.SelectProcess("select client_id, asscoiate_id from " + CONSTANT.ClientsTable + " where client_id = '" + appointment["client_id"] + "'")
				// if !ok {
				// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				// 	return
				// }

				partnerN := ""

				if clientsMap[appointment["client_id"]]["asscoiate_id"] != "" {
					// get client details
					// clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email from " + CONSTANT.ClientsTable + " where client_id = '" + clientsMap[appointment["client_id"]]["asscoiate_id"] + "'")
					// if !ok {
					// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
					// 	return
					// }

					domainName := strings.Split(extractClientFromFamilyMembersMap[clientsMap[appointment["client_id"]]["asscoiate_id"]]["email"], "@")

					// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

					partnerN = companyNameMap[domainName[1]]["partner_name"]
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
				mode,
				clientsMap[appointment["client_id"]]["first_name"] + " " + clientsMap[appointment["client_id"]]["last_name"],
				clientsMap[appointment["client_id"]]["gender"],
				clientsMap[appointment["client_id"]]["age"],
				clientsMap[appointment["client_id"]]["email"],
				clientsMap[appointment["client_id"]]["phone"],
				partnerName,
				isFamilyMember,
				location,
				counsellorsMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsMap[appointment["counsellor_id"]]["last_name"],
				UTIL.BuildOnlyDateInDDMMYYYY(appointment["date"]),
				UTIL.GetTimeFromTimeSlot(appointment["time"]),
				appointment["times_rescheduled"],
				startTime,
				endTime,
				clientStartTime,
				clientEndTime,
				duration,
				createdAt,
				status,
				isRecordFilled,
				mentalHealth,
			})
		}

		// appointmentsInPerson, status, ok := DB.SelectProcess("select * from " + CONSTANT.InPersonAppointmentsTable + " where `date` >= '" + startBy.String() + "' and `date` <= '" + endBy.String() + "' order by created_at desc")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		// get counsellor, client ids to get details
		// clientInPersonIDs := UTIL.ExtractValuesFromArrayMap(appointmentsInPerson, "client_id")
		// counsellorInPersonIDs := UTIL.ExtractValuesFromArrayMap(appointmentsInPerson, "counsellor_id")

		// // get client details
		// clientsInPerson, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientInPersonIDs, "','") + "')")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		// // get counsellor details
		// counsellorsInPerson, status, ok := DB.SelectProcess("(select counsellor_id as id, first_name, last_name, 'Counsellor' as type from " + CONSTANT.CounsellorsTable + " where counsellor_id in ('" + strings.Join(counsellorInPersonIDs, "','") + "')) union (select listener_id as id, first_name, last_name, 'Listener' as type from " + CONSTANT.ListenersTable + " where listener_id in ('" + strings.Join(counsellorInPersonIDs, "','") + "')) union (select therapist_id as id, first_name, last_name, 'Therapist' as type from " + CONSTANT.TherapistsTable + " where therapist_id in ('" + strings.Join(counsellorInPersonIDs, "','") + "'))")
		// if !ok {
		// 	UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 	return
		// }

		// clientsInPersonMap := UTIL.ConvertMapToKeyMap(clientsInPerson, "client_id")
		// counsellorsInPersonMap := UTIL.ConvertMapToKeyMap(counsellorsInPerson, "id")

		// for _, appointment := range appointmentsInPerson {

		// 	var startTime, endTime, createdAt, partnerName, statusValue, location, duration, isRecordFilled, mentalHealth string

		// 	if appointment["started_at"] == "" {
		// 		startTime = ""
		// 	} else {
		// 		startTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["started_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
		// 	}

		// 	if appointment["ended_at"] == "" {
		// 		endTime = ""
		// 	} else {
		// 		endTime = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["ended_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
		// 	}

		// 	if appointment["started_at"] != "" && appointment["ended_at"] != "" {
		// 		duration = strconv.Itoa(int(UTIL.BuildToDteTime(appointment["ended_at"]).Sub(UTIL.BuildToDteTime(appointment["started_at"])).Minutes())) + " minutes"
		// 	} else if appointment["started_at"] != "" && appointment["ended_at"] == "" {
		// 		duration = "5 minutes"
		// 	} else if appointment["started_at"] == "" && appointment["ended_at"] != "" {
		// 		duration = "5 minutes"
		// 	} else {
		// 		duration = "0 minutes"
		// 	}

		// 	if appointment["created_at"] == "" {
		// 		createdAt = ""
		// 	} else {
		// 		createdAt = UTIL.ConvertTimezone(UTIL.BuildToDteTime(appointment["created_at"]), "330").Format(CONSTANT.ReadbleDateTimeFormat)
		// 	}

		// 	statusValue = getAppointmentStatusInText("3")

		// 	domainName := strings.Split(clientsInPersonMap[appointment["client_id"]]["email"], "@")

		// 	title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

		// 	if len(title) > 0 {
		// 		partnerName = title[0]["partner_name"]
		// 	} else {
		// 		partnerName = "None"
		// 	}

		// 	if clientsInPersonMap[appointment["client_id"]]["location"] == "40.04" {
		// 		location = ""
		// 	} else {
		// 		location = clientsInPersonMap[appointment["client_id"]]["location"]
		// 	}

		// 	// get client details
		// 	counsellorRecord, status, ok := DB.SelectProcess("select * from "+CONSTANT.CounsellorRecordsTable+" where appointment_id = ? ", appointment["appointment_id"])
		// 	if !ok {
		// 		UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
		// 		return
		// 	}

		// 	if len(counsellorRecord) > 0 {
		// 		if counsellorRecord[0]["noshow"] == "1" {
		// 			statusValue = getAppointmentStatusInText("7")
		// 		}
		// 		isRecordFilled = "Yes"
		// 		mentalHealth = counsellorRecord[0]["mental_health"]
		// 	} else {
		// 		isRecordFilled = "No"
		// 		mentalHealth = ""
		// 	}

		// 	data = append(data, []string{
		// 		"In-Person",
		// 		clientsInPersonMap[appointment["client_id"]]["first_name"] + " " + clientsInPersonMap[appointment["client_id"]]["last_name"],
		// 		clientsInPersonMap[appointment["client_id"]]["gender"],
		// 		clientsInPersonMap[appointment["client_id"]]["age"],
		// 		clientsInPersonMap[appointment["client_id"]]["phone"],
		// 		clientsInPersonMap[appointment["client_id"]]["email"],
		// 		partnerName,
		// 		"No",
		// 		location,
		// 		counsellorsInPersonMap[appointment["counsellor_id"]]["first_name"] + " " + counsellorsInPersonMap[appointment["counsellor_id"]]["last_name"],
		// 		UTIL.BuildOnlyDateInDDMMYYYY(appointment["date"]),
		// 		UTIL.GetTimeFromTimeSlot(appointment["time"]),
		// 		appointment["times_rescheduled"],
		// 		startTime,
		// 		endTime,
		// 		"",
		// 		"",
		// 		duration,
		// 		createdAt,
		// 		statusValue,
		// 		isRecordFilled,
		// 		mentalHealth,
		// 	})
		// }

	case "28": // webinar request report

		fileName = "webinar_request_report_" + startBy.Format("02Jan2006") + "_to_" + endBy.Format("02Jan2006") + ".csv"

		heading = []string{"Webinar Name", "Company Name", "Client Name", "Email", "Phone", "Age", "Gender", "Location", "Department", "Request At"}
		webinars, status, ok := DB.SelectProcess("select webinar_id, client_id, created_at from " + CONSTANT.WebinarsBookTable + " where `created_at` >= '" + startBy.String() + "' and `created_at` <= '" + endBy.String() + "' order by created_at desc")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get client ids to get details
		clientIDs := UTIL.ExtractValuesFromArrayMap(webinars, "client_id")
		// get webinar ids to get details
		webinarIDs := UTIL.ExtractValuesFromArrayMap(webinars, "webinar_id")
		// get client details
		clients, status, ok := DB.SelectProcess("select client_id, first_name, last_name, email, phone, gender, year(curdate())-year(date_of_birth) as age, location, department  from " + CONSTANT.ClientsTable + " where client_id in ('" + strings.Join(clientIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get webinar Names details
		webinarNames, status, ok := DB.SelectProcess("select * from " + CONSTANT.WebinarsTable + " where webinar_id in ('" + strings.Join(webinarIDs, "','") + "')")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		// get company name details
		companyName, status, ok := DB.SelectProcess("select partner_name, domain from " + CONSTANT.CorporatePartnersTable + "")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientsMap := UTIL.ConvertMapToKeyMap(clients, "client_id")
		webinarNamesMap := UTIL.ConvertMapToKeyMap(webinarNames, "webinar_id")
		companyNameMap := UTIL.ConvertMapToKeyMap(companyName, "domain")

		for _, webinar := range webinars {
			var partnerName string

			domainName := strings.Split(clientsMap[webinar["client_id"]]["email"], "@")

			// title, _, _ := DB.SelectSQL(CONSTANT.CorporatePartnersTable, []string{"partner_name", "domain"}, map[string]string{"domain": domainName[1]})

			if len(companyNameMap[domainName[1]]) > 0 {
				partnerName = companyNameMap[domainName[1]]["partner_name"]
			} else {
				partnerName = "None"
			}
			data = append(data, []string{
				webinarNamesMap[webinar["webinar_id"]]["title"],
				partnerName,
				clientsMap[webinar["client_id"]]["first_name"] + " " + clientsMap[webinar["client_id"]]["last_name"],
				clientsMap[webinar["client_id"]]["email"],
				clientsMap[webinar["client_id"]]["phone"],
				clientsMap[webinar["client_id"]]["age"],
				clientsMap[webinar["client_id"]]["gender"],
				clientsMap[webinar["client_id"]]["location"],
				clientsMap[webinar["client_id"]]["department"],
				UTIL.ConvertTimezone(UTIL.BuildToDteTime(webinar["created_at"]), "330").Format(CONSTANT.ReadableDateTimeYearFormat),
			})
		}

	default:
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "Invalid report type", CONSTANT.ShowDialog, response)
		return

	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

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
	case CONSTANT.AppointmentIncompleteSession:
		return "Incomplete session"
	case CONSTANT.AppointmentIncompleteSessionDueToTechIssue:
		return "Incomplete session(tech issue)"
	default:
		return ""
	}
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

	appointmentTotal, clientTotal, clientFamilyAppointment, webinarTotal, appointmentsInPersonTotal, emeCaseVirtualTotal, emeCaseInPersonTotal, contentsTotal, moodsTotal, assessmentsTotal, totalRatingTotal, avgRatingTotal, appointmentsCancellationTotal, appointmentsNoShowTotal, familyMemberTotal, companyName, inpersonCafe, inpersonCafeAttend := "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""

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

			clientIDsFamily, status, ok := DB.SelectProcess("select client_id from " + CONSTANT.ClientsTable + " where asscoiate_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientIDsWithFamily := UTIL.ExtractValuesFromArrayMap(clientIDsFamily, "client_id")

			appointmentsFamily, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '3' and client_id in ('"+strings.Join(clientIDsWithFamily, "','")+"') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientFamilyAppointment = appointmentsFamily[0]["total"]

			webinar, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.WebinarsTable+whereAppointment+" and status = '1' and partner_name = '"+companyNames[0]["partner_name"]+"'", queryArgs...)
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			webinarTotal = webinar[0]["total"]

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

			familyMembers, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.ClientsTable + " where status = 1 and client_id in ('" + strings.Join(asscoiateIDs, "','") + "') and email like '%" + companyName + "'")
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

			clientIDsFamily, status, ok := DB.SelectProcess("select client_id from " + CONSTANT.ClientsTable + " where asscoiate_id in ('" + strings.Join(clientIDs, "','") + "')")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientIDsWithFamily := UTIL.ExtractValuesFromArrayMap(clientIDsFamily, "client_id")

			appointmentsFamily, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '3' and client_id in ('" + strings.Join(clientIDsWithFamily, "','") + "') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			clientFamilyAppointment = appointmentsFamily[0]["total"]

			webinar, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.WebinarsTable + " where status = '1' and partner_name = '" + companyNames[0]["partner_name"] + "'")
			if !ok {
				UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
				return
			}

			webinarTotal = webinar[0]["total"]

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

		clientIDsFamily, status, ok := DB.SelectProcess("select client_id from " + CONSTANT.ClientsTable + " where asscoiate_id != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDsWithFamily := UTIL.ExtractValuesFromArrayMap(clientIDsFamily, "client_id")

		appointmentsFamily, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.AppointmentsTable+whereAppointment+" and type = '4' and status = '3' and client_id in ('"+strings.Join(clientIDsWithFamily, "','")+"') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientFamilyAppointment = appointmentsFamily[0]["total"]

		webinar, status, ok := DB.SelectProcess("select count(*) as total from "+CONSTANT.WebinarsTable+whereAppointment+" and status = '1'", queryArgs...)
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		webinarTotal = webinar[0]["total"]

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

		clientIDsFamily, status, ok := DB.SelectProcess("select client_id from " + CONSTANT.ClientsTable + " where asscoiate_id != ''")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientIDsWithFamily := UTIL.ExtractValuesFromArrayMap(clientIDsFamily, "client_id")

		appointmentsFamily, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.AppointmentsTable + " where type = '4' and status = '3' and client_id in ('" + strings.Join(clientIDsWithFamily, "','") + "') and (client_started_at is not null or client_ended_at is not null) and (started_at is not null or ended_at is not null)")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		clientFamilyAppointment = appointmentsFamily[0]["total"]

		webinar, status, ok := DB.SelectProcess("select count(*) as total from " + CONSTANT.WebinarsTable + " where status = '1'")
		if !ok {
			UTIL.SetReponse(w, status, "", CONSTANT.ShowDialog, response)
			return
		}

		webinarTotal = webinar[0]["total"]

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
	response["total_appointment_family"] = clientFamilyAppointment
	response["total_webinar"] = webinarTotal
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
