package qualitycheck

import "github.com/gorilla/mux"

func LoadQualityCheckRoutes(router *mux.Router) {
	qualityCheckRoutes := router.PathPrefix("/qualitycheck").Subrouter()

	//login
	qualityCheckRoutes.HandleFunc("/login", Login).Methods("POST")
	qualityCheckRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")
	// verify otp
	qualityCheckRoutes.HandleFunc("/verifyotp", VerifyOTP).Methods("POST")

	// appointment
	qualityCheckRoutes.HandleFunc("/appointment", GetAppointmentdetails).Methods("POST")


	// counsellor record 
	qualityCheckRoutes.HandleFunc("/counsellor-record", GetCounsellorRecord).Methods("POST")


	// counsellor time sheet 
	qualityCheckRoutes.HandleFunc("/counsellor-time-sheet", GetCounsellorTimeSheet).Methods("POST")

	qualityCheckRoutes.HandleFunc("/send-email", SendEmail).Methods("POST")

	qualityCheckRoutes.HandleFunc("/video-comment", VideoCallComment).Methods("POST")

	// sms
	qualityCheckRoutes.HandleFunc("/send-sms", SendSMS).Methods("POST")
}
