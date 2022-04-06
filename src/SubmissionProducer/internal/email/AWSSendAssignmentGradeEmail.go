package email

import (
	"SubmissionProducer/internal/common"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"gopkg.in/gomail.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendEmailWithAttachment(sess *session.Session, data common.AssignmentSubmissions, fileName string) error {

	subject := "Completion Report for " + data.AssignmentName + " due on " + data.DueDate.Format("Jan-02-06 15:04")

	msg := gomail.NewMessage()

	var recipientEmail string
	if common.GetEnv("envEmail", "") != "" {
		recipientEmail = common.GetEnv("envEmail", data.Teacher.TeacherEmail)
	} else {
		recipientEmail = data.Teacher.TeacherEmail
	}

	msg.SetHeaders(map[string][]string{
		"From":             {msg.FormatAddress(common.Sender, "IUS AutoGrader")},
		"To":               {recipientEmail},
		"CC":               {},
		"BCC":              {common.DepHead},
		"Subject":          {subject},
		"ReplyToAddresses": {common.ReplyToAddresses},
	})

	msg.SetBody("text/html", CreateHtmlFromData(data))
	msg.SetBody("text/plain", CreateTextFromData(data))

	// Root of the project in GoLand
	// use absolute path in docker
	msg.Attach(fileName)

	var emailRaw bytes.Buffer
	_, err := msg.WriteTo(&emailRaw)

	common.CheckIfError(err)

	message := ses.RawMessage{
		Data: emailRaw.Bytes(),
	}

	sesClient := ses.New(sess)

	_, err = sesClient.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &message,
		Destinations: []*string{
			aws.String(recipientEmail),
		},
		Source:               aws.String(common.Sender),
		ConfigurationSetName: aws.String(common.ConfigurationSet),
	})

	return err
}

func SendEmailWithOutAttachment(sess *session.Session, data common.AssignmentSubmissions) error {

	subject := "Completion Report for " + data.CourseName + " " + data.AssignmentName + "(" + data.CourseNumber + ")" + " due on " + data.DueDate.Format("Jan-02-06 15:04")

	msg := gomail.NewMessage()

	var recipientEmail string
	if common.GetEnv("env", "") != "" {
		recipientEmail = common.GetEnv("envEmail", data.Teacher.TeacherEmail)
	} else {
		recipientEmail = data.Teacher.TeacherEmail
	}

	msg.SetHeaders(map[string][]string{
		"From":             {msg.FormatAddress(common.Sender, "IUS AutoGrader")},
		"To":               {recipientEmail},
		"Subject":          {subject},
		"ReplyToAddresses": {common.ReplyToAddresses},
	})

	msg.SetBody("text/html", CreateHtmlFromData(data))
	msg.SetBody("text/plain", CreateTextFromData(data))

	var emailRaw bytes.Buffer
	_, err := msg.WriteTo(&emailRaw)
	common.CheckIfErrorWithMessage(err, "Unable to write to email")

	message := ses.RawMessage{
		Data: emailRaw.Bytes(),
	}

	sesClient := ses.New(sess)

	_, err = sesClient.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &message,
		Destinations: []*string{
			aws.String(recipientEmail),
		},
		Source:               aws.String(common.Sender),
		ConfigurationSetName: aws.String(common.ConfigurationSet),
	})

	return err
}

func AWSSendAssignmentGradeEmail(data common.AssignmentSubmissions) {

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String(common.AWSRegion),
		},
	})
	common.CheckIfErrorWithMessage(err, "Failed to initialize new session.")

	err = SendEmailWithOutAttachment(sess, data)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	common.Info("\t Sent email for " + data.AssignmentName + " to " + data.Teacher.TeacherEmail + " without attachment sent successfully!")
}
