package api

import (
	"encoding/json"
	"net/http"

	AdminAPI "salbackend/api/admin"
	ClientAPI "salbackend/api/client"
	CounsellorAPI "salbackend/api/counsellor"
	CronAPI "salbackend/api/cron"
	ListenerAPI "salbackend/api/listener"
	MiscellaneousAPI "salbackend/api/miscellaneous"
	TherapistAPI "salbackend/api/therapist"

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
	CounsellorAPI.LoadCounsellorRoutes(router)
	ListenerAPI.LoadListenerRoutes(router)
	TherapistAPI.LoadTherapistRoutes(router)
	MiscellaneousAPI.LoadMiscellaneousRoutes(router)
	CronAPI.LoadCronRoutes(router)
	AdminAPI.LoadAdminRoutes(router)

	// Swagger
	sh := http.StripPrefix("/documentaion/swagger/", http.FileServer(http.Dir("./docs/")))
	router.PathPrefix("/documentaion/swagger/").Handler(sh)

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	return router
}
