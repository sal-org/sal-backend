package admin

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all admin routes with admin prefix
func LoadAdminRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/admin").Subrouter()

	// content
	adminRoutes.HandleFunc("/appointment", AppointmentGet).Methods("GET")
	adminRoutes.HandleFunc("/appointment/refund", AppointmentRefund).Queries(
		"appointment_id", "{appointment_id}",
		"refund_amount", "{refund_amount}",
	).Methods("PUT")

	// availability
	adminRoutes.HandleFunc("/availability", AvailabilityGet).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	adminRoutes.HandleFunc("/availability", AvailabilityUpdate).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("PUT")

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
	adminRoutes.HandleFunc("/content/upload", UploadContentFile).Methods("POST")
	adminRoutes.HandleFunc("/content/upload", PreSignedS3URLToUpload).Queries(
		"fileName", "{fileName}",
	).Methods("GET")

	// partner
	adminRoutes.HandleFunc("/partner", PartnerGet).Methods("GET")
	adminRoutes.HandleFunc("/partner", PartnerAdd).Methods("POST")
	adminRoutes.HandleFunc("/partner", PartnerUpdate).Queries(
		"id", "{id}",
	).Methods("PUT")
	adminRoutes.HandleFunc("/partner/address", PartnerAddressGet).Methods("GET")
	adminRoutes.HandleFunc("/partner/address", PartnerAddressAdd).Methods("POST")
	adminRoutes.HandleFunc("/partner/address", PartnerAddressUpdate).Queries(
		"id", "{id}",
	).Methods("PUT")

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

	// event
	adminRoutes.HandleFunc("/event", EventGet).Methods("GET")
	adminRoutes.HandleFunc("/event/upload", UploadEventFile).Methods("POST")
	adminRoutes.HandleFunc("/event", EventUpdate).Queries(
		"order_id", "{order_id}",
	).Methods("PUT")
	adminRoutes.HandleFunc("/event/book", EventBookGet).Methods("GET")

	// listener
	adminRoutes.HandleFunc("/listener", ListenerGet).Methods("GET")
	adminRoutes.HandleFunc("/listener", ListenerUpdate).Queries(
		"listener_id", "{listener_id}",
	).Methods("PUT")

	// login
	adminRoutes.HandleFunc("/login", Login).Methods("GET")
	adminRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")

	// add user role
	adminRoutes.HandleFunc("/profileRole", AddProfileForUsers).Methods("POST")
	adminRoutes.HandleFunc("/profileRole", UpdateProfileForUsers).Queries("id", "{id}").Methods("PUT")
	adminRoutes.HandleFunc("/profileRole", UserProfileGet).Methods("GET")
	adminRoutes.HandleFunc("/profileRole", UserProfileDelete).Queries(
		"role_id", "{role_id}",
	).Methods("DELETE")

	// add user attch role
	adminRoutes.HandleFunc("/userRole", AttachPermission).Methods("POST")
	adminRoutes.HandleFunc("/userRole", UpdateUserPermission).Queries("id", "{id}").Methods("PUT")
	adminRoutes.HandleFunc("/userRole", UserPermissionGet).Methods("GET")
	adminRoutes.HandleFunc("/userRole", UserPermissionDelete).Queries(
		"user_id", "{user_id}",
	).Methods("DELETE")

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
