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

const (
	counselorsTable = "Counselors"
)

// CounselorsRepository provides all the repository implementation for counselor repository interface
type CounselorsRepository struct {
	DB *dynamodb.DynamoDB
}

// FetchAll returns all the counselors from dynamodb
func (rep *CounselorsRepository) FetchAll() (*[]counselor.Counselor, error) {
	// create the api params
	params := &dynamodb.ScanInput{
		TableName: aws.String(counselorsTable),
	}

	// read the item
	resp, err := rep.DB.Scan(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return nil, err
	}

	counselors := new([]counselor.Counselor)
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, counselors)
	return counselors, err
}

// Save method stores the counselor in dynamodb
func (rep *CounselorsRepository) Save(counselor *counselor.Counselor) error {
	av, err := dynamodbattribute.MarshalMap(counselor)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(counselorsTable),
	}

	_, err = rep.DB.PutItem(input)
	return err
}

// CreateCounselor allows application to create a counselor from the gateway request
func (rep *CounselorsRepository) CreateCounselor(req events.APIGatewayProxyRequest) (*counselor.Counselor, error) {
	var counselor counselor.Counselor
	if err := json.Unmarshal([]byte(req.Body), &counselor); err != nil {
		return nil, err
	}

	id := xid.New()
	counselor.Identifier = id.String()

	err := rep.Save(&counselor)
	return &counselor, err
}
