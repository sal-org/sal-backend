package counsellor

import "github.com/gorilla/mux"

// LoadCounsellorRoutes - load all counsellor routes with counsellor prefix
func LoadCounsellorRoutes(router *mux.Router) {
	clientRoutes := router.PathPrefix("/counsellor").Subrouter()

	// availability
	clientRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

	// appointment
	clientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// event
	clientRoutes.HandleFunc("/events", EventsUpcoming).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// profile
	clientRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	clientRoutes.HandleFunc("", ProfileUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

}
