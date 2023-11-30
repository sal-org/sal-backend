package counsellor

import "github.com/gorilla/mux"

// LoadCounsellorRoutes - load all counsellor routes with counsellor prefix
func LoadCounsellorRoutes(router *mux.Router) {
	counsellorRoutes := router.PathPrefix("/counsellor").Subrouter()

	// availability
	counsellorRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

	// appointment
	counsellorRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/appointment", AppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")
	counsellorRoutes.HandleFunc("/appointment/agoratoken", GenerateAgoraToken).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/appointment/start", AppointmentStart).Queries(
		"appointment_id", "{appointment_id}",
		"uid", "{uid}",
	).Methods("PUT")
	counsellorRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
		"uid", "{uid}",
	).Methods("PUT")
	counsellorRoutes.HandleFunc("/appointment/comment", CounsellorComment).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	// assessment
	counsellorRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	counsellorRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	counsellorRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/assessment/download", AssessmentDownload).Queries(
		"assessment_result_id", "{assessment_result_id}",
	).Methods("GET")

	// event
	counsellorRoutes.HandleFunc("/events", EventsList).Methods("GET")
	counsellorRoutes.HandleFunc("/event", EventDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/event/booked", EventsBooked).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	counsellorRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")
	counsellorRoutes.HandleFunc("/event/block", EventsBlocked).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/event", EventUpdate).Queries(
		"order_id", "{order_id}",
		"counsellor_id", "{counsellor_id}",
		"status", "{status}",
	).Methods("PUT")
	counsellorRoutes.HandleFunc("/event/block/order", EventBlockOrderCreate).Methods("POST")
	//counsellorRoutes.HandleFunc("/event/block/paymentcomplete", EventBlockOrderPaymentComplete).Methods("POST")

	// home
	counsellorRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	counsellorRoutes.Path("/refresh-token").Queries(
		"counsellor_id", "{counsellor_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// notification
	counsellorRoutes.HandleFunc("/notification", NotificationsGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// payment
	counsellorRoutes.HandleFunc("/payment", PaymentsGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// profile
	counsellorRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("", ProfileGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	counsellorRoutes.HandleFunc("", ProfileUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

	// training
	counsellorRoutes.HandleFunc("/training", Training).Methods("GET")

	// get partner name
	counsellorRoutes.HandleFunc("/get-partner", PartnerGet).Methods("GET")

	// GetPartnerAddress
	counsellorRoutes.HandleFunc("/get-address", GetPartnerAddress).Queries(
		"client_name", "{client_name}",
	).Methods("GET")

	// add my time sheet
	counsellorRoutes.HandleFunc("/my-time-sheet", AddMyTimeSheet).Methods("POST")

	// get my time sheet
	counsellorRoutes.HandleFunc("/my-time-sheet", MyTimeSheet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// update my time sheet
	counsellorRoutes.HandleFunc("/my-time-sheet", UpdateMyTimeSheet).Queries(
		"sheet_id", "{sheet_id}",
	).Methods("PUT")
}
