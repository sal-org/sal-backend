package client

import "github.com/gorilla/mux"

// LoadClientRoutes - load all client routes with client prefix
func LoadClientRoutes(router *mux.Router) {
	clientRoutes := router.PathPrefix("/client").Subrouter()

	// appointment
	clientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/slots", AppointmentSlotsUnused).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment", AppointmentDetail).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/appointment", AppointmentBook).Methods("POST")
	clientRoutes.HandleFunc("/appointment", AppointmentReschedule).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")
	clientRoutes.HandleFunc("/appointment", AppointmentCancel).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("DELETE")
	clientRoutes.HandleFunc("/appointment/rate", AppointmentRatingAdd).Methods("POST")

	// content
	clientRoutes.HandleFunc("/content", Content).Methods("GET")
	clientRoutes.HandleFunc("/content/like", ContentLikeGet).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/content/like", ContentLikeAdd).Queries(
		"client_id", "{client_id}",
		"content_id", "{content_id}",
	).Methods("POST")
	clientRoutes.HandleFunc("/content/like", ContentLikeDelete).Queries(
		"client_id", "{client_id}",
		"content_id", "{content_id}",
	).Methods("DELETE")

	// counsellor
	clientRoutes.HandleFunc("/counsellor", CounsellorProfile).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/slots", CounsellorSlots).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/order", CounsellorOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/counsellor/paymentcomplete", CounsellorOrderPaymentComplete).Methods("POST")

	// event
	clientRoutes.HandleFunc("/events", EventsList).Methods("GET")
	clientRoutes.HandleFunc("/event", EventDetail).Queries(
		"order_id", "{order_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event/booked", EventsBooked).Queries(
		"client_id", "{client_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/event/order", EventOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/event/paymentcomplete", EventOrderPaymentComplete).Methods("POST")

	// home
	clientRoutes.HandleFunc("/home", Home).Methods("GET")

	// login
	clientRoutes.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	clientRoutes.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")
	clientRoutes.Path("/refresh-token").Queries(
		"client_id", "{client_id}",
	).HandlerFunc(RefreshToken).Methods("GET")

	// listener
	clientRoutes.HandleFunc("/listener", ListenerProfile).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/slots", ListenerSlots).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/order", ListenerOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/listener/paymentcomplete", ListenerOrderPaymentComplete).Methods("POST")

	// payment
	clientRoutes.HandleFunc("/paymentcomplete", CounsellorOrderPaymentComplete).Methods("POST")

	// profile
	clientRoutes.HandleFunc("", ProfileGet).Queries(
		"email", "{email}",
	).Methods("GET")
	clientRoutes.HandleFunc("", ProfileAdd).Methods("POST")
	clientRoutes.HandleFunc("", ProfileUpdate).Queries(
		"client_id", "{client_id}",
	).Methods("PUT")

	// search
	clientRoutes.HandleFunc("/search", ListSearch).Methods("GET")

	// therapist
	clientRoutes.HandleFunc("/therapist", TherapistProfile).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/therapist/slots", TherapistSlots).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/therapist/order", TherapistOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/therapist/paymentcomplete", TherapistOrderPaymentComplete).Methods("POST")

}
