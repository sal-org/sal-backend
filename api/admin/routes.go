package admin

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all admin routes with admin prefix
func LoadAdminRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/admin").Subrouter()

	// content
	adminRoutes.HandleFunc("/appointment", AppointmentGet).Methods("GET")

	// client
	adminRoutes.HandleFunc("/client", ClientGet).Methods("GET")
	adminRoutes.HandleFunc("/client", ClientUpdate).Queries(
		"client_id", "{client_id}",
	).Methods("PUT")

	// content
	adminRoutes.HandleFunc("/content", ContentGet).Methods("GET")
	adminRoutes.HandleFunc("/content", ContentAdd).Methods("POST")
	adminRoutes.HandleFunc("/content", ContentUpdate).Queries(
		"content_id", "{content_id}",
	).Methods("PUT")
	router.HandleFunc("/content/upload", UploadContentFile).Methods("POST")

	// counsellor
	adminRoutes.HandleFunc("/counsellor", CounsellorGet).Methods("GET")
	adminRoutes.HandleFunc("/counsellor", CounsellorUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

	// coupon
	adminRoutes.HandleFunc("/coupon", CouponGet).Methods("GET")
	adminRoutes.HandleFunc("/coupon", CouponAdd).Methods("POST")
	adminRoutes.HandleFunc("/coupon", CouponUpdate).Queries(
		"id", "{id}",
	).Methods("PUT")

	// listener
	adminRoutes.HandleFunc("/listener", ListenerGet).Methods("GET")
	adminRoutes.HandleFunc("/listener", ListenerUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// login
	adminRoutes.HandleFunc("/login", Login).Methods("GET")
	adminRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")

	// notification
	adminRoutes.HandleFunc("/notification", NotificationGet).Methods("GET")
	adminRoutes.HandleFunc("/notification", NotificationAdd).Methods("POST")

	// quote
	adminRoutes.HandleFunc("/quote", QuoteGet).Methods("GET")
	adminRoutes.HandleFunc("/quote", QuoteAdd).Methods("POST")
	adminRoutes.HandleFunc("/quote", QuoteDelete).Queries(
		"id", "{id}",
	).Methods("DELETE")

	// report
	adminRoutes.HandleFunc("/report", ReportGet).Queries(
		"id", "{id}",
	).Methods("GET")

	// therapist
	adminRoutes.HandleFunc("/therapist", TherapistGet).Methods("GET")
	adminRoutes.HandleFunc("/therapist", TherapistUpdate).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("PUT")
}
