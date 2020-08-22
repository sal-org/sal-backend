package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
)

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	storeCognitoRequest(event.Request)

	event.Response.AutoConfirmUser = true

	return event, nil
}

// TestData struct store data received from user pool during pre authentication
type TestData struct {
	ID       string                                       `json:"id" dynamodbav:"id"`
	Response events.CognitoEventUserPoolsPreSignupRequest `json:"response" dynamodbav:"response"`
}

func storeCognitoRequest(response events.CognitoEventUserPoolsPreSignupRequest) {
	var test TestData

	id := xid.New()
	test.ID = id.String()
	test.Response = response

	av, err := dynamodbattribute.MarshalMap(test)
	if err != nil {
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return
	}

	// create a dynamodb instance
	db := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Test"),
	}

	_, _ = db.PutItem(input)
}

func main() {
	lambda.Start(handler)
}
