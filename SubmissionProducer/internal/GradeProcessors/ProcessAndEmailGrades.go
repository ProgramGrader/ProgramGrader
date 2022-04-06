package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"SubmissionProducer/internal/email"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

func ProcessAndEmailGrades(data common.AssignmentSubmissions, config common.AssignmentConfig) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess,
		aws.NewConfig().WithRegion(common.AWSRegion))

	var totalTestPassRate = 0.0
	var totalDocsPassRate = 0.0
	var totalTests = 0
	var totalDocs = 0

	// Put each student's grade to Dynamo
	for _, userSubmission := range data.Data {
		totalTests = totalTests + userSubmission.FailedTests.Total
		totalDocs = totalDocs + userSubmission.MissingDocs.Total

		input := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"CourseNumber-SemesterID": {
					S: aws.String(data.CourseNumber),
				},
				"AssignmentName-UserName": {
					S: aws.String(data.AssignmentName + "-" + userSubmission.UserName),
				},
				"course": {
					S: aws.String(data.CourseNumber),
				},
				"semester": {
					S: aws.String(common.GetCurrentSemesterIdentifier()),
				},
				"username": {
					S: aws.String(userSubmission.UserName),
				},
				"teacherID": {
					S: aws.String(data.Teacher.TeacherID),
				},
				"GithubUrl": {
					S: aws.String(userSubmission.GitURL),
				},
				"CommitTimes": {
					SS: aws.StringSlice(userSubmission.CommitTimes),
				},
			},
			TableName: aws.String(common.DYTableName),
		}

		if config.StudentTestsEnabled {

			var num = userSubmission.FailedTests.Total - len(userSubmission.FailedTests.TestNames)
			var denom = userSubmission.FailedTests.Total

			if denom != 0 {
				totalTestPassRate = totalTestPassRate + float64(num)/float64(denom)
			}

			temp := input.Item

			temp["PassingTests"] = &dynamodb.AttributeValue{
				S: aws.String(
					strconv.Itoa(num) + "/" + strconv.Itoa(denom),
				),
			}

			temp["CompletedAllTests"] = &dynamodb.AttributeValue{
				BOOL: aws.Bool(userSubmission.FailedTests.Completed),
			}

			input.Item = temp

		}

		if config.GradeDocs {

			var num = userSubmission.MissingDocs.Total - len(userSubmission.MissingDocs.TestNames)
			var denom = userSubmission.MissingDocs.Total

			if denom != 0 {
				totalDocsPassRate = totalDocsPassRate + float64(num)/float64(denom)
			}

			temp := input.Item

			temp["PassingDocs"] = &dynamodb.AttributeValue{
				S: aws.String(
					strconv.Itoa(num) + "/" + strconv.Itoa(denom),
				),
			}

			input.Item = temp

		}

		_, err := svc.PutItem(input)
		common.CheckIfError(err)
	} // for each submission

	// Put Assignment totals to Dynamo
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"CourseNumber-SemesterID": {
				S: aws.String(data.CourseNumber),
			},
			"AssignmentName-UserName": {
				S: aws.String(data.AssignmentName),
			},
			"PassingTestPercentage": {
				S: aws.String(fmt.Sprint(totalTestPassRate / float64(len(data.Data)) * 100)),
			},
			"PassingDocPercentage": {
				S: aws.String(fmt.Sprint(totalDocsPassRate / float64(len(data.Data)) * 100)),
			},
			"DueDate": {
				S: aws.String(data.DueDate.UTC().String()),
			},
			"NumberOfSubmissions": {
				S: aws.String(string(rune(len(data.Data)))),
			},
			"NonCodeSubmissions": {
				BOOL: aws.Bool(config.NonCodeSubmissions),
			},
		},
		TableName: aws.String(common.DYTableName),
	}

	_, err := svc.PutItem(input)
	common.CheckIfError(err)

	// Email Professor
	email.AWSSendAssignmentGradeEmail(data)

}
