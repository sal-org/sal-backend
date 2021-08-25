package miscellaneous

import "github.com/gorilla/mux"

// LoadMiscellaneousRoutes - load all miscellaneous routes with empty prefix
func LoadMiscellaneousRoutes(router *mux.Router) {
	// categories of content
	router.HandleFunc("/content-category", ListContentCategory).Methods("GET")

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
