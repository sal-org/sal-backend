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
		"client_id", "{client_id}",
	).Methods("GET")

	// event
	therapistRoutes.HandleFunc("/events", EventsList).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	therapistRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// home
	therapistRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	therapistRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	therapistRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")
	therapistRoutes.Path("/refresh-token").Queries(
		"therapist_id", "{therapist_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// notification
	therapistRoutes.HandleFunc("/notification", NotificationsGet).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// payment
	therapistRoutes.HandleFunc("/payment", PaymentsGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")

	// profile
	therapistRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	therapistRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	therapistRoutes.HandleFunc("", ProfileUpdate).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("PUT")

	// training
	therapistRoutes.HandleFunc("/training", Training).Methods("GET")
}
