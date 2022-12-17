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
		"uid", "{uid}",
	).Methods("PUT")
	listenerRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
		"uid", "{uid}",
	).Methods("PUT")

	// assessment
	listenerRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	listenerRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	listenerRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/assessment/download", AssessmentDownload).Queries(
		"assessment_result_id", "{assessment_result_id}",
	).Methods("GET")

	// event
	listenerRoutes.HandleFunc("/events", EventsList).Methods("GET")
	listenerRoutes.HandleFunc("/event", EventDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/event/booked", EventsBooked).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	listenerRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// home
	listenerRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
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
	listenerRoutes.HandleFunc("", ProfileGet).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	listenerRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	listenerRoutes.HandleFunc("", ProfileUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// training
	listenerRoutes.HandleFunc("/training", Training).Methods("GET")
}
