package main

import (
	"math/rand"
	"time"

	API "salbackend/api"
	CONFIG "salbackend/config"
	DATABASE "salbackend/database"

	_ "salbackend/docs"

	"github.com/akrylysov/algnhsa"
)

// @title SAL Backend API
// @version 1.0
// @description This is a api for SAL client/listener/counsellor APIs
// @schemes https
// @host hwmpf9h476.execute-api.ap-south-1.amazonaws.com
// @BasePath /prod
func main() {

	rand.Seed(time.Now().UnixNano()) // seed for random generator

	CONFIG.LoadConfig()
	DATABASE.ConnectDatabase()

	// ec2 router
	// fmt.Println(http.ListenAndServe(":5000", &WithCORS{API.LoadRouter()}))

	// lambda router
	algnhsa.ListenAndServe(&WithCORS{API.LoadRouter()}, nil)
}
