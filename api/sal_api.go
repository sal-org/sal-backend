package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sal-org/sal-backend/api/routes"
)

func main() {
	r := mux.NewRouter()
	routes.ConfigureRoute(r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
