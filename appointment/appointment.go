package appointment

import (
	"fmt"
	//"time"

	//"github.com/sal-org/sal-backend/counselor"
	//"github.com/sal-org/sal-backend/user"
)

const (
	// Scheduled indicate that the appoinment is scheduled
	Scheduled = iota

	// Active indicates that the appoinment is currently going on
	Active

	// Complete indicates that the appoinment is finished
	Complete

	// Cancelled indicated that the appoinment is cancelled
	Cancelled
)

// Appointment will contain all essential details about a appointment
type Appointment struct {
	//Duration  time.Duration        `json:"duration" dynamodbav:"duration"`
	Counselor string `json:"counselor" dynamodbav:"counselor"`
	User   string           `json:"user" dynamodbav:"user"`
	//Time      time.Time            `json:"time" dynamodbav:"appointmentTime"`
	Status    string                  `json:"status" dynamodbav:"status"`
	Id    string                  `json:"id" dynamodbav:"id"`
}

// Stringify returns custom value to be printed
func (a *Appointment) String() string {
	return fmt.Sprintf("%v, %v", a.Counselor, a.User)
}

// Repository defines the CRUD functionality for Appoinment entity
type Repository interface {
	// FetchAll returns all the appoinments for the given counselor
	FetchAll(counselorId string) (*[]Appointment, error)

	// Save allows to save an appoinment in persistent storage
	Save(appointment *Appointment) error

	// Update allows user to change the appoinment
	Update(appointment *Appointment) error
}
