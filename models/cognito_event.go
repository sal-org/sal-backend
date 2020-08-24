package models

// CognitoEvent struct store data received from user pool during pre authentication
type CognitoEvent struct {
	ID      string      `json:"id" dynamodbav:"id"`
	Payload interface{} `json:"payload" dynamodbav:"payload"`
}
