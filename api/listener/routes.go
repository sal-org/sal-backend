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
	listenerRoutes.HandleFunc("/appointment/start", AppointmentStart).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")
	listenerRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	// assessment
	listenerRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	listenerRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	listenerRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"assessment_id", "{assessment_id}",
		"client_id", "{client_id}",
	).Methods("GET")

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

	// notification
	listenerRoutes.HandleFunc("/notification", NotificationsGet).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")

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
