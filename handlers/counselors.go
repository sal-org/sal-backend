package handlers

import (
	"fmt"
)

// Counselor will contain all essential details about a counselor
type Counselor struct {
	FirstName string
	LastName  string
	Mobile    string
	Address1  string
	Address2  string
	Zipcode   int
}

// Stringify returns custom value to be printed
func (c *Counselor) String() string {
	return fmt.Sprintf("%v, %v", c.FirstName, c.LastName)
}

// GetAllCounselors returns all the counselors
func GetAllCounselors() ([]*Counselor, error) {
	counselor := &Counselor{
		FirstName: "John",
		LastName:  "Doe",
		Mobile:    "23823728372",
		Address1:  "1 Market St",
		Address2:  "Suite 3000",
		Zipcode:   94105}

	return []*Counselor{counselor}, nil
}
