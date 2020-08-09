package models

import (
	"fmt"
)

// Appointment will contain all essential details about a appointment
type Appointment struct {
	Duration  float64   `json:"duration"`
	Cost      float64   `json:"cost"`
	Counselor Counselor `json:"counselor"`
	Patient   User      `json:"patient"`
	Date      string    `json:"date"`
	Time      string    `json:"time"`
}

// Stringify returns custom value to be printed
func (c *Appointment) String() string {
	return fmt.Sprintf("%v, %v", c.Counselor, c.Patient)
}
