package handlers

import (
	"github.com/sal-org/sal-backend/models"
)

// GetAllAppointments returns all the appointments
func GetAllAppointments() ([]models.Appointment, error) {
	appointment := models.Appointment{}
	return []models.Appointment{appointment}, nil
}
