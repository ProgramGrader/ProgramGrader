package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"SubmissionProducer/internal/configProcessors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"runtime"
	"strconv"
	"time"
)

func GetCurrentlyDueAssignments() []common.DueAssignment {

	configDates := GetAssignmentsAndDueDates()

	var retArray []common.DueAssignment

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess,
		aws.NewConfig().WithRegion(common.AWSRegion))

	CurrentTime := time.Now()

	for CCIndex, CurrentClasses := range configDates {
		for CDIndex, ClassDetails := range CurrentClasses {
			for ATIndex, AssignmentTime := range ClassDetails.dueDates {

				FutureAssignmentTime := AssignmentTime

				FutureAssignmentTime.AddDate(0, 1, 28)
				if CurrentTime.After(AssignmentTime) && !CurrentTime.Before(FutureAssignmentTime) {

					t := map[string]*dynamodb.AttributeValue{
						"CourseNumber-SemesterID": {
							S: aws.String(CDIndex + "-" + common.GetCurrentSemesterIdentifier()),
						},
						"AssignmentName-UserName": {
							S: aws.String("Assignment" + strconv.Itoa(ATIndex+1)),
						},
					}

					// If scores are posted in dynamodb then remove from return list
					// Query: P: CourseId-SemesterID s:AssignmentName
					result, err := svc.GetItem(&dynamodb.GetItemInput{
						TableName: aws.String(common.DYTableName),
						Key:       t,
					})

					common.CheckIfErrorWithMessage(err, "error calling GetItem")

					if result.Item == nil {

						var AssignmentLocation string
						if runtime.GOOS == "windows" {
							AssignmentLocation = fmt.Sprintf("\\src\\%v\\%v\\current\\", CCIndex, "Assignment"+strconv.Itoa(ATIndex+1))
						} else {
							AssignmentLocation = fmt.Sprintf("/src/%v/%v/current/", CCIndex, "Assignment"+strconv.Itoa(ATIndex+1))
						}

						temp := common.DueAssignment{
							CourseType:        CCIndex,
							CourseIDSemID:     CDIndex + "-" + common.GetCurrentSemesterIdentifier(),
							CourseTeacher:     ClassDetails.teacher,
							AssignmentDueDate: AssignmentTime,
							AssignmentName:    "Assignment" + strconv.Itoa(ATIndex+1),
							AssignmentNumber:  ATIndex + 1,
							AssignmentConfig:  configProcessors.ReadAssignmentConfig(common.ConfigPath() + AssignmentLocation),
						}
						retArray = append(retArray, temp)
					} else {
						common.Info("\t" + CCIndex + "-" + CDIndex + ": Assignment" + strconv.Itoa(ATIndex+1) + "  Has Already Been Graded.")
					}
				}
			}
		}
	}

	return retArray
}
