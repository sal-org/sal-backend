package admin

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all admin routes with admin prefix
func LoadAdminRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/admin").Subrouter()

	// content
	adminRoutes.HandleFunc("/appointment", AppointmentGet).Methods("GET")

	// content
	adminRoutes.HandleFunc("/content", ContentGet).Methods("GET")
	adminRoutes.HandleFunc("/content", ContentAdd).Methods("POST")
	adminRoutes.HandleFunc("/content", ContentUpdate).Queries(
		"content_id", "{content_id}",
	).Methods("PUT")
	router.HandleFunc("/content/upload", UploadContentFile).Methods("POST")

	// coupon
	adminRoutes.HandleFunc("/coupon", CouponGet).Methods("GET")
	adminRoutes.HandleFunc("/coupon", CouponAdd).Methods("POST")
	adminRoutes.HandleFunc("/coupon", CouponUpdate).Queries(
		"id", "{id}",
	).Methods("PUT")

	// login
	adminRoutes.HandleFunc("/login", Login).Methods("GET")
	adminRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")

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

}
