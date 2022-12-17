package qualitycheck

import "github.com/gorilla/mux"

func LoadQualityCheckRoutes(router *mux.Router) {
	qualityCheckRoutes := router.PathPrefix("/qualitycheck").Subrouter()

	//login
	qualityCheckRoutes.HandleFunc("/login", Login).Methods("POST")
	qualityCheckRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")

	// appointment
	qualityCheckRoutes.HandleFunc("/appointment", GetAppointmentdetails).Methods("POST")

	qualityCheckRoutes.HandleFunc("/send-email", SendEmail).Methods("POST")

	qualityCheckRoutes.HandleFunc("/video-comment", VideoCallComment).Methods("POST")
}
