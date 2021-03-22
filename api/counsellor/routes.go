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

	// event
	counsellorRoutes.HandleFunc("/events", EventsUpcoming).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	counsellorRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// login
	counsellorRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")

	// profile
	counsellorRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	counsellorRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	counsellorRoutes.HandleFunc("", ProfileUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

}
