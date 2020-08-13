package models

import (
	"fmt"
	"time"
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
	Duration  time.Duration `json:"duration" dynamodbav:"duration"`
	Counselor Counselor     `json:"counselor" dynamodbav:"counselor"`
	Patient   User          `json:"patient" dynamodbav:"patient"`
	Time      time.Time     `json:"time" dynamodbav:"appointmentTime"`
	Status    int           `json:"status" dynamodbav:"status"`
}

// Stringify returns custom value to be printed
func (c *Appointment) String() string {
	return fmt.Sprintf("%v, %v", c.Counselor, c.Patient)
}
