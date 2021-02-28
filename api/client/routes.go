package client

import "github.com/gorilla/mux"

// LoadClientRoutes - load all client routes with client prefix
func LoadClientRoutes(router *mux.Router) {
	clientRoutes := router.PathPrefix("/client").Subrouter()

	// search
	clientRoutes.HandleFunc("/search", ListSearch).Methods("GET")

	// counsellor
	clientRoutes.HandleFunc("/counsellor", CounsellorProfile).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/slots", CounsellorSlots).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/prices", CounsellorPrices).Queries(
		"counsellor_id", "{counsellor_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/counsellor/order", CounsellorOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/counsellor/paymentcomplete", CounsellorOrderPaymentComplete).Methods("POST")

	// listener
	clientRoutes.HandleFunc("/listener", ListenerProfile).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/slots", ListenerSlots).Queries(
		"listener_id", "{listener_id}",
	).Methods("GET")
	clientRoutes.HandleFunc("/listener/order", ListenerOrderCreate).Methods("POST")
	clientRoutes.HandleFunc("/listener/paymentcomplete", ListenerOrderPaymentComplete).Methods("POST")

	// payment
	clientRoutes.HandleFunc("/paymentcomplete", CounsellorOrderPaymentComplete).Methods("POST")
}
