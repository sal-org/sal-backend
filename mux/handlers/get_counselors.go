package handlers

import (
	"github.com/sal-org/sal-backend/counselor"
)

// GetAllCounselors returns all the counselors
func GetAllCounselors() ([]counselor.Counselor, error) {
	counselor := counselor.Counselor{
		FirstName: "John",
		LastName:  "Doe",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []counselor.Counselor{counselor}, nil
}
