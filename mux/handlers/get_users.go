package handlers

import (
	"github.com/sal-org/sal-backend/user"
)

// GetAllUsers returns all the users
func GetAllUsers() ([]user.User, error) {
	user := user.User{
		FirstName: "Mary",
		LastName:  "Allen",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []user.User{user}, nil
}
