package api

import (
	"encoding/json"
	"net/http"

	ClientAPI "salbackend/api/client"
	MiscellaneousAPI "salbackend/api/miscellaneous"

	"github.com/gorilla/mux"
)

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// for load balancer/beanstalk to know whether server/ec2 is healthy
	json.NewEncoder(w).Encode("ok")
}

// LoadRouter - get mux router with all the routes
func LoadRouter() *mux.Router {
	router := mux.NewRouter()

	ClientAPI.LoadClientRoutes(router)
	MiscellaneousAPI.LoadMiscellaneousRoutes(router)

	// Swagger
	sh := http.StripPrefix("/documentaion/swagger/", http.FileServer(http.Dir("./docs/")))
	router.PathPrefix("/documentaion/swagger/").Handler(sh)

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	return router
}
