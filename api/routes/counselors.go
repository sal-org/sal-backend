package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sal-org/sal-backend/constants"
	"github.com/sal-org/sal-backend/handlers"
)

// RouteBuilder interface defines the methods to be exposed by an object responsible for an api endpoint
type RouteBuilder interface {
	ConfigureRoute(r *mux.Router)
}

// CounselorsRouteBuilder manages the /counselors endpoint
type CounselorsRouteBuilder struct{}

// ConfigureRoute method sets up the router
func ConfigureRoute(r *mux.Router) {
	r.HandleFunc("/counselors", getAllCounselors).Methods("GET")
}

func getAllCounselors(w http.ResponseWriter, r *http.Request) {
	counselors, err := handlers.GetAllCounselors()

	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Add(constants.RequestContentType, constants.JSONResponse)
	json.NewEncoder(w).Encode(counselors)
}
