package handlers

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
