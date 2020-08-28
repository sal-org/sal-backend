package auth

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// StoreCognitoEvent method store the provided cognito event in database
func StoreCognitoEvent(event CognitoEvent) {
	av, err := dynamodbattribute.MarshalMap(event)
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
		TableName: aws.String("CognitoEvents"),
	}

	_, _ = db.PutItem(input)
}
