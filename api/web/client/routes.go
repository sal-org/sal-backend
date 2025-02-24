package client

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all web routes with web prefix
func LoadWebClientRoutes(router *mux.Router) {
	webClientRoutes := router.PathPrefix("/web/client").Subrouter()

	// client login
	webClientRoutes.HandleFunc("/access_code", CheckAccessCode).Methods("POST")
	webClientRoutes.HandleFunc("/send_otp", SendOTPWithCorporateEmail).Methods("POST")
	webClientRoutes.HandleFunc("/verify_otp", VerifyOTPWithCorporateEmail).Methods("POST")

	// client profile
	webClientRoutes.HandleFunc("/profile/get", ProfileGet).Methods("POST")
	webClientRoutes.HandleFunc("/profile/add", ProfileAdd).Methods("POST")
	webClientRoutes.HandleFunc("/profile/update", ProfileUpdate).Methods("POST")

	// client search
	webClientRoutes.HandleFunc("/search", ListSearch).Methods("POST")

	// client appointment
	webClientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/detail", AppointmentDetail).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/reschedule", AppointmentReschedule).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/cancel", AppointmentCancel).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/rate", AppointmentRatingAdd).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/start", AppointmentStart).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/end", AppointmentEnd).Methods("POST")

	// client agora token
	webClientRoutes.HandleFunc("/agora/token", GenerateAgoraToken).Methods("POST")

	// client therapist
	webClientRoutes.HandleFunc("/therapist/get", TherapistProfile).Methods("POST")
	webClientRoutes.HandleFunc("/therapist/slots", TherapistSlots).Methods("POST")
	webClientRoutes.HandleFunc("/therapist/order", CorporateCounsellorOrderCreate).Methods("POST")
	webClientRoutes.HandleFunc("/therapist/complete", CorporateCounsellorOrderPaymentComplete).Methods("POST")

}
