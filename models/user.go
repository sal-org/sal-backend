package models

import (
	"fmt"
)

// User will contain all essential details about a user
type User struct {
	FirstName    string `json:"firstName" dynamodbav:"firstName"`
	LastName     string `json:"lastName" dynamodbav:"lastName"`
	Mobile       string `json:"mobile" dynamodbav:"mobile"`
	Address1     string `json:"address1" dynamodbav:"address1"`
	Address2     string `json:"address2" dynamodbav:"address2"`
	Zipcode      int    `json:"zipcode" dynamodbav:"zipcode"`
	ThumbnailURL string `json:"thumbnailUrl" dynamodbav:"thumbnailUrl"`
	PhotoURL     string `json:"photoUrl" dynamodbav:"photoUrl"`
}

// Stringify returns custom value to be printed
func (u *User) String() string {
	return fmt.Sprintf("%v, %v", u.FirstName, u.LastName)
}
