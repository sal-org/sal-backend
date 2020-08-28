package counselor

import (
	"fmt"
)

// Counselor will contain all essential details about a counselor
type Counselor struct {
	Identifier   string  `json:"id" dynamodbav:"id"`
	FirstName    string  `json:"firstName" dynamodbav:"firstName"`
	LastName     string  `json:"lastName" dynamodbav:"lastName"`
	Mobile       string  `json:"mobile" dynamodbav:"mobile"`
	Email        string  `json:"email" dynamodbav:"email"`
	Address1     string  `json:"address1" dynamodbav:"address1"`
	Address2     string  `json:"address2" dynamodbav:"address2"`
	Zipcode      int     `json:"zipcode" dynamodbav:"zipcode"`
	ThumbnailURL string  `json:"thumbnailUrl" dynamodbav:"thumbnailUrl"`
	PhotoURL     string  `json:"photoUrl" dynamodbav:"photoUrl"`
	Rate         float64 `json:"rate" dynamodbav:"rate"`
}

// Stringify returns custom value to be printed
func (c *Counselor) String() string {
	return fmt.Sprintf("%v, %v", c.FirstName, c.LastName)
}

// Repository defines the CRUD functionality for counselor entity
type Repository interface {
	// FetchAll returns all the counselors registered
	FetchAll() (*[]Counselor, error)

	// Save allows to save a counselor in persistent storage
	Save(counselor *Counselor) error
}
