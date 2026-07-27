package miscellaneous

import "github.com/gorilla/mux"

// LoadMiscellaneousRoutes - load all miscellaneous routes with empty prefix
func LoadMiscellaneousRoutes(router *mux.Router) {

	// content
	router.HandleFunc("/content", Content).Methods("GET")
	router.HandleFunc("/content/name", GetContentUsedTitle).Queries(
		"content_name", "{content_name}",
		"type", "{type}",
	).Methods("GET")
	router.HandleFunc("/content/like", ContentLikeGet).Queries(
		"user_id", "{user_id}",
	).Methods("GET")
	router.HandleFunc("/content/like", ContentLikeAdd).Queries(
		"user_id", "{user_id}",
		"content_id", "{content_id}",
	).Methods("POST")
	router.HandleFunc("/cancellationreason", CancellationReason).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("PUT")
	router.HandleFunc("/content/like", ContentLikeDelete).Queries(
		"user_id", "{user_id}",
		"content_id", "{content_id}",
	).Methods("DELETE")

	router.HandleFunc("/content/count", IncreaseContentViewCount).Queries(
		"content_id", "{content_id}",
	).Methods("PUT")

	// categories of content
	router.HandleFunc("/content-category", ListContentCategory).Methods("GET")

	// user account delete
	router.HandleFunc("/delete-user", DeleteUserProfile).Queries(
		"user_id", "{user_id}",
		"type", "{type}",
	).Methods("DELETE")

	// user account restore
	router.HandleFunc("/restore-user-account", RestoreUserProfile).Methods("PUT")

	// counsellor content
	router.HandleFunc("/counsellor-content", ListCounsellorContent).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

	// counsellor content
	router.HandleFunc("/app-feedback", AppFeedback).Methods("POST")

	// notification change status
	router.HandleFunc("/notification-status", NotificationInactiveORActive).Queries(
		"user_id", "{user_id}",
	).Methods("PUT")

	// topic
	router.HandleFunc("/topic", ListTopic).Methods("GET")

	// language
	router.HandleFunc("/language", ListLanguage).Methods("GET")

	// meta
	router.HandleFunc("/meta", ListMeta).Methods("GET")

	// get counsellor record history
	router.HandleFunc("/get-last-history-for-new-counsellor-record", GetLastHistoryRecord).Methods("GET")

	// get counsellor record history
	router.HandleFunc("/get-last-history-new-version", GetCounsellorClientRecordForNewestVersion).Methods("GET")

	// new version counsellor record form
	router.HandleFunc("/counsellor-client-record-for-new-version", CounsellorClientRecordForNewestVersion).Methods("POST")

	// adsConent
	router.HandleFunc("/adscontents", AdsContent).Methods("GET")

	// counsellor record
	router.HandleFunc("/counsellor-record", GetCounsellorClientRecord).Queries(
		"client_id", "{client_id}",
	).Methods("GET")

	router.HandleFunc("/counsellor-record/check", CheckCounsellorClientRecord).Queries(
		"counsellor_id", "{counsellor_id}",
		"client_id", "{client_id}",
		"date", "{date}",
	).Methods("GET")

	router.HandleFunc("/counsellor-record-new-version/check", CheckGetCounsellorClientRecordForNewest).Queries(
		"appointment_id", "{appointment_id}",
	).Methods("GET")

	// router.HandleFunc("/counsellor-record", CounsellorClientRecord).Methods("POST")
	router.HandleFunc("/counsellor-record/category", GetCounsellorRecordFromMainCategory).Methods("GET")
	router.HandleFunc("/counsellor-record/sub_category", GetCounsellorRecordFromSubCategory).Methods("GET")
	router.HandleFunc("/document-list", GetDocumentList).Methods("GET")

	// login
	router.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	router.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")

	router.HandleFunc("/app_info", AppInfo).Methods("GET")

	// mood
	router.HandleFunc("/mood", ListMood).Methods("GET")

	// rating types
	router.HandleFunc("/rating-type", ListRatingType).Methods("GET")

	// upload
	router.HandleFunc("/upload", UploadFile).Methods("POST")

	// check access token
	router.HandleFunc("/check_access_token", CheckIfAccessTokenExpired).Methods("GET")

	// upload file using
	router.HandleFunc("/pre_signed_url", PreSignedS3URLToUpload).Queries(
		"fileName", "{fileName}",
	).Methods("GET")

}
