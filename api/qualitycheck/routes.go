package qualitycheck

import "github.com/gorilla/mux"

func LoadQualityCheckRoutes(router *mux.Router) {
	qualityCheckRoutes := router.PathPrefix("/qualitycheck").Subrouter()

	//login
	qualityCheckRoutes.HandleFunc("/login", Login).Methods("GET")
	qualityCheckRoutes.HandleFunc("/refresh-token", RefreshToken).Methods("GET")
}
