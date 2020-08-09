package models

import (
	"fmt"
)

// Counselor will contain all essential details about a counselor
type Counselor struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mobile    string `json:"mobile"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Zipcode   int    `json:"zipcode"`
}

// Stringify returns custom value to be printed
func (c *Counselor) String() string {
	return fmt.Sprintf("%v, %v", c.FirstName, c.LastName)
}
