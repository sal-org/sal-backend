package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
	"github.com/sal-org/sal-backend/counselor"
)

// CounselorsRepository provides all the repository implementation for counselor repository interface
type CounselorsRepository struct {
	DB *dynamodb.DynamoDB
}

// CounselorsRepositoryError defines all the errors related to counselor repository operations
type CounselorsRepositoryError struct {
	ErrorMsg string
}

func (e *CounselorsRepositoryError) Error() string {
	return fmt.Sprintf("counselor repository error - %v", e.ErrorMsg)
}

// FetchAll returns all the counselors from dynamodb
func (rep CounselorsRepository) FetchAll() (*[]counselor.Counselor, error) {
	// create the api params
	
	primaryKey := "COUNSELOR#";
	params := &dynamodb.QueryInput{
		TableName: aws.String(doctorsAppointmentsTable),
		KeyConditionExpression: aws.String("begins_with(PrimaryKey , :pk) AND begins_with(SortKey, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(primaryKey),
			},
			":sk": {
				S: aws.String("COUNSELOR#"),
			},
		},	
	}
	// read the item
	resp, err := rep.DB.Query(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	counselors := new([]counselor.Counselor)
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, counselors)
	return counselors, err
}

// Save method stores the counselor in dynamodb
func (rep CounselorsRepository) Save(counselor *counselor.Counselor) error {
	av, err := dynamodbattribute.MarshalMap(counselor)
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

// CreateCounselor allows application to create a counselor from the gateway request
func (rep CounselorsRepository) CreateCounselor(req events.APIGatewayProxyRequest) (counselor.Counselor, error) {
	var counselor counselor.Counselor
	id := xid.New()
	
	if err := json.Unmarshal([]byte(req.Body), &counselor); err != nil {
		return counselor , err
	}

	counselor.Identifier = id.String()

	email, ok1 := req.QueryStringParameters["email"]
	lastName, ok2 := req.QueryStringParameters["lastName"]
	name, ok3 := req.QueryStringParameters["firstName"]
	
	if !ok1 || !ok2 || !ok3  {
		return counselor , &CounselorsRepositoryError{ErrorMsg: "valid counselor details are not found"}
	}

	
	counselor.PrimaryKey = "COUNSELOR#" + id.String()
	counselor.SortKey = "COUNSELOR#" + name + "#" +  lastName
	counselor.Identifier = id.String()
	counselor.FirstName = name
	counselor.LastName = lastName
	counselor.Email = email

	err := rep.Save(&counselor)
	return counselor, err
}
