package handlers

import (
	"github.com/sal-org/sal-backend/models"
)

// GetAllUsers returns all the users
func GetAllUsers() ([]models.User, error) {
	user := models.User{
		FirstName: "Mary",
		LastName:  "Allen",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []models.User{user}, nil
}
