package listener

import "github.com/gorilla/mux"

// LoadListenerRoutes - load all listener routes with listener prefix
func LoadListenerRoutes(router *mux.Router) {
	listenerRoutes := router.PathPrefix("/listener").Subrouter()

	// availability
	listenerRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// appointment
	listenerRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/appointment", AppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")

	// home
	listenerRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	listenerRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")
	listenerRoutes.Path("/refresh-token").Queries(
		"listener_id", "{listener_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// profile
	listenerRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	listenerRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	listenerRoutes.HandleFunc("", ProfileUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// training
	listenerRoutes.HandleFunc("/training", Training).Methods("GET")
}
