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

func handler(event events.CognitoEventUserPoolsPreAuthentication) (events.CognitoEventUserPoolsPreAuthentication, error) {
	storePreAuthCognitoRequest(event.Request)
	return event, nil
}

// AuthTestData struct store data received from user pool during pre authentication
type AuthTestData struct {
	ID       string                                               `json:"id" dynamodbav:"id"`
	Response events.CognitoEventUserPoolsPreAuthenticationRequest `json:"response" dynamodbav:"response"`
}

func storePreAuthCognitoRequest(response events.CognitoEventUserPoolsPreAuthenticationRequest) {
	var test AuthTestData

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
