package handlers

import (
	"github.com/sal-org/sal-backend/models"
)

// GetAllCounselors returns all the counselors
func GetAllCounselors() ([]*models.Counselor, error) {
	counselor := &models.Counselor{
		FirstName: "John",
		LastName:  "Doe",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []*models.Counselor{counselor}, nil
}
