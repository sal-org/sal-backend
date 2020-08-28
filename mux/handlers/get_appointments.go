package handlers

import (
	"github.com/sal-org/sal-backend/appointment"
)

// GetAllAppointments returns all the appointments
func GetAllAppointments() ([]appointment.Appointment, error) {
	appointment := appointment.Appointment{}
	return []models.Appointment{appointment}, nil
}
