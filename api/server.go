package api

import (
	"fmt"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/gorilla/mux"
)

// StartServer - start server using mux
func StartServer(lambda bool) {
	if lambda {
		// lambda router
		algnhsa.ListenAndServe(&WithCORS{LoadRouter()}, nil)
	} else {
		// ec2 router
		fmt.Println(http.ListenAndServe(":5000", &WithCORS{LoadRouter()}))
	}
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// cors configuration
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Amz-Date,X-Amz-Security-Token,Authorization,X-Api-Key,X-Requested-With,Accept,Access-Control-Allow-Methods,Access-Control-Allow-Origin,Access-Control-Allow-Headers")

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(res, req)
}

// WithCORS .
type WithCORS struct {
	r *mux.Router
}
