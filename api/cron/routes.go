package cron

import "github.com/gorilla/mux"

// LoadCronRoutes - load all cron routes with cron prefix
func LoadCronRoutes(router *mux.Router) {
	cronRoutes := router.PathPrefix("/cron").Subrouter()

	// availability
	cronRoutes.HandleFunc("/no-shows", NoShowUpdates).Methods("PATCH")
}
