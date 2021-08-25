package therapist

import "github.com/gorilla/mux"

// LoadTherapistRoutes - load all therapist routes with therapist prefix
func LoadTherapistRoutes(router *mux.Router) {
	therapistRoutes := router.PathPrefix("/therapist").Subrouter()

	// availability
	therapistRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("PUT")

	// appointment
	therapistRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/appointment", AppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")
	therapistRoutes.HandleFunc("/appointment/start", AppointmentStart).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")
	therapistRoutes.HandleFunc("/appointment/end", AppointmentEnd).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")

	// assessment
	therapistRoutes.HandleFunc("/assessments", AssessmentsList).Methods("GET")
	therapistRoutes.HandleFunc("/assessment", AssessmentDetail).Queries(
		"assessment_id", "{assessment_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/assessment", AssessmentAdd).Methods("POST")
	therapistRoutes.HandleFunc("/assessment/history", AssessmentHistory).Queries(
		"assessment_id", "{assessment_id}",
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// event
	therapistRoutes.HandleFunc("/events", EventsList).Methods("GET")
	therapistRoutes.HandleFunc("/event", EventDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/event/booked", EventsBooked).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	therapistRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")
	therapistRoutes.HandleFunc("/event/block", EventsBlocked).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/event", EventUpdate).Queries(
		"order_id", "{order_id}",
		"therapist_id", "{therapist_id}",
		"status", "{status}",
	).Methods("PUT")
	therapistRoutes.HandleFunc("/event/block/order", EventBlockOrderCreate).Methods("POST")
	therapistRoutes.HandleFunc("/event/block/paymentcomplete", EventBlockOrderPaymentComplete).Methods("POST")

	// home
	therapistRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	therapistRoutes.Path("/refresh-token").Queries(
		"therapist_id", "{therapist_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// notification
	therapistRoutes.HandleFunc("/notification", NotificationsGet).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// payment
	therapistRoutes.HandleFunc("/payment", PaymentsGet).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// profile
	therapistRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	therapistRoutes.HandleFunc("", ProfileGet).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	therapistRoutes.HandleFunc("", ProfileUpdate).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("PUT")

	// training
	therapistRoutes.HandleFunc("/training", Training).Methods("GET")
}
