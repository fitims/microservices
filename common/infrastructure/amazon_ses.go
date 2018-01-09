package infrastructure

import (
	"bytes"
	"html/template"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	FROM_EMAIL = "auto@thathrowdown.com"
)

// MailMessage contains information about the email message
type MailMessage struct {
	From      string
	Subject   string
	PlainBody string
	HTMLBody  string
	To        string
}

// SendMail sends an email that is in the MailMessage struct
func (msg MailMessage) SendMail() error {

	sess, err := getS3Session()
	if err != nil {
		log.Println("error gettings AWS Session : ", err)
		return err
	}

	svc := ses.New(sess)

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(msg.To),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(msg.Subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(msg.HTMLBody),
				},
				Text: &ses.Content{
					Data: aws.String(msg.PlainBody),
				},
			},
		},
		Source: aws.String(FROM_EMAIL),
	}

	res, err := svc.SendEmail(params)
	if err != nil {
		log.Println("Error sending email : ", err.Error())
		return err
	}

	log.Println("Sent email. Response : ", res)
	return nil
}

// SendTemplatedMail sends an email using the specified template
func (msg MailMessage) SendTemplatedMail(templateName string, templateData interface{}) error {
	tpl, err := template.ParseFiles(templateName)
	if err != nil {
		log.Println("error parsing the email template : ", err.Error())
		return err
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, templateData); err != nil {
		log.Println("error executing the email template : ", err.Error())
		return err
	}

	msg.HTMLBody = buf.String()
	return msg.SendMail()
}
