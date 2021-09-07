package util

import (
	"fmt"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	DB "salbackend/database"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
