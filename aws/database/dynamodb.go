package database

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sal-org/sal-backend/models"
)

func getAllCounselors() ([]models.Counselor, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(dynamoDBEndpoint),
	}))

	// create a dynamodb instance
	db := dynamodb.New(sess)

	// create the api params
	params := &dynamodb.GetItemInput{
		TableName: aws.String(counselorsTable),
	}

	// read the item
	resp, err := db.GetItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	var counselors []models.Counselor
	err = dynamodbattribute.UnmarshalMap(resp.Item, &counselors)
	return counselors, err
}

func createAppointment(counselor models.Counselor, user models.User, duration time.Duration, time time.Time) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(dynamoDBEndpoint),
	}))

	// create a dynamodb instance
	db := dynamodb.New(sess)

	appoinment := models.Appointment{
		Counselor: counselor,
		Patient:   user,
		Duration:  duration,
		Time:      time,
		Status:    models.Scheduled,
	}

	// marshal the Appointment struct into an aws attribute value
	appoinmentAVMap, err := dynamodbattribute.MarshalMap(appoinment)
	if err != nil {
		fmt.Println("Cannot marshal Appointment into AttributeValue map")
		return err
	}

	// create the api params
	params := &dynamodb.PutItemInput{
		TableName: aws.String(appoinmentsTable),
		Item:      appoinmentAVMap,
	}

	// put the item
	resp, err := db.PutItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return err
	}

	fmt.Println(resp)
	return nil
}
