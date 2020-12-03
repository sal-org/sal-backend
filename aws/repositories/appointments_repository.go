package repositories

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	//"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sal-org/sal-backend/appointment"
	//"github.com/sal-org/sal-backend/counselor"
	//"github.com/sal-org/sal-backend/user"
)

const (
	appointmentsTable = "Appointments"
)

// AppointmentsRepository provides all the repository implementation for appoinment repository interface
type AppointmentsRepository struct {
	DB *dynamodb.DynamoDB
}

// FetchAll returns all the appoinments for the given counselor , user
func (rep AppointmentsRepository) FetchAll(counselorId string, userId string) (*[]appointment.Appointment, error) {
	// create the api params
	
	params := &dynamodb.ScanInput{
		TableName: aws.String(appointmentsTable),
		ExpressionAttributeNames: map[string]*string{
			"#AT": aws.String("counselor"),
			"#ST": aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				S: aws.String(counselorId),
			},
			":b": {
				S: aws.String(userId),
			},
		},
		FilterExpression:     aws.String("(counselor = :a and user = :b)"),
		ProjectionExpression: aws.String("#ST, #AT"),
		
	}

	// read the item
	result, err := rep.DB.Scan(params)

	if err != nil {
		return nil, err
	}

	item := new([]appointment.Appointment)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

// Save allows application to store an appoinment in dynamodb
func (r AppointmentsRepository) Save(appointment *appointment.Appointment) error {
	// marshal the Appointment struct into an aws attribute value
	appoinmentAVMap, err := dynamodbattribute.MarshalMap(appointment)
	if err != nil {
		fmt.Println("Cannot marshal Appointment into AttributeValue map")
		return err
	}

	// create the api params
	params := &dynamodb.PutItemInput{
		TableName: aws.String(appointmentsTable),
		Item:      appoinmentAVMap,
	}

	_, e := r.DB.PutItem(params)
	return e
}

// Update allows user to change an appoinment
func (r AppointmentsRepository) Update(appointment *appointment.Appointment) error {
	//TODO
	return nil
}

// CreateAppointment allows application to create an appoinement using data from the gateway request
func (r AppointmentsRepository) CreateAppointment(req events.APIGatewayProxyRequest) (*appointment.Appointment, error) {
	//TODO
	return nil, nil
}

func (r AppointmentsRepository) createAppointment(counselor_id string, user_id string, duration time.Duration, time time.Time) error {
	appointment := appointment.Appointment{
		Counselor: counselor_id,
		User:   user_id,
		//Duration:  duration,
		//Time:      time,
		//Status:    appointment.Scheduled,
	}

	return r.Save(&appointment)
}
