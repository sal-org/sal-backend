package database

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	region = "AWS_REGION"
)

// CreateDB creates a dynamodb instance
func CreateDB() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(region)),
	})

	if err != nil {
		return nil, err
	}

	// create a dynamodb instance
	return dynamodb.New(sess), nil
}
