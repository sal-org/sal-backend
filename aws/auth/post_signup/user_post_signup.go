package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sal-org/sal-backend/models"
)

// CustomError defines all the errors related to db operation
type CustomError struct {
	ErrorMsg string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("db error - %v", e.ErrorMsg)
}

func postUserConfirmationEventHandler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	err := createNewUser(event.Request)
	return event, err
}

func createNewUser(request events.CognitoEventUserPoolsPostConfirmationRequest) error {
	userStatus, ok := request.UserAttributes["cognito:user_status"]
	if userStatus != "CONFIRMED" || !ok {
		return &CustomError{ErrorMsg: "user is not confirmed"}
	}

	// ideally email needs to be verified.. that flow needs to be figured out
	// emailVerified, _ := strconv.ParseBool(request.UserAttributes["email_verified"])
	// if !emailVerified {
	// 	return &CustomError{ErrorMsg: "user email is not verified"}
	// }

	email, ok1 := request.UserAttributes["email"]
	lastName, ok2 := request.UserAttributes["family_name"]
	name, ok3 := request.UserAttributes["name"]
	id, ok4 := request.UserAttributes["sub"]

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return &CustomError{ErrorMsg: "valid user details are not found"}
	}

	var user models.User

	user.Identifier = id
	user.FirstName = name
	user.LastName = lastName
	user.Email = email

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return err
	}

	// create a dynamodb instance
	db := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Users"),
	}

	_, err = db.PutItem(input)
	return err
}

func main() {
	lambda.Start(postUserConfirmationEventHandler)
}
