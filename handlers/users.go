package handlers

import (
	"fmt"
)

// User will contain all essential details about a user
type User struct {
	FirstName string
	LastName  string
	Mobile    string
	Address1  string
	Address2  string
	Zipcode   int
}

// Stringify returns custom value to be printed
func (u *User) String() string {
	return fmt.Sprintf("%v, %v", u.FirstName, u.LastName)
}

// GetAllUsers returns all the users
func GetAllUsers() ([]*User, error) {
	user := &User{
		FirstName: "Mary",
		LastName:  "Allen",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []*User{user}, nil
}
