package main

import (
	"fmt"
	"sal-backend/handlers"
	// "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	counselors, error := handlers.GetAllCounselors()
	users, e := handlers.GetAllUsers()
	fmt.Println(counselors, error)
	fmt.Println(users, e)
	// lambda.Start(handlers.GetAllCounselors)
}
