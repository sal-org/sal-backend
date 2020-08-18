package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
	"github.com/sal-org/sal-backend/models"
)

// GetAllCounselors returns all the counselors registered
func GetAllCounselors() (*[]models.Counselor, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return nil, err
	}

	// create a dynamodb instance
	db := dynamodb.New(sess)

	// create the api params
	params := &dynamodb.ScanInput{
		TableName: aws.String(counselorsTable),
	}

	// read the item
	resp, err := db.Scan(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	counselors := new([]models.Counselor)
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, counselors)
	return counselors, err
}

// CreateCounselor allows to create a counselor with given information
func CreateCounselor(req events.APIGatewayProxyRequest) (*models.Counselor, error) {
	var counselor models.Counselor
	if err := json.Unmarshal([]byte(req.Body), &counselor); err != nil {
		return nil, err
	}

	id := xid.New()
	counselor.Identifier = id.String()

	av, err := dynamodbattribute.MarshalMap(counselor)
	if err != nil {
		return nil, errors.New("unable to store counselor due to internal error")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return nil, err
	}

	// create a dynamodb instance
	db := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(counselorsTable),
	}

	_, err = db.PutItem(input)
	if err != nil {
		return nil, errors.New("unable to store counselor due to internal db error")
	}
	return &counselor, nil
}

// CreateAppointment allows to setup an appoinment between counselor and user for a specified duration and time
func CreateAppointment(counselor models.Counselor, user models.User, duration time.Duration, time time.Time) error {
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
