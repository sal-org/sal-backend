package models

import (
	"fmt"
)

// User will contain all essential details about a user
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mobile    string `json:"mobile"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Zipcode   int    `json:"zipcode"`
}

// Stringify returns custom value to be printed
func (u *User) String() string {
	return fmt.Sprintf("%v, %v", u.FirstName, u.LastName)
}
