package main

import (
	"fmt"

	"github.com/sal-org/sal-backend/aws/handlers"
)

func main() {
	counselors, error := handlers.GetAllCounselors()
	users, e := handlers.GetAllUsers()
	appointments, err := handlers.GetAllAppointments()

	fmt.Println(counselors, error)
	fmt.Println(users, e)
	fmt.Println(appointments, err)

	// lambda.Start(handlers.GetAllCounselors)
}
