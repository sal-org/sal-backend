package repositories

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
	"github.com/sal-org/sal-backend/user"
	//"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	doctorsAppointmentsTable = "DoctorsAppointments"
)

// UsersRepositoryError defines all the errors related to user repository operations
type UsersRepositoryError struct {
	ErrorMsg string
}

func (e *UsersRepositoryError) Error() string {
	return fmt.Sprintf("user repository error - %v", e.ErrorMsg)
}

// UsersRepository provides all the repository implementation for user repository interface
type UsersRepository struct {
	DB *dynamodb.DynamoDB
}

// FetchUser returns a user from dynamodb with the given id
func (rep UsersRepository) FetchUser(id string) ([]user.User, error) {
	// create the api params
	//PK is USER#<userid> SK is USER#firstname#lastname
	primaryKey := "USER#" + id;
	params := &dynamodb.QueryInput{
		TableName: aws.String(doctorsAppointmentsTable),
		KeyConditionExpression: aws.String("PrimaryKey = :pk AND begins_with(SortKey, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(primaryKey),
			},
			":sk": {
				S: aws.String("USER#"),
			},
		},	
	}

	// read the item
	// TODO
	resp, err := rep.DB.Query(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	//user := new(user.User)
	//TODO
	//err = dynamodbattribute.UnmarshalMap(resp.Items, user)
	user := []user.User{}
    err = dynamodbattribute.UnmarshalListOfMaps(resp.Items,  &user)
	return user, err
}

// Save method stores the user inside dynamodb
func (rep UsersRepository) Save(user *user.User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(doctorsAppointmentsTable),
	}

	_, err = rep.DB.PutItem(input)
	return err
}

// CreateUser method creates an user from Cognito request and saves it inside dynamodb
func (rep UsersRepository) CreateUser(request events.CognitoEventUserPoolsPostConfirmationRequest) error {
	// ideally email needs to be verified.. that flow needs to be figured out
	// emailVerified, _ := strconv.ParseBool(request.UserAttributes["email_verified"])
	// if !emailVerified {
	// 	return &CustomError{ErrorMsg: "user email is not verified"}
	// }
    id := xid.New()
	email, ok1 := request.UserAttributes["email"]
	lastName, ok2 := request.UserAttributes["lastName"]
	name, ok3 := request.UserAttributes["firstName"]
	
	if !ok1 || !ok2 || !ok3  {
		return &UsersRepositoryError{ErrorMsg: "valid user details are not found"}
	}

	var user user.User

	user.PrimaryKey = "USER#" + id.String()
	user.SortKey = "USER#" + name + "#" +  lastName
	user.Identifier = id.String()
	user.FirstName = name
	user.LastName = lastName
	user.Email = email

	return rep.Save(&user)
}
