package repositories

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sal-org/sal-backend/appointment"
	"github.com/sal-org/sal-backend/counselor"
	"github.com/sal-org/sal-backend/user"
)

const (
	appointmentsTable = "Appointments"
)

// AppointmentsRepository provides all the repository implementation for appoinment repository interface
type AppointmentsRepository struct {
	DB *dynamodb.DynamoDB
}

// FetchAll returns all the appoinments for the given counselor
func (r AppointmentsRepository) FetchAll(counselor *counselor.Counselor) (*[]appointment.Appointment, error) {
	return nil, nil
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
	return nil
}

// CreateAppointment method creates an appoinment between counselor and user for a specified duration and time
func (r AppointmentsRepository) createAppointment(counselor *counselor.Counselor, user *user.User, duration time.Duration, time time.Time) error {
	appointment := appointment.Appointment{
		Counselor: counselor,
		Patient:   user,
		Duration:  duration,
		Time:      time,
		Status:    appointment.Scheduled,
	}

	return r.Save(&appointment)
}
