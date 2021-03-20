package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	API "salbackend/api"
	CONFIG "salbackend/config"
	DATABASE "salbackend/database"

	_ "salbackend/docs"
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

	if strings.EqualFold(os.Getenv("lambda"), "1") {
		API.StartServer(true)
	} else {
		API.StartServer(false)
	}
}
