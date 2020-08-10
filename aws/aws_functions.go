package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sal-org/sal-backend/handlers"
	"github.com/sal-org/sal-backend/models"
)

// LambdaFunction contains details about the function to be executed
type LambdaFunction struct {
	Name string `json:"name"`
}

// HandleLambdaRequest receives the context and LambdaFunction object to handle the event
func HandleLambdaRequest(ctx context.Context, name LambdaFunction) ([]models.Counselor, error) {
	counselors, err := handlers.GetAllCounselors()
	return counselors, err
}

func main() {
	lambda.Start(HandleLambdaRequest)
}
