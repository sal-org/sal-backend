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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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

func sendSESMail(title, body, email string) {
	// start a new aws session
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	// start a new ses session
	svc := ses.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		Region:      aws.String("ap-south-1"),
	})

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			ToAddresses: []*string{
				aws.String(email), // Required
			},
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data:    aws.String(body), // Required
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(title), // Required
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(CONFIG.FromEmailID),
	}

	//end email
	output, err := svc.SendEmail(params)
	fmt.Println(err, output.String())
}

// SendEmail - send email using SES. now : true - send now without background workers
func SendEmailForQuality(title string, body string, emailfrom string, emailto, emailcc, emailbcc []*string, todata string, now bool) {
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

func sendSESMailForQualityCheck(title string, body string, emailfrom string, emailto, emailcc, emailbcc []*string) {
	// start a new aws session
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	// start a new ses session
	svc := ses.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		Region:      aws.String("ap-south-1"),
	})

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			CcAddresses:  emailcc,
			ToAddresses:  emailto,
			BccAddresses: emailbcc,
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data:    aws.String(body), // Required
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(title), // Required
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(emailfrom),
	}

	//end email
	output, err := svc.SendEmail(params)
	fmt.Println(err, output.String())
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

	// Create a buffer to hold the raw email
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// Create an AWS session and SES client
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	// start a new ses session
	svc := ses.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccesKey, CONFIG.AWSSecretKey, ""),
		Region:      aws.String("ap-south-1"),
	})

	// Send the email with SES
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: emailRaw.Bytes(),
		},
		Source: aws.String(CONFIG.FromEmailID),
		Destinations: []*string{
			aws.String(recipients[0]),
		},
	}
	_, err = svc.SendRawEmail(input)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func downloadDocument(url string) ([]byte, error) {
	// Create an HTTP client
	client := http.Client{}

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
