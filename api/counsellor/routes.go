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
	counsellorRoutes.HandleFunc("/appointment/start", AppointmentStart).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")
	counsellorRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	// assessment
	counsellorRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	counsellorRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	counsellorRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"assessment_id", "{assessment_id}",
		"client_id", "{client_id}",
	).Methods("GET")

	// event
	counsellorRoutes.HandleFunc("/events", EventsList).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	counsellorRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// home
	counsellorRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	counsellorRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")
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
	counsellorRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	counsellorRoutes.HandleFunc("", ProfileUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

	// training
	counsellorRoutes.HandleFunc("/training", Training).Methods("GET")
}
