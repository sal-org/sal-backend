package models

import (
	"fmt"
)

const (
	// Bad indicates that the patient not satisfied at all
	Bad = iota

	// Ok indicates that the patient felt not so impressed
	Ok

	// Good indicates that the patient felt satisfied
	Good

	// VeryGood indicates that the patient felt very happy
	VeryGood

	// Excellent indicates that the patient felt the session went extremely well and would likely see the counselor again
	Excellent
)

// Feedback captures information regarding how session went
type Feedback struct {
	Appointment Appointment `json:"appoinment" dynamodbav:"appoinment"`
	Rating      int         `json:"rating" dynamodbav:"rating"`
}

// Stringify returns custom value to be printed
func (f *Feedback) String() string {
	return fmt.Sprintf("Counselor %v, Patient %v, Rating %v", f.Appointment.Counselor.FirstName, f.Appointment.Patient.FirstName, f.Rating)
}
