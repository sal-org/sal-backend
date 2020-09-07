package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sal-org/sal-backend/aws/utils"
	"github.com/sal-org/sal-backend/constants"
	"github.com/sal-org/sal-backend/models"
)

const (
	region = "AWS_REGION"

	maxAllowedFileSize = 4 << 20

	userPhotosBucket = "sal-user-photos"

	pngContentType = "image/png"

	jpegContentType = "image/jpeg"

	photosFolder = "photos"
)

type photoUploadReq struct {
	FileName string `json:"fileName"`
}

// HandleUserPhotoRequest will handle all counselor gate api requests
func HandleUserPhotoRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		//TODO
		return utils.GatewayResponse(http.StatusOK, make(map[string]string))
	case "POST":
		session, err := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv(region)),
		})

		if err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String("Internal server error")})
		}

		var uploadReq photoUploadReq
		if err := json.Unmarshal([]byte(req.Body), &uploadReq); err != nil {
			return utils.GatewayResponse(http.StatusBadRequest, err)
		}

		s3c := s3.New(session)
		// fileName := fmt.Sprintf("%v/%v", photosFolder, uploadReq.FileName)
		req, putErr := s3c.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(userPhotosBucket),
			Key:    aws.String(uploadReq.FileName),
		})

		if putErr != nil {
			return utils.GatewayResponse(http.StatusBadRequest, fmt.Sprintf("Unable to create put object req for file %v", uploadReq.FileName))
		}

		presurl, presignerr := req.Presign(5 * time.Minute)
		if presignerr != nil {
			return utils.GatewayResponse(http.StatusBadRequest, presignerr)
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				constants.RequestContentType:  constants.JSONResponse,
				"Access-Control-Allow-Origin": "*",
			},
			Body: presurl,
		}, nil
	default:
		return utils.GatewayResponse(http.StatusBadRequest, models.ErrorBody{ErrorMsg: aws.String(constants.ErrorMethodNotAllowed)})
	}
}
