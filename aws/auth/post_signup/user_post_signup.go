package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sal-org/sal-backend/aws/database"
	"github.com/sal-org/sal-backend/aws/repositories"
)

func postUserConfirmationEventHandler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	db, err := database.CreateDB()

	if err != nil {
		return event, err
	}

	repository := repositories.UsersRepository{DB: db}

	err = repository.CreateUser(event.Request)
	return event, err
}

func main() {
	lambda.Start(postUserConfirmationEventHandler)
}
