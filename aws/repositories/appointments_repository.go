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

// AppointmentsRepository provides all the repository implementation for appoinment repository interface
type AppointmentsRepository struct {
	DB *dynamodb.DynamoDB
}

// FetchAll returns all the appoinments for the given counselor , user
func (rep AppointmentsRepository) FetchAll(counselorId string, userId string) (*[]appointment.Appointment, error) {
	// create the api params
	//PK APPOINTMENT#<id> SK APPOINTMENT#COUNSELOR#<id>#USER#<id>
	primaryKey := "APPOINTMENT#";
	secKey := "APPOINTMENT#COUNSELOR#" + counselorId + "USER#" + userId
	params := &dynamodb.QueryInput{
		TableName: aws.String(doctorsAppointmentsTable),
		KeyConditionExpression: aws.String("begins_with(PrimaryKey , :pk) AND begins_with(SortKey, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(primaryKey),
			},
			":sk": {
				S: aws.String(secKey),
			},
		},	
	}

	// read the item
	result, err := rep.DB.Query(params)

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
		TableName: aws.String(doctorsAppointmentsTable),
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

func (r AppointmentsRepository) createAppointment(counselor_id string, user_id string, duration time.Duration, time time.Time, status string) error {
	//TODO
	appointment := appointment.Appointment{
		Counselor: counselor_id,
		User:   user_id,
		Duration:  duration,
		Time:      time,
		Status:    status,
	}

	return r.Save(&appointment)
}
