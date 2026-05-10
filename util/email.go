package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	Model "salbackend/model"
	"strings"
	"time"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"
)

// SendEmail - send email using SES. now : true - send now without background workers
func SendEmail(title, body, email string, now bool) {
	if strings.Contains(title, "###") || strings.Contains(body, "###") { // check if mail variables are replaced
		return
	}

	// add data to mails
	mail := map[string]string{}
	mail["title"] = title
	mail["body"] = body
	mail["email"] = email
	if now {
		// set mail sent status as sent if now is true
		mail["status"] = CONSTANT.EmailSent
		sendSESMail(title, body, email)
	} else {
		mail["status"] = CONSTANT.EmailInProgress
	}
	mail["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, CONSTANT.EmailsDigits, mail, "email_id")

}

func SendEmailWithCCBB(title, body, email string, ccemail []string, now bool) {
	if strings.Contains(title, "###") || strings.Contains(body, "###") { // check if mail variables are replaced
		return
	}

	// add data to mails
	mail := map[string]string{}
	mail["title"] = title
	mail["body"] = body
	mail["email"] = email
	if now {
		// set mail sent status as sent if now is true
		mail["status"] = CONSTANT.EmailSent
		sendSESMailWithCCBB(title, body, email, ccemail)
	} else {
		mail["status"] = CONSTANT.EmailInProgress
	}
	mail["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, CONSTANT.EmailsDigits, mail, "email_id")

}

func sendSESMailWithCCBB(title, body, email string, ccemail []string) {

	// old version of aws sdk

	// // start a new aws session
	// sess, err := session.NewSession()
	// if err != nil {
	// 	fmt.Println("failed to create session,", err)
	// 	return
	// }

	// // start a new ses session
	// svc := ses.New(sess, &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
	// 	Region:      aws.String("ap-south-1"), // ap-south-1
	// })

	// params := &ses.SendEmailInput{
	// 	Destination: &ses.Destination{ // Required
	// 		ToAddresses: []*string{
	// 			aws.String(email), // Required
	// 		},
	// 		CcAddresses: ccAddressess,
	// 	},
	// 	Message: &ses.Message{ // Required
	// 		Body: &ses.Body{ // Required
	// 			Html: &ses.Content{
	// 				Data:    aws.String(body), // Required
	// 				Charset: aws.String("UTF-8"),
	// 			},
	// 		},
	// 		Subject: &ses.Content{ // Required
	// 			Data:    aws.String(title), // Required
	// 			Charset: aws.String("UTF-8"),
	// 		},
	// 	},
	// 	Source: aws.String(CONFIG.FromEmailID),
	// }

	// //end email
	// output, err := svc.SendEmail(params)

	// new version of aws sdk

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(CONFIG.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		),
	)
	if err != nil {
		fmt.Println("config load failed: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	ccAddressess := []*string{}
	for _, e := range ccemail {
		ccAddressess = append(ccAddressess, aws.String(e))
	}

	input := &ses.SendEmailInput{
		Source: aws.String(CONFIG.FromEmailID),
		Destination: &types.Destination{
			ToAddresses: []string{email},
			CcAddresses: ccemail,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(title),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	output, err := client.SendEmail(ctx, input)
	if err != nil {
		fmt.Println("failed to send email: %w", err)
	}

	fmt.Println("Email sent:", *output.MessageId)

}

func sendSESMail(title, body, email string) {
	// // start a new aws session
	// sess, err := session.NewSession()
	// if err != nil {
	// 	fmt.Println("failed to create session,", err)
	// 	return
	// }

	// // start a new ses session
	// svc := ses.New(sess, &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
	// 	Region:      aws.String("ap-south-1"), // ap-south-1
	// })

	// params := &ses.SendEmailInput{
	// 	Destination: &ses.Destination{ // Required
	// 		ToAddresses: []*string{
	// 			aws.String(email), // Required
	// 		},
	// 	},
	// 	Message: &ses.Message{ // Required
	// 		Body: &ses.Body{ // Required
	// 			Html: &ses.Content{
	// 				Data:    aws.String(body), // Required
	// 				Charset: aws.String("UTF-8"),
	// 			},
	// 		},
	// 		Subject: &ses.Content{ // Required
	// 			Data:    aws.String(title), // Required
	// 			Charset: aws.String("UTF-8"),
	// 		},
	// 	},
	// 	Source: aws.String(CONFIG.FromEmailID),
	// }

	// //end email
	// output, err := svc.SendEmail(params)
	// fmt.Println(err, output.String())

	// new version of aws sdk
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(CONFIG.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		),
	)
	if err != nil {
		fmt.Println("config load failed: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Source: aws.String(CONFIG.FromEmailID),
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(title),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	output, err := client.SendEmail(ctx, input)
	if err != nil {
		fmt.Println("failed to send email: %w", err)
	}

	fmt.Println("Email sent:", *output.MessageId)
}

// SendEmail - send email using SES. now : true - send now without background workers
func SendEmailWithCcRefrence(title string, body string, emailfrom string, emailto, emailcc, emailbcc []string, todata string, now bool) {
	if strings.Contains(title, "###") || strings.Contains(body, "###") { // check if mail variables are replaced
		return
	}

	// add data to mails
	mail := map[string]string{}
	mail["title"] = title
	mail["body"] = body
	mail["email_from"] = emailfrom
	mail["email_to"] = todata
	if now {
		// set mail sent status as sent if now is true
		mail["status"] = CONSTANT.EmailSent
		sendSESMailForQualityCheck(title, body, emailfrom, emailto, emailcc, emailbcc)
	} else {
		mail["status"] = CONSTANT.EmailInProgress
	}
	mail["created_at"] = GetCurrentTime().String()
	DB.InsertWithUniqueID(CONSTANT.QualityCheckEmailTable, CONSTANT.EmailsDigits, mail, "email_id")

}

func sendSESMailForQualityCheck(title string, body string, emailfrom string, emailto, emailcc, emailbcc []string) {
	// // start a new aws session
	// sess, err := session.NewSession()
	// if err != nil {
	// 	fmt.Println("failed to create session,", err)
	// 	return
	// }

	// // start a new ses session
	// svc := ses.New(sess, &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
	// 	Region:      aws.String("ap-south-1"), // ap-south-1
	// })

	// params := &ses.SendEmailInput{
	// 	Destination: &ses.Destination{ // Required
	// 		CcAddresses:  emailcc,
	// 		ToAddresses:  emailto,
	// 		BccAddresses: emailbcc,
	// 	},
	// 	Message: &ses.Message{ // Required
	// 		Body: &ses.Body{ // Required
	// 			Html: &ses.Content{
	// 				Data:    aws.String(body), // Required
	// 				Charset: aws.String("UTF-8"),
	// 			},
	// 		},
	// 		Subject: &ses.Content{ // Required
	// 			Data:    aws.String(title), // Required
	// 			Charset: aws.String("UTF-8"),
	// 		},
	// 	},
	// 	Source: aws.String(emailfrom),
	// }

	// //end email
	// output, err := svc.SendEmail(params)
	// fmt.Println(err, output.String())

	// new version of aws sdk
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(CONFIG.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		),
	)
	if err != nil {
		fmt.Println("config load failed: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Source: aws.String(CONFIG.FromEmailID),
		Destination: &types.Destination{
			CcAddresses:  emailcc,
			ToAddresses:  emailto,
			BccAddresses: emailbcc,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(title),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	output, err := client.SendEmail(ctx, input)
	if err != nil {
		fmt.Println("failed to send email: %w", err)
	}

	fmt.Println("Email sent:", *output.MessageId)
}

func IsValidEmail(email_id string) string {
	_, err := mail.ParseAddress(email_id)
	if err != nil {
		return ""
	}

	return email_id
}

func SendEmailWithDocument(toemail string, body string, title string, documentB []Model.DocumentList) {
	recipients := []string{toemail}

	// Build email with attachment
	msg := gomail.NewMessage()
	msg.SetHeader("From", CONFIG.FromEmailID)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", body)

	for _, value := range documentB {
		documentInBytes, _ := downloadDocument(value.DocumentLink)
		msg.Attach(value.DocumentName+".pdf", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(documentInBytes)
			return err
		}))
	}

	// // Create a buffer to hold the raw email
	// var emailRaw bytes.Buffer
	// msg.WriteTo(&emailRaw)

	// // Create an AWS session and SES client
	// sess, err := session.NewSession()
	// if err != nil {
	// 	fmt.Println("failed to create session,", err)
	// 	return
	// }

	// // start a new ses session
	// svc := ses.New(sess, &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
	// 	Region:      aws.String("ap-south-1"), // ap-south-1
	// })

	// // Send the email with SES
	// input := &ses.SendRawEmailInput{
	// 	RawMessage: &ses.RawMessage{
	// 		Data: emailRaw.Bytes(),
	// 	},
	// 	Source: aws.String(CONFIG.FromEmailID),
	// 	Destinations: []*string{
	// 		aws.String(recipients[0]),
	// 	},
	// }
	// _, err = svc.SendRawEmail(input)
	// if err != nil {
	// 	fmt.Println("Error sending email:", err)
	// 	return
	// }

	// fmt.Println("Email sent successfully!")

	// new version of aws sdk

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// Context with timeout (recommended)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(CONFIG.AWSRegion), // set your region
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		),
	)
	if err != nil {
		fmt.Println("failed to load config: %w", err)
	}

	// Create SES client
	client := ses.NewFromConfig(cfg)

	// Prepare raw email input
	input := &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: emailRaw.Bytes(),
		},
		Source:       aws.String(CONFIG.FromEmailID),
		Destinations: recipients, // ✅ v2 uses []string (not []*string)
	}

	// Send email
	output, err := client.SendRawEmail(ctx, input)
	if err != nil {
		fmt.Println("error sending email: %w", err)
	}

	fmt.Println("Email sent successfully:", *output.MessageId)
}

func downloadDocument(url string) ([]byte, error) {
	// Create an HTTP client
	client := http.Client{}

	preSignedUrl := PreSignedS3URLToGetTheData(CONFIG.S3Bucket, url, CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, CONFIG.AWSRegion)
	_, endPointURL := GetBaseURLAndEndpointFromURL(preSignedUrl)
	url = endPointURL

	urlWithBaseURL := CONFIG.MediaURL + url
	// Send a GET request to the URL
	response, err := client.Get(urlWithBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error sending GET request: %w", err)
	}
	defer response.Body.Close()

	// Check for successful response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response body into a byte slice
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
