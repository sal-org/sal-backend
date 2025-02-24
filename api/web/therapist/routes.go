package therapist

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all web routes with web prefix
func LoadWebTherapistRoutes(router *mux.Router) {
	webClientRoutes := router.PathPrefix("/web/therapist").Subrouter()

	// client profile
	webClientRoutes.HandleFunc("/profile/get", ProfileGet).Methods("POST")
	webClientRoutes.HandleFunc("/profile/add", ProfileAdd).Methods("POST")
	webClientRoutes.HandleFunc("/profile/update", ProfileUpdate).Methods("POST")

	// client search
	webClientRoutes.HandleFunc("/availability/get", AvailabilityGet).Methods("POST")
	webClientRoutes.HandleFunc("/availability/update", AvailabilityUpdate).Methods("POST")

	// client appointment
	webClientRoutes.HandleFunc("/appointment/upcoming", AppointmentsUpcoming).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/past", AppointmentsPast).Methods("POST")
	webClientRoutes.HandleFunc("/appointment/cancel", AppointmentCancel).Methods("POST")



}