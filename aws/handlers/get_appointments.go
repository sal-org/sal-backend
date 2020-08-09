package handlers

import (
	"github.com/sal/backend/models"
)

// GetAllAppointments returns all the appointments
func GetAllAppointments() ([]*Appointment, error) {
	appointment := &Appointment{}

	return []*Appointment{appointment}, nil
}
