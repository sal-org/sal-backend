package listener

import "github.com/gorilla/mux"

// LoadListenerRoutes - load all listener routes with listener prefix
func LoadListenerRoutes(router *mux.Router) {
	clientRoutes := router.PathPrefix("/listener").Subrouter()

	// availability
	clientRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// appointment
	clientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")

	// profile
	clientRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	clientRoutes.HandleFunc("", ProfileUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

}
