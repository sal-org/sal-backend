package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func signUpEventHandler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	// ideally, we would like to add some checks, before allowing user to be registered
	event.Response.AutoConfirmUser = true
	return event, nil
}

func main() {
	lambda.Start(signUpEventHandler)
}
