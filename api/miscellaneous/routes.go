package miscellaneous

import "github.com/gorilla/mux"

// LoadMiscellaneousRoutes - load all miscellaneous routes with empty prefix
func LoadMiscellaneousRoutes(router *mux.Router) {

	// content
	router.HandleFunc("/content", Content).Methods("GET")
	router.HandleFunc("/content/like", ContentLikeGet).Queries(
		"user_id", "{user_id}",
	).Methods("GET")
	router.HandleFunc("/content/like", ContentLikeAdd).Queries(
		"user_id", "{user_id}",
		"content_id", "{content_id}",
	).Methods("POST")
	router.HandleFunc("/content/like", ContentLikeDelete).Queries(
		"user_id", "{user_id}",
		"content_id", "{content_id}",
	).Methods("DELETE")

	// categories of content
	router.HandleFunc("/content-category", ListContentCategory).Methods("GET")

	// counsellor account delete
	router.HandleFunc("/delete-user", DeleteUserProfile).Queries(
		"user_id", "{user_id}",
		"type", "{type}",
	).Methods("DELETE")

	// counsellor content
	router.HandleFunc("/counsellor-content", ListCounsellorContent).Queries(
		"therapist_id", "{therapist_id}",
	).Methods("GET")

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

	// login
	router.HandleFunc("/sendotp", SendOTP).Queries(
		"phone", "{phone}",
	).Methods("GET")
	router.HandleFunc("/verifyotp", VerifyOTP).Queries(
		"phone", "{phone}",
		"otp", "{otp}",
	).Methods("GET")

	// mood
	router.HandleFunc("/mood", ListMood).Methods("GET")

	// rating types
	router.HandleFunc("/rating-type", ListRatingType).Methods("GET")

	// upload
	router.HandleFunc("/upload", UploadFile).Methods("POST")

}
